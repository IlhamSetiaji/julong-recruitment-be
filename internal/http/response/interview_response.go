package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type InterviewResponse struct {
	ID                         uuid.UUID              `json:"id"`
	JobPostingID               uuid.UUID              `json:"job_posting_id"`
	ProjectPicID               uuid.UUID              `json:"project_pic_id"`
	ProjectRecruitmentHeaderID uuid.UUID              `json:"project_recruitment_header_id"`
	ProjectRecruitmentLineID   uuid.UUID              `json:"project_recruitment_line_id"`
	Name                       string                 `json:"name"`
	DocumentNumber             string                 `json:"document_number"`
	ScheduleDate               time.Time              `json:"schedule_date"`
	StartTime                  time.Time              `json:"start_time"`
	EndTime                    time.Time              `json:"end_time"`
	LocationLink               string                 `json:"location_link"`
	Description                string                 `json:"description"`
	RangeDuration              *int                   `json:"range_duration"`
	TotalCandidate             int                    `json:"total_candidate"`
	Status                     entity.InterviewStatus `json:"status"`
	CreatedAt                  time.Time              `json:"created_at"`
	UpdatedAt                  time.Time              `json:"updated_at"`

	JobPosting               *JobPostingResponse               `json:"job_posting"`
	ProjectPic               *ProjectPicResponse               `json:"project_pic"`
	ProjectRecruitmentHeader *ProjectRecruitmentHeaderResponse `json:"project_recruitment_header"`
	ProjectRecruitmentLine   *ProjectRecruitmentLineResponse   `json:"project_recruitment_line"`
	InterviewApplicants      []InterviewApplicantResponse      `json:"interview_applicants"`
	InterviewAssessors       []InterviewAssessorResponse       `json:"interview_assessors"`
}

type InterviewMyselfResponse struct {
	ID                         uuid.UUID              `json:"id"`
	JobPostingID               uuid.UUID              `json:"job_posting_id"`
	ProjectPicID               uuid.UUID              `json:"project_pic_id"`
	ProjectRecruitmentHeaderID uuid.UUID              `json:"project_recruitment_header_id"`
	ProjectRecruitmentLineID   uuid.UUID              `json:"project_recruitment_line_id"`
	Name                       string                 `json:"name"`
	DocumentNumber             string                 `json:"document_number"`
	ScheduleDate               time.Time              `json:"schedule_date"`
	StartTime                  time.Time              `json:"start_time"`
	EndTime                    time.Time              `json:"end_time"`
	LocationLink               string                 `json:"location_link"`
	Description                string                 `json:"description"`
	RangeDuration              *int                   `json:"range_duration"`
	TotalCandidate             int                    `json:"total_candidate"`
	Status                     entity.InterviewStatus `json:"status"`
	CreatedAt                  time.Time              `json:"created_at"`
	UpdatedAt                  time.Time              `json:"updated_at"`

	JobPosting               *JobPostingResponse               `json:"job_posting"`
	ProjectPic               *ProjectPicResponse               `json:"project_pic"`
	ProjectRecruitmentHeader *ProjectRecruitmentHeaderResponse `json:"project_recruitment_header"`
	ProjectRecruitmentLine   *ProjectRecruitmentLineResponse   `json:"project_recruitment_line"`
	InterviewApplicant       *InterviewApplicantResponse       `json:"interview_applicant"`
	InterviewAssessors       []InterviewAssessorResponse       `json:"interview_assessors"`
}

type InterviewMyselfForAssessorResponse struct {
	ID                         uuid.UUID              `json:"id"`
	JobPostingID               uuid.UUID              `json:"job_posting_id"`
	ProjectPicID               uuid.UUID              `json:"project_pic_id"`
	ProjectRecruitmentHeaderID uuid.UUID              `json:"project_recruitment_header_id"`
	ProjectRecruitmentLineID   uuid.UUID              `json:"project_recruitment_line_id"`
	Name                       string                 `json:"name"`
	DocumentNumber             string                 `json:"document_number"`
	ScheduleDate               time.Time              `json:"schedule_date"`
	StartTime                  time.Time              `json:"start_time"`
	EndTime                    time.Time              `json:"end_time"`
	LocationLink               string                 `json:"location_link"`
	Description                string                 `json:"description"`
	RangeDuration              *int                   `json:"range_duration"`
	TotalCandidate             int                    `json:"total_candidate"`
	Status                     entity.InterviewStatus `json:"status"`
	CreatedAt                  time.Time              `json:"created_at"`
	UpdatedAt                  time.Time              `json:"updated_at"`

	JobPosting               *JobPostingResponse               `json:"job_posting"`
	ProjectPic               *ProjectPicResponse               `json:"project_pic"`
	ProjectRecruitmentHeader *ProjectRecruitmentHeaderResponse `json:"project_recruitment_header"`
	ProjectRecruitmentLine   *ProjectRecruitmentLineResponse   `json:"project_recruitment_line"`
	InterviewApplicants      []InterviewApplicantResponse      `json:"interview_applicant"`
	InterviewAssessor        *InterviewAssessorResponse        `json:"interview_assessor"`
}
