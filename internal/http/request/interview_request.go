package request

type CreateInterviewRequest struct {
	JobPostingID               string `json:"job_posting_id" validate:"required,uuid"`
	ProjectPicID               string `json:"project_pic_id" validate:"required,uuid"`
	ProjectRecruitmentHeaderID string `json:"project_recruitment_header_id" validate:"required,uuid"`
	ProjectRecruitmentLineID   string `json:"project_recruitment_line_id" validate:"required,uuid"`
	Name                       string `json:"name" validate:"required,max=255"`
	DocumentNumber             string `json:"document_number" validate:"required,max=255"`
	ScheduleDate               string `json:"schedule_date" validate:"required,datetime=2006-01-02"`
	StartTime                  string `json:"start_time" validate:"required,datetime=2006-01-02 15:04:05"`
	EndTime                    string `json:"end_time" validate:"required,datetime=2006-01-02 15:04:05"`
	LocationLink               string `json:"location_link" validate:"omitempty"`
	Description                string `json:"description" validate:"omitempty"`
	RangeDuration              *int   `json:"range_duration" validate:"omitempty"`
	TotalCandidate             int    `json:"total_candidate" validate:"required"`
	MeetingLink                string `json:"meeting_link" validate:"omitempty"`
	Status                     string `json:"status" validate:"required,interview_status_validation"`
	InterviewAssessors         []struct {
		EmployeeID string `json:"employee_id" validate:"required,uuid"`
	} `json:"interview_assessors" validate:"required,dive"`
}

type UpdateInterviewRequest struct {
	ID                         string `json:"id" validate:"required,uuid"`
	JobPostingID               string `json:"job_posting_id" validate:"required,uuid"`
	ProjectPicID               string `json:"project_pic_id" validate:"required,uuid"`
	ProjectRecruitmentHeaderID string `json:"project_recruitment_header_id" validate:"required,uuid"`
	ProjectRecruitmentLineID   string `json:"project_recruitment_line_id" validate:"required,uuid"`
	Name                       string `json:"name" validate:"required,max=255"`
	DocumentNumber             string `json:"document_number" validate:"required,max=255"`
	ScheduleDate               string `json:"schedule_date" validate:"required,datetime=2006-01-02"`
	StartTime                  string `json:"start_time" validate:"required,datetime=2006-01-02 15:04:05"`
	EndTime                    string `json:"end_time" validate:"required,datetime=2006-01-02 15:04:05"`
	LocationLink               string `json:"location_link" validate:"omitempty"`
	Description                string `json:"description" validate:"omitempty"`
	RangeDuration              *int   `json:"range_duration" validate:"omitempty"`
	TotalCandidate             int    `json:"total_candidate" validate:"required"`
	MeetingLink                string `json:"meeting_link" validate:"omitempty"`
	Status                     string `json:"status" validate:"required,interview_status_validation"`
	InterviewAssessors         []struct {
		EmployeeID string `json:"employee_id" validate:"required,uuid"`
	} `json:"interview_assessors" validate:"required,dive"`
}

type UpdateStatusInterviewRequest struct {
	ID     string `json:"id" validate:"required,uuid"`
	Status string `json:"status" validate:"required,interview_status_validation"`
}
