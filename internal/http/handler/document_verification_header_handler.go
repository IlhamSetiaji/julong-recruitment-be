package handler

import (
	"net/http"
	"strconv"

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

type IDocumentVerificationHeaderHandler interface {
	CreateDocumentVerificationHeader(ctx *gin.Context)
	UpdateDocumentVerificationHeader(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	DeleteDocumentVerificationHeader(ctx *gin.Context)
}

type DocumentVerificationHeaderHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IDocumentVerificationHeaderUseCase
}

func NewDocumentVerificationHeaderHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDocumentVerificationHeaderUseCase,
) IDocumentVerificationHeaderHandler {
	return &DocumentVerificationHeaderHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func DocumentVerificationHeaderHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IDocumentVerificationHeaderHandler {
	useCase := usecase.DocumentVerificationHeaderUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewDocumentVerificationHeaderHandler(log, viper, validate, useCase)
}

// CreateDocumentVerificationHeader create document verification header
//
//		@Summary		Create document verification header
//		@Description	Create document verification header
//		@Tags			Document Verification Header
//		@Accept			json
//	 @Produce		json
//	 @Param payload body request.CreateDocumentVerificationHeaderRequest true "Create Document Verification Header"
//		@Success		200	{object} response.DocumentVerificationHeaderResponse
//		@Router			/document-verification-headers [post]
func (h *DocumentVerificationHeaderHandler) CreateDocumentVerificationHeader(ctx *gin.Context) {
	req := new(request.CreateDocumentVerificationHeaderRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	res, err := h.UseCase.CreateDocumentVerificationHeader(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when creating document verification header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Document verification header created", res)
}

// UpdateDocumentVerificationHeader update document verification header
//
//	@Summary		Update document verification header
//	@Description	Update document verification header
//	@Tags			Document Verification Header
//	@Accept			json
//	@Produce		json
//	@Param			payload body request.UpdateDocumentVerificationHeaderRequest true "Update Document Verification Header"
//	@Success		200	{object} response.DocumentVerificationHeaderResponse
//	@Router			/document-verification-headers/update [put]
func (h *DocumentVerificationHeaderHandler) UpdateDocumentVerificationHeader(ctx *gin.Context) {
	req := new(request.UpdateDocumentVerificationHeaderRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	res, err := h.UseCase.UpdateDocumentVerificationHeader(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when updating document verification header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document verification header updated", res)
}

// FindByID find document verification header by id
//
//	@Success		200	{object} response.DocumentVerificationHeaderResponse
//	@Security BearerAuth
//	@Router			/document-verification-headers/{id} [get]
func (h *DocumentVerificationHeaderHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.FindByID] error when parsing id: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Error when parsing id", err.Error())
		return
	}

	res, err := h.UseCase.FindByID(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.FindByID] error when finding document verification header by id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding document verification header by id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get document verification header by id", res)
}

// FindAllPaginated find all document verification header with pagination
//
//	@Summary		Find all document verification header with pagination
//	@Description	Find all document verification header with pagination
//	@Tags			Document Verification Header
//	@Accept			json
//	@Produce		json
//	@Param			page query int false "Page"
//	@Param			page_size query int false "Page Size"
//	@Param			search query string false "Search"
//	@Param			sort query string false "Sort"
//	@Success		200	{object} response.DocumentVerificationHeaderResponse
//	@Security BearerAuth
//	@Router			/document-verification-headers [get]
func (h *DocumentVerificationHeaderHandler) FindAllPaginated(ctx *gin.Context) {
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
		h.Log.Errorf("[DocumentVerificationHeaderHandler.FindAllPaginated] error when finding all document verification header: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding all document verification header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get all document verification header", gin.H{
		"document_verification_headers": res,
		"total":                         total,
	})
}

// DeleteDocumentVerificationHeader delete document verification header by id
//
//	@Summary		Delete document verification header by id
//	@Description	Delete document verification header by id
//	@Tags			Document Verification Header
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "ID"
//	@Success		200	{string} string
//	@Security BearerAuth
//	@Router			/document-verification-headers/{id} [delete]
func (h *DocumentVerificationHeaderHandler) DeleteDocumentVerificationHeader(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.DeleteDocumentVerificationHeader] error when parsing id: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Error when parsing id", err.Error())
		return
	}

	err = h.UseCase.DeleteDocumentVerificationHeader(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.DeleteDocumentVerificationHeader] error when deleting document verification header by id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when deleting document verification header by id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document verification header deleted", nil)
}
