package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IRecruitmentTypeHandler interface {
	FindAll(ctx *gin.Context)
}

type RecruitmentTypeHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IRecruitmentTypeUseCase
}

func NewRecruitmentTypeHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IRecruitmentTypeUseCase,
) IRecruitmentTypeHandler {
	return &RecruitmentTypeHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func RecruitmentTypeHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IRecruitmentTypeHandler {
	useCase := usecase.RecruitmentTypeUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewRecruitmentTypeHandler(log, viper, validate, useCase)
}

// FindAll find all recruitment types
//
//		@Summary		Find all recruitment types
//		@Description	Find all recruitment types
//		@Tags			Recruitment Types
//		@Accept			json
//		@Produce		json
//		@Success		200	{object} response.RecruitmentTypeResponse
//	 @Security BearerAuth
//		@Router			/recruitment-types [get]
func (h *RecruitmentTypeHandler) FindAll(ctx *gin.Context) {
	recruitmentTypes, err := h.UseCase.FindAll()
	if err != nil {
		h.Log.Errorf("[RecruitmentTypeHandler.FindAll] error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all recruitment types", err.Error())
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success find all recruitment types", recruitmentTypes)
}
