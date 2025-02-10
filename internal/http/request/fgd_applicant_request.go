package request

type CreateOrUpdateFgdApplicantsRequest struct {
	FgdID         string `json:"fgd_id" validate:"required,uuid"`
	FgdApplicants []struct {
		ID            string `json:"id" validate:"omitempty,uuid"`
		ApplicantID   string `json:"applicant_id" validate:"required,uuid"`
		UserProfileID string `json:"user_profile_id" validate:"required,uuid"`
		StartTime     string `json:"start_time" validate:"required,datetime=2006-01-02 15:04:05"`
		EndTime       string `json:"end_time" validate:"required,datetime=2006-01-02 15:04:05"`
		FinalResult   string `json:"final_result" validate:"required,final_result_status_validation"`
	} `json:"fgd_applicants" validate:"required,dive"`
	DeletedFgdApplicantIDs []string `json:"deleted_fgd_applicant_ids" validate:"omitempty,dive,uuid"`
}

type UpdateStatusFgdApplicantRequest struct {
	ID     string `json:"id" validate:"required,uuid"`
	Status string `json:"status" validate:"required,assessment_status_validation"`
}

type UpdateFinalResultFgdApplicantRequest struct {
	ID          string `json:"id" validate:"required,uuid"`
	FinalResult string `json:"final_result" validate:"required,final_result_status_validation"`
}
