package request

import "mime/multipart"

type CreateDocumentAgreementRequest struct {
	DocumentSendingID string                `form:"document_sending_id" validate:"required,uuid"`
	ApplicantID       string                `form:"applicant_id" validate:"required,uuid"`
	File              *multipart.FileHeader `form:"file" validate:"required"`
	Path              string                `form:"path" validate:"omitempty"`
}

type UpdateDocumentAgreementRequest struct {
	ID                string                `form:"id" validate:"required,uuid"`
	DocumentSendingID string                `form:"document_sending_id" validate:"required,uuid"`
	ApplicantID       string                `form:"applicant_id" validate:"required,uuid"`
	File              *multipart.FileHeader `form:"file" validate:"required"`
	Path              string                `form:"path" validate:"omitempty"`
}

type UpdateStatusDocumentAgreementRequest struct {
	ID     string `json:"id" validate:"required,uuid"`
	Status string `json:"status" validate:"required,document_agreement_status_validation"`
}
