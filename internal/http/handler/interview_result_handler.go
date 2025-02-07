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
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type IInterviewResultHandler interface {
	FillInterviewResult(ctx *gin.Context)
	FindByInterviewApplicantAndAssessorID(ctx *gin.Context)
}

type InterviewResultHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IInterviewResultUseCase
	UserHelper helper.IUserHelper
	DB         *gorm.DB
}

func NewInterviewResultHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IInterviewResultUseCase,
	userHelper helper.IUserHelper,
	db *gorm.DB,
) IInterviewResultHandler {
	return &InterviewResultHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
		DB:         db,
	}
}

func InterviewResultHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IInterviewResultHandler {
	validate := config.NewValidator(viper)
	useCase := usecase.InterviewResultUseCaseFactory(log, viper)
	userHelper := helper.UserHelperFactory(log)
	db := config.NewDatabase()
	return NewInterviewResultHandler(log, viper, validate, useCase, userHelper, db)
}

// FillInterviewResult fill interview result
//
// @Summary Fill interview result
// @Description Fill interview result
// @Tags Interview Result
// @Accept json
// @Produce json
// @Param interview_result body request.FillInterviewResultRequest true "Interview Result"
// @Success 200 {object} response.InterviewResultResponse
// @Security BearerAuth
// @Router /interview-result [post]
func (h *InterviewResultHandler) FillInterviewResult(ctx *gin.Context) {
	var req request.FillInterviewResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.FillInterviewResult(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fill interview result", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Interview result filled successfully", res)
}

// FindByInterviewApplicantAndAssessorID find interview result by interview applicant and assessor id
//
// @Summary Find interview result by interview applicant and assessor id
// @Description Find interview result by interview applicant and assessor id
// @Tags Interview Result
// @Accept json
// @Produce json
// @Param interview_applicant_id query string true "Interview Applicant ID"
// @Param interview_assessor_id query string true "Interview Assessor ID"
// @Success 200 {object} response.InterviewResultResponse
// @Security BearerAuth
// @Router /interview-result [get]
func (h *InterviewResultHandler) FindByInterviewApplicantAndAssessorID(ctx *gin.Context) {
	interviewApplicantID := ctx.Query("interview_applicant_id")
	interviewAssessorID := ctx.Query("interview_assessor_id")

	if interviewApplicantID == "" {
		utils.BadRequestResponse(ctx, "Interview applicant ID is required", nil)
		return
	}

	if interviewAssessorID == "" {
		utils.BadRequestResponse(ctx, "Interview assessor ID is required", nil)
		return
	}

	parsedInterviewApplicantID, err := uuid.Parse(interviewApplicantID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	parsedInterviewAssessorID, err := uuid.Parse(interviewAssessorID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.FindByInterviewApplicantAndAssessorID(parsedInterviewApplicantID, parsedInterviewAssessorID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find interview result", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Interview result not found", "Interview result not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Interview result found successfully", res)
}
