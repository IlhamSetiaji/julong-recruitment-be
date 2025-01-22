package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type TestTypeResponse struct {
	ID              uuid.UUID                     `json:"id"`
	Name            string                        `json:"name"`
	RecruitmentType entity.ProjectRecruitmentType `json:"recruitment_type"`
	Status          entity.TestTypeStatus         `json:"status"`
	CreatedAt       time.Time                     `json:"created_at"`
	UpdatedAt       time.Time                     `json:"updated_at"`
}
