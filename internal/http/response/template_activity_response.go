package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type TemplateActivityResponse struct {
	ID              uuid.UUID                     `json:"id"`
	Name            string                        `json:"name"`
	Description     string                        `json:"description"`
	RecruitmentType entity.ProjectRecruitmentType `json:"recruitment_type"`
	Status          entity.TemplateActivityStatus `json:"status"`
	CreatedAt       time.Time                     `json:"created_at"`
	UpdatedAt       time.Time                     `json:"updated_at"`

	TemplateActivityLines *[]TemplateActivityLineResponse `json:"template_activity_lines"`
}
