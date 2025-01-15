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

type IDocumentSetupHandler interface {
	CreateDocumentSetup(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateDocumentSetup(ctx *gin.Context)
	DeleteDocumentSetup(ctx *gin.Context)
	FindByDocumentTypeName(ctx *gin.Context)
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

func (h *DocumentSetupHandler) FindAllPaginated(ctx *gin.Context) {
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

	documentSetups, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Errorf("[DocumentSetupHandler.FindAllPaginated] error when getting document setups: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when getting document setups", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find all document setups", gin.H{
		"document_setups": documentSetups,
		"total":           total,
	})
}

func (h *DocumentSetupHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	documentSetup, err := h.UseCase.FindByID(parsedUUID)
	if err != nil {
		h.Log.Errorf("[DocumentSetupHandler.FindByID] error when getting document setup: %v", err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when getting document setup", err.Error())
		return
	}

	if documentSetup == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "document setup not found", "document setup not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find document setup", documentSetup)
}

func (h *DocumentSetupHandler) UpdateDocumentSetup(ctx *gin.Context) {
	var payload request.UpdateDocumentSetupRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[DocumentSetupHandler.UpdateDocumentSetup] error when binding request: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error when binding request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[DocumentSetupHandler.UpdateDocumentSetup] error when validating request: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error when validating request", err.Error())
		return
	}

	exist, err := h.UseCase.FindByID(uuid.MustParse(payload.ID))
	if err != nil {
		h.Log.Errorf("[DocumentSetupHandler.UpdateDocumentSetup] error when getting document setup: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when getting document setup", err.Error())
		return
	}

	if exist == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "document setup not found", "document setup not found")
		return
	}

	documentSetup, err := h.UseCase.UpdateDocumentSetup(&payload)
	if err != nil {
		h.Log.Errorf("[DocumentSetupHandler.UpdateDocumentSetup] error when updating document setup: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when updating document setup", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success update document setup", documentSetup)
}

func (h *DocumentSetupHandler) DeleteDocumentSetup(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	exist, err := h.UseCase.FindByID(parsedUUID)
	if err != nil {
		h.Log.Errorf("[DocumentSetupHandler.DeleteDocumentSetup] error when getting document setup: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when getting document setup", err.Error())
		return
	}

	if exist == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "document setup not found", "document setup not found")
		return
	}

	err = h.UseCase.DeleteDocumentSetup(parsedUUID)
	if err != nil {
		h.Log.Errorf("[DocumentSetupHandler.DeleteDocumentSetup] error when deleting document setup: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when deleting document setup", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success delete document setup", nil)
}

func (h *DocumentSetupHandler) FindByDocumentTypeName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "name is required", "name is required")
		return
	}

	documentSetups, err := h.UseCase.FindByDocumentTypeName(name)
	if err != nil {
		h.Log.Errorf("[DocumentSetupHandler.FindByDocumentTypeName] error when getting document setups: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when getting document setups", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success find document setups by document type name", documentSetups)
}
