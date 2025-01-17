package response

import (
	"time"

	"github.com/google/uuid"
)

type ProjectRecruitmentLineResponse struct {
	ID                         uuid.UUID `json:"id"`
	ProjectRecruitmentHeaderID uuid.UUID `json:"project_recruitment_header_id"`
	TemplateActivityLineID     uuid.UUID `json:"template_activity_id"`
	StartDate                  time.Time `json:"start_date"`
	EndDate                    time.Time `json:"end_date"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`

	ProjectPics          []ProjectPicResponse          `json:"project_pics"`
	TemplateActivityLine *TemplateActivityLineResponse `json:"template_activity"`
}
