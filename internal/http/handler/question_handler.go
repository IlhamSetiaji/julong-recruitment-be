package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase.go"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IQuestionHandler interface {
	CreateOrUpdateQuestions(ctx *gin.Context)
}

type QuestionHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IQuestionUseCase
}

func NewQuestionHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IQuestionUseCase,
) IQuestionHandler {
	return &QuestionHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func QuestionHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IQuestionHandler {
	useCase := usecase.QuestionUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewQuestionHandler(log, viper, validate, useCase)
}

func (h *QuestionHandler) CreateOrUpdateQuestions(ctx *gin.Context) {
	var payload request.CreateOrUpdateQuestions
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("[QuestionHandler.CreateOrUpdateQuestions] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[QuestionHandler.CreateOrUpdateQuestions] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	tq, err := h.UseCase.CreateOrUpdateQuestions(&payload)
	if err != nil {
		h.Log.Error("[QuestionHandler.CreateOrUpdateQuestions] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success", tq)
}
