package response

import "github.com/google/uuid"

type DocumentVerificationResponse struct {
	ID                 uuid.UUID `json:"id"`
	TemplateQuestionID uuid.UUID `json:"template_question_id"`
	Name               string    `json:"name"`
	Format             string    `json:"format"`

	TemplateQuestion *TemplateQuestionResponse `json:"template_question"`
}
