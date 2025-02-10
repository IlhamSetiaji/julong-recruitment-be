package request

type FillFgdResultRequest struct {
	FgdApplicantID string `json:"fgd_applicant_id" validate:"required,uuid"`
	FgdAssessorID  string `json:"fgd_assessor_id" validate:"required,uuid"`
	Status         string `json:"status" validate:"required,oneof=ACCEPTED REJECTED"`
}
