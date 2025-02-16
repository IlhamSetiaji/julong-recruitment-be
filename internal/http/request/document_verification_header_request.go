package request

type CreateDocumentVerificationHeaderRequest struct {
	ProjectRecruitmentLineID string `json:"project_recruitment_line_id" validate:"required,uuid"`
	ApplicantID              string `json:"applicant_id" validate:"required,uuid"`
	JobPostingID             string `json:"job_posting_id" validate:"required,uuid"`
	Status                   string `json:"status" validate:"required,document_verification_header_status_validation"`
}

type UpdateDocumentVerificationHeaderRequest struct {
	ID                       string  `json:"id" validate:"required,uuid"`
	ProjectRecruitmentLineID string  `json:"project_recruitment_line_id" validate:"required,uuid"`
	ApplicantID              string  `json:"applicant_id" validate:"required,uuid"`
	JobPostingID             string  `json:"job_posting_id" validate:"required,uuid"`
	VerifiedBy               *string `json:"verified_by" validate:"omitempty,uuid"`
	Status                   string  `json:"status" validate:"required,document_verification_header_status_validation"`
}
