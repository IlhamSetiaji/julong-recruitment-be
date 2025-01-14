package response

import "github.com/google/uuid"

type QuestionResponse struct {
	ID                 uuid.UUID `json:"id"`
	TemplateQuestionID uuid.UUID `json:"template_question_id"`
	AnswerTypeID       uuid.UUID `json:"answer_type_id"`
	Name               string    `json:"name"`

	AnswerTypeResponse *AnswerTypeResponse `json:"answer_types"`
}
