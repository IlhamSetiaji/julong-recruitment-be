package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type DocumentSetupResponse struct {
	ID              uuid.UUID                     `json:"id"`
	Title           string                        `json:"title"`
	DocumentTypeID  uuid.UUID                     `json:"document_type_id"`
	RecruitmentType entity.ProjectRecruitmentType `json:"recruitment_type"`
	Header          string                        `json:"header"`
	Body            string                        `json:"body"`
	Footer          string                        `json:"footer"`
	CreatedAt       time.Time                     `json:"created_at"`
	UpdatedAt       time.Time                     `json:"updated_at"`

	DocumentType *DocumentTypeResponse `json:"document_type"`
}
