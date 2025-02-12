package request

type CreateDocumentSendingRequest struct {
	ProjectRecruitmentLineID string  `json:"project_recruitment_line_id" validate:"required,uuid"`
	ApplicantID              string  `json:"applicant_id" validate:"required,uuid"`
	DocumentSetupID          string  `json:"document_setup_id" validate:"required,uuid"`
	DocumentDate             string  `json:"document_date" validate:"required,datetime=2006-01-02"`
	DocumentNumber           string  `json:"document_number" validate:"required"`
	Status                   string  `json:"status" validate:"required,oneof=DRAFT PENDING SENT FAILED"`
	BasicWage                float64 `json:"basic_wage" validate:"required"`
	PositionalAllowance      float64 `json:"positional_allowance" validate:"required"`
	OperationalAllowance     float64 `json:"operational_allowance" validate:"required"`
	MealAllowance            float64 `json:"meal_allowance" validate:"required"`
	JobLocation              string  `json:"job_location" validate:"required"`
	HometripTicket           string  `json:"hometrip_ticket" validate:"required"`
	PeriodAgreement          string  `json:"period_agreement" validate:"required"`
	HomeLocation             string  `json:"home_location" validate:"required"`
	JobLevelID               string  `json:"job_level_id" validate:"required,uuid"`
	// JobID                    string  `json:"job_id" validate:"uuid"`
	JobPostingID      string `json:"job_posting_id" validate:"required,uuid"`
	JoinedDate        string `json:"joined_date" validate:"datetime=2006-01-02"`
	ForOrganizationID string `json:"for_organization_id" validate:"required,uuid"`
	DetailContent     string `json:"detail_content" validate:"required"`
}

type UpdateDocumentSendingRequest struct {
	ID                       string  `json:"id" validate:"required,uuid"`
	ProjectRecruitmentLineID string  `json:"project_recruitment_line_id" validate:"required,uuid"`
	ApplicantID              string  `json:"applicant_id" validate:"required,uuid"`
	DocumentSetupID          string  `json:"document_setup_id" validate:"required,uuid"`
	DocumentDate             string  `json:"document_date" validate:"required,datetime=2006-01-02"`
	DocumentNumber           string  `json:"document_number" validate:"required"`
	Status                   string  `json:"status" validate:"required,oneof=DRAFT PENDING SENT FAILED"`
	BasicWage                float64 `json:"basic_wage" validate:"required"`
	PositionalAllowance      float64 `json:"positional_allowance" validate:"required"`
	OperationalAllowance     float64 `json:"operational_allowance" validate:"required"`
	MealAllowance            float64 `json:"meal_allowance" validate:"required"`
	JobLocation              string  `json:"job_location" validate:"required"`
	HometripTicket           string  `json:"hometrip_ticket" validate:"required"`
	PeriodAgreement          string  `json:"period_agreement" validate:"required"`
	HomeLocation             string  `json:"home_location" validate:"required"`
	JobLevelID               string  `json:"job_level_id" validate:"required,uuid"`
	// JobID                    string  `json:"job_id" validate:"uuid"`
	JobPostingID      string `json:"job_posting_id" validate:"required,uuid"`
	JoinedDate        string `json:"joined_date" validate:"datetime=2006-01-02"`
	ForOrganizationID string `json:"for_organization_id" validate:"required,uuid"`
	DetailContent     string `json:"detail_content" validate:"required"`
}
