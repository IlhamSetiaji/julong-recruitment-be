package handler

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase.go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITestScheduleHeaderHandler interface {
	CreateTestScheduleHeader(ctx *gin.Context)
}

type TestScheduleHeaderHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITestScheduleHeaderUsecase
}

func NewTestScheduleHeaderHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITestScheduleHeaderUsecase,
) ITestScheduleHeaderHandler {
	return &TestScheduleHeaderHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TestScheduleHeaderHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITestScheduleHeaderHandler {
	useCase := usecase.TestScheduleHeaderUsecaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewTestScheduleHeaderHandler(log, viper, validate, useCase)
}

func (h *TestScheduleHeaderHandler) CreateTestScheduleHeader(ctx *gin.Context) {
	
}
