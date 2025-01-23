package request

type CreateTestScheduleHeaderRequest struct {
	JobPostingID           string `json:"job_posting_id" validate:"required,uuid"`
	TestTypeID             string `json:"test_type_id" validate:"required,uuid"`
	ProjectPicID           string `json:"project_pic_id" validate:"required,uuid"`
	TemplateActivityLineID string `json:"template_activity_line_id" validate:"required,uuid"`
	JobID                  string `json:"job_id" validate:"required,uuid"`
	Name                   string `json:"name" validate:"required,max=255"`
	DocumentNumber         string `json:"document_number" validate:"required,max=255"`
	StartDate              string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate                string `json:"end_date" validate:"required,datetime=2006-01-02"`
	StartTime              string `json:"start_time" validate:"required,datetime=2006-01-02 15:04:05"`
	EndTime                string `json:"end_time" validate:"required,datetime=2006-01-02 15:04:05"`
	Link                   string `json:"link" validate:"required"`
	Location               string `json:"location" validate:"required"`
	Description            string `json:"description" validate:"required"`
	TotalCandidate         int    `json:"total_candidate" validate:"required"`
	Status                 string `json:"status" validate:"required,test_schedule_header_status_validation"`
}

type UpdateTestScheduleHeaderRequest struct {
	ID                     string `json:"id" validate:"required,uuid"`
	JobPostingID           string `json:"job_posting_id" validate:"required,uuid"`
	TestTypeID             string `json:"test_type_id" validate:"required,uuid"`
	ProjectPicID           string `json:"project_pic_id" validate:"required,uuid"`
	TemplateActivityLineID string `json:"template_activity_line_id" validate:"required,uuid"`
	JobID                  string `json:"job_id" validate:"required,uuid"`
	Name                   string `json:"name" validate:"required,max=255"`
	DocumentNumber         string `json:"document_number" validate:"required,max=255"`
	StartDate              string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate                string `json:"end_date" validate:"required,datetime=2006-01-02"`
	StartTime              string `json:"start_time" validate:"required,datetime=2006-01-02 15:04:05"`
	EndTime                string `json:"end_time" validate:"required,datetime=2006-01-02 15:04:05"`
	Link                   string `json:"link" validate:"required"`
	Location               string `json:"location" validate:"required"`
	Description            string `json:"description" validate:"required"`
	TotalCandidate         int    `json:"total_candidate" validate:"required"`
	Status                 string `json:"status" validate:"required,test_schedule_header_status_validation"`
}
