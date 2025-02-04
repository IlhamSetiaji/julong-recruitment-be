package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IProjectPicHandler interface {
	FindByProjectRecruitmentLineIDAndEmployeeID(ctx *gin.Context)
}

type ProjectPicHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IProjectPicUseCase
	UserHelper helper.IUserHelper
}

func NewProjectPicHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IProjectPicUseCase,
	userHelper helper.IUserHelper,
) IProjectPicHandler {
	return &ProjectPicHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func ProjectPicHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IProjectPicHandler {
	useCase := usecase.ProjectPicUseCaseFactory(log)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewProjectPicHandler(log, viper, validate, useCase, userHelper)
}

// FindByProjectRecruitmentLineIDAndEmployeeID find project pic by project recruitment line id and employee id
//
// @Summary Find project pic by project recruitment line id and employee id
// @Description Find project pic by project recruitment line id and employee id
// @Tags Project PIC
// @Accept json
// @Produce json
// @Param project_recruitment_line_id path string true "Project Recruitment Line ID"
// @Param employee_id path string true "Employee ID"
// @Success 200 {object} response.ProjectPicResponse
// @Security BearerAuth
// @Router /project-pic/project-recruitment-line/{project_recruitment_line_id} [get]
func (h *ProjectPicHandler) FindByProjectRecruitmentLineIDAndEmployeeID(ctx *gin.Context) {
	projectRecruitmentLineID := ctx.Param("project_recruitment_line_id")

	parsedPrlID, err := uuid.Parse(projectRecruitmentLineID)
	if err != nil {
		h.Log.Errorf("[ProjectPicHandler.FindByProjectRecruitmentLineIDAndEmployeeID] error when parsing project_recruitment_line_id: %v", err)
		utils.BadRequestResponse(ctx, "project_recruitment_line_id is not a valid UUID", err)
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

	res, err := h.UseCase.FindByProjectRecruitmentLineIDAndEmployeeID(parsedPrlID, employeUUID)
	if err != nil {
		h.Log.Errorf("[ProjectPicHandler.FindByProjectRecruitmentLineIDAndEmployeeID] error when getting project pic: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, 404, "error", "Project PIC not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Successfully find project pic", res)
}
