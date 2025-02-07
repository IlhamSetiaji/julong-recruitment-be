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
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IInterviewApplicantHandler interface {
	CreateOrUpdateInterviewApplicants(ctx *gin.Context)
	UpdateStatusInterviewApplicants(ctx *gin.Context)
	FindByUserProfileIDAndInterviewID(ctx *gin.Context)
	FindAllByInterviewIDPaginated(ctx *gin.Context)
	UpdateFinalResultStatusInterviewApplicants(ctx *gin.Context)
}

type InterviewApplicantHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IInterviewApplicantUseCase
}

func NewInterviewApplicantHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IInterviewApplicantUseCase,
) IInterviewApplicantHandler {
	return &InterviewApplicantHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func InterviewApplicantHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IInterviewApplicantHandler {
	useCase := usecase.InterviewApplicantUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewInterviewApplicantHandler(log, viper, validate, useCase)
}

// CreateOrUpdateInterviewApplicants create or update interview applicants
//
//		@Summary		Create or update interview applicants
//		@Description	Create or update interview applicants
//		@Tags			Interview Applicants
//		@Accept			json
//	 @Produce		json
//	 @Param interview_applicants body request.CreateOrUpdateInterviewApplicantsRequest true "Create or update interview applicants"
//		@Success		200	{object} response.InterviewResponse
//	 @Security BearerAuth
//	 @Router /api/interview-applicants [post]
func (h *InterviewApplicantHandler) CreateOrUpdateInterviewApplicants(ctx *gin.Context) {
	var req request.CreateOrUpdateInterviewApplicantsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.CreateOrUpdateInterviewApplicants] error when binding request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.CreateOrUpdateInterviewApplicants] error when validating request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.CreateOrUpdateInterviewApplicants(&req)
	if err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.CreateOrUpdateInterviewApplicants] error when creating or updating interview applicants: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success create or update interview applicants", res)
}

// UpdateStatusInterviewApplicants update status interview applicants
//
//	@Summary		Update status interview applicants
//	@Description	Update status interview applicants
//	@Tags			Interview Applicants
//	@Accept			json
//	@Produce		json
//	@Param			interview_applicants	body		request.UpdateStatusInterviewApplicantRequest	true	"Update status interview applicants"
//	@Success		200			{object}	response.InterviewResponse
//	@Security		BearerAuth
//	@Router			/api/interview-applicants/update-status [put]
func (h *InterviewApplicantHandler) UpdateStatusInterviewApplicants(ctx *gin.Context) {
	var req request.UpdateStatusInterviewApplicantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.UpdateStatusInterviewApplicants] error when binding request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.UpdateStatusInterviewApplicants] error when validating request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	err := h.UseCase.UpdateStatusInterviewApplicant(&req)
	if err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.UpdateStatusInterviewApplicants] error when updating status interview applicants: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success update status interview applicants", nil)
}

// FindByUserProfileIDAndInterviewID find interview applicant by user profile id and interview id
//
//	@Summary		Find interview applicant by user profile id and interview id
//	@Description	Find interview applicant by user profile id and interview id
//	@Tags			Interview Applicants
//	@Accept			json
//	@Produce		json
//	@Param			user_profile_id	query	string	true	"User profile id"
//	@Param			interview_id	query	string	true	"Interview id"
//	@Success		200			{object}	response.InterviewApplicantResponse
//	@Security		BearerAuth
//	@Router			/api/interview-applicants/me [get]
func (h *InterviewApplicantHandler) FindByUserProfileIDAndInterviewID(ctx *gin.Context) {
	userProfileID := ctx.Query("user_profile_id")
	interviewID := ctx.Query("interview_id")

	res, err := h.UseCase.FindByUserProfileIDAndInterviewID(userProfileID, interviewID)
	if err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.FindByUserProfileIDAndInterviewID] error when finding interview applicant by user profile id and interview id: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find interview applicant by user profile id and interview id", res)
}

// FindAllByInterviewIDPaginated find all interview applicants by interview id paginated
//
//	@Summary		Find all interview applicants by interview id paginated
//	@Description	Find all interview applicants by interview id paginated
//	@Tags			Interview Applicants
//	@Accept			json
//	@Produce		json
//	@Param			interview_id	path	string	true	"Interview id"
//	@Param			page	query	int	true	"Page"
//	@Param			page_size	query	int	true	"Page size"
//	@Param			search	query	string	false	"Search"
//	@Param			sort	query	string	false	"Sort"
//	@Param			filter	query	string	false	"Filter"
//	@Success		200			{object}	response.InterviewApplicantResponse
//	@Security		BearerAuth
//	@Router			/api/interview-applicants [get]
func (h *InterviewApplicantHandler) FindAllByInterviewIDPaginated(ctx *gin.Context) {
	interviewID := ctx.Param("interview_id")
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
		if status != "" {
			filter["status"] = status
		}
	}

	res, total, err := h.UseCase.FindAllByInterviewIDPaginated(interviewID, page, pageSize, search, sort, filter)
	if err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.FindAllByInterviewIDPaginated] error when finding all interview applicants by interview id paginated: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find all interview applicants by interview id paginated", gin.H{
		"interview_applicants": res,
		"total":                total,
	})
}

// UpdateFinalResultStatusInterviewApplicants update final result status interview applicants
//
//	@Summary		Update final result status interview applicants
//	@Description	Update final result status interview applicants
//	@Tags			Interview Applicants
//	@Accept			json
//	@Produce		json
//	@Param			interview_applicants	body		request.UpdateFinalResultInterviewApplicantRequest	true	"Update final result status interview applicants"
//	@Success		200			{object}	response.InterviewResponse
//	@Security		BearerAuth
//	@Router			/api/interview-applicants/update-final-result [put]
func (h *InterviewApplicantHandler) UpdateFinalResultStatusInterviewApplicants(ctx *gin.Context) {
	var req request.UpdateFinalResultInterviewApplicantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.UpdateFinalResultStatusInterviewApplicants] error when binding request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.UpdateFinalResultStatusInterviewApplicants] error when validating request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	err := h.UseCase.UpdateFinalResultStatusInterviewApplicant(&req)
	if err != nil {
		h.Log.Errorf("[InterviewApplicantHandler.UpdateFinalResultStatusInterviewApplicants] error when updating final result status interview applicants: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success update final result status interview applicants", nil)
}
