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

type IFgdResultHandler interface {
	FillFgdResult(ctx *gin.Context)
	FindByFgdApplicantAndAssessorID(ctx *gin.Context)
}

type FgdResultHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IFgdResultUseCase
	UserHelper helper.IUserHelper
	DB         *gorm.DB
}

func NewFgdResultHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IFgdResultUseCase,
	userHelper helper.IUserHelper,
	db *gorm.DB,
) IFgdResultHandler {
	return &FgdResultHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
		DB:         db,
	}
}

func FgdResultHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IFgdResultHandler {
	validate := config.NewValidator(viper)
	useCase := usecase.FgdResultUseCaseFactory(log, viper)
	userHelper := helper.UserHelperFactory(log)
	db := config.NewDatabase()
	return NewFgdResultHandler(log, viper, validate, useCase, userHelper, db)
}

// FillFgdResult fill Fgd result
//
// @Summary Fill Fgd result
// @Description Fill Fgd result
// @Tags Fgd Result
// @Accept json
// @Produce json
// @Param Fgd_result body request.FillFgdResultRequest true "Fgd Result"
// @Success 200 {object} response.FgdResultResponse
// @Security BearerAuth
// @Router /fgd-results [post]
func (h *FgdResultHandler) FillFgdResult(ctx *gin.Context) {
	var req request.FillFgdResultRequest
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

	res, err := h.UseCase.FillFgdResult(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fill Fgd result", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Fgd result filled successfully", res)
}

// FindByFgdApplicantAndAssessorID find Fgd result by Fgd applicant and assessor id
//
// @Summary Find Fgd result by Fgd applicant and assessor id
// @Description Find Fgd result by Fgd applicant and assessor id
// @Tags Fgd Result
// @Accept json
// @Produce json
// @Param fgd_applicant_id query string true "Fgd Applicant ID"
// @Param fgd_assessor_id query string true "Fgd Assessor ID"
// @Success 200 {object} response.FgdResultResponse
// @Security BearerAuth
// @Router /fgd-results/find [get]
func (h *FgdResultHandler) FindByFgdApplicantAndAssessorID(ctx *gin.Context) {
	FgdApplicantID := ctx.Query("fgd_applicant_id")
	FgdAssessorID := ctx.Query("fgd_assessor_id")

	if FgdApplicantID == "" {
		utils.BadRequestResponse(ctx, "Fgd applicant ID is required", nil)
		return
	}

	if FgdAssessorID == "" {
		utils.BadRequestResponse(ctx, "Fgd assessor ID is required", nil)
		return
	}

	parsedFgdApplicantID, err := uuid.Parse(FgdApplicantID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	parsedFgdAssessorID, err := uuid.Parse(FgdAssessorID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.FindByFgdApplicantAndAssessorID(parsedFgdApplicantID, parsedFgdAssessorID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find Fgd result", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Fgd result not found", "Fgd result not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Fgd result found successfully", res)
}
