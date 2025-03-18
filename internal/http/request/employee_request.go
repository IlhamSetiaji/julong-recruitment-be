package request

type SendFindEmployeeByIDMessageRequest struct {
	ID string `json:"id" binding:"required"`
}

type SendCreateEmployeeMessageRequest struct {
	UserID                  string `json:"user_id" validate:"required,uuid"`
	Name                    string `json:"name" validate:"required"`
	Email                   string `json:"email" validate:"required,email"`
	JobID                   string `json:"job_id" validate:"required,uuid"`
	OrganizationID          string `json:"organization_id" validate:"required,uuid"`
	OrganizationLocationID  string `json:"organization_location_id" validate:"required,uuid"`
	OrganizationStructureID string `json:"organization_structure_id" validate:"required,uuid"`
}

type SendCreateEmployeeTaskMessageRequest struct {
	EmployeeID       string `json:"employee_id" validate:"required,uuid"`
	JoinedDate       string `json:"joined_date" validate:"required,datetime=2006-01-02"`
	OrganizationType string `json:"organization_type" validate:"required"`
}

type AdOrgId struct {
	ID         string `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcGender struct {
	ID string `json:"id" binding:"required"`
}

type HcMaritalStatus struct {
	ID         string `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcReligionId struct {
	ID         string `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcStatus struct {
	ID         string `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcBasicAcceptance struct {
	ID         string `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcRecruitmentTypeId struct {
	ID         string `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type SyncEmployeeMidsuitRequest struct {
	AdOrgId           AdOrgId             `json:"AD_Org_ID" binding:"required"`
	Name              string              `json:"Name" binding:"required"`
	Birthday          string              `json:"Birthday" binding:"required"`
	City              string              `json:"City" binding:"required"`
	Email             string              `json:"EMail" binding:"required"`
	HcGender          HcGender            `json:"HC_Gender" binding:"required"`
	HcMaritalStatus   HcMaritalStatus     `json:"HC_MaritalStatus" binding:"required"`
	HcNationalID1     string              `json:"HC_NationalID1" binding:"required"`
	HcReligionID      HcReligionId        `json:"HC_Religion_ID" binding:"required"`
	HcStatus          HcStatus            `json:"HC_Status" binding:"required"`
	HcBasicAcceptance HcBasicAcceptance   `json:"HC_BasicAcceptance" binding:"required"`
	HcRecruitmentType HcRecruitmentTypeId `json:"HC_RecruitmentType" binding:"required"`
	HCWorkStartDate   string              `json:"HC_WorkStartDate" binding:"required"`
}
