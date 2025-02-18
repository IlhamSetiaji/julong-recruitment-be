package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type ApplicantResponse struct {
	ID               uuid.UUID                     `json:"id"`
	UserProfileID    uuid.UUID                     `json:"user_profile_id"`
	JobPostingID     uuid.UUID                     `json:"job_posting_id"`
	AppliedDate      time.Time                     `json:"applied_date"`
	Order            int                           `json:"order"`
	Status           entity.ApplicantStatus        `json:"status"`
	ProcessStatus    entity.ApplicantProcessStatus `json:"process_status"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
	TemplateQuestion *TemplateQuestionResponse     `json:"template_question"`
	JobPosting       *JobPostingResponse           `json:"job_posting"`
	UserProfile      *UserProfileResponse          `json:"user_profile"`
}
