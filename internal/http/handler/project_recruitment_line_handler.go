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

type IProjectRecruitmentLineHandler interface {
	CreateOrUpdateProjectRecruitmentLines(ctx *gin.Context)
	FindAllByProjectRecruitmentHeaderID(ctx *gin.Context)
}

type ProjectRecruitmentLineHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IProjectRecruitmentLineUseCase
}

func NewProjectRecruitmentLineHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IProjectRecruitmentLineUseCase,
) IProjectRecruitmentLineHandler {
	return &ProjectRecruitmentLineHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func ProjectRecruitmentLineHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IProjectRecruitmentLineHandler {
	useCase := usecase.ProjectRecruitmentLineUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewProjectRecruitmentLineHandler(log, viper, validate, useCase)
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
//	@Summary		Find all project recruitment lines by project recruitment header id
//	@Description	Find all project recruitment lines by project recruitment header id
//	@Tags			Project Recruitment Lines
//	@Accept			json
//	@Produce		json
//	@Param			project_recruitment_header_id path string true "project recruitment header id"
//	@Success		200	{array} response.ProjectRecruitmentLineResponse
//	@Security BearerAuth
//	@Router			/project-recruitment-lines/header/{project_recruitment_header_id} [get]
func (h *ProjectRecruitmentLineHandler) FindAllByProjectRecruitmentHeaderID(ctx *gin.Context) {
	projectRecruitmentHeaderID := ctx.Param("project_recruitment_header_id")
	parsedID, err := uuid.Parse(projectRecruitmentHeaderID)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByProjectRecruitmentHeaderID] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.GetAllByKeys(map[string]interface{}{
		"project_recruitment_header_id": parsedID,
	})
	if err != nil {
		h.Log.Error("[ProjectRecruitmentLineHandler.FindAllByProjectRecruitmentHeaderID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}
