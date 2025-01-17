package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type ProjectRecruitmentHeaderResponse struct {
	ID                 uuid.UUID                             `json:"id"`
	TemplateActivityID uuid.UUID                             `json:"template_activity_id"`
	Name               string                                `json:"name"`
	Description        string                                `json:"description"`
	DocumentDate       time.Time                             `json:"document_date"`
	DocumentNumber     string                                `json:"document_number"`
	RecruitmentType    entity.ProjectRecruitmentType         `json:"recruitment_type"`
	StartDate          time.Time                             `json:"start_date"`
	EndDate            time.Time                             `json:"end_date"`
	Status             entity.ProjectRecruitmentHeaderStatus `json:"status"`

	TemplateActivity *TemplateActivityResponse `json:"template_activity"`
}
