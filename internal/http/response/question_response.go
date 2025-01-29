package response

import (
	"time"

	"github.com/google/uuid"
)

type QuestionResponse struct {
	ID                 uuid.UUID `json:"id"`
	TemplateQuestionID uuid.UUID `json:"template_question_id"`
	AnswerTypeID       uuid.UUID `json:"answer_type_id"`
	Name               string    `json:"name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	AnswerTypeResponse *AnswerTypeResponse         `json:"answer_types"`
	QuestionOptions    *[]QuestionOptionResponse   `json:"question_options"`
	QuestionResponses  *[]QuestionResponseResponse `json:"question_responses"`
}
