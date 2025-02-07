package request

type FillInterviewResultRequest struct {
	InterviewApplicantID string `json:"interview_applicant_id" validate:"required,uuid"`
	InterviewAssessorID  string `json:"interview_assessor_id" validate:"required,uuid"`
	Status               string `json:"status" validate:"required,oneof=ACCEPTED REJECTED"`
}
