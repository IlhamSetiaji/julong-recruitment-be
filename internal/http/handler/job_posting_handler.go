package handler

import (
	"net/http"
	"strconv"
	"time"

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

type IJobPostingHandler interface {
	CreateJobPosting(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateJobPosting(ctx *gin.Context)
	DeleteJobPosting(ctx *gin.Context)
	GenerateDocumentNumber(ctx *gin.Context)
}

type JobPostingHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IJobPostingUseCase
}

func NewJobPostingHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IJobPostingUseCase,
) IJobPostingHandler {
	return &JobPostingHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func JobPostingHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IJobPostingHandler {
	useCase := usecase.JobPostingUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewJobPostingHandler(log, viper, validate, useCase)
}

func (h *JobPostingHandler) CreateJobPosting(ctx *gin.Context) {
	var req request.CreateJobPostingRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("failed to bind request: ", err)
		utils.BadRequestResponse(ctx, "invalid requeest payload", err)
		return
	}

	if err := h.Validate.Struct(&req); err != nil {
		h.Log.Error("validation error: ", err)
		utils.BadRequestResponse(ctx, "invalid request payload", err)
		return
	}

	// Handle file uploads
	if req.OrganizationLogo != nil {
		timestamp := time.Now().UnixNano()
		organizationLogoPath := "storage/job_posting/logos/" + strconv.FormatInt(timestamp, 10) + "_" + req.OrganizationLogo.Filename
		if err := ctx.SaveUploadedFile(req.OrganizationLogo, organizationLogoPath); err != nil {
			h.Log.Error("failed to save organization logo: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save organization logo", err.Error())
			return
		}
		req.OrganizationLogo = nil
		req.OrganizationLogoPath = organizationLogoPath
	}

	if req.Poster != nil {
		timestamp := time.Now().UnixNano()
		posterPath := "storage/job_posting/posters/" + strconv.FormatInt(timestamp, 10) + "_" + req.Poster.Filename
		if err := ctx.SaveUploadedFile(req.Poster, posterPath); err != nil {
			h.Log.Error("failed to save poster: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save poster", err.Error())
			return
		}
		req.Poster = nil
		req.PosterPath = posterPath
	}

	// Create job posting
	res, err := h.UseCase.CreateJobPosting(&req)
	if err != nil {
		h.Log.Error("failed to create job posting: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to create job posting", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "job posting created", res)
}

func (h *JobPostingHandler) FindAllPaginated(ctx *gin.Context) {
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

	res, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Error("failed to find all paginated job postings: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find all paginated job postings", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", gin.H{
		"job_postings": res,
		"total":        total,
	})
}

func (h *JobPostingHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("failed to parse id: ", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "failed to parse id", err.Error())
		return
	}

	res, err := h.UseCase.FindByID(parsedUUID)
	if err != nil {
		h.Log.Error("failed to find job posting by id: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find job posting by id", err.Error())
		return
	}

	if res == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "data not found", "job posting not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

func (h *JobPostingHandler) UpdateJobPosting(ctx *gin.Context) {
	var req request.UpdateJobPostingRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("failed to bind request: ", err)
		utils.BadRequestResponse(ctx, "invalid request payload", err)
		return
	}

	if err := h.Validate.Struct(&req); err != nil {
		h.Log.Error("validation error: ", err)
		utils.BadRequestResponse(ctx, "invalid request payload", err)
		return
	}

	// Handle file uploads
	if req.OrganizationLogo != nil {
		timestamp := time.Now().UnixNano()
		organizationLogoPath := "storage/job_posting/logos/" + strconv.FormatInt(timestamp, 10) + "_" + req.OrganizationLogo.Filename
		if err := ctx.SaveUploadedFile(req.OrganizationLogo, organizationLogoPath); err != nil {
			h.Log.Error("failed to save organization logo: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save organization logo", err.Error())
			return
		}
		req.OrganizationLogo = nil
		req.OrganizationLogoPath = organizationLogoPath
	}

	if req.Poster != nil {
		timestamp := time.Now().UnixNano()
		posterPath := "storage/job_posting/posters/" + strconv.FormatInt(timestamp, 10) + "_" + req.Poster.Filename
		if err := ctx.SaveUploadedFile(req.Poster, posterPath); err != nil {
			h.Log.Error("failed to save poster: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save poster", err.Error())
			return
		}
		req.Poster = nil
		req.PosterPath = posterPath
	}

	// Update job posting
	res, err := h.UseCase.UpdateJobPosting(&req)
	if err != nil {
		h.Log.Error("failed to update job posting: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update job posting", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "job posting updated", res)
}

func (h *JobPostingHandler) DeleteJobPosting(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("failed to parse id: ", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "failed to parse id", err.Error())
		return
	}

	err = h.UseCase.DeleteJobPosting(parsedUUID)
	if err != nil {
		h.Log.Error("failed to delete job posting: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to delete job posting", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "job posting deleted", nil)
}

func (h *JobPostingHandler) GenerateDocumentNumber(ctx *gin.Context) {
	dateNow := time.Now()
	documentNumber, err := h.UseCase.GenerateDocumentNumber(dateNow)
	if err != nil {
		h.Log.Error("failed to generate document number: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to generate document number", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", documentNumber)
}
