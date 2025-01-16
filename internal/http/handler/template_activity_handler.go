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

type ITemplateActivityHandler interface {
	CreateTemplateActivity(ctx *gin.Context)
}

type TemplateActivityHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITemplateActivityUseCase
}

func NewTemplateActivityHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITemplateActivityUseCase,
) ITemplateActivityHandler {
	return &TemplateActivityHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TemplateActivityHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITemplateActivityHandler {
	useCase := usecase.TemplateActivityUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewTemplateActivityHandler(log, viper, validate, useCase)
}

func (h *TemplateActivityHandler) CreateTemplateActivity(ctx *gin.Context) {
	var req request.CreateTemplateActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[TemplateActivityHandler.CreateTemplateActivity] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[TemplateActivityHandler.CreateTemplateActivity] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	resp, err := h.UseCase.CreateTemplateActivity(&req)
	if err != nil {
		h.Log.Error("[TemplateActivityHandler.CreateTemplateActivity] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success", resp)
}
