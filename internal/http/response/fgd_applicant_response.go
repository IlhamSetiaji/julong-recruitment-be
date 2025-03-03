package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type FgdApplicantResponse struct {
	ID               uuid.UUID                `json:"id"`
	FgdScheduleID    uuid.UUID                `json:"fgd_schedule_id"`
	ApplicantID      uuid.UUID                `json:"applicant_id"`
	UserProfileID    uuid.UUID                `json:"user_profile_id"`
	StartTime        time.Time                `json:"start_time"`
	EndTime          time.Time                `json:"end_time"`
	StartedTime      *time.Time               `json:"started_time"`
	EndedTime        *time.Time               `json:"ended_time"`
	AssessmentStatus entity.AssessmentStatus  `json:"assessment_status"`
	FinalResult      entity.FinalResultStatus `json:"final_result"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`

	UserProfile       *UserProfileResponse `json:"user_profile"`
	Applicant         *ApplicantResponse   `json:"applicant"`
	FgdResultAssessor *FgdResultResponse   `json:"fgd_result_assessor"`
	FgdResults        []FgdResultResponse  `json:"fgd_results"`
}
