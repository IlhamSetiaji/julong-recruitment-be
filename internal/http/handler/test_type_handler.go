package handler

import (
	"net/http"
	"strconv"

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

type ITestTypeHandler interface {
	CreateTestType(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateTestType(ctx *gin.Context)
	DeleteTestType(ctx *gin.Context)
}

type TestTypeHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.ITestTypeUseCase
}

func NewTestTypeHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITestTypeUseCase,
) ITestTypeHandler {
	return &TestTypeHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func TestTypeHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITestTypeHandler {
	useCase := usecase.TestTypeUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewTestTypeHandler(log, viper, validate, useCase)
}

func (h *TestTypeHandler) CreateTestType(ctx *gin.Context) {
	var payload request.CreateTestTypeRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[TestTypeHandler.CreateTestType] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[TestTypeHandler.CreateTestType] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.CreateTestType(&payload)
	if err != nil {
		h.Log.Errorf("[TestTypeHandler.CreateTestType] error when creating test type: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to create test type", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "test type created", res)
}

// FindAllPaginated find all test types paginated
func (h *TestTypeHandler) FindAllPaginated(ctx *gin.Context) {
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

	filter := make(map[string]interface{})
	name := ctx.Query("name")
	if name != "" {
		filter["name"] = name
	}
	testTypes, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort, filter)
	if err != nil {
		h.Log.Errorf("[TestTypeHandler.FindAllPaginated] error when finding all test types paginated: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find all test types paginated", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "test types found", gin.H{
		"test_types": testTypes,
		"total":      total,
	})
}

func (h *TestTypeHandler) FindAll(ctx *gin.Context) {
	res, err := h.UseCase.FindAll()
	if err != nil {
		h.Log.Errorf("[TestTypeHandler.FindAll] error when finding all test type: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find all test type", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "test type found", res)
}

func (h *TestTypeHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[TestTypeHandler.FindByID] error when parsing id: %v", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.FindByID(parsedID)
	if err != nil {
		h.Log.Errorf("[TestTypeHandler.FindByID] error when finding test type by id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find test type by id", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "test type not found", "test type not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "test type found", res)
}

func (h *TestTypeHandler) UpdateTestType(ctx *gin.Context) {
	var payload request.UpdateTestTypeRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[TestTypeHandler.UpdateTestType] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[TestTypeHandler.UpdateTestType] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	res, err := h.UseCase.UpdateTestType(&payload)
	if err != nil {
		h.Log.Errorf("[TestTypeHandler.UpdateTestType] error when updating test type: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update test type", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "test type not found", "test type not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "test type updated", res)
}

func (h *TestTypeHandler) DeleteTestType(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[TestTypeHandler.DeleteTestType] error when parsing id: %v", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	err = h.UseCase.DeleteTestType(parsedID)
	if err != nil {
		h.Log.Errorf("[TestTypeHandler.DeleteTestType] error when deleting test type: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to delete test type", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "test type deleted", nil)
}
