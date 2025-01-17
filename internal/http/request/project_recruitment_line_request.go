package request

type CreateOrUpdateProjectRecruitmentLinesRequest struct {
	ProjectRecruitmentHeaderID string `json:"project_recruitment_header_id" validate:"required,uuid"`
	ProjectRecruitmentLines    []struct {
		ID                     string `json:"id" validate:"omitempty,uuid"`
		TemplateActivityLineID string `json:"template_activity_line_id" validate:"required,uuid"`
		StartDate              string `json:"start_date" validate:"required,datetime=2006-01-02"`
		EndDate                string `json:"end_date" validate:"required,datetime=2006-01-02"`
		ProjectPics            []struct {
			EmployeeID          string `json:"employee_id" validate:"required,uuid"`
			AdministrativeTotal int    `json:"administrative_total" validate:"required"`
		} `json:"project_pics" validate:"required,dive"`
	} `json:"project_recruitment_lines" validate:"required,dive"`
	DeletedProjectRecruitmentLineIDs []string `json:"deleted_project_recruitment_line_ids" validate:"omitempty,dive,uuid"`
}
