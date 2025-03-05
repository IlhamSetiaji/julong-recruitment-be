package handler

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IProjectRecruitmentLineHandler interface {
	CreateOrUpdateProjectRecruitmentLines(ctx *gin.Context)
	FindAllByProjectRecruitmentHeaderID(ctx *gin.Context)
	FindAllByFormType(ctx *gin.Context)
	FindAllByProjectRecruitmentHeaderIDAndEmployeeID(ctx *gin.Context)
	FindAllByMonthAndYear(ctx *gin.Context)
}

type ProjectRecruitmentLineHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IProjectRecruitmentLineUseCase
	UserHelper helper.IUserHelper
}

func NewProjectRecruitmentLineHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IProjectRecruitmentLineUseCase,
	userHelper helper.IUserHelper,
) IProjectRecruitmentLineHandler {
	return &ProjectRecruitmentLineHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func ProjectRecruitmentLineHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IProjectRecruitmentLineHandler {
	useCase := usecase.ProjectRecruitmentLineUseCaseFactory(log)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewProjectRecruitmentLineHandler(log, viper, validate, useCase, userHelper)
}

// CreateOrUpdateProjectRecruitmentLines create or update project recruitment lines
//
//	@Summary		Create or update project recruitment lines
//	@Description	Create or update project recruitment lines
//	@Tags			Project Recruitment Lines
//	@Accept			json
//	@Produce		json
//	@Param			payload body request.CreateOrUpdateProjectRecruitmentLinesRequest true "payload"
//	@Success		201	{object} response.ProjectRecruitmentHeaderResponse
//	@Security BearerAuth
//	@Router			/project-recruitment-lines [post]
func (h *ProjectRecruitmentLineHandler) CreateOrUpdateProjectRecruitmentLines(ctx *gin.Context) {
	var payload request.CreateOrUpdateProjectRecruitmentLinesRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.CreateOrUpdateProjectRecruitmentLines] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.CreateOrUpdateProjectRecruitmentLines] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.CreateOrUpdateProjectRecruitmentLines(&payload)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.CreateOrUpdateProjectRecruitmentLines] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success", res)
}

// FindAllByProjectRecruitmentHeaderID find all project recruitment lines by project recruitment header id
//
//		@Summary		Find all project recruitment lines by project recruitment header id
//		@Description	Find all project recruitment lines by project recruitment header id
//		@Tags			Project Recruitment Lines
//		@Accept			json
//		@Produce		json
//		@Param			project_recruitment_header_id path string true "project recruitment header id"
//	 	@Param form_type query string false "form type"
//		@Success		200	{array} response.ProjectRecruitmentLineResponse
//		@Security BearerAuth
//		@Router			/project-recruitment-lines/header/{project_recruitment_header_id} [get]
func (h *ProjectRecruitmentLineHandler) FindAllByProjectRecruitmentHeaderID(ctx *gin.Context) {
	projectRecruitmentHeaderID := ctx.Param("project_recruitment_header_id")
	parsedID, err := uuid.Parse(projectRecruitmentHeaderID)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByProjectRecruitmentHeaderID] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	formType := ctx.Query("form_type")
	if formType != "" {
		res, err := h.UseCase.FindAllByHeaderIDAndFormType(parsedID, entity.TemplateQuestionFormType(formType))
		if err != nil {
			h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByProjectRecruitmentHeaderID] " + err.Error())
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}

		utils.SuccessResponse(ctx, http.StatusOK, "success", res)
		return
	}

	res, err := h.UseCase.GetAllByKeysWithoutPic(map[string]interface{}{
		"project_recruitment_header_id": parsedID,
	})
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByProjectRecruitmentHeaderID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

// FindAllByFormType find all project recruitment lines by form type
//
//	@Summary		Find all project recruitment lines by form type
//	@Description	Find all project recruitment lines by form type
//	@Tags			Project Recruitment Lines
//	@Accept			json
//	@Produce		json
//	@Param			form_type query string true "form type"
//	@Success		200	{array} response.ProjectRecruitmentLineResponse
//	@Security BearerAuth
//	@Router			/project-recruitment-lines/form-type/{form_type} [get]
func (h *ProjectRecruitmentLineHandler) FindAllByFormType(ctx *gin.Context) {
	formType := ctx.Query("form_type")
	if formType == "" {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByFormType] form type is required")
		utils.BadRequestResponse(ctx, "bad request", "form type is required")
		return
	}

	res, err := h.UseCase.FindAllByFormType(entity.TemplateQuestionFormType(formType))
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByFormType] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

// FindAllByProjectRecruitmentHeaderIDAndEmployeeID find all project recruitment lines by project recruitment header id and employee id
//
//	@Summary		Find all project recruitment lines by project recruitment header id and employee id
//	@Description	Find all project recruitment lines by project recruitment header id and employee id
//	@Tags			Project Recruitment Lines
//	@Accept			json
//	@Produce		json
//	@Param			project_recruitment_header_id path string true "project recruitment header id"
//	@Success		200	{array} response.ProjectRecruitmentLineResponse
//	@Security BearerAuth
//	@Router			/project-recruitment-lines/header-pic/{project_recruitment_header_id} [get]
func (h *ProjectRecruitmentLineHandler) FindAllByProjectRecruitmentHeaderIDAndEmployeeID(ctx *gin.Context) {
	projectRecruitmentHeaderID := ctx.Param("project_recruitment_header_id")
	parsedProjectRecruitmentHeaderID, err := uuid.Parse(projectRecruitmentHeaderID)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByProjectRecruitmentHeaderIDAndEmployeeID] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
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
	employeUUID, err := h.UserHelper.GetEmployeeId(user)
	if err != nil {
		h.Log.Errorf("Error when getting employee id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	res, err := h.UseCase.FindAllByProjectRecruitmentHeaderIDAndEmployeeID(parsedProjectRecruitmentHeaderID, employeUUID)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByProjectRecruitmentHeaderIDAndEmployeeID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

// FindAllByMonthAndYear find all project recruitment lines by month and year
//
//	@Summary		Find all project recruitment lines by month and year
//	@Description	Find all project recruitment lines by month and year
//	@Tags			Project Recruitment Lines
//	@Accept			json
//	@Produce		json
//	@Param			month query string true "month"
//	@Param			year query string true "year"
//	@Success		200	{array} response.ProjectRecruitmentLineResponse
//	@Security BearerAuth
//	@Router			/project-recruitment-lines/calendar [get]
func (h *ProjectRecruitmentLineHandler) FindAllByMonthAndYear(ctx *gin.Context) {
	month := ctx.Query("month")
	year := ctx.Query("year")
	if month == "" || year == "" {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByMonthAndYear] month and year is required")
		utils.BadRequestResponse(ctx, "bad request", "month and year is required")
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
	employeeID, err := h.UserHelper.GetEmployeeId(user)
	if err != nil {
		h.Log.Errorf("Error when getting user id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	intMonth, err := strconv.Atoi(month)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByMonthAndYear] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	intYear, err := strconv.Atoi(year)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByMonthAndYear] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	res, err := h.UseCase.FindAllByMonthAndYear(intMonth, intYear, employeeID)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByMonthAndYear] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}
