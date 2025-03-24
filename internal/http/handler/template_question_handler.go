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

type ITemplateQuestionHandler interface {
	FindAllPaginated(ctx *gin.Context)
	CreateTemplateQuestion(ctx *gin.Context)
	FindAllFormTypes(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateTemplateQuestion(ctx *gin.Context)
	DeleteTemplateQuestion(ctx *gin.Context)
}

type TemplateQuestionHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITemplateQuestionUseCase
}

func NewTemplateQuestionHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITemplateQuestionUseCase,
) ITemplateQuestionHandler {
	return &TemplateQuestionHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TemplateQuestionHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITemplateQuestionHandler {
	useCase := usecase.TemplateQuestionUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewTemplateQuestionHandler(log, viper, validate, useCase)
}

// FindAllPaginated find all template questions paginated
//
//	@Summary		Find all template questions paginated
//	@Description	Find all template questions paginated
//	@Tags			Template Questions
//	@Accept			json
//	@Produce		json
//	@Param			page	query	int	false	"Page"
//	@Param			page_size	query	int	false	"Page Size"
//	@Param			search	query	string	false	"Search"
//	@Param			created_at	query	string	false	"Created At"
//	@Success		200	{object}	response.TemplateQuestionResponse
//	@Security BearerAuth
//	@Router			/api/template-questions	[get]
func (h *TemplateQuestionHandler) FindAllPaginated(ctx *gin.Context) {
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
	// Filter DocumentSetup name

	sort := map[string]interface{}{
		"created_at": createdAt,
	}
	// Filter DocumentSetup name
	filter := map[string]interface{}{}
	// Filter DocumentSetup name
	if ctx.Query("document_setup_title") != "" {
		filter["document_setup_title"] = ctx.Query("document_setup_title")
	}

	templateQuestions, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort, filter)
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", gin.H{
		"template_questions": templateQuestions,
		"total":              total,
	})
}

// CreateTemplateQuestion create template question
//
//	@Summary		Create template question
//	@Description	Create template question
//	@Tags			Template Questions
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	request.CreateTemplateQuestion	true	"Payload"
//	@Success		201	{object}	response.TemplateQuestionResponse
//	@Security BearerAuth
//	@Router			/api/template-questions	[post]
func (h *TemplateQuestionHandler) CreateTemplateQuestion(ctx *gin.Context) {
	var payload request.CreateTemplateQuestion
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("[TemplateQuestionHandler.CreateTemplateQuestion] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[TemplateQuestionHandler.CreateTemplateQuestion] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	response, err := h.UseCase.CreateTemplateQuestion(&payload)
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.CreateTemplateQuestion] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success", response)
}

// FindAllFormTypes find all form types
//
//	@Summary		Find all form types
//	@Description	Find all form types
//	@Tags			Template Questions
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.FormTypeResponse
//	@Security BearerAuth
//	@Router			/api/template-questions/form-types	[get]
func (h *TemplateQuestionHandler) FindAllFormTypes(ctx *gin.Context) {
	formTypes, err := h.UseCase.FindAllFormTypes()
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.FindAllFormTypes] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", formTypes)
}

// FindByID find template question by id
//
//	@Summary		Find template question by id
//	@Description	Find template question by id
//	@Tags			Template Questions
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200	{object}	response.TemplateQuestionResponse
//	@Security BearerAuth
//	@Router			/api/template-questions/{id}	[get]
func (h *TemplateQuestionHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Error("[TemplateQuestionHandler.FindByID] id is required")
		utils.BadRequestResponse(ctx, "bad request", "id is required")
		return
	}

	templateQuestionID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	response, err := h.UseCase.FindByID(templateQuestionID)
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if response == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "error not found", "template question not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response)
}

// UpdateTemplateQuestion update template question
//
//	@Summary		Update template question
//	@Description	Update template question
//	@Tags			Template Questions
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Param			payload	body	request.UpdateTemplateQuestion	true	"Payload"
//	@Success		200	{object}	response.TemplateQuestionResponse
//	@Security BearerAuth
//	@Router			/api/template-questions/{id}	[put]
func (h *TemplateQuestionHandler) UpdateTemplateQuestion(ctx *gin.Context) {
	var payload request.UpdateTemplateQuestion
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("[TemplateQuestionHandler.UpdateTemplateQuestion] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[TemplateQuestionHandler.UpdateTemplateQuestion] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	exist, err := h.UseCase.FindByID(uuid.MustParse(payload.ID))
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.UpdateTemplateQuestion] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if exist == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "error not found", "template question not found")
		return
	}

	response, err := h.UseCase.UpdateTemplateQuestion(&payload)
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.UpdateTemplateQuestion] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response)
}

// DeleteTemplateQuestion delete template question
//
//	@Summary		Delete template question
//	@Description	Delete template question
//	@Tags			Template Questions
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200	{string}	string
//	@Security BearerAuth
//	@Router			/api/template-questions/{id}	[delete]
func (h *TemplateQuestionHandler) DeleteTemplateQuestion(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Error("[TemplateQuestionHandler.DeleteTemplateQuestion] id is required")
		utils.BadRequestResponse(ctx, "bad request", "id is required")
		return
	}

	templateQuestionID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.DeleteTemplateQuestion] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	exist, err := h.UseCase.FindByID(templateQuestionID)
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.DeleteTemplateQuestion] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	if exist == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "error not found", "template question not found")
		return
	}

	err = h.UseCase.DeleteTemplateQuestion(templateQuestionID)
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.DeleteTemplateQuestion] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", nil)
}
