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

type IDocumentSetupHandler interface {
	CreateDocumentSetup(ctx *gin.Context)
}

type DocumentSetupHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IDocumentSetupUseCase
}

func NewDocumentSetupHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDocumentSetupUseCase,
) IDocumentSetupHandler {
	return &DocumentSetupHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func DocumentSetupHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IDocumentSetupHandler {
	useCase := usecase.DocumentSetupUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewDocumentSetupHandler(log, viper, validate, useCase)
}

func (h *DocumentSetupHandler) CreateDocumentSetup(ctx *gin.Context) {
	var payload request.CreateDocumentSetupRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[DocumentSetupHandler.CreateDocumentSetup] error when binding request: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error when binding request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[DocumentSetupHandler.CreateDocumentSetup] error when validating request: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error when validating request", err.Error())
		return
	}

	documentSetup, err := h.UseCase.CreateDocumentSetup(&payload)
	if err != nil {
		h.Log.Errorf("[DocumentSetupHandler.CreateDocumentSetup] error when creating document setup: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when creating document setup", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success create document setup", documentSetup)
}
