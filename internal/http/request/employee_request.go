package request

type SendFindEmployeeByIDMessageRequest struct {
	ID string `json:"id" binding:"required"`
}

type SendCreateEmployeeMessageRequest struct {
	UserID                 string `json:"user_id" validate:"required,uuid"`
	Name                   string `json:"name" validate:"required"`
	Email                  string `json:"email" validate:"required,email"`
	JobID                  string `json:"job_id" validate:"required,uuid"`
	OrganizationID         string `json:"organization_id" validate:"required,uuid"`
	OrganizationLocationID string `json:"organization_location_id" validate:"required,uuid"`
}
