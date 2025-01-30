package handler

import (
	"net/http"

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

type ITemplateActivityLineHandler interface {
	CreateOrUpdateTemplateActivityLine(ctx *gin.Context)
	FindByTemplateActivityID(ctx *gin.Context)
	FindByID(ctx *gin.Context)
}

type TemplateActivityLineHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITemplateActivityLineUseCase
}

func NewTemplateActivityLineHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITemplateActivityLineUseCase,
) ITemplateActivityLineHandler {
	return &TemplateActivityLineHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TemplateActivityLineHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITemplateActivityLineHandler {
	useCase := usecase.TemplateActivityLineUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewTemplateActivityLineHandler(log, viper, validate, useCase)
}

// CreateOrUpdateTemplateActivityLine create or update template activity line
//
//	@Summary		Create or update template activity line
//	@Description	Create or update template activity line
//	@Tags			Template Activity Lines
//	@Accept		json
//	@Produce		json
//	@Param			payload body request.CreateOrUpdateTemplateActivityLineRequest true "payload"
//	@Success		201	{object} response.TemplateActivityResponse
//	@Security BearerAuth
//	@Router			/template-activity-lines [post]
func (h *TemplateActivityLineHandler) CreateOrUpdateTemplateActivityLine(ctx *gin.Context) {
	var payload request.CreateOrUpdateTemplateActivityLineRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("[TemplateActivityLineHandler.CreateOrUpdateTemplateActivityLine] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[TemplateActivityLineHandler.CreateOrUpdateTemplateActivityLine] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	ta, err := h.UseCase.CreateOrUpdateTemplateActivityLine(&payload)
	if err != nil {
		h.Log.Error("[TemplateActivityLineHandler.CreateOrUpdateTemplateActivityLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", ta)
}

// FindByTemplateActivityID find template activity line by template activity id
//
//	@Summary		Find template activity line by template activity id
//	@Description	Find template activity line by template activity id
//	@Tags			Template Activity Lines
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "Template Activity ID"
//	@Success		200	{object} response.TemplateActivityLineResponse
//	@Security BearerAuth
//	@Router			/template-activity-lines/template-activity/{id} [get]
func (h *TemplateActivityLineHandler) FindByTemplateActivityID(ctx *gin.Context) {
	templateActivityID := ctx.Param("id")

	tal, err := h.UseCase.FindByTemplateActivityID(templateActivityID)
	if err != nil {
		h.Log.Error("[TemplateActivityLineHandler.FindByTemplateActivityID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", tal)
}

// FindByID find template activity line by id
//
//	@Summary		Find template activity line by id
//	@Description	Find template activity line by id
//	@Tags			Template Activity Lines
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "Template Activity Line ID"
//	@Success		200	{object} response.TemplateActivityLineResponse
//	@Security BearerAuth
//	@Router			/template-activity-lines/{id} [get]
func (h *TemplateActivityLineHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[TemplateActivityLineHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}
	tal, err := h.UseCase.FindByID(parsedID)
	if err != nil {
		h.Log.Error("[TemplateActivityLineHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", tal)
}
