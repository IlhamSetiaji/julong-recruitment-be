package request

type SendFindEmployeeByIDMessageRequest struct {
	ID string `json:"id" binding:"required"`
}

type SendCreateEmployeeMessageRequest struct {
	UserID                  string `json:"user_id" validate:"required,uuid"`
	Name                    string `json:"name" validate:"required"`
	Email                   string `json:"email" validate:"required,email"`
	JobID                   string `json:"job_id" validate:"required,uuid"`
	JobLevelID              string `json:"job_level_id" validate:"required,uuid"`
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
	ID         int    `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcGender struct {
	ID string `json:"id" binding:"required"`
}

type HcMaritalStatus struct {
	ID         int    `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcReligionId struct {
	ID         int    `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"omitempty"`
}

type HcStatus struct {
	ID         string `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcBasicAcceptance struct {
	ID         int    `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcRecruitmentTypeId struct {
	ID         int    `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcEmployeeId struct {
	ID         int    `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"omitempty"`
}

type HcEmployeeCategoryId struct {
	ID         int    `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcEmployeeGradeId struct {
	ID         int    `json:"id" binding:"omitempty"`
	Identifier string `json:"identifier" binding:"required"`
}

type HcJobId struct {
	ID         int    `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"midsuit"`
}

type HcJobLevelId struct {
	ID         int    `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"omitempty"`
}

type HcOrgId struct {
	ID         int    `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"omitempty"`
}

type HcWorkSiteId struct {
	ID         int    `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"omitempty"`
}

type CDocTypeID struct {
	ID int `json:"id" binding:"required"`
}

type CPeriodID struct {
	Identifier string `json:"identifier" binding:"required"`
	ModelName  string `json:"model-name" binding:"required"`
}

type HCAllowanceType struct {
	ID         string `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"omitempty"`
	ModelName  string `json:"model-name" binding:"omitempty"`
}

type JobLevelCategory struct {
	ID         int    `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"required"`
	ModelName  string `json:"model-name" binding:"required"`
}

type HCUOM struct {
	ID         string `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"omitempty"`
	ModelName  string `json:"model-name" binding:"omitempty"`
}

type HCProvisionType struct {
	ID         string `json:"id" binding:"required"`
	Identifier string `json:"identifier" binding:"omitempty"`
	ModelName  string `json:"model-name" binding:"omitempty"`
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
	HcRecruitmentType HcRecruitmentTypeId `json:"HC_RecruitmentType_ID" binding:"required"`
	HCWorkStartDate   string              `json:"HC_WorkStartDate" binding:"required"`
}

type SyncEmployeeJobMidsuitRequest struct {
	AdOrgId              AdOrgId               `json:"AD_Org_ID" binding:"required"`
	HCCompensation1      int                   `json:"HC_Compensation1" binding:"required"`
	HCEmployeeID         HcEmployeeId          `json:"HC_Employee_ID" binding:"required"`
	HCEmployeeCategoryID *HcEmployeeCategoryId `json:"HC_EmployeeCategory_ID" binding:"required"`
	HCEmployeeGradeID    HcEmployeeGradeId     `json:"HC_EmployeeGrade_ID" binding:"required"`
	HCJobID              HcJobId               `json:"HC_Job_ID" binding:"required"`
	HCJobLevelID         HcJobLevelId          `json:"HC_JobLevel_ID" binding:"required"`
	HCOrgID              HcOrgId               `json:"HC_Org_ID" binding:"required"`
	HCWorkStartDate      string                `json:"HC_WorkStartDate" binding:"required"`
	HCRecruitmentTypeID  HcRecruitmentTypeId   `json:"HC_RecruitmentType_ID" binding:"required"`
	ADEmploymentOrgID    AdOrgId               `json:"AD_EmploymentOrg_ID" binding:"required"`
	HCWorkSiteID         HcWorkSiteId          `json:"HC_WorkSite_ID" binding:"required"`
	IsPrimary            bool                  `json:"IsPrimary" binding:"required"`
	ModelName            string                `json:"model-name" binding:"required"`
}

type SyncEmployeeWorkExperienceMidsuitRequest struct {
	AdOrgId        AdOrgId      `json:"AD_Org_ID" binding:"required"`
	Name           string       `json:"Name" binding:"required"`
	Description    string       `json:"Description" binding:"required"`
	HCEmployeeID   HcEmployeeId `json:"HC_Employee_ID" binding:"required"`
	YearExperience string       `json:"YearExperience" binding:"required"`
	ModelName      string       `json:"model-name" binding:"required"`
}

type SyncEmployeeEducationMidsuitRequest struct {
	AdOrgId               AdOrgId           `json:"AD_Org_ID" binding:"required"`
	BidangPendidikanAkhir string            `json:"BidangPendidikanAkhir" binding:"required"`
	HcEducationInstitute  string            `json:"hc_educationinstitute" binding:"required"`
	HCEmployeeID          HcEmployeeId      `json:"HC_Employee_ID" binding:"required"`
	HcGpaScore            int               `json:"hc_gpascore" binding:"required"`
	SeqNo                 int               `json:"SeqNo" binding:"required"`
	HCBasicAcceptance     HcBasicAcceptance `json:"HC_BasicAcceptance" binding:"required"`
	ModelName             string            `json:"model-name" binding:"required"`
}

type SyncEmployeeAllowanceMidsuitRequest struct {
	AdOrgId          AdOrgId          `json:"AD_Org_ID" binding:"required"`
	CDocTypeID       CDocTypeID       `json:"C_DocType_ID" binding:"omitempty"`
	DateDoc          string           `json:"DateDoc" binding:"required"`
	HCEmployeeID     HcEmployeeId     `json:"HC_Employee_ID" binding:"required"`
	HCNIK            string           `json:"HC_NIK" binding:"required"`
	HCJobID          HcJobId          `json:"HC_Job_ID" binding:"required"`
	HCOrgID          HcOrgId          `json:"HC_Org_ID" binding:"required"`
	CPeriodID        CPeriodID        `json:"C_Period_ID" binding:"omitempty"`
	HCAllowanceType  HCAllowanceType  `json:"HC_AllowanceType" binding:"required"`
	HCEmployee2ID    HcEmployeeId     `json:"HC_Employee2_ID" binding:"omitempty"`
	HCNIK2           string           `json:"HC_NIK2" binding:"omitempty"`
	HCJob2ID         HcJobId          `json:"HC_Job2_ID" binding:"omitempty"`
	HCOrg2ID         HcOrgId          `json:"HC_Org2_ID" binding:"omitempty"`
	HCJobLevel2ID    HcJobLevelId     `json:"HC_JobLevel2_ID" binding:"omitempty"`
	JobLevelCategory JobLevelCategory `json:"JobLevelCategory" binding:"omitempty"`
	Distance         int              `json:"Distance" binding:"required"`
	Amount           int              `json:"Amount" binding:"required"`
	HCUOM            HCUOM            `json:"HC_UOM" binding:"required"`
	CPeriod2ID       CPeriodID        `json:"C_Period2_ID" binding:"required"`
	IsUseDate        bool             `json:"IsUseDate" binding:"required"`
	HCProvisionType  HCProvisionType  `json:"HC_ProvisionType" binding:"required"`
	IsGenerated      bool             `json:"IsGenerated" binding:"required"`
	ModelName        string           `json:"model-name" binding:"required"`
}
