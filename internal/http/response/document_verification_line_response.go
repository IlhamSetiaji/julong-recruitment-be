package response

import (
	"time"

	"github.com/google/uuid"
)

type DocumentVerificationLineResponse struct {
	ID                           uuid.UUID `json:"id"`
	DocumentVerificationHeaderID uuid.UUID `json:"document_verification_header_id"`
	DocumentVerificationID       uuid.UUID `json:"document_verification_id"`
	Path                         string    `json:"path"`
	Answer                       string    `json:"answer"`
	CreatedAt                    time.Time `json:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at"`

	DocumentVerificationHeader *DocumentVerificationHeaderResponse `json:"document_verification_header"`
	DocumentVerification       *DocumentVerificationResponse       `json:"document_verification"`
}
