package request

import "mime/multipart"

type QuestionResponseRequest struct {
	QuestionID       string          `form:"question_id" validate:"required,uuid"`
	Answers          []AnswerRequest `form:"answers" validate:"omitempty,dive"`
	DeletedAnswerIDs []string        `form:"deleted_answer_ids" validate:"omitempty,dive,uuid"`
}

type AnswerRequest struct {
	ID            *string               `form:"id" validate:"omitempty,uuid"`
	JobPostingID  string                `form:"job_posting_id" validate:"required,uuid"`
	UserProfileID string                `form:"user_profile_id" validate:"required,uuid"`
	Answer        string                `form:"answer" validate:"omitempty"`
	AnswerFile    *multipart.FileHeader `form:"answer_file" validate:"omitempty"`
	AnswerPath    string                `form:"answer_path" validate:"omitempty"`
}
