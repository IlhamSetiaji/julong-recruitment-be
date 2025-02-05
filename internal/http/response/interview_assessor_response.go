package response

import (
	"time"

	"github.com/google/uuid"
)

type InterviewAssessorResponse struct {
	ID           uuid.UUID  `json:"id"`
	InterviewID  uuid.UUID  `json:"interview_id"`
	EmployeeID   *uuid.UUID `json:"employee_id"`
	EmployeeName string     `json:"employee_name"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
