package handler

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IAdministrativeResultHandler interface {
	CreateOrUpdateAdministrativeResults(ctx *gin.Context)
	FindAllByAdministrativeSelectionID(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateStatusAdministrativeResult(ctx *gin.Context)
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
//			@Summary		Find all administrative results by administrative selection id
//			@Description	Find all administrative results by administrative selection id
//		 @Tags Administrative Result
//			@Accept		json
//			@Produce		json
//		 @Param administrativeSelectionID path string true "Administrative Selection ID"
//	  @Param			page	query	int	false	"Page"
//			@Param			page_size	query	int	false	"Page Size"
//			@Param			search	query	string	false	"Search"
//			@Param			created_at	query	string	false	"Created At"
//			@Param			status	query	string	false	"Status"
//		 @Success 200 {object} response.AdministrativeResultResponse
//		 @Security BearerAuth
//		 @Router /api/administrative-results/administrative-selection/{administrativeSelectionID} [get]
func (h *AdministrativeResultHandler) FindAllByAdministrativeSelectionID(ctx *gin.Context) {
	administrativeSelectionID := ctx.Param("administrative_selection_id")

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

	filter := make(map[string]interface{})
	status := ctx.Query("status")
	if status != "" {
		filter["status"] = status
	}

	res, total, err := h.UseCase.FindAllByAdministrativeSelectionID(administrativeSelectionID, page, pageSize, search, sort, filter)
	if err != nil {
		h.Log.Error("[AdministrativeResultHandler.FindAllByAdministrativeSelectionID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when find all administrative results by administrative selection id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", gin.H{
		"administrative_results": res,
		"total":                  total,
	})
}

// FindByID find administrative result by id
//
//		@Summary		Find administrative result by id
//		@Description	Find administrative result by id
//	 @Tags Administrative Result
//		@Accept		json
//		@Produce		json
//	 @Param id path string true "ID"
//	 @Success 200 {object} response.AdministrativeResultResponse
//	 @Security BearerAuth
//	 @Router /api/administrative-results/{id} [get]
func (h *AdministrativeResultHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[AdministrativeResultHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.FindByID(parsedID)
	if err != nil {
		h.Log.Error("[AdministrativeResultHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when find administrative result by id", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "administrative result not found", "administrative result not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

// UpdateStatusAdministrativeResult update status administrative result
//
//		@Summary		Update status administrative result
//		@Description	Update status administrative result
//	 @Tags Administrative Result
//		@Accept		json
//		@Produce		json
//	 @Param id path string true "ID"
//	 @Param status query string true "Status"
//	 @Success 200 {object} response.AdministrativeResultResponse
//	 @Security BearerAuth
//	 @Router /api/administrative-results/{id}/update-status [get]
func (h *AdministrativeResultHandler) UpdateStatusAdministrativeResult(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Query("status")

	if status == "" {
		utils.BadRequestResponse(ctx, "bad request", "status is required")
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[AdministrativeResultHandler.UpdateStatusAdministrativeResult] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.UpdateStatusAdministrativeResult(parsedID, entity.AdministrativeResultStatus(status))
	if err != nil {
		h.Log.Error("[AdministrativeResultHandler.UpdateStatusAdministrativeResult] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error when update status administrative result", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}
