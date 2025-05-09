package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/service"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentAgreementHandler interface {
	CreateDocumentAgreement(ctx *gin.Context)
	UpdateDocumentAgreement(ctx *gin.Context)
	UpdateStatusDocumentAgreement(ctx *gin.Context)
	FindByDocumentSendingIDAndApplicantID(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
}

type DocumentAgreementHandler struct {
	Log                           *logrus.Logger
	Viper                         *viper.Viper
	Validate                      *validator.Validate
	UseCase                       usecase.IDocumentAgreementUseCase
	ProjectRecruitmentLineUseCase usecase.IProjectRecruitmentLineUseCase
	UserMessage                   messaging.IUserMessage
	NotificationService           service.INotificationService
}

func NewDocumentAgreementHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDocumentAgreementUseCase,
	projectRecruitmentLineUseCase usecase.IProjectRecruitmentLineUseCase,
	userMessage messaging.IUserMessage,
	notificationService service.INotificationService,
) IDocumentAgreementHandler {
	return &DocumentAgreementHandler{
		Log:                           log,
		Viper:                         viper,
		Validate:                      validate,
		UseCase:                       useCase,
		ProjectRecruitmentLineUseCase: projectRecruitmentLineUseCase,
		UserMessage:                   userMessage,
		NotificationService:           notificationService,
	}
}

func DocumentAgreementHandlerFactory(log *logrus.Logger, viper *viper.Viper) IDocumentAgreementHandler {
	useCase := usecase.DocumentAgreementUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	projectRecruitmentLineUseCase := usecase.ProjectRecruitmentLineUseCaseFactory(log)
	userMessage := messaging.UserMessageFactory(log)
	notificationService := service.NotificationServiceFactory(viper, log)
	return NewDocumentAgreementHandler(
		log,
		viper,
		validate,
		useCase,
		projectRecruitmentLineUseCase,
		userMessage,
		notificationService,
	)
}

// CreateDocumentAgreement create document agreement
//
// @Summary Create document agreement
// @Description Create document agreement
// @Tags Document Agreement
// @Accept json
// @Produce json
// @Param document_agreement body request.CreateDocumentAgreementRequest true "Document Agreement"
// @Success 200 {object} response.DocumentAgreementResponse
// @Security BearerAuth
// @Router /document-agreement [post]
func (h *DocumentAgreementHandler) CreateDocumentAgreement(ctx *gin.Context) {
	var req request.CreateDocumentAgreementRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[DocumentAgreementHandler.CreateDocumentAgreement] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[DocumentAgreementHandler.CreateDocumentAgreement] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	// handle file uploads
	if req.File != nil {
		timestamp := time.Now().UnixNano()
		filePath := "storage/document-agreement/" + strconv.FormatInt(timestamp, 10) + "_" + req.File.Filename
		if err := ctx.SaveUploadedFile(req.File, filePath); err != nil {
			h.Log.Error("[DocumentAgreementHandler.CreateDocumentAgreement] " + err.Error())
			utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Failed to upload file", err.Error())
			return
		}
		req.Path = filePath
	}

	res, err := h.UseCase.CreateDocumentAgreement(&req)
	if err != nil {
		h.Log.Error("[DocumentAgreementHandler.CreateDocumentAgreement] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create document agreement", err.Error())
		return
	}

	prl, err := h.ProjectRecruitmentLineUseCase.FindByID(res.DocumentSending.ProjectRecruitmentLineID)
	if err != nil {
		h.Log.Error("[DocumentAgreementHandler.CreateDocumentAgreement] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find project recruitment line", err.Error())
		return
	}

	if prl == nil {
		h.Log.Error("[DocumentAgreementHandler.CreateDocumentAgreement] " + errors.New("project recruitment line not found").Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Project recruitment line not found", err.Error())
		return
	}

	var userIDs []string
	for _, pp := range prl.ProjectPics {
		userResp, err := h.UserMessage.SendFindUserByEmployeeIDMessage(pp.EmployeeID.String())
		if err != nil {
			h.Log.Error("[DocumentAgreementHandler.CreateDocumentAgreement] " + err.Error())
			continue
			// utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find user by employee id", err.Error())
			// return
		}
		if userResp != nil {
			userIDs = append(userIDs, userResp.ID)
			h.Log.Infof("[DocumentAgreementHandler.CreateDocumentAgreement] Appended user ID: %s", userIDs)
		}
	}

	if len(userIDs) > 0 {
		err = h.NotificationService.CreateDocumentAgreementNotification(res.Applicant.UserProfile.UserID.String(), userIDs, res.DocumentSending.DocumentSetup.Title)
		if err != nil {
			h.Log.Error("[DocumentAgreementHandler.CreateDocumentAgreement] " + err.Error())
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create document agreement notification", err.Error())
			return
		}
		h.Log.Infof("[DocumentAgreementHandler.CreateDocumentAgreement] Document agreement notification created successfully")
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Document agreement created", res)
}

// UpdateDocumentAgreement update document agreement
//
// @Summary Update document agreement
// @Description Update document agreement
// @Tags Document Agreement
// @Accept json
// @Produce json
// @Param document_agreement body request.UpdateDocumentAgreementRequest true "Document Agreement"
// @Success 200 {object} response.DocumentAgreementResponse
// @Security BearerAuth
// @Router /document-agreement/update [put]
func (h *DocumentAgreementHandler) UpdateDocumentAgreement(ctx *gin.Context) {
	var req request.UpdateDocumentAgreementRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[DocumentAgreementHandler.UpdateDocumentAgreement] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[DocumentAgreementHandler.UpdateDocumentAgreement] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	// handle file uploads
	if req.File != nil {
		timestamp := time.Now().UnixNano()
		filePath := "storage/document-agreement/" + strconv.FormatInt(timestamp, 10) + "_" + req.File.Filename
		if err := ctx.SaveUploadedFile(req.File, filePath); err != nil {
			h.Log.Error("[DocumentAgreementHandler.UpdateDocumentAgreement] " + err.Error())
			utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Failed to upload file", err.Error())
			return
		}
		req.Path = filePath
	}

	res, err := h.UseCase.UpdateDocumentAgreement(&req)
	if err != nil {
		h.Log.Error("[DocumentAgreementHandler.UpdateDocumentAgreement] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update document agreement", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document agreement updated", res)
}

// UpdateStatusDocumentAgreement update status document agreement
//
// @Summary Update status document agreement
// @Description Update status document agreement
// @Tags Document Agreement
// @Accept json
// @Produce json
// @Param status body request.UpdateStatusDocumentAgreementRequest true "Document Agreement Status"
// @Success 200 {object} response.DocumentAgreementResponse
// @Security BearerAuth
// @Router /document-agreement/update-status [put]
func (h *DocumentAgreementHandler) UpdateStatusDocumentAgreement(ctx *gin.Context) {
	var req request.UpdateStatusDocumentAgreementRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[DocumentAgreementHandler.UpdateStatusDocumentAgreement] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[DocumentAgreementHandler.UpdateStatusDocumentAgreement] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.UpdateStatusDocumentAgreement(&req)
	if err != nil {
		h.Log.Error("[DocumentAgreementHandler.UpdateStatusDocumentAgreement] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update status document agreement", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Status document agreement updated", res)
}

// FindByDocumentSendingIDAndApplicantID find document agreement by document sending id and applicant id
//
// @Summary Find document agreement by document sending id and applicant id
// @Description Find document agreement by document sending id and applicant id
// @Tags Document Agreement
// @Accept json
// @Produce json
// @Param document_sending_id query string true "Document Sending ID"
// @Param applicant_id query string true "Applicant ID"
// @Success 200 {object} response.DocumentAgreementResponse
// @Security BearerAuth
// @Router /document-agreement/find [get]
func (h *DocumentAgreementHandler) FindByDocumentSendingIDAndApplicantID(ctx *gin.Context) {
	documentSendingID := ctx.Query("document_sending_id")
	applicantID := ctx.Query("applicant_id")

	res, err := h.UseCase.FindByDocumentSendingIDAndApplicantID(documentSendingID, applicantID)
	if err != nil {
		h.Log.Error("[DocumentAgreementHandler.FindByDocumentSendingIDAndApplicantID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find document agreement", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document agreement found", res)
}

// FindAllPaginated find all document agreement paginated
//
// @Summary Find all document agreement paginated
// @Description Find all document agreement paginated
// @Tags Document Agreement
// @Accept json
// @Produce json
// @Param page query int true "Page"
// @Param page_size query int true "Page Size"
// @Param search query string false "Search"
// @Param sort query string false "Sort"
// @Param status query string false "Status"
// @Param document_type_id query string false "Document Type ID"
// @Success 200 {object} response.DocumentAgreementResponse
// @Security BearerAuth
// @Router /document-agreement [get]
func (h *DocumentAgreementHandler) FindAllPaginated(ctx *gin.Context) {
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

	status := ctx.Query("status")
	if status != "" {
		status = string(entity.DocumentAgreementStatus(status))
	}

	filter := map[string]interface{}{
		"status": status,
	}

	documentTypeID := ctx.Query("document_type_id")
	if documentTypeID == "" {
		documentTypeID = ""
	}

	res, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort, filter, documentTypeID)
	if err != nil {
		h.Log.Error("[DocumentAgreementHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all document agreement", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document agreement found", gin.H{
		"document_agreements": res,
		"total":               total,
	})
}

// FindByID find document agreement by id
//
// @Summary Find document agreement by id
// @Description Find document agreement by id
// @Tags Document Agreement
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} response.DocumentAgreementResponse
// @Security BearerAuth
// @Router /document-agreement/{id} [get]
func (h *DocumentAgreementHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[DocumentAgreementHandler.FindByID] " + err.Error())
		utils.BadRequestResponse(ctx, "Invalid ID", err)
		return
	}

	res, err := h.UseCase.FindByID(parsedID)
	if err != nil {
		h.Log.Error("[DocumentAgreementHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find document agreement", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document agreement found", res)
}
