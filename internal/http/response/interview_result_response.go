package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type InterviewResultResponse struct {
	ID                   uuid.UUID                    `json:"id"`
	InterviewApplicantID uuid.UUID                    `json:"interview_applicant_id"`
	InterviewAssessorID  uuid.UUID                    `json:"interview_assessor_id"`
	Status               entity.InterviewResultStatus `json:"status"`
	CreatedAt            time.Time                    `json:"created_at"`
	UpdatedAt            time.Time                    `json:"updated_at"`

	InterviewApplicant *InterviewApplicantResponse `json:"interview_applicant"`
	InterviewAssessor  *InterviewAssessorResponse  `json:"interview_assessor"`
}
