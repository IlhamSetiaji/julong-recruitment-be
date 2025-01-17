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
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IJobPostingHandler interface {
	CreateJobPosting(ctx *gin.Context)
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
	useCase := usecase.JobPostingUseCaseFactory(log)
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
