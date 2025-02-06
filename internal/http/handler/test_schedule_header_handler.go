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

type ITestScheduleHeaderHandler interface {
	CreateTestScheduleHeader(ctx *gin.Context)
	UpdateTestScheduleHeader(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	DeleteTestScheduleHeader(ctx *gin.Context)
	GenerateDocumentNumber(ctx *gin.Context)
	UpdateStatusTestScheduleHeader(ctx *gin.Context)
	FindMySchedule(ctx *gin.Context)
}

type TestScheduleHeaderHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.ITestScheduleHeaderUsecase
	UserHelper helper.IUserHelper
}

func NewTestScheduleHeaderHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITestScheduleHeaderUsecase,
	userHelper helper.IUserHelper,
) ITestScheduleHeaderHandler {
	return &TestScheduleHeaderHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func TestScheduleHeaderHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITestScheduleHeaderHandler {
	useCase := usecase.TestScheduleHeaderUsecaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewTestScheduleHeaderHandler(log, viper, validate, useCase, userHelper)
}

// CreateTestScheduleHeader create test schedule header
//
//	@Summary		Create test schedule header
//	@Description	Create test schedule header
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Param			test_schedule_header	body		request.CreateTestScheduleHeaderRequest	true	"Create test schedule header"
//	@Success		201			{object}	response.TestScheduleHeaderResponse
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers [post]
func (h *TestScheduleHeaderHandler) CreateTestScheduleHeader(ctx *gin.Context) {
	var req request.CreateTestScheduleHeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.CreateTestScheduleHeader(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create test schedule header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Test schedule header created", res)
}

// UpdateTestScheduleHeader update test schedule header
//
//	@Summary		Update test schedule header
//	@Description	Update test schedule header
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Param			test_schedule_header	body		request.UpdateTestScheduleHeaderRequest	true	"Update test schedule header"
//	@Success		200			{object}	response.TestScheduleHeaderResponse
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/update [put]
func (h *TestScheduleHeaderHandler) UpdateTestScheduleHeader(ctx *gin.Context) {
	var req request.UpdateTestScheduleHeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.UpdateTestScheduleHeader(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update test schedule header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test schedule header updated", res)
}

// FindAllPaginated find all test schedule headers paginated
//
//		@Summary		Find all test schedule headers paginated
//		@Description	Find all test schedule headers paginated
//		@Tags			Test Schedule Headers
//		@Accept			json
//		@Produce		json
//	 	@Param			page	query	int	false	"Page"
//		@Param			page_size	query	int	false	"Page Size"
//		@Param			search	query	string	false	"Search"
//		@Param			created_at	query	string	false	"Created At"
//		@Success		200			{object}	response.TestScheduleHeaderResponse
//		@Security		BearerAuth
//		@Router			/api/test-schedule-headers [get]
func (h *TestScheduleHeaderHandler) FindAllPaginated(ctx *gin.Context) {
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

	testScheduleHeaders, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find test schedule headers", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test schedule headers found", gin.H{
		"test_schedule_headers": testScheduleHeaders,
		"total":                 total,
	})
}

// FindByID find test schedule header by id
//
//	@Summary		Find test schedule header by id
//	@Description	Find test schedule header by id
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Test schedule header ID"
//	@Success		200			{object}	response.TestScheduleHeaderResponse
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/{id} [get]
func (h *TestScheduleHeaderHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Error("Test schedule header ID is required")
		utils.BadRequestResponse(ctx, "Test schedule header ID is required", nil)
		return
	}

	testScheduleHeaderID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid test schedule header ID", err)
		return
	}

	res, err := h.UseCase.FindByID(testScheduleHeaderID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find test schedule header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test schedule header found", res)
}

// DeleteTestScheduleHeader delete test schedule header
//
//	@Summary		Delete test schedule header
//	@Description	Delete test schedule header
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Test schedule header ID"
//	@Success		200			{string}	string
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/{id} [delete]
func (h *TestScheduleHeaderHandler) DeleteTestScheduleHeader(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Error("Test schedule header ID is required")
		utils.BadRequestResponse(ctx, "Test schedule header ID is required", nil)
		return
	}

	testScheduleHeaderID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid test schedule header ID", err)
		return
	}

	err = h.UseCase.DeleteTestScheduleHeader(testScheduleHeaderID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete test schedule header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test schedule header deleted", nil)
}

// GenerateDocumentNumber generate document number
//
//	@Summary		Generate document number
//	@Description	Generate document number
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Success		200			{string}	string
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/document-number [get]
func (h *TestScheduleHeaderHandler) GenerateDocumentNumber(ctx *gin.Context) {
	dateNow := time.Now()
	documentNumber, err := h.UseCase.GenerateDocumentNumber(dateNow)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate document number", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document number generated", documentNumber)
}

// UpdateStatusTestScheduleHeader update status test schedule header
//
//	@Summary		Update status test schedule header
//	@Description	Update status test schedule header
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Param			id		body	string	true	"Test schedule header ID"
//	@Param			status	body	string	true	"Status"
//	@Success		200			{object}	string
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/update-status [put]
func (h *TestScheduleHeaderHandler) UpdateStatusTestScheduleHeader(ctx *gin.Context) {
	var req request.UpdateStatusTestScheduleHeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	err := h.UseCase.UpdateStatusTestScheduleHeader(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update status test schedule header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Status test schedule header updated", "Status test schedule header updated")
}

// FindMySchedule find my schedule
//
//		@Summary		Find my schedule
//		@Description	Find my schedule
//		@Tags			Test Schedule Headers
//		@Accept			json
//		@Produce		json
//		@Param			project_recruitment_line_id	query	string	false	"Project Recruitment Line ID"
//	 @Param      job_posting_id	query	string	false	"Job Posting ID"
//		@Success		200			{object}	response.TestScheduleHeaderResponse
//		@Security		BearerAuth
//		@Router			/api/test-schedule-headers/my-schedule [get]
func (h *TestScheduleHeaderHandler) FindMySchedule(ctx *gin.Context) {
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
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find my schedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "My schedule found", res)
}
