package request

type CreateOrUpdateAdministrativeResults struct {
	AdministrativeSelectionID string `json:"administrative_selection_id" validate:"required,uuid"`
	AdministrativeResults     []struct {
		ID            string `json:"id" validate:"omitempty,uuid"`
		UserProfileID string `json:"user_profile_id" validate:"required,uuid"`
		Status        string `json:"status" validate:"required,administrative_result_status_validation"`
	} `json:"administrative_results" validate:"required,dive"`
	DeletedAdministrativeResultIDs []string `json:"deleted_administrative_result_ids" validate:"omitempty,dive,uuid"`
}
