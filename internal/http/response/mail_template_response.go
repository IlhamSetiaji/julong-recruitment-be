package response

import (
	"time"

	"github.com/google/uuid"
)

type MailTemplateResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	DocumentTypeID uuid.UUID `json:"document_type_id"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	DocumentType *DocumentTypeResponse `json:"document_type"`
}
