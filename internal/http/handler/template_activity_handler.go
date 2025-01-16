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

type ITemplateActivityHandler interface {
	CreateTemplateActivity(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
}

type TemplateActivityHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITemplateActivityUseCase
}

func NewTemplateActivityHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITemplateActivityUseCase,
) ITemplateActivityHandler {
	return &TemplateActivityHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TemplateActivityHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITemplateActivityHandler {
	useCase := usecase.TemplateActivityUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewTemplateActivityHandler(log, viper, validate, useCase)
}

func (h *TemplateActivityHandler) CreateTemplateActivity(ctx *gin.Context) {
	var req request.CreateTemplateActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[TemplateActivityHandler.CreateTemplateActivity] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[TemplateActivityHandler.CreateTemplateActivity] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	resp, err := h.UseCase.CreateTemplateActivity(&req)
	if err != nil {
		h.Log.Error("[TemplateActivityHandler.CreateTemplateActivity] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success", resp)
}

func (h *TemplateActivityHandler) FindAllPaginated(ctx *gin.Context) {
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

	templateActivities, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Error("[TemplateActivityHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", gin.H{
		"template_activities": templateActivities,
		"total":               total,
	})
}

func (h *TemplateActivityHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[TemplateActivityHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "bad request", err.Error())
		return
	}

	templateActivity, err := h.UseCase.FindByID(parsedUUID)
	if err != nil {
		h.Log.Error("[TemplateActivityHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	if templateActivity == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "not found", "template activity not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", templateActivity)
}
