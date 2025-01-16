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

type IDocumentVerificationHandler interface {
	CreateDocumentVerification(ctx *gin.Context)
}

type DocumentVerificationHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IDocumentVerificationUseCase
}

func NewDocumentVerificationHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDocumentVerificationUseCase,
) IDocumentVerificationHandler {
	return &DocumentVerificationHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func DocumentVerificationHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IDocumentVerificationHandler {
	useCase := usecase.DocumentVerificationUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewDocumentVerificationHandler(log, viper, validate, useCase)
}

func (h *DocumentVerificationHandler) CreateDocumentVerification(ctx *gin.Context) {
	var payload request.CreateDocumentVerificationRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.CreateDocumentVerification] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, "Invalid request payload", err)
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.CreateDocumentVerification] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, "Invalid request payload", err)
		return
	}

	res, err := h.UseCase.CreateDocumentVerification(&payload)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.CreateDocumentVerification] error when creating document verification: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when creating document verification", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Document verification created successfully", res)
}
