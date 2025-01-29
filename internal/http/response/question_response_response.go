package response

import (
	"time"

	"github.com/google/uuid"
)

type QuestionResponseResponse struct {
	ID            uuid.UUID `json:"id"`
	QuestionID    uuid.UUID `json:"question_id"`
	UserProfileID uuid.UUID `json:"user_profile_id"`
	JobPostingID  uuid.UUID `json:"job_posting_id"`
	Answer        string    `json:"answer"`
	AnswerFile    string    `json:"answer_file"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Question    *QuestionResponse    `json:"question"`
	UserProfile *UserProfileResponse `json:"user_profile"`
	JobPosting  *JobPostingResponse  `json:"job_posting"`
}

type QuestionAnswerJSON struct {
	Answers []struct {
		QuestionID uuid.UUID `json:"question_id"`
		Answer     string    `json:"answer"`
	} `json:"answers"`
}
