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

type IDocumentVerificationLineHandler interface {
	CreateOrUpdateDocumentVerificationLine(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindAllByDocumentVerificationHeaderID(ctx *gin.Context)
}

type DocumentVerificationLineHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IDocumentVerificationLineUsecase
}

func NewDocumentVerificationLineHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDocumentVerificationLineUsecase,
) IDocumentVerificationLineHandler {
	return &DocumentVerificationLineHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func DocumentVerificationLineHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IDocumentVerificationLineHandler {
	useCase := usecase.DocumentVerificationLineFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewDocumentVerificationLineHandler(log, viper, validate, useCase)
}

// CreateOrUpdateDocumentVerificationLine create or update document verification line
//
// @Summary Create or update document verification line
// @Description Create or update document verification line
// @Tags Document Verification Lines
// @Accept json
// @Produce json
// @Param payload body request.CreateOrUpdateDocumentVerificationLine true "Create or update document verification line"
// @Success 200 {object} response.DocumentVerificationHeaderResponse
// @Security BearerAuth
// @Router /document-verification-lines [post]
func (h *DocumentVerificationLineHandler) CreateOrUpdateDocumentVerificationLine(ctx *gin.Context) {
	var payload request.CreateOrUpdateDocumentVerificationLine
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	res, err := h.UseCase.CreateOrUpdateDocumentVerificationLine(&payload)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success", res)
}

// FindByID find document verification line by id
//
// @Summary Find document verification line by id
// @Description Find document verification line by id
// @Tags Document Verification Lines
// @Accept json
// @Produce json
// @Param id path string true "Document Verification Line ID"
// @Success 200 {object} response.DocumentVerificationLineResponse
// @Security BearerAuth
// @Router /document-verification-lines/{id} [get]
func (h *DocumentVerificationLineHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := h.UseCase.FindByID(id)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Not Found", "Document Verification Line not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success", res)
}

// FindAllByDocumentVerificationHeaderID find all document verification lines by document verification header id
//
// @Summary Find all document verification lines by document verification header id
// @Description Find all document verification lines by document verification header id
// @Tags Document Verification Lines
// @Accept json
// @Produce json
// @Param document_verification_header_id path string true "Document Verification Header ID"
// @Success 200 {array} response.DocumentVerificationLineResponse
// @Security BearerAuth
// @Router /document-verification-lines/document-verification-header/{document_verification_header_id} [get]
func (h *DocumentVerificationLineHandler) FindAllByDocumentVerificationHeaderID(ctx *gin.Context) {
	documentVerificationHeaderID := ctx.Param("document_verification_header_id")

	res, err := h.UseCase.FindAllByDocumentVerificationHeaderID(documentVerificationHeaderID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success", res)
}
