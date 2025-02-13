package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentAgreementHandler interface {
	CreateDocumentAgreement(ctx *gin.Context)
	UpdateDocumentAgreement(ctx *gin.Context)
	UpdateStatusDocumentAgreement(ctx *gin.Context)
	FindByDocumentSendingIDAndApplicantID(ctx *gin.Context)
}

type DocumentAgreementHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IDocumentAgreementUseCase
}

func NewDocumentAgreementHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IDocumentAgreementUseCase,
) IDocumentAgreementHandler {
	return &DocumentAgreementHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func DocumentAgreementHandlerFactory(log *logrus.Logger, viper *viper.Viper) IDocumentAgreementHandler {
	useCase := usecase.DocumentAgreementUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	return NewDocumentAgreementHandler(log, viper, validate, useCase)
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
