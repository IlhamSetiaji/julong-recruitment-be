package handler

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase.go"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentVerificationHandler interface {
	CreateDocumentVerification(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateDocumentVerification(ctx *gin.Context)
	FindByTemplateQuestionID(ctx *gin.Context)
	DeleteDocumentVerification(ctx *gin.Context)
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

func (h *DocumentVerificationHandler) FindAllPaginated(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	createdAt := ctx.Query("created_at")
	if createdAt == "" {
		createdAt = "DESC"
	}

	sort := map[string]interface{}{
		"created_at": createdAt,
	}

	res, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.FindAllPaginated] error when finding all document verification: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding all document verification", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get all document verification", gin.H{
		"document_verifications": res,
		"total":                  total,
	})
}

func (h *DocumentVerificationHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.FindByID] error when parsing id: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Error when parsing id", err.Error())
		return
	}

	res, err := h.UseCase.FindByID(parsedUUID)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.FindByID] error when finding document verification by id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding document verification by id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get document verification by id", res)
}

func (h *DocumentVerificationHandler) UpdateDocumentVerification(ctx *gin.Context) {
	var payload request.UpdateDocumentVerificationRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.UpdateDocumentVerification] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, "Invalid request payload", err)
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.UpdateDocumentVerification] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, "Invalid request payload", err)
		return
	}

	res, err := h.UseCase.UpdateDocumentVerification(&payload)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.UpdateDocumentVerification] error when updating document verification: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when updating document verification", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document verification updated successfully", res)
}

func (h *DocumentVerificationHandler) FindByTemplateQuestionID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.FindByTemplateQuestionID] error when parsing id: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Error when parsing id", err.Error())
		return
	}

	res, err := h.UseCase.FindByTemplateQuestionID(parsedUUID)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.FindByTemplateQuestionID] error when finding document verification by template question id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding document verification by template question id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get document verification by template question id", res)
}

func (h *DocumentVerificationHandler) DeleteDocumentVerification(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.DeleteDocumentVerification] error when parsing id: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Error when parsing id", err.Error())
		return
	}

	err = h.UseCase.DeleteDocumentVerification(parsedUUID)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHandler.DeleteDocumentVerification] error when deleting document verification: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when deleting document verification", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document verification deleted successfully", nil)
}
