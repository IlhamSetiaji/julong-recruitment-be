package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
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
}

type InterviewHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IInterviewUseCase
	UserHelper helper.IUserHelper
}

func NewInterviewHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IInterviewUseCase,
	userHelper helper.IUserHelper,
) IInterviewHandler {
	return &InterviewHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func InterviewHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IInterviewHandler {
	useCase := usecase.InterviewUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewInterviewHandler(log, viper, validate, useCase, userHelper)
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
// @Router /api/interview/generate-document-number [get]
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

	res, err := h.UseCase.FindScheduleForApplicant(applicantUUID, projectRecruitmentLineUUID, jobPostingUUID)
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
