package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type IFgdScheduleHandler interface {
	CreateFgdSchedule(ctx *gin.Context)
	UpdateFgdSchedule(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	DeleteFgdSchedule(ctx *gin.Context)
	GenerateDocumentNumber(ctx *gin.Context)
	UpdateStatusFgdSchedule(ctx *gin.Context)
	FindMySchedule(ctx *gin.Context)
	FindMyScheduleForAssessor(ctx *gin.Context)
	FindApplicantSchedule(ctx *gin.Context)
	ExportFgdScheduleAnswer(ctx *gin.Context)
	ExportResultTemplate(ctx *gin.Context)
	ReadResultTemplate(ctx *gin.Context)
}

type FgdScheduleHandler struct {
	Log                           *logrus.Logger
	Viper                         *viper.Viper
	Validate                      *validator.Validate
	UseCase                       usecase.IFgdScheduleUseCase
	UserHelper                    helper.IUserHelper
	UserProfileUseCase            usecase.IUserProfileUseCase
	ProjectRecruitmentLineUseCase usecase.IProjectRecruitmentLineUseCase
	DB                            *gorm.DB
	FgdApplicantUseCase           usecase.IFgdApplicantUseCase
}

func NewFgdScheduleHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IFgdScheduleUseCase,
	userHelper helper.IUserHelper,
	upUseCase usecase.IUserProfileUseCase,
	prlUseCase usecase.IProjectRecruitmentLineUseCase,
	db *gorm.DB,
	iaUseCase usecase.IFgdApplicantUseCase,
) IFgdScheduleHandler {
	return &FgdScheduleHandler{
		Log:                           log,
		Viper:                         viper,
		Validate:                      validate,
		UseCase:                       useCase,
		UserHelper:                    userHelper,
		UserProfileUseCase:            upUseCase,
		ProjectRecruitmentLineUseCase: prlUseCase,
		DB:                            db,
		FgdApplicantUseCase:           iaUseCase,
	}
}

func FgdScheduleHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IFgdScheduleHandler {
	useCase := usecase.FgdScheduleUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	upUseCase := usecase.UserProfileUseCaseFactory(log, viper)
	prlUseCase := usecase.ProjectRecruitmentLineUseCaseFactory(log)
	db := config.NewDatabase()
	iaUseCase := usecase.FgdApplicantUseCaseFactory(log, viper)
	return NewFgdScheduleHandler(log, viper, validate, useCase, userHelper, upUseCase, prlUseCase, db, iaUseCase)
}

// CreateFgdSchedule creates a new FgdSchedule
//
// @Summary Create a new FgdSchedule
// @Description Create a new FgdSchedule
// @Tags FgdSchedule
// @Accept json
// @Produce json
// @Param FgdSchedule body request.CreateFgdScheduleRequest true "Create FgdSchedule Request"
// @Success 201 {object} response.FgdScheduleResponse
// @Security BearerAuth
// @Router /api/fgd-schedules [post]
func (h *FgdScheduleHandler) CreateFgdSchedule(ctx *gin.Context) {
	var req request.CreateFgdScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[FgdScheduleHandler.CreateFgdSchedule] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[FgdScheduleHandler.CreateFgdSchedule] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.CreateFgdSchedule(&req)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.CreateFgdSchedule] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create FgdSCreateFgdSchedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "FgdSCreateFgdSchedule created", res)
}

// UpdateFgdSchedule updates an FgdSchedule
//
// @Summary Update an FgdSchedule
// @Description Update an FgdSchedule
// @Tags FgdSchedule
// @Accept json
// @Produce json
// @Param FgdSchedule body request.UpdateFgdScheduleRequest true "Update FgdSchedule Request"
// @Success 200 {object} response.FgdScheduleResponse
// @Security BearerAuth
// @Router /api/fgd-schedules/update [put]
func (h *FgdScheduleHandler) UpdateFgdSchedule(ctx *gin.Context) {
	var req request.UpdateFgdScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[FgdScheduleHandler.UpdateFgdSchedule] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[FgdScheduleHandler.UpdateFgdSchedule] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.UpdateFgdSchedule(&req)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.UpdateFgdSchedule] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update FgdSchedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "FgdSchedule updated", res)
}

// FindAllPaginated finds all FgdSchedules paginated
//
// @Summary Find all FgdSchedules paginated
// @Description Find all FgdSchedules paginated
// @Tags FgdSchedule
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param search query string false "Search"
// @Param created_at query string false "Created At"
// @Success 200 {object} response.FgdScheduleResponse
// @Security BearerAuth
// @Router /api/fgd-schedules [get]
func (h *FgdScheduleHandler) FindAllPaginated(ctx *gin.Context) {
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

	FgdSchedules, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find FgdSchedules", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "FgdSchedules found", gin.H{
		"fgd_schedules": FgdSchedules,
		"total":         total,
	})
}

// FindByID finds an FgdSchedule by ID
//
// @Summary Find an FgdSchedule by ID
// @Description Find an FgdSchedule by ID
// @Tags FgdSchedule
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} response.FgdScheduleResponse
// @Security BearerAuth
// @Router /api/FgdSchedule/{id} [get]
func (h *FgdScheduleHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid ID", err)
		return
	}

	res, err := h.UseCase.FindByID(parsedID)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find FgdSchedule", err.Error())
		return
	}

	if res == nil {
		h.Log.Error("[FgdScheduleHandler.FindByID] FgdSchedule not found")
		utils.ErrorResponse(ctx, http.StatusNotFound, "Failed to find FgdSchedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "FgdSchedule found", res)
}

// DeleteFgdSchedule deletes an FgdSchedule by ID
//
// @Summary Delete an FgdSchedule by ID
// @Description Delete an FgdSchedule by ID
// @Tags FgdSchedule
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {string} string
// @Security BearerAuth
// @Router /api/fgd-schedules/{id} [delete]
func (h *FgdScheduleHandler) DeleteFgdSchedule(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.DeleteFgdSchedule] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid ID", err)
		return
	}

	err = h.UseCase.DeleteByID(parsedID)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.DeleteFgdSchedule] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete FgdSchedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "FgdSchedule deleted", nil)
}

// GenerateDocumentNumber generates a document number for an FgdSchedule
//
// @Summary Generate a document number for an FgdSchedule
// @Description Generate a document number for an FgdSchedule
// @Tags FgdSchedule
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Security BearerAuth
// @Router /api/fgd-schedules/generate-document-number [get]
func (h *FgdScheduleHandler) GenerateDocumentNumber(ctx *gin.Context) {
	dateNow := time.Now()
	documentNumber, err := h.UseCase.GenerateDocumentNumber(dateNow)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.GenerateDocumentNumber] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate document number", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document number generated", documentNumber)
}

// UpdateStatusFgdSchedule updates an FgdSchedule status
//
// @Summary Update an FgdSchedule status
// @Description Update an FgdSchedule status
// @Tags FgdSchedule
// @Accept json
// @Produce json
// @Param status body request.UpdateStatusFgdScheduleRequest true "Update Status FgdSchedule Request"
// @Success 200 {object} response.FgdScheduleResponse
// @Security BearerAuth
// @Router /api/fgd-schedules/update-status [put]
func (h *FgdScheduleHandler) UpdateStatusFgdSchedule(ctx *gin.Context) {
	var req request.UpdateStatusFgdScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[FgdScheduleHandler.UpdateStatusFgdSchedule] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[FgdScheduleHandler.UpdateStatusFgdSchedule] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.UpdateStatusFgdSchedule(&req)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.UpdateStatusFgdSchedule] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update FgdSchedule status", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "FgdSchedule status updated", res)
}

// FindMySchedule finds my FgdSchedule schedule
//
// @Summary Find my FgdSchedule schedule
// @Description Find my FgdSchedule schedule
// @Tags FgdSchedule
// @Accept json
// @Produce json
// @Param			project_recruitment_line_id	query	string	false	"Project Recruitment Line ID"
// @Param      job_posting_id	query	string	false	"Job Posting ID"
// @Success 200 {object} response.FgdScheduleMyselfResponse
// @Security BearerAuth
// @Router /api/fgd-schedules/my-schedule [get]
func (h *FgdScheduleHandler) FindMySchedule(ctx *gin.Context) {
	projectRecruitmentLineID := ctx.Query("project_recruitment_line_id")
	if projectRecruitmentLineID == "" {
		h.Log.Error("Project recruitment line ID is required")
		utils.BadRequestResponse(ctx, "Project recruitment line ID is required", nil)
		return
	}

	jobPostingID := ctx.Query("job_posting_id")
	if jobPostingID == "" {
		h.Log.Error("Job posting ID is required")
		utils.BadRequestResponse(ctx, "Job posting ID is required", nil)
		return
	}

	projectRecruitmentLineUUID, err := uuid.Parse(projectRecruitmentLineID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid project recruitment line ID", err)
		return
	}

	jobPostingUUID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid job posting ID", err)
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

	res, err := h.UseCase.FindMySchedule(userUUID, projectRecruitmentLineUUID, jobPostingUUID)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.FindMySchedule] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find my FgdSchedule schedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "My FgdSchedule schedule found", res)
}

// FindApplicantSchedule finds applicant FgdSchedule schedule
//
// @Summary Find applicant FgdSchedule schedule
// @Description Find applicant FgdSchedule schedule
// @Tags FgdSchedule
// @Accept json
// @Produce json
// @Param			project_recruitment_line_id	query	string	false	"Project Recruitment Line ID"
// @Param      job_posting_id	query	string	false	"Job Posting ID"
// @Param      applicant_id	query	string	false	"Applicant ID"
// @Success 200 {object} response.FgdScheduleMyselfResponse
// @Security BearerAuth
// @Router /api/fgd-schedules/applicant-schedule [get]
func (h *FgdScheduleHandler) FindApplicantSchedule(ctx *gin.Context) {
	projectRecruitmentLineID := ctx.Query("project_recruitment_line_id")
	if projectRecruitmentLineID == "" {
		h.Log.Error("Project recruitment line ID is required")
		utils.BadRequestResponse(ctx, "Project recruitment line ID is required", nil)
		return
	}

	jobPostingID := ctx.Query("job_posting_id")
	if jobPostingID == "" {
		h.Log.Error("Job posting ID is required")
		utils.BadRequestResponse(ctx, "Job posting ID is required", nil)
		return
	}

	applicantID := ctx.Query("applicant_id")
	if applicantID == "" {
		h.Log.Error("Applicant ID is required")
		utils.BadRequestResponse(ctx, "Applicant ID is required", nil)
		return
	}

	projectRecruitmentLineUUID, err := uuid.Parse(projectRecruitmentLineID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid project recruitment line ID", err)
		return
	}

	jobPostingUUID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid job posting ID", err)
		return
	}

	applicantUUID, err := uuid.Parse(applicantID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid applicant ID", err)
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
	employeeUUID, err := h.UserHelper.GetEmployeeId(user)
	if err != nil {
		h.Log.Errorf("Error when getting employee id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	res, err := h.UseCase.FindScheduleForApplicant(applicantUUID, projectRecruitmentLineUUID, jobPostingUUID, employeeUUID)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.FindApplicantSchedule] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find applicant FgdSchedule schedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Applicant FgdSchedule schedule found", res)
}

// FindMyScheduleForAssessor finds my FgdSchedule schedule for assessor
//
// @Summary Find my FgdSchedule schedule for assessor
// @Description Find my FgdSchedule schedule for assessor
// @Tags FgdSchedule
// @Accept json
// @Produce json
// @Param			project_recruitment_line_id	query	string	false	"Project Recruitment Line ID"
// @Param      job_posting_id	query	string	false	"Job Posting ID"
// @Success 200 {object} response.FgdScheduleMyselfForAssessorResponse
// @Security BearerAuth
// @Router /api/FgdSchedule/assessor-schedule [get]
func (h *FgdScheduleHandler) FindMyScheduleForAssessor(ctx *gin.Context) {
	projectRecruitmentLineID := ctx.Query("project_recruitment_line_id")
	if projectRecruitmentLineID == "" {
		h.Log.Error("Project recruitment line ID is required")
		utils.BadRequestResponse(ctx, "Project recruitment line ID is required", nil)
		return
	}

	jobPostingID := ctx.Query("job_posting_id")
	if jobPostingID == "" {
		h.Log.Error("Job posting ID is required")
		utils.BadRequestResponse(ctx, "Job posting ID is required", nil)
		return
	}

	projectRecruitmentLineUUID, err := uuid.Parse(projectRecruitmentLineID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid project recruitment line ID", err)
		return
	}

	jobPostingUUID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid job posting ID", err)
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
	employeeUUID, err := h.UserHelper.GetEmployeeId(user)
	if err != nil {
		h.Log.Errorf("Error when getting employee id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	res, err := h.UseCase.FindMyScheduleForAssessor(employeeUUID, projectRecruitmentLineUUID, jobPostingUUID)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.FindMyScheduleForAssessor] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find my FgdSchedule schedule for assessor", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "My FgdSchedule schedule for assessor found", res)
}

// ExportFgdScheduleAnswer exports FgdSchedule schedule answer
//
// @Summary Export FgdSchedule schedule answer
// @Description Export FgdSchedule schedule answer
// @Tags FgdSchedule
// @Accept json
//
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Param			id					query	string	true	"FgdSchedule ID"
//	@Param			job_posting_id		query	string	true	"Job Posting ID"
//	@Success		200					{file}		file
//	@Security		BearerAuth
//	@Router			/api/fgd-schedules/export-answer [get]
func (h *FgdScheduleHandler) ExportFgdScheduleAnswer(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		h.Log.Error("FgdSchedule ID is required")
		utils.BadRequestResponse(ctx, "FgdSchedule ID is required", nil)
		return
	}

	jobPostingID := ctx.Query("job_posting_id")
	if jobPostingID == "" {
		h.Log.Error("Job posting ID is required")
		utils.BadRequestResponse(ctx, "Job posting ID is required", nil)
		return
	}

	testScheduleHeaderID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid test schedule header ID", err)
		return
	}

	jobPostingUUID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid job posting ID", err)
		return
	}

	FgdSchedule, err := h.UseCase.FindByIDForAnswer(testScheduleHeaderID, jobPostingUUID)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.ExportFgdScheduleScheduleAnswer] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find FgdSchedule", err.Error())
		return
	}

	if FgdSchedule == nil {
		h.Log.Error("[FgdScheduleHandler.ExportFgdScheduleScheduleAnswer] FgdSchedule not found")
		utils.ErrorResponse(ctx, http.StatusNotFound, "Failed to find FgdSchedule", err.Error())
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
			return
		}
	}()
	// loop FgdSchedule assessors for sheet name
	for _, FgdScheduleData := range FgdSchedule.FgdAssessors {
		sheetName := FgdScheduleData.EmployeeName
		newSheet, err := f.NewSheet(sheetName)
		if err != nil {
			h.Log.Error("[FgdScheduleHandler.ExportFgdScheduleScheduleAnswer] " + err.Error())
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
			return
		}
		f.SetActiveSheet(newSheet)
		f.SetCellValue(sheetName, "A1", "Applicant ID")
		f.SetCellValue(sheetName, "B1", "Applicant Name")

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
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
			return
		}

		// Set the style to the header
		f.SetCellStyle(sheetName, "A1", "A1", headerStyle)
		f.SetCellStyle(sheetName, "B1", "B1", headerStyle)

		// loop through FgdSchedules -> FgdSchedule applicants
		for i, ia := range FgdSchedule.FgdApplicants {
			f.SetCellValue(sheetName, "A"+strconv.Itoa(i+2), ia.ApplicantID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(i+2), ia.UserProfile.Name)

			prl, err := h.ProjectRecruitmentLineUseCase.FindByIDForAnswerFgd(FgdSchedule.ProjectRecruitmentLineID, jobPostingUUID, ia.UserProfileID, FgdScheduleData.ID)
			if err != nil {
				h.Log.Error(err)
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find project recruitment line", err.Error())
				return
			}

			if prl == nil {
				utils.ErrorResponse(ctx, http.StatusNotFound, "Project recruitment line not found", "Project recruitment line not found")
				return
			}

			// loop through project recruitment line -> template activity line -> template question -> questions for header
			for i, questionData := range *prl.TemplateActivityLine.TemplateQuestion.Questions {
				var concatenatedValue string
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", string(rune(i+67)), 1), questionData.Name)

				// Set the style to the header
				f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", string(rune(i+67)), 1), fmt.Sprintf("%s%d", string(rune(i+67)), 1), headerStyle)

				// loop through project recruitment line -> template activity line -> template question -> questions -> question responses for line
				if questionData.QuestionResponses == nil {
					continue
				}
				for _, questionResponse := range *questionData.QuestionResponses {
					var cellValue string
					if questionResponse.AnswerFile == "" {
						cellValue = questionResponse.Answer
					} else {
						cellValue = h.Viper.GetString("app.url") + questionResponse.AnswerFile
					}

					if concatenatedValue != "" {
						concatenatedValue += ", "
					}
					concatenatedValue += cellValue
				}

				// Set the concatenated value to the cell
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", string(rune(i+67)), 2), concatenatedValue)

				// Set the width of the columns
				f.SetColWidth(sheetName, fmt.Sprintf("%s", string(rune(i+67))), string(rune(i+67)), 20)
			}
		}
		// Set the width of the columns
		f.SetColWidth(sheetName, "A", "A", 20)
		f.SetColWidth(sheetName, "B", "B", 20)
	}

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=FgdSchedule_answers.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(ctx.Writer); err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
		return
	}
}

// ExportResultTemplate exports result template
//
// @Summary Export result template
// @Description Export result template
// @Tags FgdSchedule
// @Accept json
//
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Param			id					query	string	true	"FgdSchedule ID"
//	@Param			job_posting_id		query	string	true	"Job Posting ID"
//	@Success		200					{file}		file
//	@Security		BearerAuth
//	@Router			/api/fgd-schedules/export-result-template [get]
func (h *FgdScheduleHandler) ExportResultTemplate(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		h.Log.Error("FgdSchedule ID is required")
		utils.BadRequestResponse(ctx, "FgdSchedule ID is required", nil)
		return
	}

	jobPostingID := ctx.Query("job_posting_id")
	if jobPostingID == "" {
		h.Log.Error("Job posting ID is required")
		utils.BadRequestResponse(ctx, "Job posting ID is required", nil)
		return
	}

	testScheduleHeaderID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid FgdSchedule schedule header ID", err)
		return
	}

	jobPostingUUID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid job posting ID", err)
		return
	}

	FgdSchedule, err := h.UseCase.FindByIDForAnswer(testScheduleHeaderID, jobPostingUUID)
	if err != nil {
		h.Log.Error("[FgdScheduleHandler.ExportResultTemplate] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find FgdSchedule", err.Error())
		return
	}

	if FgdSchedule == nil {
		h.Log.Error("[FgdScheduleHandler.ExportResultTemplate] FgdSchedule not found")
		utils.ErrorResponse(ctx, http.StatusNotFound, "Failed to find FgdSchedule", err.Error())
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
			return
		}
	}()

	f.SetSheetName("Sheet1", "Applicants")
	// Set value of a cell.
	f.SetCellValue("Applicants", "A1", "FgdSchedule Applicant ID")
	f.SetCellValue("Applicants", "B1", "Applicant Name")
	f.SetCellValue("Applicants", "C1", "Final Result")

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
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
		return
	}

	// Set the style to the header
	f.SetCellStyle("Applicants", "A1", "A1", headerStyle)
	f.SetCellStyle("Applicants", "B1", "B1", headerStyle)
	f.SetCellStyle("Applicants", "C1", "C1", headerStyle)

	// loop through test schedule header -> test applicants
	for i, ta := range FgdSchedule.FgdApplicants {
		f.SetCellValue("Applicants", fmt.Sprintf("A%d", i+2), ta.ID)
		f.SetCellValue("Applicants", fmt.Sprintf("B%d", i+2), ta.UserProfile.Name)
	}

	// Set the width of the columns
	f.SetColWidth("Applicants", "A", "A", 20)
	f.SetColWidth("Applicants", "B", "B", 20)
	f.SetColWidth("Applicants", "C", "C", 20)

	// Set active sheet of the workbook.
	// f.SetActiveSheet(index)

	// Write the file to the response body
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=FgdSchedule_applicants.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(ctx.Writer); err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
		return
	}
}

// ReadResultTemplate reads result template
//
// @Summary Read result template
// @Description Read result template
// @Tags FgdSchedule
//
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"File"
//	@Success		200			{object}	string
//	@Security		BearerAuth
//	@Router			/api/FgdSchedules/read-result-template [post]
func (h *FgdScheduleHandler) ReadResultTemplate(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "File is required", err)
		return
	}

	timestamp := time.Now().UnixNano()
	filePath := fmt.Sprintf("storage/tests/results/%s", strconv.FormatInt(timestamp, 10)+"_"+file.Filename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to save file", err.Error())
		return
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to open file", err.Error())
		return
	}

	rows, err := f.GetRows("Applicants")
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get rows", err.Error())
		return
	}

	tx := h.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		h.Log.Warnf("Failed begin transaction : %+v", tx.Error)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to begin transaction", tx.Error.Error())
		return
	}
	defer tx.Rollback()

	for i, row := range rows {
		if i == 0 {
			continue
		}

		testApplicantID := row[0]
		var finalResult string
		if len(row) > 2 {
			finalResult = row[2]
			h.Log.Info("finalResult: ", finalResult)
		} else {
			h.Log.Warn("finalResult not found for row: ", i)
			tx.Rollback()
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Final result not found for row "+strconv.Itoa(i), "Final result not found for row "+strconv.Itoa(i))
			return
		}

		if finalResult != string(entity.FINAL_RESULT_STATUS_ACCEPTED) && finalResult != string(entity.FINAL_RESULT_STATUS_REJECTED) {
			h.Log.Warn("finalResult not valid for row: ", i)
			tx.Rollback()
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Final result not valid for row "+strconv.Itoa(i), "Final result not valid for row "+strconv.Itoa(i))
			return
		}

		h.Log.Info("finalResult: ", finalResult)

		testApplicantUUID, err := uuid.Parse(testApplicantID)
		if err != nil {
			h.Log.Error(err)
			tx.Rollback()
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid FgdSchedule applicant ID for row "+strconv.Itoa(i), err.Error())
			return
		}

		_, err = h.FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant(ctx, testApplicantUUID, entity.FinalResultStatus(finalResult))
		if err != nil {
			h.Log.Error(err)
			tx.Rollback()
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update final result status FgdSchedule applicant for row "+strconv.Itoa(i), err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("Failed commit transaction : %+v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Result template read", "Result template read")
}
