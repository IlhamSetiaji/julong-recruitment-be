package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type TestApplicantResponse struct {
	ID                   uuid.UUID                `json:"id"`
	TestScheduleHeaderID uuid.UUID                `json:"test_schedule_header_id"`
	UserProfileID        uuid.UUID                `json:"user_profile_id"`
	StartTime            time.Time                `json:"start_time"`
	EndTime              time.Time                `json:"end_time"`
	FinalResult          entity.FinalResultStatus `json:"final_result"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            time.Time                `json:"updated_at"`

	UserProfile        *UserProfileResponse        `json:"user_profile"`
	Applicant          *ApplicantResponse          `json:"applicant"`
	TestScheduleHeader *TestScheduleHeaderResponse `json:"test_schedule_header"`
}
