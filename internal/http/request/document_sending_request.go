package request

type CreateDocumentSendingRequest struct {
	ProjectRecruitmentLineID string  `json:"project_recruitment_line_id" validate:"required,uuid"`
	ApplicantID              string  `json:"applicant_id" validate:"required,uuid"`
	DocumentSetupID          string  `json:"document_setup_id" validate:"required,uuid"`
	DocumentDate             string  `json:"document_date" validate:"required,datetime=2006-01-02"`
	DocumentNumber           string  `json:"document_number" validate:"required"`
	Status                   string  `json:"status" validate:"required,document_sending_status_validation"`
	BasicWage                float64 `json:"basic_wage" validate:"omitempty"`
	PositionalAllowance      float64 `json:"positional_allowance" validate:"omitempty"`
	OperationalAllowance     float64 `json:"operational_allowance" validate:"omitempty"`
	MealAllowance            float64 `json:"meal_allowance" validate:"omitempty"`
	JobLocation              string  `json:"job_location" validate:"omitempty"`
	HometripTicket           string  `json:"hometrip_ticket" validate:"omitempty"`
	PeriodAgreement          string  `json:"period_agreement" validate:"omitempty"`
	HomeLocation             string  `json:"home_location" validate:"omitempty"`
	JobLevelID               string  `json:"job_level_id" validate:"omitempty"`
	JobID                    string  `json:"job_id" validate:"omitempty"`
	JobPostingID             string  `json:"job_posting_id" validate:"required,uuid"`
	JoinedDate               string  `json:"joined_date" validate:"omitempty,datetime=2006-01-02"`
	ForOrganizationID        string  `json:"for_organization_id" validate:"omitempty"`
	OrganizationLocationID   string  `json:"organization_location_id" validate:"omitempty"`
	DetailContent            string  `json:"detail_content" validate:"required"`
	RecruitmentType          string  `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	SyncMidsuit              string  `json:"sync_midsuit" validate:"omitempty,oneof=YES NO"`
	HiredStatus              string  `json:"hired_status" validate:"omitempty,document_sending_hired_status_validation"`
}

type UpdateDocumentSendingRequest struct {
	ID                       string  `json:"id" validate:"required,uuid"`
	ProjectRecruitmentLineID string  `json:"project_recruitment_line_id" validate:"required,uuid"`
	ApplicantID              string  `json:"applicant_id" validate:"required,uuid"`
	DocumentSetupID          string  `json:"document_setup_id" validate:"required,uuid"`
	DocumentDate             string  `json:"document_date" validate:"required,datetime=2006-01-02"`
	DocumentNumber           string  `json:"document_number" validate:"required"`
	Status                   string  `json:"status" validate:"required,document_sending_status_validation"`
	BasicWage                float64 `json:"basic_wage" validate:"omitempty"`
	PositionalAllowance      float64 `json:"positional_allowance" validate:"omitempty"`
	OperationalAllowance     float64 `json:"operational_allowance" validate:"omitempty"`
	MealAllowance            float64 `json:"meal_allowance" validate:"omitempty"`
	JobLocation              string  `json:"job_location" validate:"omitempty"`
	HometripTicket           string  `json:"hometrip_ticket" validate:"omitempty"`
	PeriodAgreement          string  `json:"period_agreement" validate:"omitempty"`
	HomeLocation             string  `json:"home_location" validate:"omitempty"`
	JobLevelID               string  `json:"job_level_id" validate:"omitempty"`
	JobID                    string  `json:"job_id" validate:"omitempty"`
	JobPostingID             string  `json:"job_posting_id" validate:"required,uuid"`
	JoinedDate               string  `json:"joined_date" validate:"omitempty,datetime=2006-01-02"`
	ForOrganizationID        string  `json:"for_organization_id" validate:"omitempty"`
	OrganizationLocationID   string  `json:"organization_location_id" validate:"omitempty"`
	DetailContent            string  `json:"detail_content" validate:"required"`
	RecruitmentType          string  `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	SyncMidsuit              string  `json:"sync_midsuit" validate:"omitempty,oneof=YES NO"`
	HiredStatus              string  `json:"hired_status" validate:"omitempty,document_sending_hired_status_validation"`
}

type GeneratePdfBufferFromHTMLRequest struct {
	HTML string `json:"html" validate:"required"`
}
