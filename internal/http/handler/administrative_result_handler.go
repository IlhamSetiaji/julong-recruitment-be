package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IAdministrativeResultHandler interface {
	CreateOrUpdateAdministrativeResults(ctx *gin.Context)
	FindAllByAdministrativeSelectionID(ctx *gin.Context)
}

type AdministrativeResultHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IAdministrativeResultUseCase
	UserHelper helper.IUserHelper
}

func NewAdministrativeResultHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IAdministrativeResultUseCase,
	userHelper helper.IUserHelper,
) IAdministrativeResultHandler {
	return &AdministrativeResultHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func AdministrativeResultHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IAdministrativeResultHandler {
	useCase := usecase.AdministrativeResultUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewAdministrativeResultHandler(log, viper, validate, useCase, userHelper)
}

// CreateOrUpdateAdministrativeResults create or update administrative results
//
//		@Summary		Create or update administrative results
//		@Description	Create or update administrative results
//	 @Tags Administrative Result
//		@Accept		json
//		@Produce		json
//	 @Param data body request.CreateOrUpdateAdministrativeResults true "Create or update administrative results"
//	 @Success 201 {object} response.AdministrativeSelectionResponse
//	 @Security BearerAuth
//	 @Router /api/administrative-results [post]
func (h *AdministrativeResultHandler) CreateOrUpdateAdministrativeResults(ctx *gin.Context) {
	var payload request.CreateOrUpdateAdministrativeResults
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("[AdministrativeResultHandler.CreateOrUpdateAdministrativeResults] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[AdministrativeResultHandler.CreateOrUpdateAdministrativeResults] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.CreateOrUpdateAdministrativeResults(&payload)
	if err != nil {
		h.Log.Error("[AdministrativeResultHandler.CreateOrUpdateAdministrativeResults] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when create or update administrative results", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success", res)
}

// FindAllByAdministrativeSelectionID find all administrative results by administrative selection id
//
//		@Summary		Find all administrative results by administrative selection id
//		@Description	Find all administrative results by administrative selection id
//	 @Tags Administrative Result
//		@Accept		json
//		@Produce		json
//	 @Param administrativeSelectionID path string true "Administrative Selection ID"
//	 @Success 200 {object} response.AdministrativeResultResponse
//	 @Security BearerAuth
//	 @Router /api/administrative-results/administrative-selection/{administrativeSelectionID} [get]
func (h *AdministrativeResultHandler) FindAllByAdministrativeSelectionID(ctx *gin.Context) {
	administrativeSelectionID := ctx.Param("administrativeSelectionID")

	res, err := h.UseCase.FindAllByAdministrativeSelectionID(administrativeSelectionID)
	if err != nil {
		h.Log.Error("[AdministrativeResultHandler.FindAllByAdministrativeSelectionID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when find all administrative results by administrative selection id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}
