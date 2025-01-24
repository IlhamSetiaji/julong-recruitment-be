package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
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
	FindAllPaginatedWithoutUserID(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateJobPosting(ctx *gin.Context)
	DeleteJobPosting(ctx *gin.Context)
	GenerateDocumentNumber(ctx *gin.Context)
	FindAllAppliedJobPostingByUserID(ctx *gin.Context)
	InsertSavedJob(ctx *gin.Context)
	FindAllSavedJobsByUserID(ctx *gin.Context)
}

type JobPostingHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IJobPostingUseCase
	UserHelper helper.IUserHelper
}

func NewJobPostingHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IJobPostingUseCase,
	userHelper helper.IUserHelper,
) IJobPostingHandler {
	return &JobPostingHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func JobPostingHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IJobPostingHandler {
	useCase := usecase.JobPostingUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewJobPostingHandler(log, viper, validate, useCase, userHelper)
}

// CreateJobPosting create job posting
//
//		@Summary		Create job posting
//		@Description	Create job posting
//		@Tags			Job Postings
//		@Accept			multipart/form-data
//		@Produce		json
//		@Param			project_recruitment_header_id formData string true "Project Recruitment Header ID"
//		@Param			mp_request_id formData string true "MP Request ID"
//		@Param			job_id formData string true "Job ID"
//		@Param			for_organization_id formData string true "For Organization ID"
//		@Param			for_organization_location_id formData string true "For Organization Location ID"
//		@Param			document_number formData string true "Document Number"
//		@Param			document_date formData string true "Document Date"
//		@Param			recruitment_type formData string true "Recruitment Type"
//		@Param			start_date formData string true "Start Date"
//		@Param			end_date formData string true "End Date"
//		@Param			status formData string true "Status"
//		@Param			salary_min formData string true "Salary Min"
//		@Param			salary_max formData string true "Salary Max"
//		@Param			content_description formData string false "Content Description"
//		@Param			organization_logo formData file false "Organization Logo"
//		@Param			poster formData file false "Poster"
//		@Param			link formData string false "Link"
//		@Success		201	{object} response.JobPostingResponse
//	 @Security BearerAuth
//		@Router			/job-postings [post]
func (h *JobPostingHandler) CreateJobPosting(ctx *gin.Context) {
	var req request.CreateJobPostingRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("failed to bind request: ", err)
		utils.BadRequestResponse(ctx, "invalid requeest payload", err.Error())
		return
	}

	if err := h.Validate.Struct(&req); err != nil {
		h.Log.Error("validation error: ", err)
		utils.BadRequestResponse(ctx, "invalid request payload", err.Error())
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

// FindAllPaginated find all job postings
//
//		@Summary		Find all job postings
//		@Description	Find all job postings
//		@Tags			Job Postings
//		@Accept			json
//		@Produce		json
//		@Param			page	query	int	false	"Page"
//		@Param			page_size	query	int	false	"Page Size"
//		@Param			search	query	string	false	"Search"
//		@Param			created_at	query	string	false	"Created At"
//		@Param			status	query	string	false	"Status"
//		@Success		200	{object} response.JobPostingResponse
//	 @Security BearerAuth
//		@Router			/job-postings [get]
func (h *JobPostingHandler) FindAllPaginated(ctx *gin.Context) {
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
	status := ctx.Query("status")
	if status != "" {
		filter["status"] = status
	}

	res, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort, filter, userUUID)
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

// FindAllPaginatedWithoutUserID find all job postings without user id
//
//		@Summary		Find all job postings without user id
//		@Description	Find all job postings without user id
//		@Tags			Job Postings
//		@Accept			json
//		@Produce		json
//		@Param			page	query	int	false	"Page"
//		@Param			page_size	query	int	false	"Page Size"
//		@Param			search	query	string	false	"Search"
//		@Param			created_at	query	string	false	"Created At"
//		@Param			status	query	string	false	"Status"
//		@Success		200	{object} response.JobPostingResponse
//	 @Security BearerAuth
//		@Router			/api/no-auth/job-postings [get]
func (h *JobPostingHandler) FindAllPaginatedWithoutUserID(ctx *gin.Context) {
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
	status := ctx.Query("status")
	if status != "" {
		filter["status"] = status
	}

	res, total, err := h.UseCase.FindAllPaginatedWithoutUserID(page, pageSize, search, sort, filter)
	if err != nil {
		h.Log.Error("failed to find all paginated job postings without user id: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find all paginated job postings without user id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", gin.H{
		"job_postings": res,
		"total":        total,
	})
}

// FindByID find job posting by id
//
//		@Summary		Find job posting by id
//		@Description	Find job posting by id
//		@Tags			Job Postings
//		@Accept			json
//		@Produce		json
//		@Param			id	path	string	true	"ID"
//		@Success		200	{object} response.JobPostingResponse
//	 @Security BearerAuth
//		@Router			/job-postings/{id} [get]
func (h *JobPostingHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("failed to parse id: ", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "failed to parse id", err.Error())
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
	userUUID, err := h.UserHelper.GetUserId(user)
	if err != nil {
		h.Log.Errorf("Error when getting user id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	res, err := h.UseCase.FindByID(parsedUUID, userUUID)
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

// UpdateJobPosting update job posting
//
//		@Summary		Update job posting
//		@Description	Update job posting
//		@Tags			Job Postings
//		@Accept			multipart/form-data
//		@Produce		json
//		@Param			id	formData	string	true	"ID"
//		@Param			project_recruitment_header_id formData string true "Project Recruitment Header ID"
//		@Param			mp_request_id formData string true "MP Request ID"
//		@Param			job_id formData string true "Job ID"
//		@Param			for_organization_id formData string true "For Organization ID"
//		@Param			for_organization_location_id formData string true "For Organization Location ID"
//		@Param			document_number formData string true "Document Number"
//		@Param			document_date formData string true "Document Date"
//		@Param			recruitment_type formData string true "Recruitment Type"
//		@Param			start_date formData string true "Start Date"
//		@Param			end_date formData string true "End Date"
//		@Param			status formData string true "Status"
//		@Param			salary_min formData string true "Salary Min"
//		@Param			salary_max formData string true "Salary Max"
//		@Param			content_description formData string false "Content Description"
//		@Param			organization_logo formData file false "Organization Logo"
//		@Param			poster formData file false "Poster"
//		@Param			link formData string false "Link"
//		@Param			organization_logo_path formData string false "Organization Logo Path"
//		@Param			poster_path formData string false "Poster Path"
//		@Param			deleted_organization_logo formData string false "Deleted Organization Logo"
//		@Param			deleted_poster formData string false "Deleted Poster"
//		@Success		200	{object} response.JobPostingResponse
//	 @Security BearerAuth
//		@Router			/job-postings/update [put]
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

	parsedUUID, err := uuid.Parse(req.ID)
	if err != nil {
		h.Log.Error("failed to parse id: ", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "failed to parse id", err.Error())
		return
	}

	parsedBoolDeletedOrganizationLogo, err := strconv.ParseBool(req.DeletedOrganizationLogo)
	if err != nil {
		h.Log.Error("failed to parse deleted organization logo: ", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "failed to parse deleted organization logo", err.Error())
		return
	}

	if parsedBoolDeletedOrganizationLogo != false {
		err := h.UseCase.UpdateJobPostingOrganizationLogoToNull(parsedUUID)
		if err != nil {
			h.Log.Error("failed to update organization logo to null: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update organization logo to null", err.Error())
			return
		}
	}

	parsedBoolDeletedPoster, err := strconv.ParseBool(req.DeletedPoster)
	if parsedBoolDeletedPoster != false {
		err := h.UseCase.UpdateJobPostingPosterToNull(parsedUUID)
		if err != nil {
			h.Log.Error("failed to update poster to null: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update poster to null", err.Error())
			return
		}
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

// DeleteJobPosting delete job posting
//
//		@Summary		Delete job posting
//		@Description	Delete job posting
//		@Tags			Job Postings
//		@Accept			json
//		@Produce		json
//		@Param			id	path	string	true	"ID"
//		@Success		200	{string}	string
//	 @Security BearerAuth
//		@Router			/job-postings/{id} [delete]
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

// GenerateDocumentNumber generate document number
//
//		@Summary		Generate document number
//		@Description	Generate document number
//		@Tags			Job Postings
//		@Accept			json
//		@Produce		json
//		@Success		200	{string}	string
//	 @Security BearerAuth
//		@Router			/job-postings/generate-document-number [get]
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

// FindAllAppliedJobPostingByUserID find all applied job posting by user id
//
//		@Summary		Find all applied job posting by user id
//		@Description	Find all applied job posting by user id
//		@Tags			Job Postings
//		@Accept			json
//		@Produce		json
//		@Success		200	{object} response.JobPostingResponse
//	 @Security BearerAuth
//		@Router			/job-postings/applied [get]
func (h *JobPostingHandler) FindAllAppliedJobPostingByUserID(ctx *gin.Context) {
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

	res, err := h.UseCase.FindAllAppliedJobPostingByUserID(userUUID)
	if err != nil {
		h.Log.Error("failed to find all applied job posting by user id: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find all applied job posting by user id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

// InsertSavedJob insert saved job
//
//		@Summary		Insert saved job
//		@Description	Insert saved job
//		@Tags			Job Postings
//		@Accept			json
//		@Produce		json
//		@Param			job_posting_id	query	string	true	"Job Posting ID"
//		@Success		200	{string}	string
//	 @Security BearerAuth
//		@Router			/job-postings/save [get]
func (h *JobPostingHandler) InsertSavedJob(ctx *gin.Context) {
	jobPostingID := ctx.Query("job_posting_id")

	parsedJobPostingID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error("failed to parse job posting id: ", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "failed to parse job posting id", err.Error())
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
	userUUID, err := h.UserHelper.GetUserId(user)
	if err != nil {
		h.Log.Errorf("Error when getting user id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	err = h.UseCase.InsertSavedJob(userUUID, parsedJobPostingID)
	if err != nil {
		h.Log.Error("failed to insert saved job: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to insert saved job", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", nil)
}

// FindAllSavedJobsByUserID find all saved jobs by user id
//
//		@Summary		Find all saved jobs by user id
//		@Description	Find all saved jobs by user id
//		@Tags			Job Postings
//		@Accept			json
//		@Produce		json
//		@Success		200	{object} response.JobPostingResponse
//	 @Security BearerAuth
//		@Router			/job-postings/saved [get]
func (h *JobPostingHandler) FindAllSavedJobsByUserID(ctx *gin.Context) {
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

	res, err := h.UseCase.FindAllSavedJobsByUserID(userUUID)
	if err != nil {
		h.Log.Error("failed to find all saved jobs by user id: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find all saved jobs by user id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}
