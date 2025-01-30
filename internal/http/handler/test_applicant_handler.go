package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITestApplicantHandler interface {
	CreateOrUpdateTestApplicants(ctx *gin.Context)
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
