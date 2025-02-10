package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type FgdResultResponse struct {
	ID             uuid.UUID              `json:"id"`
	FgdApplicantID uuid.UUID              `json:"fgd_applicant_id"`
	FgdAssessorID  uuid.UUID              `json:"fgd_assessor_id"`
	Status         entity.FgdResultStatus `json:"status"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`

	FgdApplicant *FgdApplicantResponse `json:"fgd_applicant"`
	FgdAssessor  *FgdAssessorResponse  `json:"fgd_assessor"`
}
