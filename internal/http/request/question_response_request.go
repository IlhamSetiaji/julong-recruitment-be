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

type InterviewQuestionResponseRequest struct {
	TemplateQuestionID string `json:"template_question_id" validate:"required,uuid"`
	Questions          []struct {
		ID      string                   `json:"id" validate:"required,uuid"`
		Answers []InterviewAnswerRequest `json:"answers" validate:"required,dive"`
	} `json:"questions" validate:"required,dive"`
	DeletedAnswerIDs []string `json:"deleted_answer_ids" validate:"omitempty,dive,uuid"`
}

type InterviewAnswerRequest struct {
	ID            string `json:"id" validate:"omitempty,uuid"`
	JobPostingID  string `json:"job_posting_id" validate:"required,uuid"`
	UserProfileID string `json:"user_profile_id" validate:"required,uuid"`
	Answer        string `json:"answer" validate:"required"`
}
