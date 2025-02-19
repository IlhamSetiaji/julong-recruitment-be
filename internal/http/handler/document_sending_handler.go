package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-pdf/fpdf"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentSendingHandler interface {
	CreateDocumentSending(ctx *gin.Context)
	UpdateDocumentSending(ctx *gin.Context)
	FindAllPaginatedByDocumentTypeID(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	DeleteDocumentSending(ctx *gin.Context)
	GenerateDocumentNumber(ctx *gin.Context)
	FindAllByDocumentSetupID(ctx *gin.Context)
	FindByDocumentTypeIDAndApplicantID(ctx *gin.Context)
	TestGeneratePDF(ctx *gin.Context)
	TestSendEmail(ctx *gin.Context)
}

type DocumentSendingHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IDocumentSendingUseCase
}

func NewDocumentSendingHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDocumentSendingUseCase,
) *DocumentSendingHandler {
	return &DocumentSendingHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func DocumentSendingHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IDocumentSendingHandler {
	validate := config.NewValidator(viper)
	useCase := usecase.DocumentSendingUseCaseFactory(log, viper)
	return NewDocumentSendingHandler(log, viper, validate, useCase)
}

// CreateDocumentSending  create document sending
//
// @Summary create document sending
// @Description create document sending
// @Tags Document Sendings
// @Accept json
// @Produce json
// @Param document_sending body request.CreateDocumentSendingRequest true "Document Sending"
// @Success 200 {object} response.DocumentSendingResponse "Success"
// @Security BearerAuth
// @Router /document-sending [post]
func (h *DocumentSendingHandler) CreateDocumentSending(ctx *gin.Context) {
	var payload request.CreateDocumentSendingRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[DocumentSendingHandler.CreateDocumentSending] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[DocumentSendingHandler.CreateDocumentSending] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	res, err := h.UseCase.CreateDocumentSending(&payload)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.CreateDocumentSending] error when creating document sending: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create document sending", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Document sending created", res)
}

// UpdateDocumentSending  update document sending
//
// @Summary update document sending
// @Description update document sending
// @Tags Document Sendings
// @Accept json
// @Produce json
// @Param document_sending body request.UpdateDocumentSendingRequest true "Document Sending"
// @Success 200 {object} response.DocumentSendingResponse "Success"
// @Security BearerAuth
// @Router /document-sending/update [put]
func (h *DocumentSendingHandler) UpdateDocumentSending(ctx *gin.Context) {
	var payload request.UpdateDocumentSendingRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[DocumentSendingHandler.UpdateDocumentSending] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[DocumentSendingHandler.UpdateDocumentSending] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	res, err := h.UseCase.UpdateDocumentSending(&payload)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.UpdateDocumentSending] error when updating document sending: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update document sending", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document sending updated", res)
}

// FindAllPaginatedByDocumentTypeID  find all paginated by document type id
//
// @Summary find all paginated by document type id
// @Description find all paginated by document type id
// @Tags Document Sendings
// @Accept json
// @Produce json
// @Param document_type_id query string true "Document Type ID"
// @Param page query int true "Page"
// @Param page_size query int true "Page Size"
// @Param search query string false "Search"
// @Param sort query string false "Sort"
// @Success 200 {object} response.DocumentSendingResponse "Success"
// @Security BearerAuth
// @Router /document-sending [get]
func (h *DocumentSendingHandler) FindAllPaginatedByDocumentTypeID(ctx *gin.Context) {
	documentTypeID := ctx.Query("document_type_id")
	if documentTypeID == "" {
		utils.BadRequestResponse(ctx, "Document type id is required", nil)
		return
	}
	parsedDocumentTypeID, err := uuid.Parse(documentTypeID)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid document type id", err)
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

	res, total, err := h.UseCase.FindAllPaginatedByDocumentTypeID(parsedDocumentTypeID, page, pageSize, search, sort)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.FindAllPaginatedByDocumentTypeID] error when finding all paginated by document type id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all paginated by document type id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document sending found", gin.H{
		"document_sendings": res,
		"total":             total,
	})
}

// FindByID  find by id
//
// @Summary find by id
// @Description find by id
// @Tags Document Sendings
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} response.DocumentSendingResponse "Success"
// @Security BearerAuth
// @Router /document-sending/{id} [get]
func (h *DocumentSendingHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.BadRequestResponse(ctx, "ID is required", nil)
		return
	}

	res, err := h.UseCase.FindByID(id)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.FindByID] error when finding by id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find by id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document sending found", res)
}

// DeleteDocumentSending  delete document sending
//
// @Summary delete document sending
// @Description delete document sending
// @Tags Document Sendings
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {string} string "Success"
// @Security BearerAuth
// @Router /document-sending/{id} [delete]
func (h *DocumentSendingHandler) DeleteDocumentSending(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.BadRequestResponse(ctx, "ID is required", nil)
		return
	}

	err := h.UseCase.DeleteDocumentSending(id)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.DeleteDocumentSending] error when deleting document sending: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete document sending", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document sending deleted", nil)
}

// GenerateDocumentNumber  generate document number
//
// @Summary generate document number
// @Description generate document number
// @Tags Document Sendings
// @Accept json
// @Produce json
// @Success 200 {string} string "Success"
// @Security BearerAuth
// @Router /document-sending/document-number [get]
func (h *DocumentSendingHandler) GenerateDocumentNumber(ctx *gin.Context) {
	dateNow := time.Now()
	documentNumber, err := h.UseCase.GenerateDocumentNumber(dateNow)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.GenerateDocumentNumber] error when generating document number: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate document number", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document number generated", documentNumber)
}

// FindAllByDocumentSetupID  find all by document setup id
//
// @Summary find all by document setup id
// @Description find all by document setup id
// @Tags Document Sendings
// @Accept json
// @Produce json
// @Param document_setup_id path string true "Document Setup ID"
// @Success 200 {object} response.DocumentSendingResponse "Success"
// @Security BearerAuth
// @Router /document-sending/document-setup/{document_setup_id} [get]
func (h *DocumentSendingHandler) FindAllByDocumentSetupID(ctx *gin.Context) {
	documentSetupID := ctx.Param("document_setup_id")
	if documentSetupID == "" {
		utils.BadRequestResponse(ctx, "Document setup id is required", nil)
		return
	}

	res, err := h.UseCase.FindAllByDocumentSetupID(documentSetupID)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.FindAllByDocumentSetupID] error when finding all by document setup id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all by document setup id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document sending found", res)
}

// FindByDocumentTypeIDAndApplicantID  find by document type id and applicant id
//
// @Summary find by document type id and applicant id
// @Description find by document type id and applicant id
// @Tags Document Sendings
// @Accept json
// @Produce json
// @Param document_type_id query string true "Document Type ID"
// @Param applicant_id query string true "Applicant ID"
// @Success 200 {object} response.DocumentSendingResponse "Success"
// @Security BearerAuth
// @Router /document-sending/applicant [get]
func (h *DocumentSendingHandler) FindByDocumentTypeIDAndApplicantID(ctx *gin.Context) {
	documentTypeID := ctx.Query("document_type_id")
	if documentTypeID == "" {
		utils.BadRequestResponse(ctx, "Document type id is required", nil)
		return
	}
	parsedDocumentTypeID, err := uuid.Parse(documentTypeID)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid document type id", err)
		return
	}

	applicantID := ctx.Query("applicant_id")
	if applicantID == "" {
		utils.BadRequestResponse(ctx, "Applicant id is required", nil)
		return
	}
	parsedApplicantID, err := uuid.Parse(applicantID)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid applicant id", err)
		return
	}

	res, err := h.UseCase.FindByDocumentTypeIDAndApplicantID(parsedDocumentTypeID, parsedApplicantID)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.FindByDocumentTypeIDAndApplicantID] error when finding by document type id and applicant id: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find by document type id and applicant id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document sending found", res)
}

// TestGeneratePDF  test generate pdf
//
// @Summary test generate pdf
// @Description test generate pdf
// @Tags Document Sendings
// @Accept json
// @Produce json
// @Success 200 {string} string "Success"
// @Security BearerAuth
// @Router /document-sending/test-generate-pdf [get]
func (h *DocumentSendingHandler) TestGeneratePDF(ctx *gin.Context) {
	// Create a new PDF document
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")

	// Define the file path
	filePath := "storage/hello.pdf"

	// Ensure the directory exists
	err := os.MkdirAll("storage", os.ModePerm)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.TestGeneratePDF] error when creating directory: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create directory", err.Error())
		return
	}

	// Save the PDF to the file
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.TestGeneratePDF] error when generating pdf: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate pdf", err.Error())
		return
	}

	// Return the generated PDF file as a response
	ctx.File(filePath)
}

// TestSendEmail  test send email
//
// @Summary test send email
// @Description test send email
// @Tags Document Sendings
// @Accept json
// @Produce json
// @Success 200 {string} string "Success"
// @Security BearerAuth
// @Router /document-sending/test-send-email [get]
func (h *DocumentSendingHandler) TestSendEmail(ctx *gin.Context) {
	// Send an email
	err := h.UseCase.TestSendEmail()
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.TestSendEmail] error when sending email: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to send email", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Email sent", nil)
}
