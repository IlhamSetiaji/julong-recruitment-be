package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type TemplateActivityLineResponse struct {
	ID                 uuid.UUID                         `json:"id"`
	TemplateActivityID uuid.UUID                         `json:"template_activity_id"`
	Description        string                            `json:"description"`
	Status             entity.TemplateActivityLineStatus `json:"status"`
	TemplateQuestionID uuid.UUID                         `json:"question_template_id"`
	ColorHexCode       string                            `json:"color_hex_code"`
	CreatedAt          time.Time                         `json:"created_at"`
	UpdatedAt          time.Time                         `json:"updated_at"`

	TemplateActivity *TemplateActivityResponse `json:"template_activity"`
	TemplateQuestion *TemplateQuestionResponse `json:"template_question"`
}
