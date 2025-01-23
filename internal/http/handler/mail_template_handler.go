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

type IMailTemplateHandler interface {
	CreateMailTemplate(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateMailTemplate(ctx *gin.Context)
	DeleteMailTemplate(ctx *gin.Context)
	FindAllByDocumentTypeID(ctx *gin.Context)
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

// CreateMailTemplate create mail template
//
//	@Summary		Create mail template
//	@Description	Create mail template
//	@Tags			Mail Templates
//	@Accept			json
//	@Produce		json
//	@Param			payload body request.CreateMailTemplateRequest true "Create Mail Template"
//	@Success		201	{object} response.MailTemplateResponse
//
// @Security BearerAuth
//
//	@Router			/mail-templates [post]
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

// FindAllPaginated find all mail templates paginated
//
//	@Summary		Find all mail templates paginated
//	@Description	Find all mail templates paginated
//	@Tags			Mail Templates
//	@Accept			json
//	@Produce		json
//	@Param			page	query	int	false	"Page"
//	@Param			page_size	query	int	false	"Page Size"
//	@Param			search	query	string	false	"Search"
//	@Param			created_at	query	string	false	"Created At"
//
// @Security BearerAuth
//
//	@Router			/mail-templates [get]
func (h *MailTemplateHandler) FindAllPaginated(ctx *gin.Context) {
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
		h.Log.Errorf("[MailTemplateHandler.FindAllPaginated] error when getting mail templates: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get mail templates", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get mail templates", gin.H{
		"mail_templates": res,
		"total":          total,
	})
}

// FindByID find mail template by ID
//
//	@Summary		Find mail template by ID
//	@Description	Find mail template by ID
//	@Tags			Mail Templates
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//
// @Security BearerAuth
//
//	@Router			/mail-templates/{id} [get]
func (h *MailTemplateHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	mtID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[MailTemplateHandler.FindByID] error when parsing UUID: %v", err)
		utils.BadRequestResponse(ctx, "Bad request", err)
		return
	}

	res, err := h.UseCase.FindByID(mtID)
	if err != nil {
		h.Log.Errorf("[MailTemplateHandler.FindByID] error when getting mail template: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get mail template", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Mail template not found", "Mail template not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get mail template", res)
}

func (h *MailTemplateHandler) UpdateMailTemplate(ctx *gin.Context) {
	var payload request.UpdateMailTemplateRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[MailTemplateHandler.UpdateMailTemplate] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, "Bad request", err)
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[MailTemplateHandler.UpdateMailTemplate] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, "Bad request", err)
		return
	}

	exist, err := h.UseCase.FindByID(uuid.MustParse(payload.ID))
	if err != nil {
		h.Log.Errorf("[MailTemplateHandler.UpdateMailTemplate] error when finding mail template by ID: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find mail template by ID", err.Error())
		return
	}

	if exist == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Mail template not found", "Mail template not found")
		return
	}

	res, err := h.UseCase.UpdateMailTemplate(&payload)
	if err != nil {
		h.Log.Errorf("[MailTemplateHandler.UpdateMailTemplate] error when updating mail template: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update mail template", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success update mail template", res)
}

// DeleteMailTemplate delete mail template
//
//	@Summary		Delete mail template
//	@Description	Delete mail template
//	@Tags			Mail Templates
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200			{string}	string
//
// @Security BearerAuth
//
//	@Router			/mail-templates/{id} [delete]
func (h *MailTemplateHandler) DeleteMailTemplate(ctx *gin.Context) {
	id := ctx.Param("id")
	mtID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[MailTemplateHandler.DeleteMailTemplate] error when parsing UUID: %v", err)
		utils.BadRequestResponse(ctx, "Bad request", err)
		return
	}

	exist, err := h.UseCase.FindByID(mtID)
	if err != nil {
		h.Log.Errorf("[MailTemplateHandler.DeleteMailTemplate] error when finding mail template by ID: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find mail template by ID", err.Error())
		return
	}

	if exist == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Mail template not found", "Mail template not found")
		return
	}

	err = h.UseCase.DeleteMailTemplate(mtID)
	if err != nil {
		h.Log.Errorf("[MailTemplateHandler.DeleteMailTemplate] error when deleting mail template: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete mail template", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success delete mail template", nil)
}

// FindAllByDocumentTypeID find all mail templates by document type ID
//
//	@Summary		Find all mail templates by document type ID
//	@Description	Find all mail templates by document type ID
//	@Tags			Mail Templates
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200	{object} response.MailTemplateResponse
//
// @Security BearerAuth
//
//	@Router			/mail-templates/document-type/{id} [get]
func (h *MailTemplateHandler) FindAllByDocumentTypeID(ctx *gin.Context) {
	id := ctx.Param("id")
	documentTypeID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[MailTemplateHandler.FindAllByDocumentTypeID] error when parsing UUID: %v", err)
		utils.BadRequestResponse(ctx, "Bad request", err)
		return
	}

	res, err := h.UseCase.FindAllByDocumentTypeID(documentTypeID)
	if err != nil {
		h.Log.Errorf("[MailTemplateHandler.FindAllByDocumentTypeID] error when finding mail templates by document type ID: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find mail templates by document type ID", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get mail templates by document type ID", res)
}
