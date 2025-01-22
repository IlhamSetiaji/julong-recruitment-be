package response

import (
	"time"

	"github.com/google/uuid"
)

type TestApplicantResponse struct {
	ID                   uuid.UUID `json:"id"`
	TestScheduleHeaderID uuid.UUID `json:"test_schedule_header_id"`
	UserProfileID        uuid.UUID `json:"user_profile_id"`
	StartTime            string    `json:"start_time"`
	EndTime              string    `json:"end_time"`
	FinalResult          string    `json:"final_result"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	UserProfile *UserProfileResponse `json:"user_profile"`
}
