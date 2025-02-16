package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type DocumentVerificationHeaderResponse struct {
	ID                       uuid.UUID                               `json:"id"`
	ProjectRecruitmentLineID uuid.UUID                               `json:"project_recruitment_line_id"`
	ApplicantID              uuid.UUID                               `json:"applicant_id"`
	JobPostingID             uuid.UUID                               `json:"job_posting_id"`
	VerifiedBy               *uuid.UUID                              `json:"verified_by"`
	Status                   entity.DocumentVerificationHeaderStatus `json:"status"`
	CreatedAt                time.Time                               `json:"created_at"`
	UpdatedAt                time.Time                               `json:"updated_at"`

	ProjectRecruitmentLine    *ProjectRecruitmentLineResponse    `json:"project_recruitment_line"`
	Applicant                 *ApplicantResponse                 `json:"applicant"`
	JobPosting                *JobPostingResponse                `json:"job_posting"`
	DocumentVerificationLines []DocumentVerificationLineResponse `json:"document_verification_lines"`

	EmployeeName string `json:"employee_name"`
}
