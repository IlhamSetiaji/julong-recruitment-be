package response

import (
	"time"

	"github.com/google/uuid"
)

type ProjectPicResponse struct {
	ID                       uuid.UUID  `json:"id"`
	ProjectRecruitmentLineID uuid.UUID  `json:"project_recruitment_line_id"`
	EmployeeID               *uuid.UUID `json:"employee_id"`
	EmployeeName             string     `json:"employee_name"`
	AdministrativeTotal      int        `json:"administrative_total"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}
