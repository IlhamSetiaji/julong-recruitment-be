package request

type CreateOrUpdateTestApplicantsRequest struct {
	TestScheduleHeaderID string `json:"test_schedule_header_id" validate:"required,uuid"`
	TestApplicants       []struct {
		ID            string `json:"id" validate:"omitempty,uuid"`
		ApplicantID   string `json:"applicant_id" validate:"required,uuid"`
		UserProfileID string `json:"user_profile_id" validate:"required,uuid"`
		StartTime     string `json:"start_time" validate:"required,datetime=2006-01-02 15:04:05"`
		EndTime       string `json:"end_time" validate:"required,datetime=2006-01-02 15:04:05"`
		FinalResult   string `json:"final_result" validate:"required,test_applicant_status_validation"`
	} `json:"test_applicants" validate:"required,dive"`
	DeletedTestApplicantIDs []string `json:"deleted_test_applicant_ids" validate:"omitempty,dive,uuid"`
}
