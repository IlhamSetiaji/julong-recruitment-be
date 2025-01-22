package request

type CreateTestTypeRequest struct {
	Name            string `json:"name" validate:"required"`
	RecruitmentType string `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	Status          string `json:"status" validate:"required,test_type_status_validation"`
}

type UpdateTestTypeRequest struct {
	ID              string `json:"id" validate:"required,uuid"`
	Name            string `json:"name" validate:"required"`
	RecruitmentType string `json:"recruitment_type" validate:"required,recruitment_type_validation"`
	Status          string `json:"status" validate:"required,test_type_status_validation"`
}
