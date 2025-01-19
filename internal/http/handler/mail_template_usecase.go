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

type IMailTemplateHandler interface {
	CreateMailTemplate(ctx *gin.Context)
}

type MailTemplateHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IMailTemplateUseCase
}

func NewMailTemplateHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IMailTemplateUseCase,
) IMailTemplateHandler {
	return &MailTemplateHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func MailTemplateHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IMailTemplateHandler {
	useCase := usecase.MailTemplateUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewMailTemplateHandler(log, viper, validate, useCase)
}

func (h *MailTemplateHandler) CreateMailTemplate(ctx *gin.Context) {
	var payload request.CreateMailTemplateRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[MailTemplateHandler.CreateMailTemplate] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, "Bad request", err)
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[MailTemplateHandler.CreateMailTemplate] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, "Bad request", err)
		return
	}

	res, err := h.UseCase.CreateMailTemplate(&payload)
	if err != nil {
		h.Log.Errorf("[MailTemplateHandler.CreateMailTemplate] error when creating mail template: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create mail template", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Success created mail template", res)
}
