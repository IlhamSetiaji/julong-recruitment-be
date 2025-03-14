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

type IInterviewHandler interface {
	CreateInterview(ctx *gin.Context)
	UpdateInterview(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	DeleteInterview(ctx *gin.Context)
	GenerateDocumentNumber(ctx *gin.Context)
	UpdateStatusInterview(ctx *gin.Context)
	FindMySchedule(ctx *gin.Context)
	FindMyScheduleForAssessor(ctx *gin.Context)
	FindApplicantSchedule(ctx *gin.Context)
	ExportInterviewScheduleAnswer(ctx *gin.Context)
	ExportResultTemplate(ctx *gin.Context)
	ReadResultTemplate(ctx *gin.Context)
}

type InterviewHandler struct {
	Log                           *logrus.Logger
	Viper                         *viper.Viper
	Validate                      *validator.Validate
	UseCase                       usecase.IInterviewUseCase
	UserHelper                    helper.IUserHelper
	UserProfileUseCase            usecase.IUserProfileUseCase
	ProjectRecruitmentLineUseCase usecase.IProjectRecruitmentLineUseCase
	DB                            *gorm.DB
	InterviewApplicantUseCase     usecase.IInterviewApplicantUseCase
}

func NewInterviewHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IInterviewUseCase,
	userHelper helper.IUserHelper,
	upUseCase usecase.IUserProfileUseCase,
	prlUseCase usecase.IProjectRecruitmentLineUseCase,
	db *gorm.DB,
	iaUseCase usecase.IInterviewApplicantUseCase,
) IInterviewHandler {
	return &InterviewHandler{
		Log:                           log,
		Viper:                         viper,
		Validate:                      validate,
		UseCase:                       useCase,
		UserHelper:                    userHelper,
		UserProfileUseCase:            upUseCase,
		ProjectRecruitmentLineUseCase: prlUseCase,
		DB:                            db,
		InterviewApplicantUseCase:     iaUseCase,
	}
}

func InterviewHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IInterviewHandler {
	useCase := usecase.InterviewUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	upUseCase := usecase.UserProfileUseCaseFactory(log, viper)
	prlUseCase := usecase.ProjectRecruitmentLineUseCaseFactory(log)
	db := config.NewDatabase()
	iaUseCase := usecase.InterviewApplicantUseCaseFactory(log, viper)
	return NewInterviewHandler(log, viper, validate, useCase, userHelper, upUseCase, prlUseCase, db, iaUseCase)
}

// CreateInterview creates a new interview
//
// @Summary Create a new interview
// @Description Create a new interview
// @Tags Interview
// @Accept json
// @Produce json
// @Param interview body request.CreateInterviewRequest true "Create Interview Request"
// @Success 201 {object} response.InterviewResponse
// @Security BearerAuth
// @Router /api/interview [post]
func (h *InterviewHandler) CreateInterview(ctx *gin.Context) {
	var req request.CreateInterviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[InterviewHandler.CreateInterview] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[InterviewHandler.CreateInterview] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.CreateInterview(&req)
	if err != nil {
		h.Log.Error("[InterviewHandler.CreateInterview] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create interview", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Interview created", res)
}

// UpdateInterview updates an interview
//
// @Summary Update an interview
// @Description Update an interview
// @Tags Interview
// @Accept json
// @Produce json
// @Param interview body request.UpdateInterviewRequest true "Update Interview Request"
// @Success 200 {object} response.InterviewResponse
// @Security BearerAuth
// @Router /api/interview/update [put]
func (h *InterviewHandler) UpdateInterview(ctx *gin.Context) {
	var req request.UpdateInterviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[InterviewHandler.UpdateInterview] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[InterviewHandler.UpdateInterview] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.UpdateInterview(&req)
	if err != nil {
		h.Log.Error("[InterviewHandler.UpdateInterview] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update interview", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Interview updated", res)
}

// FindAllPaginated finds all interviews paginated
//
// @Summary Find all interviews paginated
// @Description Find all interviews paginated
// @Tags Interview
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param search query string false "Search"
// @Param created_at query string false "Created At"
// @Success 200 {object} response.InterviewResponse
// @Security BearerAuth
// @Router /api/interview [get]
func (h *InterviewHandler) FindAllPaginated(ctx *gin.Context) {
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

	interviews, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find interviews", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "interviews found", gin.H{
		"interviews": interviews,
		"total":      total,
	})
}

// FindByID finds an interview by ID
//
// @Summary Find an interview by ID
// @Description Find an interview by ID
// @Tags Interview
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} response.InterviewResponse
// @Security BearerAuth
// @Router /api/interview/{id} [get]
func (h *InterviewHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[InterviewHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid ID", err)
		return
	}

	res, err := h.UseCase.FindByID(parsedID)
	if err != nil {
		h.Log.Error("[InterviewHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find interview", err.Error())
		return
	}

	if res == nil {
		h.Log.Error("[InterviewHandler.FindByID] Interview not found")
		utils.ErrorResponse(ctx, http.StatusNotFound, "Failed to find interview", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Interview found", res)
}

// DeleteInterview deletes an interview by ID
//
// @Summary Delete an interview by ID
// @Description Delete an interview by ID
// @Tags Interview
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {string} string
// @Security BearerAuth
// @Router /api/interview/{id} [delete]
func (h *InterviewHandler) DeleteInterview(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[InterviewHandler.DeleteInterview] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid ID", err)
		return
	}

	err = h.UseCase.DeleteByID(parsedID)
	if err != nil {
		h.Log.Error("[InterviewHandler.DeleteInterview] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete interview", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Interview deleted", nil)
}

// GenerateDocumentNumber generates a document number for an interview
//
// @Summary Generate a document number for an interview
// @Description Generate a document number for an interview
// @Tags Interview
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Security BearerAuth
// @Router /api/interview/document-number [get]
func (h *InterviewHandler) GenerateDocumentNumber(ctx *gin.Context) {
	dateNow := time.Now()
	documentNumber, err := h.UseCase.GenerateDocumentNumber(dateNow)
	if err != nil {
		h.Log.Error("[InterviewHandler.GenerateDocumentNumber] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate document number", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document number generated", documentNumber)
}

// UpdateStatusInterview updates an interview status
//
// @Summary Update an interview status
// @Description Update an interview status
// @Tags Interview
// @Accept json
// @Produce json
// @Param status body request.UpdateStatusInterviewRequest true "Update Status Interview Request"
// @Success 200 {object} response.InterviewResponse
// @Security BearerAuth
// @Router /api/interview/update-status [put]
func (h *InterviewHandler) UpdateStatusInterview(ctx *gin.Context) {
	var req request.UpdateStatusInterviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[InterviewHandler.UpdateStatusInterview] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[InterviewHandler.UpdateStatusInterview] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.UpdateStatusInterview(&req)
	if err != nil {
		h.Log.Error("[InterviewHandler.UpdateStatusInterview] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update interview status", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Interview status updated", res)
}

// FindMySchedule finds my interview schedule
//
// @Summary Find my interview schedule
// @Description Find my interview schedule
// @Tags Interview
// @Accept json
// @Produce json
// @Param			project_recruitment_line_id	query	string	false	"Project Recruitment Line ID"
// @Param      job_posting_id	query	string	false	"Job Posting ID"
// @Success 200 {object} response.InterviewMyselfResponse
// @Security BearerAuth
// @Router /api/interview/my-schedule [get]
func (h *InterviewHandler) FindMySchedule(ctx *gin.Context) {
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
		h.Log.Error("[InterviewHandler.FindMySchedule] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find my interview schedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "My interview schedule found", res)
}

// FindApplicantSchedule finds applicant interview schedule
//
// @Summary Find applicant interview schedule
// @Description Find applicant interview schedule
// @Tags Interview
// @Accept json
// @Produce json
// @Param			project_recruitment_line_id	query	string	false	"Project Recruitment Line ID"
// @Param      job_posting_id	query	string	false	"Job Posting ID"
// @Param      applicant_id	query	string	false	"Applicant ID"
// @Success 200 {object} response.InterviewMyselfResponse
// @Security BearerAuth
// @Router /api/interview/applicant-schedule [get]
func (h *InterviewHandler) FindApplicantSchedule(ctx *gin.Context) {
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
		h.Log.Error("[InterviewHandler.FindApplicantSchedule] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find applicant interview schedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Applicant interview schedule found", res)
}

// FindMyScheduleForAssessor finds my interview schedule for assessor
//
// @Summary Find my interview schedule for assessor
// @Description Find my interview schedule for assessor
// @Tags Interview
// @Accept json
// @Produce json
// @Param			project_recruitment_line_id	query	string	false	"Project Recruitment Line ID"
// @Param      job_posting_id	query	string	false	"Job Posting ID"
// @Success 200 {object} response.InterviewMyselfForAssessorResponse
// @Security BearerAuth
// @Router /api/interview/assessor-schedule [get]
func (h *InterviewHandler) FindMyScheduleForAssessor(ctx *gin.Context) {
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
		h.Log.Error("[InterviewHandler.FindMyScheduleForAssessor] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find my interview schedule for assessor", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "My interview schedule for assessor found", res)
}

// ExportInterviewScheduleAnswer exports interview schedule answer
//
// @Summary Export interview schedule answer
// @Description Export interview schedule answer
// @Tags Interview
// @Accept json
//
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Param			id					query	string	true	"Interview ID"
//	@Param			job_posting_id		query	string	true	"Job Posting ID"
//	@Success		200					{file}		file
//	@Security		BearerAuth
//	@Router			/api/interviews/export-answer [get]
func (h *InterviewHandler) ExportInterviewScheduleAnswer(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		h.Log.Error("interview ID is required")
		utils.BadRequestResponse(ctx, "Interview ID is required", nil)
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

	interview, err := h.UseCase.FindByIDForAnswer(testScheduleHeaderID, jobPostingUUID)
	if err != nil {
		h.Log.Error("[InterviewHandler.ExportInterviewScheduleAnswer] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find interview", err.Error())
		return
	}

	if interview == nil {
		h.Log.Error("[InterviewHandler.ExportInterviewScheduleAnswer] Interview not found")
		utils.ErrorResponse(ctx, http.StatusNotFound, "Failed to find interview", err.Error())
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
	// loop interview assessors for sheet name
	for _, interviewData := range interview.InterviewAssessors {
		sheetName := interviewData.EmployeeName
		newSheet, err := f.NewSheet(sheetName)
		if err != nil {
			h.Log.Error("[InterviewHandler.ExportInterviewScheduleAnswer] " + err.Error())
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

		// loop through interviews -> interview applicants
		for i, ia := range interview.InterviewApplicants {
			f.SetCellValue(sheetName, "A"+strconv.Itoa(i+2), ia.ApplicantID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(i+2), ia.UserProfile.Name)

			prl, err := h.ProjectRecruitmentLineUseCase.FindByIDForAnswerInterview(interview.ProjectRecruitmentLineID, jobPostingUUID, ia.UserProfileID, interviewData.ID)
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
	ctx.Header("Content-Disposition", "attachment; filename=interview_answers.xlsx")
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
// @Tags Interview
// @Accept json
//
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Param			id					query	string	true	"Interview ID"
//	@Param			job_posting_id		query	string	true	"Job Posting ID"
//	@Success		200					{file}		file
//	@Security		BearerAuth
//	@Router			/api/interviews/export-result-template [get]
func (h *InterviewHandler) ExportResultTemplate(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		h.Log.Error("interview ID is required")
		utils.BadRequestResponse(ctx, "Interview ID is required", nil)
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
		utils.BadRequestResponse(ctx, "Invalid interview schedule header ID", err)
		return
	}

	jobPostingUUID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid job posting ID", err)
		return
	}

	interview, err := h.UseCase.FindByIDForAnswer(testScheduleHeaderID, jobPostingUUID)
	if err != nil {
		h.Log.Error("[InterviewHandler.ExportResultTemplate] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find interview", err.Error())
		return
	}

	if interview == nil {
		h.Log.Error("[InterviewHandler.ExportResultTemplate] Interview not found")
		utils.ErrorResponse(ctx, http.StatusNotFound, "Failed to find interview", err.Error())
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
	f.SetCellValue("Applicants", "A1", "Interview Applicant ID")
	f.SetCellValue("Applicants", "B1", "Applicant Name")
	f.SetCellValue("Applicants", "C1", "Final Result (ACCEPTED/REJECTED)")

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
	for i, ta := range interview.InterviewApplicants {
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
	ctx.Header("Content-Disposition", "attachment; filename=interview_applicants.xlsx")
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
// @Tags Interview
//
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"File"
//	@Success		200			{object}	string
//	@Security		BearerAuth
//	@Router			/api/interviews/read-result-template [post]
func (h *InterviewHandler) ReadResultTemplate(ctx *gin.Context) {
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
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid interview applicant ID for row "+strconv.Itoa(i), err.Error())
			return
		}

		_, err = h.InterviewApplicantUseCase.UpdateFinalResultStatusTestApplicant(ctx, testApplicantUUID, entity.FinalResultStatus(finalResult))
		if err != nil {
			h.Log.Error(err)
			tx.Rollback()
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update final result status interview applicant for row "+strconv.Itoa(i), err.Error())
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
