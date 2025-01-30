package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentTypeHandler interface {
	FindAll(ctx *gin.Context)
}

type DocumentTypeHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IDocumentTypeUseCase
}

func NewDocumentTypeHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDocumentTypeUseCase,
) IDocumentTypeHandler {
	return &DocumentTypeHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func DocumentTypeHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IDocumentTypeHandler {
	useCase := usecase.DocumentTypeUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewDocumentTypeHandler(log, viper, validate, useCase)
}

// FindAll find all document types
//
//		@Summary		Find all document types
//		@Description	Find all document types
//		@Tags			Document Types
//		@Accept			json
//		@Produce		json
//		@Success		200	{object} response.DocumentTypeResponse
//	 @Security BearerAuth
//		@Router			/document-types [get]
func (h *DocumentTypeHandler) FindAll(ctx *gin.Context) {
	documentTypes, err := h.UseCase.GetAllDocumentType()
	if err != nil {
		h.Log.Errorf("[DocumentTypeHandler.FindAll] error when getting document types: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when getting document types", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find all document types", documentTypes)
}
