package request

import "mime/multipart"

type QuestionResponseRequest struct {
	QuestionID string `json:"question_id" validate:"required,uuid"`
	Answers    []struct {
		ID            string                `json:"id" validate:"omitempty,uuid"`
		JobPostingID  string                `json:"job_posting_id" validate:"required,uuid"`
		UserProfileID string                `json:"user_profile_id" validate:"required,uuid"`
		Answer        string                `json:"answer" validate:"omitempty"`
		AnswerFile    *multipart.FileHeader `json:"answer_file" validate:"omitempty"`
		AnswerPath    string                `json:"answer_path" validate:"omitempty"`
	} `json:"answers" validate:"required,dive"`
	DeletedAnswerIDs []string `json:"deleted_answer_ids" validate:"omitempty,dive,uuid"`
}
