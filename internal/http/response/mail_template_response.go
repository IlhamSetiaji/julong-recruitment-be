package response

import "github.com/google/uuid"

type MailTemplateResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	DocumentTypeID uuid.UUID `json:"document_type_id"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body"`

	DocumentType *DocumentTypeResponse `json:"document_type"`
}
