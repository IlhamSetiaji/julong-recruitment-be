package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentVerificationLineHandler interface {
	CreateOrUpdateDocumentVerificationLine(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindAllByDocumentVerificationHeaderID(ctx *gin.Context)
	UploadDocumentVerificationLine(ctx *gin.Context)
	UpdateAnswer(ctx *gin.Context)
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

// UploadDocumentVerificationLine upload document verification line
//
// @Summary Upload document verification line
// @Description Upload document verification line
// @Tags Document Verification Lines
// @Accept multipart/form-data
// @Produce json
// @Param id formData string true "Document Verification Line ID"
// @Param file formData file true "File"
// @Success 200 {object} response.DocumentVerificationLineResponse
// @Security BearerAuth
// @Router /document-verification-lines/upload [post]
func (h *DocumentVerificationLineHandler) UploadDocumentVerificationLine(ctx *gin.Context) {
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		h.Log.Error("Failed to parse form-data: ", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	id := ctx.PostForm("id")
	if id == "" {
		utils.BadRequestResponse(ctx, "bad request", "id is required")
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		h.Log.Error("Failed to get file: ", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	timestamp := time.Now().UnixNano()
	filePath := "storage/document_verification_lines/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		h.Log.Error("Failed to save file: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	res, err := h.UseCase.UploadDocumentVerificationLine(&request.UploadDocumentVerificationLine{
		ID:   id,
		Path: filePath,
	})
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success", res)
}

// Update Answer
//
// @Summary Update Answer
// @Description Update Answer
// @Tags Document Verification Lines
// @Accept json
// @Produce json
// @Param id path string true "Document Verification Line ID"
// @Param payload body request.UpdateAnswer true "Update Answer"
// @Success 200 {object} response.DocumentVerificationLineResponse
// @Security BearerAuth
// @Router /document-verification-lines/{id}/answer [put]
func (h *DocumentVerificationLineHandler) UpdateAnswer(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("Invalid UUID: ", err)
		utils.BadRequestResponse(ctx, "bad request", "Invalid UUID format")
		return
	}

	var payload request.UpdateAnswer
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

	res, err := h.UseCase.UpdateAnswer(parsedID, &payload)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success", res)
}
