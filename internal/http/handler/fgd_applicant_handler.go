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

type IFgdApplicantHandler interface {
	CreateOrUpdateFgdApplicants(ctx *gin.Context)
	UpdateStatusFgdApplicants(ctx *gin.Context)
	FindByUserProfileIDAndFgdID(ctx *gin.Context)
	FindAllByFgdIDPaginated(ctx *gin.Context)
	UpdateFinalResultStatusFgdApplicants(ctx *gin.Context)
}

type FgdApplicantHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IFgdApplicantUseCase
}

func NewFgdApplicantHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IFgdApplicantUseCase,
) IFgdApplicantHandler {
	return &FgdApplicantHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func FgdApplicantHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IFgdApplicantHandler {
	useCase := usecase.FgdApplicantUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewFgdApplicantHandler(log, viper, validate, useCase)
}

// CreateOrUpdateFgdApplicants create or update Fgd applicants
//
//		@Summary		Create or update Fgd applicants
//		@Description	Create or update Fgd applicants
//		@Tags			Fgd Applicants
//		@Accept			json
//	 @Produce		json
//	 @Param Fgd_applicants body request.CreateOrUpdateFgdApplicantsRequest true "Create or update Fgd applicants"
//		@Success		200	{object} response.FgdScheduleResponse
//	 @Security BearerAuth
//	 @Router /api/fgd-applicants [post]
func (h *FgdApplicantHandler) CreateOrUpdateFgdApplicants(ctx *gin.Context) {
	var req request.CreateOrUpdateFgdApplicantsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[FgdApplicantHandler.CreateOrUpdateFgdApplicants] error when binding request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[FgdApplicantHandler.CreateOrUpdateFgdApplicants] error when validating request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.CreateOrUpdateFgdApplicants(&req)
	if err != nil {
		h.Log.Errorf("[FgdApplicantHandler.CreateOrUpdateFgdApplicants] error when creating or updating Fgd applicants: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success create or update Fgd applicants", res)
}

// UpdateStatusFgdApplicants update status Fgd applicants
//
//	@Summary		Update status Fgd applicants
//	@Description	Update status Fgd applicants
//	@Tags			Fgd Applicants
//	@Accept			json
//	@Produce		json
//	@Param			Fgd_applicants	body		request.UpdateStatusFgdApplicantRequest	true	"Update status Fgd applicants"
//	@Success		200			{object}	response.FgdApplicantResponse
//	@Security		BearerAuth
//	@Router			/api/fgd-applicants/update-status [put]
func (h *FgdApplicantHandler) UpdateStatusFgdApplicants(ctx *gin.Context) {
	var req request.UpdateStatusFgdApplicantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[FgdApplicantHandler.UpdateStatusFgdApplicants] error when binding request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[FgdApplicantHandler.UpdateStatusFgdApplicants] error when validating request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	err := h.UseCase.UpdateStatusFgdApplicant(&req)
	if err != nil {
		h.Log.Errorf("[FgdApplicantHandler.UpdateStatusFgdApplicants] error when updating status Fgd applicants: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success update status Fgd applicants", nil)
}

// FindByUserProfileIDAndFgdID find Fgd applicant by user profile id and Fgd id
//
//	@Summary		Find Fgd applicant by user profile id and Fgd id
//	@Description	Find Fgd applicant by user profile id and Fgd id
//	@Tags			Fgd Applicants
//	@Accept			json
//	@Produce		json
//	@Param			user_profile_id	query	string	true	"User profile id"
//	@Param			Fgd_id	query	string	true	"Fgd id"
//	@Success		200			{object}	response.FgdApplicantResponse
//	@Security		BearerAuth
//	@Router			/api/fgd-applicants/me [get]
func (h *FgdApplicantHandler) FindByUserProfileIDAndFgdID(ctx *gin.Context) {
	userProfileID := ctx.Query("user_profile_id")
	FgdID := ctx.Query("Fgd_id")

	res, err := h.UseCase.FindByUserProfileIDAndFgdID(userProfileID, FgdID)
	if err != nil {
		h.Log.Errorf("[FgdApplicantHandler.FindByUserProfileIDAndFgdID] error when finding Fgd applicant by user profile id and Fgd id: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find Fgd applicant by user profile id and Fgd id", res)
}

// FindAllByFgdIDPaginated find all Fgd applicants by Fgd id paginated
//
//	@Summary		Find all Fgd applicants by Fgd id paginated
//	@Description	Find all Fgd applicants by Fgd id paginated
//	@Tags			Fgd Applicants
//	@Accept			json
//	@Produce		json
//	@Param			fgd_id	path	string	true	"Fgd id"
//	@Param			page	query	int	true	"Page"
//	@Param			page_size	query	int	true	"Page size"
//	@Param			search	query	string	false	"Search"
//	@Param			sort	query	string	false	"Sort"
//	@Param			filter	query	string	false	"Filter"
//	@Success		200			{object}	response.FgdApplicantResponse
//	@Security		BearerAuth
//	@Router			/api/fgd-applicants/fgd-schedule/{fgd_id} [get]
func (h *FgdApplicantHandler) FindAllByFgdIDPaginated(ctx *gin.Context) {
	FgdID := ctx.Param("fgd_id")
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

	res, total, err := h.UseCase.FindAllByFgdIDPaginated(FgdID, page, pageSize, search, sort, filter)
	if err != nil {
		h.Log.Errorf("[FgdApplicantHandler.FindAllByFgdIDPaginated] error when finding all Fgd applicants by Fgd id paginated: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find all Fgd applicants by Fgd id paginated", gin.H{
		"fgd_applicants": res,
		"total":          total,
	})
}

// UpdateFinalResultStatusFgdApplicants update final result status Fgd applicants
//
//	@Summary		Update final result status Fgd applicants
//	@Description	Update final result status Fgd applicants
//	@Tags			Fgd Applicants
//	@Accept			json
//	@Produce		json
//	@Param			Fgd_applicants	body		request.UpdateFinalResultFgdApplicantRequest	true	"Update final result status Fgd applicants"
//	@Success		200			{object}	response.FgdApplicantResponse
//	@Security		BearerAuth
//	@Router			/api/fgd-applicants/update-final-result [put]
func (h *FgdApplicantHandler) UpdateFinalResultStatusFgdApplicants(ctx *gin.Context) {
	var req request.UpdateFinalResultFgdApplicantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[FgdApplicantHandler.UpdateFinalResultStatusFgdApplicants] error when binding request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[FgdApplicantHandler.UpdateFinalResultStatusFgdApplicants] error when validating request: %s", err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	err := h.UseCase.UpdateFinalResultStatusFgdApplicant(&req)
	if err != nil {
		h.Log.Errorf("[FgdApplicantHandler.UpdateFinalResultStatusFgdApplicants] error when updating final result status Fgd applicants: %s", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success update final result status Fgd applicants", nil)
}
