package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type TestScheduleHeaderResponse struct {
	ID                         uuid.UUID                 `json:"id"`
	JobPostingID               uuid.UUID                 `json:"job_posting_id"`
	TestTypeID                 uuid.UUID                 `json:"test_type_id"`
	ProjectPicID               uuid.UUID                 `json:"project_pic_id"`
	ProjectRecruitmentHeaderID uuid.UUID                 `json:"project_recruitment_header_id"`
	ProjectRecruitmentLineID   uuid.UUID                 `json:"project_recruitment_line_id"`
	JobID                      *uuid.UUID                `json:"job_id"`
	Name                       string                    `json:"name"`
	DocumentNumber             string                    `json:"document_number"`
	StartDate                  time.Time                 `json:"start_date"`
	EndDate                    time.Time                 `json:"end_date"`
	StartTime                  time.Time                 `json:"start_time"`
	EndTime                    time.Time                 `json:"end_time"`
	Link                       string                    `json:"link"`
	Location                   string                    `json:"location"`
	Description                string                    `json:"description"`
	TotalCandidate             int                       `json:"total_candidate"`
	Status                     entity.TestScheduleStatus `json:"status"`
	ScheduleDate               time.Time                 `json:"schedule_date"`
	Platform                   string                    `json:"platform"`
	CreatedAt                  time.Time                 `json:"created_at"`
	UpdatedAt                  time.Time                 `json:"updated_at"`

	JobPosting               *JobPostingResponse               `json:"job_posting"`
	TestType                 *TestTypeResponse                 `json:"test_type"`
	ProjectPic               *ProjectPicResponse               `json:"project_pic"`
	TestApplicants           []TestApplicantResponse           `json:"test_applicants"`
	ProjectRecruitmentHeader *ProjectRecruitmentHeaderResponse `json:"project_recruitment_header"`
	ProjectRecruitmentLine   *ProjectRecruitmentLineResponse   `json:"project_recruitment_line"`
}

type TestScheduleHeaderMyselfResponse struct {
	ID             uuid.UUID                 `json:"id"`
	JobPostingID   uuid.UUID                 `json:"job_posting_id"`
	TestTypeID     uuid.UUID                 `json:"test_type_id"`
	ProjectPicID   uuid.UUID                 `json:"project_pic_id"`
	JobID          *uuid.UUID                `json:"job_id"`
	Name           string                    `json:"name"`
	DocumentNumber string                    `json:"document_number"`
	StartDate      time.Time                 `json:"start_date"`
	EndDate        time.Time                 `json:"end_date"`
	StartTime      time.Time                 `json:"start_time"`
	EndTime        time.Time                 `json:"end_time"`
	Link           string                    `json:"link"`
	Location       string                    `json:"location"`
	Description    string                    `json:"description"`
	TotalCandidate int                       `json:"total_candidate"`
	Status         entity.TestScheduleStatus `json:"status"`
	ScheduleDate   time.Time                 `json:"schedule_date"`
	Platform       string                    `json:"platform"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`

	JobPosting               *JobPostingResponse               `json:"job_posting"`
	TestType                 *TestTypeResponse                 `json:"test_type"`
	ProjectPic               *ProjectPicResponse               `json:"project_pic"`
	TestApplicant            *TestApplicantResponse            `json:"test_applicants"`
	ProjectRecruitmentHeader *ProjectRecruitmentHeaderResponse `json:"project_recruitment_header"`
	ProjectRecruitmentLine   *ProjectRecruitmentLineResponse   `json:"project_recruitment_line"`
}

type TestApplicantsPayload struct {
	ApplicantIDs   []uuid.UUID `json:"applicant_ids"`
	UserProfileIDs []uuid.UUID `json:"user_profile_ids"`
	Total          int         `json:"total"`
}
