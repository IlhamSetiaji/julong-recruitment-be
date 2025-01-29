package response

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type AdministrativeResultResponse struct {
	ID                        uuid.UUID                         `json:"id"`
	AdministrativeSelectionID uuid.UUID                         `json:"administrative_selection_id"`
	ApplicantID               uuid.UUID                         `json:"applicant_id"`
	Status                    entity.AdministrativeResultStatus `json:"status"`
	CreatedAt                 string                            `json:"created_at"`
	UpdatedAt                 string                            `json:"updated_at"`

	Applicant *ApplicantResponse `json:"applicant"`
}
