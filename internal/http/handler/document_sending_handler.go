package handler

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
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
	TestGenerateHTMLPDF(ctx *gin.Context)
	GeneratePdfBufferFromHTML(ctx *gin.Context)
	GeneratePdfBufferForDocumentSending(ctx *gin.Context)
}

type DocumentSendingHandler struct {
	Log                 *logrus.Logger
	Viper               *viper.Viper
	Validate            *validator.Validate
	UseCase             usecase.IDocumentSendingUseCase
	OrganizationMessage messaging.IOrganizationMessage
	UserMessage         messaging.IUserMessage
	UserHelper          helper.IUserHelper
}

func NewDocumentSendingHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDocumentSendingUseCase,
	orgMessage messaging.IOrganizationMessage,
	userMessage messaging.IUserMessage,
	userHelper helper.IUserHelper,
) *DocumentSendingHandler {
	return &DocumentSendingHandler{
		Log:                 log,
		Viper:               viper,
		Validate:            validate,
		UseCase:             useCase,
		OrganizationMessage: orgMessage,
		UserMessage:         userMessage,
		UserHelper:          userHelper,
	}
}

func DocumentSendingHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IDocumentSendingHandler {
	validate := config.NewValidator(viper)
	useCase := usecase.DocumentSendingUseCaseFactory(log, viper)
	orgMessage := messaging.OrganizationMessageFactory(log)
	userMessage := messaging.UserMessageFactory(log)
	userHelper := helper.UserHelperFactory(log)
	return NewDocumentSendingHandler(log, viper, validate, useCase, orgMessage, userMessage, userHelper)
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

func (h *DocumentSendingHandler) TestGenerateHTMLPDF(ctx *gin.Context) {
	docSendingID := ctx.Query("doc_sending_id")
	if docSendingID == "" {
		utils.BadRequestResponse(ctx, "Document sending id is required", nil)
		return
	}

	parsedDocSendingID, err := uuid.Parse(docSendingID)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid document sending id", err)
		return
	}

	filepath, err := h.UseCase.TestGenerateHTMLPDF(parsedDocSendingID)
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.TestGenerateHTMLPDF] error when generating html pdf: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate html pdf", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "HTML PDF generated", filepath)
}

func (h *DocumentSendingHandler) GeneratePdfBufferFromHTML(ctx *gin.Context) {
	var payload request.GeneratePdfBufferFromHTMLRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[DocumentSendingHandler.GeneratePdfBufferFromHTML] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[DocumentSendingHandler.GeneratePdfBufferFromHTML] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	c, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	c, cancel = context.WithTimeout(c, 30*time.Second)
	defer cancel()

	cssStyles := `
<style>
body {
	font-size: 20px;
}
.tiptap h1 {
  font-size: 1.4rem;
}

.tiptap h2 {
  font-size: 1.2rem;
}

.tiptap h3 {
  font-size: 1.1rem;
}

.tiptap {
  ul,
  ol {
    padding: 0 1rem;
    margin: 1.25rem 1rem 1.25rem 0.4rem;
  }
  li p {
    margin-top: 0.25em;
    margin-bottom: 0.25em;
  }
  code {
    background-color: var(--purple-light);
    border-radius: 0.4rem;
    color: var(--black);
    font-size: 0.85rem;
    padding: 0.25em 0.3em;
  }

  pre {
    background: var(--black);
    border-radius: 0.5rem;
    color: var(--white);
    font-family: "JetBrainsMono", monospace;
    margin: 1.5rem 0;
    padding: 0.75rem 1rem;

    code {
      background: none;
      color: inherit;
      font-size: 0.8rem;
      padding: 0;
    }
  }

  blockquote {
    border-left: 3px solid var(--gray-3);
    margin: 1.5rem 0;
    padding-left: 1rem;
  }

  hr {
    border: none;
    border-top: 1px solid var(--gray-2);
    margin: 2rem 0;
  }
}

.tiptap table {
  border-collapse: collapse;
  margin: 0;
  overflow: hidden;
  table-layout: fixed;
  width: 100%;
}

.tiptap td,
.tiptap th {
  border: 1px solid var(--primary);
  box-sizing: border-box;
  min-width: 1em;
  padding: 6px 8px;
  position: relative;
  vertical-align: top;
}

.tiptap th {
  background-color: var(--second);
  font-weight: normal !important;
  text-align: left;
}

.tiptap .selectedCell:after {
  background: var(--selectGray);
  content: "";
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  pointer-events: none;
  position: absolute;
  z-index: 2;
}

.tiptap .column-resize-handle {
  background-color: var(--gray);
  bottom: -2px;
  pointer-events: none;
  position: absolute;
  right: -2px;
  top: 0;
  width: 1px;
}

.tiptap .tableWrapper {
  margin: 1.5rem 0;
  overflow-x: auto;
}

.tiptap.resize-cursor {
  cursor: ew-resize;
  cursor: col-resize;
}

.tiptap-border-none {
  border: 0px solid transparent !important;
  background-color: transparent !important;
}

.tiptap-border-none th {
  border: 0px solid transparent !important;
  background-color: transparent !important;
}

.tiptap-border-none td {
  border: 0px solid transparent !important;
  background-color: transparent !important;
}
</style>
`

	htmlContent := `<html><head><meta charset="UTF-8">` + cssStyles + `</head><body><div class="tiptap">` + payload.HTML + `</div></body></html>`
	dataURL := "data:text/html," + url.PathEscape(htmlContent)

	var pdfBuffer []byte

	err := chromedp.Run(c,
		chromedp.Navigate(dataURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuffer, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(1.0).
				WithMarginRight(1.0).
				WithMarginBottom(1.0).
				WithMarginLeft(1.0).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		h.Log.Errorf("Failed to generate PDF: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate PDF", err.Error())
		return
	}

	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "attachment; filename=document.pdf")
	ctx.Header("Content-Length", strconv.Itoa(len(pdfBuffer)))

	// Write the PDF buffer to the response
	_, err = ctx.Writer.Write(pdfBuffer)
	if err != nil {
		h.Log.Errorf("Failed to write PDF to response: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to write PDF to response", err.Error())
		return
	}
}

func (h *DocumentSendingHandler) GeneratePdfBufferForDocumentSending(ctx *gin.Context) {
	var payload request.GeneratePdfBufferFromHTMLRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("[DocumentSendingHandler.GeneratePdfBufferFromHTML] error when binding request: %v", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[DocumentSendingHandler.GeneratePdfBufferFromHTML] error when validating request: %v", err)
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if payload.DocumentSendingID == "" {
		utils.BadRequestResponse(ctx, "Document sending id is required", nil)
		return
	}

	parsedDocSendingID, err := uuid.Parse(payload.DocumentSendingID)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid document sending id", err)
		return
	}

	documentSending, err := h.UseCase.FindByID(parsedDocSendingID.String())
	if err != nil {
		h.Log.Errorf("[DocumentSendingHandler.GeneratePdfBufferForDocumentSending] error when finding document sending: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find document sending", err.Error())
		return
	}

	if documentSending == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Document sending not found", "Document sending not found")
		return
	}

	pdfBuffer, err := h.UseCase.GeneratePdfBuffer(documentSending.ID, payload.HTML)

	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "attachment; filename=document.pdf")
	ctx.Header("Content-Length", strconv.Itoa(len(pdfBuffer)))

	// Write the PDF buffer to the response
	_, err = ctx.Writer.Write(pdfBuffer)
	if err != nil {
		h.Log.Errorf("Failed to write PDF to response: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to write PDF to response", err.Error())
		return
	}
}
