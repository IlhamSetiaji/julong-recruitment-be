package handler

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IMPRequestHandler interface {
	FindAllPaginated(ctx *gin.Context)
}

type MPRequestHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IMPRequestUseCase
}

func NewMPRequestHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IMPRequestUseCase,
) IMPRequestHandler {
	return &MPRequestHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func MPRequestHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IMPRequestHandler {
	useCase := usecase.MPRequestUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewMPRequestHandler(log, viper, validate, useCase)
}

// FindAllPaginated find all mp requests paginated
//
//	@Summary		Find all mp requests paginated
//	@Description	Find all mp requests paginated
//	@Tags			MP Requests
//	@Accept			json
//	@Produce		json
//	@Success		200	{object} response.MPRequestPaginatedResponse
//	@Security BearerAuth
//	@Router			/mp-requests [get]
func (h *MPRequestHandler) FindAllPaginated(ctx *gin.Context) {
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

	filter := make(map[string]interface{})

	status := ctx.Query("status")
	if status != "" {
		filter["status"] = status
	}

	res, err := h.UseCase.FindAllPaginated(page, pageSize, search, filter)
	if err != nil {
		h.Log.Errorf("[MPRequestHandler.FindAllPaginated] error when find all paginated: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all paginated", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success find all paginated", res)
}
