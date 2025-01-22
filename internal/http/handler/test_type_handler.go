package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase.go"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITestTypeHandler interface {
	CreateTestType(ctx *gin.Context)
}

type TestTypeHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITestTypeUseCase
}

func NewTestTypeHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITestTypeUseCase,
) ITestTypeHandler {
	return &TestTypeHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TestTypeHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITestTypeHandler {
	useCase := usecase.TestTypeUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewTestTypeHandler(log, viper, validate, useCase)
}

func (h *TestTypeHandler) CreateTestType(ctx *gin.Context) {
	var payload request.CreateTestTypeRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[TestTypeHandler.CreateTestType] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[TestTypeHandler.CreateTestType] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.CreateTestType(&payload)
	if err != nil {
		h.Log.Errorf("[TestTypeHandler.CreateTestType] error when creating test type: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to create test type", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "test type created", res)
}
