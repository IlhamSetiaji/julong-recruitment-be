package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IProjectRecruitmentHeaderHandler interface {
	CreateProjectRecruitmentHeader(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateProjectRecruitmentHeader(ctx *gin.Context)
	DeleteProjectRecruitmentHeader(ctx *gin.Context)
	GenerateDocumentNumber(ctx *gin.Context)
}

type ProjectRecruitmentHeaderHandler struct {
	Log                     *logrus.Logger
	Viper                   *viper.Viper
	Validate                *validator.Validate
	UseCase                 usecase.IProjectRecruitmentHeaderUseCase
	TemplateActivityUseCase usecase.ITemplateActivityUseCase
}

func NewProjectRecruitmentHeaderHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IProjectRecruitmentHeaderUseCase,
	taUseCase usecase.ITemplateActivityUseCase,
) IProjectRecruitmentHeaderHandler {
	return &ProjectRecruitmentHeaderHandler{
		Log:                     log,
		Viper:                   viper,
		Validate:                validate,
		UseCase:                 useCase,
		TemplateActivityUseCase: taUseCase,
	}
}

func ProjectRecruitmentHeaderHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IProjectRecruitmentHeaderHandler {
	useCase := usecase.ProjectRecruitmentHeaderUseCaseFactory(log)
	validate := config.NewValidator(viper)
	taUseCase := usecase.TemplateActivityUseCaseFactory(log)
	return NewProjectRecruitmentHeaderHandler(log, viper, validate, useCase, taUseCase)
}

// CreateProjectRecruitmentHeader create project recruitment header
//
//	@Summary		Create project recruitment header
//	@Description	Create project recruitment header
//	@Tags			Project Recruitment Headers
//	@Accept			json
//	@Produce		json
//	@Param			payload			body	request.CreateProjectRecruitmentHeader	true	"Create Project Recruitment Header"
//	@Success		201	{object} response.ProjectRecruitmentHeaderResponse
//	@Security BearerAuth
//	@Router /api/project-recruitment-headers [post]
func (h *ProjectRecruitmentHeaderHandler) CreateProjectRecruitmentHeader(ctx *gin.Context) {
	var req request.CreateProjectRecruitmentHeader
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	taID, err := uuid.Parse(req.TemplateActivityID)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	ta, err := h.TemplateActivityUseCase.FindByID(taID)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusNotFound, "error", err.Error())
		return
	}

	if ta == nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader] template activity not found")
		utils.ErrorResponse(ctx, http.StatusNotFound, "template activity not found", "template activity not found")
		return
	}

	response, err := h.UseCase.CreateProjectRecruitmentHeader(&req)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success", response)
}

// FindAllPaginated find all project recruitment headers
//
// @Summary		Find all project recruitment headers
// @Description	Find all project recruitment headers
// @Tags			Project Recruitment Headers
// @Accept			json
// @Produce		json
// @Param			page	query	int	false	"Page"
// @Param			page_size	query	int	false	"Page Size"
// @Param			search	query	string	false	"Search"
// @Param			created_at	query	string	false	"Created At"
// @Success		200	{object} response.ProjectRecruitmentHeaderResponse
// @Security BearerAuth
// @Router			/api/project-recruitment-headers [get]
func (h *ProjectRecruitmentHeaderHandler) FindAllPaginated(ctx *gin.Context) {
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

	responses, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", gin.H{
		"project_recruitment_headers": responses,
		"total":                       total,
	})
}

// FindByID find project recruitment header by id
//
// @Summary		Find project recruitment header by id
// @Description	Find project recruitment header by id
// @Tags			Project Recruitment Headers
// @Accept			json
// @Produce		json
// @Param			id	path	string	true	"ID"
// @Success		200	{object} response.ProjectRecruitmentHeaderResponse
// @Security BearerAuth
// @Router			/api/project-recruitment-headers/{id} [get]
func (h *ProjectRecruitmentHeaderHandler) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	response, err := h.UseCase.FindByID(id)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	if response == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "project recruitment header not found", "project recruitment header not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response)
}

// UpdateProjectRecruitmentHeader update project recruitment header
//
// @Summary		Update project recruitment header
// @Description	Update project recruitment header
// @Tags			Project Recruitment Headers
// @Accept			json
// @Produce		json
// @Param			payload	body	request.UpdateProjectRecruitmentHeader	true	"Update Project Recruitment Header"
// @Success		200	{object} response.ProjectRecruitmentHeaderResponse
// @Security BearerAuth
// @Router			/api/project-recruitment-headers/update [put]
func (h *ProjectRecruitmentHeaderHandler) UpdateProjectRecruitmentHeader(ctx *gin.Context) {
	var req request.UpdateProjectRecruitmentHeader
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.UpdateProjectRecruitmentHeader] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.UpdateProjectRecruitmentHeader] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	taID, err := uuid.Parse(req.TemplateActivityID)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	ta, err := h.TemplateActivityUseCase.FindByID(taID)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusNotFound, "error", err.Error())
		return
	}

	if ta == nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader] template activity not found")
		utils.ErrorResponse(ctx, http.StatusNotFound, "template activity not found", "template activity not found")
		return
	}

	response, err := h.UseCase.UpdateProjectRecruitmentHeader(&req)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.UpdateProjectRecruitmentHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response)
}

// DeleteProjectRecruitmentHeader delete project recruitment header
//
// @Summary		Delete project recruitment header
// @Description	Delete project recruitment header
// @Tags			Project Recruitment Headers
// @Accept			json
// @Produce		json
// @Param			id	path	string	true	"ID"
// @Success		200	{string}	string
// @Security BearerAuth
// @Router			/api/project-recruitment-headers/{id} [delete]
func (h *ProjectRecruitmentHeaderHandler) DeleteProjectRecruitmentHeader(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.DeleteProjectRecruitmentHeader] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	exist, err := h.UseCase.FindByID(id)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.DeleteProjectRecruitmentHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	if exist == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "error not found", "project recruitment header not found")
		return
	}

	err = h.UseCase.DeleteProjectRecruitmentHeader(id)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.DeleteProjectRecruitmentHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", nil)
}

// GenerateDocumentNumber generate document number
//
// @Summary		Generate document number
// @Description	Generate document number
// @Tags			Project Recruitment Headers
// @Accept			json
// @Produce		json
// @Success		200	{string}	string
// @Security BearerAuth
// @Router			/api/project-recruitment-headers/document-number [get]
func (h *ProjectRecruitmentHeaderHandler) GenerateDocumentNumber(ctx *gin.Context) {
	dateNow := time.Now()
	documentNumber, err := h.UseCase.GenerateDocumentNumber(dateNow)
	if err != nil {
		h.Log.Error("[ProjectRecruitmentHeaderHandler.GenerateDocumentNumber] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", documentNumber)
}
