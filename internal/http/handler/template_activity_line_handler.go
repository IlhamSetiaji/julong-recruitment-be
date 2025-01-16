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

type ITemplateActivityLineHandler interface {
	CreateOrUpdateTemplateActivityLine(ctx *gin.Context)
}

type TemplateActivityLineHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITemplateActivityLineUseCase
}

func NewTemplateActivityLineHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITemplateActivityLineUseCase,
) ITemplateActivityLineHandler {
	return &TemplateActivityLineHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TemplateActivityLineHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITemplateActivityLineHandler {
	useCase := usecase.TemplateActivityLineUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewTemplateActivityLineHandler(log, viper, validate, useCase)
}

func (h *TemplateActivityLineHandler) CreateOrUpdateTemplateActivityLine(ctx *gin.Context) {
	var payload request.CreateOrUpdateTemplateActivityLineRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("[TemplateActivityLineHandler.CreateOrUpdateTemplateActivityLine] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[TemplateActivityLineHandler.CreateOrUpdateTemplateActivityLine] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	ta, err := h.UseCase.CreateOrUpdateTemplateActivityLine(&payload)
	if err != nil {
		h.Log.Error("[TemplateActivityLineHandler.CreateOrUpdateTemplateActivityLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", ta)
}
