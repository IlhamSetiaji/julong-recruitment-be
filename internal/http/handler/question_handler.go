package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
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

type IQuestionHandler interface {
	CreateOrUpdateQuestions(ctx *gin.Context)
	FindByIDAndUserID(ctx *gin.Context)
	FindAllByProjectRecruitmentLineIDAndJobPostingID(ctx *gin.Context)
}

type QuestionHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IQuestionUseCase
	UserHelper helper.IUserHelper
}

func NewQuestionHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IQuestionUseCase,
	userHelper helper.IUserHelper,
) IQuestionHandler {
	return &QuestionHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func QuestionHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IQuestionHandler {
	useCase := usecase.QuestionUseCaseFactory(log)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewQuestionHandler(log, viper, validate, useCase, userHelper)
}

// CreateOrUpdateQuestions create or update questions
//
//	@Summary		Create or update questions
//	@Description	Create or update questions
//	@Tags			Questions
//	@Accept			json
//	@Produce		json
//	@Param			employee	body		request.CreateOrUpdateQuestions	true	"Create employee"
//	@Success		201			{object}	response.TemplateQuestionResponse
//	@Security		BearerAuth
//	@Router			/api/questions [post]
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

// FindByIDAndUserID find by id and user id
//
//	@Summary		Find by id and user id
//	@Description	Find by id and user id
//	@Tags			Questions
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string	true	"ID"
//	@Success		200		{object}	response.QuestionResponse
//	@Security		BearerAuth
//	@Router			/api/questions/{id} [get]
func (h *QuestionHandler) FindByIDAndUserID(ctx *gin.Context) {
	id := ctx.Param("id")

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
	userUUID, err := h.UserHelper.GetUserId(user)
	if err != nil {
		h.Log.Errorf("Error when getting user id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	qr, err := h.UseCase.FindByIDAndUserID(id, userUUID.String())
	if err != nil {
		h.Log.Error("[QuestionHandler.FindByIDAndUserID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", qr)
}

// FindAllByProjectRecruitmentLineIDAndJobPostingID find all by project recruitment line id and job posting id
//
//	@Summary		Find all by project recruitment line id and job posting id
//	@Description	Find all by project recruitment line id and job posting id
//	@Tags			Questions
//	@Accept			json
//	@Produce		json
//	@Param			project_recruitment_line_id	query	string	true	"Project Recruitment Line ID"
//	@Param			job_posting_id				query	string	true	"Job Posting ID"
//	@Success		200							{object}	response.QuestionResponse
//	@Security		BearerAuth
//	@Router			/api/questions/result [get]
func (h *QuestionHandler) FindAllByProjectRecruitmentLineIDAndJobPostingID(ctx *gin.Context) {
	projectRecruitmentLineID := ctx.Query("project_recruitment_line_id")
	if projectRecruitmentLineID == "" {
		h.Log.Error("[QuestionHandler.FindAllByProjectRecruitmentLineIDAndJobPostingID] project_recruitment_line_id is required")
		utils.BadRequestResponse(ctx, "bad request", "project_recruitment_line_id is required")
		return
	}

	jobPostingID := ctx.Query("job_posting_id")
	if jobPostingID == "" {
		h.Log.Error("[QuestionHandler.FindAllByProjectRecruitmentLineIDAndJobPostingID] job_posting_id is required")
		utils.BadRequestResponse(ctx, "bad request", "job_posting_id is required")
		return
	}

	parsedPrhID, err := uuid.Parse(projectRecruitmentLineID)
	if err != nil {
		h.Log.Error("[QuestionHandler.FindAllByProjectRecruitmentLineIDAndJobPostingID] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	parsedJpID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error("[QuestionHandler.FindAllByProjectRecruitmentLineIDAndJobPostingID] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	qr, err := h.UseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID(parsedPrhID, parsedJpID)
	if err != nil {
		h.Log.Error("[QuestionHandler.FindAllByProjectRecruitmentLineIDAndJobPostingID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", qr)
}
