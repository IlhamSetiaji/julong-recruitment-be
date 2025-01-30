package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type AdministrativeSelectionResponse struct {
	ID           uuid.UUID                            `json:"id"`
	JobPostingID uuid.UUID                            `json:"job_posting_id"`
	ProjectPicID uuid.UUID                            `json:"project_pic_id"`
	Status       entity.AdministrativeSelectionStatus `json:"status"`
	// VerifiedAt      time.Time                            `json:"verified_at"`
	// VerifiedBy      *uuid.UUID                           `json:"verified_by"`
	DocumentDate    time.Time `json:"document_date"`
	DocumentNumber  string    `json:"document_number"`
	TotalApplicants int       `json:"total_applicants"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	JobPosting            *JobPostingResponse            `json:"job_posting"`
	ProjectPIC            *ProjectPicResponse            `json:"project_pic"`
	AdministrativeResults []AdministrativeResultResponse `json:"administrative_results"`
}
