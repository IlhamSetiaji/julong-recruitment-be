package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type DocumentAgreementResponse struct {
	ID                uuid.UUID                      `json:"id"`
	DocumentSendingID uuid.UUID                      `json:"document_sending_id"`
	ApplicantID       uuid.UUID                      `json:"applicant_id"`
	Status            entity.DocumentAgreementStatus `json:"status"`
	Path              string                         `json:"path"`
	CreatedAt         time.Time                      `json:"created_at"`
	UpdatedAt         time.Time                      `json:"updated_at"`

	DocumentSending *DocumentSendingResponse `json:"document_sending"`
	Applicant       *ApplicantResponse       `json:"applicant"`
}
