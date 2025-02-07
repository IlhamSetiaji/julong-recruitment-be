package request

type CreateOrUpdateInterviewApplicantsRequest struct {
	InterviewID         string `json:"interview_id" validate:"required,uuid"`
	InterviewApplicants []struct {
		ID            string `json:"id" validate:"omitempty,uuid"`
		ApplicantID   string `json:"applicant_id" validate:"required,uuid"`
		UserProfileID string `json:"user_profile_id" validate:"required,uuid"`
		StartTime     string `json:"start_time" validate:"required,datetime=2006-01-02 15:04:05"`
		EndTime       string `json:"end_time" validate:"required,datetime=2006-01-02 15:04:05"`
		FinalResult   string `json:"final_result" validate:"required,final_result_status_validation"`
	} `json:"interview_applicants" validate:"required,dive"`
	DeletedInterviewApplicantIDs []string `json:"deleted_interview_applicant_ids" validate:"omitempty,dive,uuid"`
}

type UpdateStatusInterviewApplicantRequest struct {
	ID     string `json:"id" validate:"required,uuid"`
	Status string `json:"status" validate:"required,assessment_status_validation"`
}

type UpdateFinalResultInterviewApplicantRequest struct {
	ID          string `json:"id" validate:"required,uuid"`
	FinalResult string `json:"final_result" validate:"required,final_result_status_validation"`
}
