package request

type CreateAdministrativeSelectionRequest struct {
	JobPostingID   string `json:"job_posting_id" validate:"required,uuid"`
	ProjectPicID   string `json:"project_pic_id" validate:"required,uuid"`
	Status         string `json:"status" validate:"required,administrative_selection_status_validation"`
	DocumentDate   string `json:"document_date" validate:"required,datetime=2006-01-02"`
	DocumentNumber string `json:"document_number" validate:"required"`
}

type UpdateAdministrativeSelectionRequest struct {
	ID             string `json:"id" validate:"required,uuid"`
	JobPostingID   string `json:"job_posting_id" validate:"required,uuid"`
	ProjectPicID   string `json:"project_pic_id" validate:"required,uuid"`
	Status         string `json:"status" validate:"required,administrative_selection_status_validation"`
	DocumentDate   string `json:"document_date" validate:"required,datetime=2006-01-02"`
	DocumentNumber string `json:"document_number" validate:"required"`
}
