package handler

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IAdministrativeSelectionHandler interface {
	CreateAdministrativeSelection(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateAdministrativeSelection(ctx *gin.Context)
	DeleteAdministrativeSelection(ctx *gin.Context)
	VerifyAdministrativeSelection(ctx *gin.Context)
}

type AdministrativeSelectionHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IAdministrativeSelectionUsecase
	UserHelper helper.IUserHelper
}

func NewAdministrativeSelectionHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IAdministrativeSelectionUsecase,
	userHelper helper.IUserHelper,
) IAdministrativeSelectionHandler {
	return &AdministrativeSelectionHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func AdministrativeSelectionHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IAdministrativeSelectionHandler {
	useCase := usecase.AdministrativeSelectionUsecaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewAdministrativeSelectionHandler(log, viper, validate, useCase, userHelper)
}

// CreateAdministrativeSelection create administrative selection
//
//		@Summary		Create administrative selection
//		@Description	Create administrative selection
//	 @Tags			Administrative Selection
//		@Accept		json
//		@Produce		json
//		@Param			payload body request.CreateAdministrativeSelectionRequest true "Create Document Verification"
//		@Success		201	{object} response.AdministrativeSelectionResponse
//		@Security BearerAuth
//		@Router			/administrative-selections [post]
func (h *AdministrativeSelectionHandler) CreateAdministrativeSelection(ctx *gin.Context) {
	var payload request.CreateAdministrativeSelectionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.CreateAdministrativeSelection] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid request payload", err)
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.CreateAdministrativeSelection] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid request payload", err)
		return
	}

	res, err := h.UseCase.CreateAdministrativeSelection(&payload)
	if err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.CreateAdministrativeSelection] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when creating administrative selection", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Administrative selection created successfully", res)
}

// FindAllPaginated find all administrative selection paginated
//
//	@Summary		Find all administrative selection paginated
//	@Description	Find all administrative selection paginated
//	@Tags			Administrative Selection
//	@Accept			json
//	@Produce		json
//	@Param			page query int false "Page"
//	@Param			pageSize query int false "Page Size"
//	@Param			search query string false "Search"
//	@Param			sort query string false "Sort"
//	@Success		200 {object} response.AdministrativeSelectionResponse
//	@Security BearerAuth
//	@Router			/administrative-selections [get]
func (h *AdministrativeSelectionHandler) FindAllPaginated(ctx *gin.Context) {
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
		h.Log.Error("[AdministrativeSelectionHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding all administrative selection", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Administrative selection found successfully", gin.H{
		"administrative_selections": res,
		"total":                     total,
	})
}

// FindByID find administrative selection by id
//
//	@Summary		Find administrative selection by id
//	@Description	Find administrative selection by id
//	@Tags			Administrative Selection
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "ID"
//	@Success		200 {object} response.AdministrativeSelectionResponse
//	@Security BearerAuth
//	@Router			/administrative-selections/{id} [get]
func (h *AdministrativeSelectionHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := h.UseCase.FindByID(id)
	if err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding administrative selection", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Administrative selection not found", "Administrative selection not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Administrative selection found successfully", res)
}

// UpdateAdministrativeSelection update administrative selection
//
//	@Summary		Update administrative selection
//	@Description	Update administrative selection
//	@Tags			Administrative Selection
//	@Accept			json
//	@Produce		json
//	@Param			payload body request.UpdateAdministrativeSelectionRequest true "Update Administrative Selection"
//	@Success		200 {object} response.AdministrativeSelectionResponse
//	@Security BearerAuth
//	@Router			/administrative-selections/{id} [put]
func (h *AdministrativeSelectionHandler) UpdateAdministrativeSelection(ctx *gin.Context) {
	var payload request.UpdateAdministrativeSelectionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.UpdateAdministrativeSelection] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid request payload", err)
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.UpdateAdministrativeSelection] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid request payload", err)
		return
	}

	exist, err := h.UseCase.FindByID(payload.ID)
	if err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.UpdateAdministrativeSelection] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding administrative selection", err.Error())
		return
	}

	if exist == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Administrative selection not found", "Administrative selection not found")
		return
	}

	res, err := h.UseCase.UpdateAdministrativeSelection(&payload)
	if err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.UpdateAdministrativeSelection] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when updating administrative selection", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Administrative selection updated successfully", res)
}

// DeleteAdministrativeSelection delete administrative selection
//
//	@Summary		Delete administrative selection
//	@Description	Delete administrative selection
//	@Tags			Administrative Selection
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "ID"
//	@Success		204
//	@Security BearerAuth
//	@Router			/administrative-selections/{id} [delete]
func (h *AdministrativeSelectionHandler) DeleteAdministrativeSelection(ctx *gin.Context) {
	id := ctx.Param("id")

	exist, err := h.UseCase.FindByID(id)
	if err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.DeleteAdministrativeSelection] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding administrative selection", err.Error())
		return
	}

	if exist == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Administrative selection not found", "Administrative selection not found")
		return
	}

	err = h.UseCase.DeleteAdministrativeSelection(id)
	if err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.DeleteAdministrativeSelection] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when deleting administrative selection", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusNoContent, "Administrative selection deleted successfully", nil)
}

// VerifyAdministrativeSelection verify administrative selection
//
//	@Summary		Verify administrative selection
//	@Description	Verify administrative selection
//	@Tags			Administrative Selection
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "ID"
//	@Success		200
//	@Security BearerAuth
//	@Router			/administrative-selections/verify/{id} [get]
func (h *AdministrativeSelectionHandler) VerifyAdministrativeSelection(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "ID is required", "ID is required")
		return
	}

	user, err := middleware.GetUser(ctx, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		return
	}

	verifiedBy, err := h.UserHelper.GetEmployeeId(user)
	if err != nil {
		h.Log.Errorf("Error when getting employee id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	err = h.UseCase.VerifyAdministrativeSelection(id, verifiedBy.String())
	if err != nil {
		h.Log.Error("[AdministrativeSelectionHandler.VerifyAdministrativeSelection] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when verifying administrative selection", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Administrative selection verified successfully", nil)
}
