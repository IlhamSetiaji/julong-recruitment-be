package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDashboardHandler interface {
	GetDashboard(ctx *gin.Context)
}

type DashboardHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IDashboardUseCase
}

func NewDashboardHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDashboardUseCase,
) IDashboardHandler {
	return &DashboardHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func DashboardHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IDashboardHandler {
	useCase := usecase.DashboardUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewDashboardHandler(log, viper, validate, useCase)
}

// GetDashboard get dashboard data
//
// @Summary Get dashboard data
// @Description Get dashboard data
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.DashboardResponse
// @Router /dashboard [get]
func (h *DashboardHandler) GetDashboard(ctx *gin.Context) {
	res, err := h.UseCase.GetDashboard()
	if err != nil {
		h.Log.Errorf("[DashboardHandler.GetDashboard] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Successfully get dashboard data", res)
}
