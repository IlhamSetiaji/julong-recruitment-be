package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/service"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
)

type IApplicantHandler interface {
	ApplyJobPosting(ctx *gin.Context)
	GetApplicantsByJobPostingID(ctx *gin.Context)
	FindApplicantByJobPostingIDAndUserID(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	ExportApplicantsByJobPosting(ctx *gin.Context)
	GetApplicantsForCoverLetter(ctx *gin.Context)
}

type ApplicantHandler struct {
	Log                 *logrus.Logger
	Viper               *viper.Viper
	Validate            *validator.Validate
	UseCase             usecase.IApplicantUseCase
	UserProfileUseCase  usecase.IUserProfileUseCase
	UserHelper          helper.IUserHelper
	NotificationService service.INotificationService
}

func NewApplicantHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IApplicantUseCase,
	upUseCase usecase.IUserProfileUseCase,
	userHelper helper.IUserHelper,
	notificationService service.INotificationService,
) IApplicantHandler {
	return &ApplicantHandler{
		Log:                 log,
		Viper:               viper,
		Validate:            validate,
		UseCase:             useCase,
		UserProfileUseCase:  upUseCase,
		UserHelper:          userHelper,
		NotificationService: notificationService,
	}
}

func ApplicantHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IApplicantHandler {
	useCase := usecase.ApplicantUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	upUseCase := usecase.UserProfileUseCaseFactory(log, viper)
	userHelper := helper.UserHelperFactory(log)
	notificationService := service.NotificationServiceFactory(viper, log)
	return NewApplicantHandler(log, viper, validate, useCase, upUseCase, userHelper, notificationService)
}

// ApplyJobPosting apply job posting
//
// @Summary apply job posting
// @Description apply job posting
// @Tags Applicants
// @Accept json
// @Produce json
// @Param job_posting_id query string true "Job Posting ID"
// @Success 200 {object} response.ApplicantResponse
// @Security BearerAuth
// @Router /applicants/apply [post]
func (h *ApplicantHandler) ApplyJobPosting(ctx *gin.Context) {
	jobPostingID, err := uuid.Parse(ctx.Query("job_posting_id"))
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.ApplyJobPosting] error when parsing job_posting_id: %v", err)
		utils.BadRequestResponse(ctx, "job_posting_id is not a valid UUID", err)
		return
	}

	user, err := middleware.GetUser(ctx, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		return
	}
	userUUID, err := h.UserHelper.GetUserId(user)
	if err != nil {
		h.Log.Errorf("Error when getting user id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	userProfile, err := h.UserProfileUseCase.FindByUserID(userUUID)
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.ApplyJobPosting] error when getting user profile: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get user profile", err.Error())
		return
	}

	if userProfile == nil {
		h.Log.Errorf("[ApplicantHandler.ApplyJobPosting] user profile not found")
		utils.ErrorResponse(ctx, http.StatusNotFound, "User profile not found", "")
		return
	}

	if userProfile.Status != entity.USER_ACTIVE {
		h.Log.Errorf("[ApplicantHandler.ApplyJobPosting] user profile is not active")
		utils.ErrorResponse(ctx, http.StatusForbidden, "User profile is not active", "")
		return
	}

	applicant, err := h.UseCase.ApplyJobPosting(userProfile.ID, jobPostingID)
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.ApplyJobPosting] error when applying job posting: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to apply job posting", err.Error())
		return
	}

	if err := h.NotificationService.ApplicantAppliedNotification(userUUID.String()); err != nil {
		h.Log.Errorf("[ApplicantHandler.ApplyJobPosting] error when sending notification: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to send notification", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Successfully applied job posting", applicant)
}

// GetApplicantsByJobPostingID get applicants by job posting ID
//
// @Summary get applicants by job posting ID
// @Description get applicants by job posting ID
// @Tags Applicants
// @Accept json
// @Produce json
// @Param job_posting_id path string true "Job Posting ID"
// @Param order query string false "Order"
// @Param total query string false "Total"
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param search query string false "Search"
// @Param created_at query string false "Created At"
// @Success 200 {array} response.ApplicantResponse
// @Security BearerAuth
// @Router /applicants/job-posting/{job_posting_id} [get]
func (h *ApplicantHandler) GetApplicantsByJobPostingID(ctx *gin.Context) {
	jobPostingID, err := uuid.Parse(ctx.Param("job_posting_id"))
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.GetApplicantsByJobPostingID] error when parsing job_posting_id: %v", err)
		utils.BadRequestResponse(ctx, "job_posting_id is not a valid UUID", err)
		return
	}

	orderStr := ctx.Query("order")
	if orderStr == "" {
		orderStr = ""
	}

	totalStr := ctx.Query("total")
	if totalStr == "" {
		totalStr = "0"
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	createdAt := ctx.Query("created_at")
	if createdAt == "" {
		createdAt = "DESC"
	}

	sort := map[string]interface{}{
		"created_at": createdAt,
	}

	// order, err := strconv.Atoi(orderStr)
	// if err != nil {
	// 	h.Log.Errorf("[ApplicantHandler.GetApplicantsByJobPostingID] error when converting order to int: %v", err)
	// 	utils.BadRequestResponse(ctx, "order is not a valid integer", err)
	// 	return
	// }

	total, err := strconv.Atoi(totalStr)
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.GetApplicantsByJobPostingID] error when converting total to int: %v", err)
		utils.BadRequestResponse(ctx, "total is not a valid integer", err)
		return
	}

	applicants, totalData, err := h.UseCase.GetApplicantsByJobPostingID(jobPostingID, orderStr, total, page, pageSize, search, sort)
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.GetApplicantsByJobPostingID] error when getting applicants: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get applicants", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Successfully get applicants", gin.H{
		"applicants": applicants,
		"total":      totalData,
	})
}

// FindApplicantByJobPostingIDAndUserID find applicant by job posting ID and user ID
//
// @Summary find applicant by job posting ID and user ID
// @Description find applicant by job posting ID and user ID
// @Tags Applicants
// @Accept json
// @Produce json
// @Param job_posting_id path string true "Job Posting ID"
// @Success 200 {object} response.ApplicantResponse
// @Security BearerAuth
// @Router /applicants/me [get]
func (h *ApplicantHandler) FindApplicantByJobPostingIDAndUserID(ctx *gin.Context) {
	jobPostingID, err := uuid.Parse(ctx.Param("job_posting_id"))
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.FindApplicantByJobPostingIDAndUserID] error when parsing job_posting_id: %v", err)
		utils.BadRequestResponse(ctx, "job_posting_id is not a valid UUID", err)
		return
	}

	user, err := middleware.GetUser(ctx, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		return
	}
	userUUID, err := h.UserHelper.GetUserId(user)
	if err != nil {
		h.Log.Errorf("Error when getting user id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	applicant, err := h.UseCase.FindApplicantByJobPostingIDAndUserID(jobPostingID, userUUID)
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.FindApplicantByJobPostingIDAndUserID] error when finding applicant: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find applicant", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Successfully find applicant", applicant)
}

// FindByID find applicant by ID
//
// @Summary find applicant by ID
// @Description find applicant by ID
// @Tags Applicants
// @Accept json
// @Produce json
// @Param id path string true "Applicant ID"
// @Success 200 {object} response.ApplicantResponse
// @Security BearerAuth
// @Router /applicants/{id} [get]
func (h *ApplicantHandler) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.FindByID] error when parsing id: %v", err)
		utils.BadRequestResponse(ctx, "id is not a valid UUID", err)
		return
	}

	applicant, err := h.UseCase.FindByID(id)
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.FindByID] error when finding applicant: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find applicant", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Successfully find applicant", applicant)
}

// ExportApplicantsByJobPosting export applicants by job posting
//
// @Summary export applicants by job posting
// @Description export applicants by job posting
// @Tags Applicants
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param job_posting_id path string true "Job Posting ID"
// @Security BearerAuth
// @Router /applicants/job-posting/{job_posting_id}/export [get]
func (h *ApplicantHandler) ExportApplicantsByJobPosting(ctx *gin.Context) {
	jobPostingID, err := uuid.Parse(ctx.Param("job_posting_id"))
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.ExportApplicantsByJobPosting] error when parsing job_posting_id: %v", err)
		utils.BadRequestResponse(ctx, "job_posting_id is not a valid UUID", err)
		return
	}

	applicants, err := h.UseCase.GetApplicantsByJobPostingIDForExport(jobPostingID)
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.ExportApplicantsByJobPosting] error when exporting applicants: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export applicants", err.Error())
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export applicants", err.Error())
			return
		}
	}()

	f.SetSheetName("Sheet1", "Applicants")
	f.SetCellValue("Applicants", "A1", "Applicant Name")
	f.SetCellValue("Applicants", "B1", "Job Name")
	f.SetCellValue("Applicants", "C1", "Work Experience")
	f.SetCellValue("Applicants", "D1", "Applied Date")
	f.SetCellValue("Applicants", "E1", "Phone Number")
	f.SetCellValue("Applicants", "F1", "Expected Salary")
	f.SetCellValue("Applicants", "G1", "Current Salary")

	// Create a style for the header
	headerStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#00FF00"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export applicants", err.Error())
		return
	}

	// Set the style to the header
	f.SetCellStyle("Applicants", "A1", "A1", headerStyle)
	f.SetCellStyle("Applicants", "B1", "B1", headerStyle)
	f.SetCellStyle("Applicants", "C1", "C1", headerStyle)
	f.SetCellStyle("Applicants", "D1", "D1", headerStyle)
	f.SetCellStyle("Applicants", "E1", "E1", headerStyle)
	f.SetCellStyle("Applicants", "F1", "F1", headerStyle)
	f.SetCellStyle("Applicants", "G1", "G1", headerStyle)

	for i, applicant := range *applicants {
		f.SetCellValue("Applicants", fmt.Sprintf("A%d", i+2), applicant.UserProfile.Name)
		f.SetCellValue("Applicants", fmt.Sprintf("B%d", i+2), applicant.JobPosting.JobName)
		if applicant.UserProfile.WorkExperiences != nil {
			var concatenatedValue string
			for _, we := range *applicant.UserProfile.WorkExperiences {
				cellValue := we.Name

				if concatenatedValue != "" {
					concatenatedValue += ", "
				}
				concatenatedValue += cellValue
			}

			f.SetCellValue("Applicants", fmt.Sprintf("C%d", i+2), concatenatedValue)
		}
		f.SetCellValue("Applicants", fmt.Sprintf("D%d", i+2), applicant.AppliedDate)
		f.SetCellValue("Applicants", fmt.Sprintf("E%d", i+2), applicant.UserProfile.PhoneNumber)
		f.SetCellValue("Applicants", fmt.Sprintf("F%d", i+2), applicant.UserProfile.ExpectedSalary)
		f.SetCellValue("Applicants", fmt.Sprintf("G%d", i+2), applicant.UserProfile.CurrentSalary)
	}

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=applicants_exported.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(ctx.Writer); err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export applicants", err.Error())
		return
	}
}

// GetApplicantsForCoverLetter get applicants for cover letter
//
// @Summary get applicants for cover letter
// @Description get applicants for cover letter
// @Tags Applicants
// @Accept json
// @Produce json
// @Param job_posting_id query string true "Job Posting ID"
// @Param project_recruitment_line_id query string false "Project Recruitment Line ID"
// @Param hired_status query string false "Hired Status"
// @Security BearerAuth
// @Success 200 {array} response.ApplicantResponse
// @Router /applicants/cover-letter/{job_posting_id} [get]
func (h *ApplicantHandler) GetApplicantsForCoverLetter(ctx *gin.Context) {
	jobPostingID, err := uuid.Parse(ctx.Query("job_posting_id"))
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.GetApplicantsForCoverLetter] error when parsing job_posting_id: %v", err)
		utils.BadRequestResponse(ctx, "job_posting_id is not a valid UUID", err)
		return
	}

	projectRecruitmentLineID, err := uuid.Parse(ctx.Query("project_recruitment_line_id"))
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.GetApplicantsForCoverLetter] error when parsing project_recruitment_line_id: %v", err)
		utils.BadRequestResponse(ctx, "project_recruitment_line_id is not a valid UUID", err)
		return
	}

	// hiredStatus := ctx.Query("hired_status")
	// if hiredStatus == "" {
	// 	h.Log.Errorf("[ApplicantHandler.GetApplicantsForCoverLetter] hired_status is required")
	// 	utils.BadRequestResponse(ctx, "hired_status is required", nil)
	// 	return
	// }

	applicants, err := h.UseCase.GetApplicantsForCoverLetter(jobPostingID, projectRecruitmentLineID)
	if err != nil {
		h.Log.Errorf("[ApplicantHandler.GetApplicantsForCoverLetter] error when getting applicants for cover letter: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get applicants for cover letter", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Successfully get applicants for cover letter", applicants)
}
