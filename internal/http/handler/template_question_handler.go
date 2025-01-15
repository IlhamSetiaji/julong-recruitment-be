package handler

import (
	"net/http"

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

type ITemplateQuestionHandler interface {
	CreateTemplateQuestion(ctx *gin.Context)
	FindAllFormTypes(ctx *gin.Context)
	FindByID(ctx *gin.Context)
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

func (h *TemplateQuestionHandler) FindAllFormTypes(ctx *gin.Context) {
	formTypes, err := h.UseCase.FindAllFormTypes()
	if err != nil {
		h.Log.Error("[TemplateQuestionHandler.FindAllFormTypes] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", formTypes)
}

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
