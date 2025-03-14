package handler

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
)

type IDocumentVerificationHeaderHandler interface {
	CreateDocumentVerificationHeader(ctx *gin.Context)
	UpdateDocumentVerificationHeader(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	DeleteDocumentVerificationHeader(ctx *gin.Context)
	FindByJobPostingAndApplicant(ctx *gin.Context)
	ExportBPJSTenagaKerja(ctx *gin.Context)
	ImportBPJSTenagaKerja(ctx *gin.Context)
}

type DocumentVerificationHeaderHandler struct {
	Log                *logrus.Logger
	Viper              *viper.Viper
	Validate           *validator.Validate
	UseCase            usecase.IDocumentVerificationHeaderUseCase
	ApplicantUseCase   usecase.IApplicantUseCase
	UserHelper         helper.IUserHelper
	EmployeeMessage    messaging.IEmployeeMessage
	UserProfileUseCase usecase.IUserProfileUseCase
}

func NewDocumentVerificationHeaderHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDocumentVerificationHeaderUseCase,
	applicantUseCase usecase.IApplicantUseCase,
	userHelper helper.IUserHelper,
	employeeMessage messaging.IEmployeeMessage,
	userProfileUseCase usecase.IUserProfileUseCase,
) IDocumentVerificationHeaderHandler {
	return &DocumentVerificationHeaderHandler{
		Log:                log,
		Viper:              viper,
		Validate:           validate,
		UseCase:            useCase,
		ApplicantUseCase:   applicantUseCase,
		UserHelper:         userHelper,
		EmployeeMessage:    employeeMessage,
		UserProfileUseCase: userProfileUseCase,
	}
}

func DocumentVerificationHeaderHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IDocumentVerificationHeaderHandler {
	useCase := usecase.DocumentVerificationHeaderUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	applicantUseCase := usecase.ApplicantUseCaseFactory(log, viper)
	userHelper := helper.UserHelperFactory(log)
	employeeMessage := messaging.EmployeeMessageFactory(log)
	userProfileUseCase := usecase.UserProfileUseCaseFactory(log, viper)
	return NewDocumentVerificationHeaderHandler(
		log,
		viper,
		validate,
		useCase,
		applicantUseCase,
		userHelper,
		employeeMessage,
		userProfileUseCase,
	)
}

// CreateDocumentVerificationHeader create document verification header
//
//		@Summary		Create document verification header
//		@Description	Create document verification header
//		@Tags			Document Verification Header
//		@Accept			json
//	 @Produce		json
//	 @Param payload body request.CreateDocumentVerificationHeaderRequest true "Create Document Verification Header"
//		@Success		200	{object} response.DocumentVerificationHeaderResponse
//		@Router			/document-verification-headers [post]
func (h *DocumentVerificationHeaderHandler) CreateDocumentVerificationHeader(ctx *gin.Context) {
	req := new(request.CreateDocumentVerificationHeaderRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	res, err := h.UseCase.CreateDocumentVerificationHeader(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when creating document verification header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Document verification header created", res)
}

// UpdateDocumentVerificationHeader update document verification header
//
//	@Summary		Update document verification header
//	@Description	Update document verification header
//	@Tags			Document Verification Header
//	@Accept			json
//	@Produce		json
//	@Param			payload body request.UpdateDocumentVerificationHeaderRequest true "Update Document Verification Header"
//	@Success		200	{object} response.DocumentVerificationHeaderResponse
//	@Router			/document-verification-headers/update [put]
func (h *DocumentVerificationHeaderHandler) UpdateDocumentVerificationHeader(ctx *gin.Context) {
	req := new(request.UpdateDocumentVerificationHeaderRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	res, err := h.UseCase.UpdateDocumentVerificationHeader(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when updating document verification header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document verification header updated", res)
}

// FindByID find document verification header by id
//
//	@Success		200	{object} response.DocumentVerificationHeaderResponse
//	@Security BearerAuth
//	@Router			/document-verification-headers/{id} [get]
func (h *DocumentVerificationHeaderHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.FindByID] error when parsing id: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Error when parsing id", err.Error())
		return
	}

	res, err := h.UseCase.FindByID(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.FindByID] error when finding document verification header by id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding document verification header by id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get document verification header by id", res)
}

// FindAllPaginated find all document verification header with pagination
//
//	@Summary		Find all document verification header with pagination
//	@Description	Find all document verification header with pagination
//	@Tags			Document Verification Header
//	@Accept			json
//	@Produce		json
//	@Param			page query int false "Page"
//	@Param			page_size query int false "Page Size"
//	@Param			search query string false "Search"
//	@Param			sort query string false "Sort"
//	@Success		200	{object} response.DocumentVerificationHeaderResponse
//	@Security BearerAuth
//	@Router			/document-verification-headers [get]
func (h *DocumentVerificationHeaderHandler) FindAllPaginated(ctx *gin.Context) {
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
		h.Log.Errorf("[DocumentVerificationHeaderHandler.FindAllPaginated] error when finding all document verification header: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding all document verification header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get all document verification header", gin.H{
		"document_verification_headers": res,
		"total":                         total,
	})
}

// DeleteDocumentVerificationHeader delete document verification header by id
//
//	@Summary		Delete document verification header by id
//	@Description	Delete document verification header by id
//	@Tags			Document Verification Header
//	@Accept			json
//	@Produce		json
//	@Param			id path string true "ID"
//	@Success		200	{string} string
//	@Security BearerAuth
//	@Router			/document-verification-headers/{id} [delete]
func (h *DocumentVerificationHeaderHandler) DeleteDocumentVerificationHeader(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.DeleteDocumentVerificationHeader] error when parsing id: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Error when parsing id", err.Error())
		return
	}

	err = h.UseCase.DeleteDocumentVerificationHeader(id)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.DeleteDocumentVerificationHeader] error when deleting document verification header by id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when deleting document verification header by id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document verification header deleted", nil)
}

// FindByJobPostingAndApplicant find document verification header by job posting and applicant
//
//		@Summary		Find document verification header by job posting and applicant
//		@Description	Find document verification header by job posting and applicant
//		@Tags			Document Verification Header
//		@Accept			json
//		@Produce		json
//		@Param			job_posting_id query string true "Job Posting ID"
//		@Param			applicant_id query string true "Applicant ID"
//		@Success		200	{object} response.DocumentVerificationHeaderResponse
//		@Security BearerAuth
//	 @Router			/document-verification-headers/find [get]
func (h *DocumentVerificationHeaderHandler) FindByJobPostingAndApplicant(ctx *gin.Context) {
	jobPostingID := ctx.Query("job_posting_id")
	applicantID := ctx.Query("applicant_id")

	parsedJpID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.FindByJobPostingAndApplicant] error when parsing job posting id: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Error when parsing job posting id", err.Error())
		return
	}

	parsedAppId, err := uuid.Parse(applicantID)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.FindByJobPostingAndApplicant] error when parsing applicant id: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Error when parsing applicant id", err.Error())
		return
	}

	res, err := h.UseCase.FindByJobPostingAndApplicant(parsedJpID, parsedAppId)
	if err != nil {
		h.Log.Errorf("[DocumentVerificationHeaderHandler.FindByJobPostingAndApplicant] error when finding document verification header by job posting and applicant: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding document verification header by job posting and applicant", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success get document verification header by job posting and applicant", res)
}

// ExportBPJSTenagaKerja exports Excel file for BPJS Tenaga Kerja
//
// @Summary Export Excel file for BPJS Tenaga Kerja
// @Description Export Excel file for BPJS Tenaga Kerja
// @Tags Document Verification Header
// @Accept json
//
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Success		200					{file}		file
//	@Security		BearerAuth
//	@Router			/document-verification-headers/bpjs-tk [get]
func (h *DocumentVerificationHeaderHandler) ExportBPJSTenagaKerja(ctx *gin.Context) {
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
	userName, err := h.UserHelper.GetUserName(user)
	if err != nil {
		h.Log.Errorf("Error when getting user name: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	userId, err := h.UserHelper.GetUserId(user)
	if err != nil {
		h.Log.Errorf("Error when getting user id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	userEmail, err := h.UserHelper.GetUserEmail(user)
	if err != nil {
		h.Log.Errorf("Error when getting user email: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	userProfile, err := h.UserProfileUseCase.FindByUserID(userId)
	if err != nil {
		h.Log.Errorf("Error when finding user profile by user id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	employeeId, err := h.UserHelper.GetEmployeeId(user)
	if err != nil {
		h.Log.Errorf("Error when getting employee id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	_, err = h.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
		ID: employeeId.String(),
	})
	if err != nil {
		h.Log.Errorf("Error when sending find employee by id message: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	f, err := excelize.OpenFile("./storage/template_tk_14005857.xlsx")
	if err != nil {
		h.Log.Errorf("Error when opening file: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			h.Log.Errorf("Error when closing file: %v", err)
			utils.ErrorResponse(ctx, 500, "error", err.Error())
			return
		}
	}()

	// Set value to the cell.
	sheetName := "data_tk_baru"
	// Nama Lengkap
	f.SetCellValue(sheetName, "B2", userName)
	// HP
	f.SetCellValue(sheetName, "I2", userProfile.PhoneNumber)
	// Email
	f.SetCellValue(sheetName, "J2", userEmail)
	// Birth Place
	f.SetCellValue(sheetName, "K2", userProfile.BirthPlace)
	// Birth Date
	f.SetCellValue(sheetName, "L2", userProfile.BirthDate.Format("02-01-2006"))
	// Marital Status
	f.SetCellValue(sheetName, "T2", userProfile.MaritalStatus)

	// Save the spreadsheet by the given path.
	fileName := "bpjs_tk_" + userId.String() + ".xlsx"
	if err := f.SaveAs("./storage/bpjs" + fileName); err != nil {
		h.Log.Errorf("Error when saving file: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=template_tk_14005857.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(ctx.Writer); err != nil {
		h.Log.Errorf("Error when writing file: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
		return
	}
}

// ImportBPJSTenagaKerja import Excel file for BPJS Tenaga Kerja
//
// @Summary Import Excel file for BPJS Tenaga Kerja
// @Description Import Excel file for BPJS Tenaga Kerja
// @Tags Document Verification Header
// @Accept json
//
//	@Produce		json
//	@Success		200	{string} string
//	@Security		BearerAuth
//	@Router			/document-verification-headers/bpjs-tk [post]
func (h *DocumentVerificationHeaderHandler) ImportBPJSTenagaKerja(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		h.Log.Errorf("Error when getting file: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	err = ctx.SaveUploadedFile(file, "./storage/"+file.Filename)
	if err != nil {
		h.Log.Errorf("Error when saving file: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, "Success import file", nil)
}
