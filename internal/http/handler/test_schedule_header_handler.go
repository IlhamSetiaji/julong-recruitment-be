package handler

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
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
}

type TestScheduleHeaderHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITestScheduleHeaderUsecase
}

func NewTestScheduleHeaderHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITestScheduleHeaderUsecase,
) ITestScheduleHeaderHandler {
	return &TestScheduleHeaderHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TestScheduleHeaderHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITestScheduleHeaderHandler {
	useCase := usecase.TestScheduleHeaderUsecaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewTestScheduleHeaderHandler(log, viper, validate, useCase)
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
		utils.BadRequestResponse(ctx, "Invalid request body", err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid request body", err)
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
		utils.BadRequestResponse(ctx, "Invalid request body", err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid request body", err)
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
