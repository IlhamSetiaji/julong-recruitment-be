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

type ITestApplicantHandler interface {
	CreateOrUpdateTestApplicants(ctx *gin.Context)
	UpdateStatusTestApplicants(ctx *gin.Context)
	FindByUserProfileIDAndTestScheduleHeaderID(ctx *gin.Context)
	FindAllByTestScheduleHeaderIDPaginated(ctx *gin.Context)
}

type TestApplicantHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITestApplicantUseCase
}

func NewTestApplicantHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITestApplicantUseCase,
) ITestApplicantHandler {
	return &TestApplicantHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TestApplicantHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITestApplicantHandler {
	useCase := usecase.TestApplicantUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewTestApplicantHandler(log, viper, validate, useCase)
}

// CreateOrUpdateTestApplicants create or update test applicants
//
//	@Summary		Create or update test applicants
//	@Description	Create or update test applicants
//	@Tags			Test Applicants
//	@Accept			json
//	@Produce		json
//	@Param			employee	body		request.CreateOrUpdateTestApplicantsRequest	true	"Create test applicants"
//	@Success		201			{object}	response.TestScheduleHeaderResponse
//	@Security		BearerAuth
//	@Router			/api/test-applicants [post]
func (h *TestApplicantHandler) CreateOrUpdateTestApplicants(ctx *gin.Context) {
	var req request.CreateOrUpdateTestApplicantsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[TestApplicantHandler.CreateOrUpdateTestApplicants] error when binding request: %s", err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[TestApplicantHandler.CreateOrUpdateTestApplicants] error when validating request: %s", err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.CreateOrUpdateTestApplicants(&req)
	if err != nil {
		h.Log.Errorf("[TestApplicantHandler.CreateOrUpdateTestApplicants] error when creating or updating test applicants: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success create or update test applicants", res)
}

// UpdateStatusTestApplicants update status test applicants
//
//	@Summary		Update status test applicants
//	@Description	Update status test applicants
//	@Tags			Test Applicants
//	@Accept			json
//	@Produce		json
//	@Param			employee	body		request.UpdateStatusTestApplicantRequest	true	"Update status test applicants"
//	@Success		200			{object}	response.TestScheduleHeaderResponse
//	@Security		BearerAuth
//	@Router			/api/test-applicants/update-status [put]
func (h *TestApplicantHandler) UpdateStatusTestApplicants(ctx *gin.Context) {
	var req request.UpdateStatusTestApplicantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[TestApplicantHandler.UpdateStatusTestApplicants] error when binding request: %s", err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[TestApplicantHandler.UpdateStatusTestApplicants] error when validating request: %s", err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.UpdateStatusTestApplicant(&req)
	if err != nil {
		h.Log.Errorf("[TestApplicantHandler.UpdateStatusTestApplicants] error when updating status test applicants: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success update status test applicants", res)
}

// FindByUserProfileIDAndTestScheduleHeaderID find test applicant by user profile id and test schedule header id
//
//	@Summary		Find test applicant by user profile id and test schedule header id
//	@Description	Find test applicant by user profile id and test schedule header id
//	@Tags			Test Applicants
//	@Accept			json
//	@Produce		json
//	@Param			user_profile_id	query	string	true	"User profile id"
//	@Param			test_schedule_header_id	query	string	true	"Test schedule header id"
//	@Success		200			{object}	response.TestApplicantResponse
//	@Security		BearerAuth
//	@Router			/api/test-applicants/me [get]
func (h *TestApplicantHandler) FindByUserProfileIDAndTestScheduleHeaderID(ctx *gin.Context) {
	userProfileID := ctx.Query("user_profile_id")
	testScheduleHeaderID := ctx.Query("test_schedule_header_id")

	if userProfileID == "" || testScheduleHeaderID == "" {
		h.Log.Errorf("[TestApplicantHandler.FindByUserProfileIDAndTestScheduleHeaderID] error when binding request: user profile id or test schedule header id is empty")
		utils.BadRequestResponse(ctx, "bad request", "user profile id or test schedule header id is empty")
		return
	}

	parsedUserProfileID, err := uuid.Parse(userProfileID)
	if err != nil {
		h.Log.Errorf("[TestApplicantHandler.FindByUserProfileIDAndTestScheduleHeaderID] error when parsing user profile id: %s", err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	parsedTestScheduleHeaderID, err := uuid.Parse(testScheduleHeaderID)
	if err != nil {
		h.Log.Errorf("[TestApplicantHandler.FindByUserProfileIDAndTestScheduleHeaderID] error when parsing test schedule header id: %s", err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.FindByUserProfileIDAndTestScheduleHeaderID(parsedUserProfileID, parsedTestScheduleHeaderID)
	if err != nil {
		h.Log.Errorf("[TestApplicantHandler.FindByUserProfileIDAndTestScheduleHeaderID] error when finding test applicant by user profile id and test schedule header id: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find test applicant by user profile id and test schedule header id", res)
}

// FindAllByTestScheduleHeaderIDPaginated find all test applicants by test schedule header id paginated
//
//	@Summary		Find all test applicants by test schedule header id paginated
//	@Description	Find all test applicants by test schedule header id paginated
//	@Tags			Test Applicants
//	@Accept			json
//	@Produce		json
//	@Param			test_schedule_header_id	path	string	true	"Test schedule header id"
//	@Param			page	query	int	true	"Page"
//	@Param			page_size	query	int	true	"Page size"
//	@Param			search	query	string	false	"Search"
//	@Param			sort	query	string	false	"Sort"
//	@Param			filter	query	string	false	"Filter"
//	@Success		200			{object}	response.TestApplicantResponse
//	@Security		BearerAuth
//	@Router			/api/test-applicants [get]
func (h *TestApplicantHandler) FindAllByTestScheduleHeaderIDPaginated(ctx *gin.Context) {
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

	filter := make(map[string]interface{})
	status := ctx.Query("status")
	if status != "" {
		filter["status"] = status
	}

	testScheduleHeaderID := ctx.Param("test_schedule_header_id")

	parsedTestScheduleHeaderID, err := uuid.Parse(testScheduleHeaderID)
	if err != nil {
		h.Log.Errorf("[TestApplicantHandler.FindAllByTestScheduleHeaderIDPaginated] error when parsing test schedule header id: %s", err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, total, err := h.UseCase.FindAllByTestScheduleHeaderIDPaginated(parsedTestScheduleHeaderID, page, pageSize, search, sort, filter)
	if err != nil {
		h.Log.Errorf("[TestApplicantHandler.FindAllByTestScheduleHeaderIDPaginated] error when finding all test applicants by test schedule header id paginated: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find all test applicants by test schedule header id paginated", gin.H{
		"test_applicants": res,
		"total":           total,
	})
}
