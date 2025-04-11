package helper

import (
	"errors"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IMPRequestHelper interface {
	ConvertMapInterfaceToResponse(mprMap map[string]interface{}) (*response.MPRequestHeaderResponse, error)
	ConvertMapInterfaceToResponseMinimal(mprMap map[string]interface{}) (*response.MPRequestHeaderResponse, error)
}

type MPRequestHelper struct {
	Log *logrus.Logger
}

func NewMPRequestHelper(
	log *logrus.Logger,
) IMPRequestHelper {
	return &MPRequestHelper{
		Log: log,
	}
}

func MPRequestHelperFactory(
	log *logrus.Logger,
) IMPRequestHelper {
	return NewMPRequestHelper(log)
}

func (d *MPRequestHelper) ConvertMapInterfaceToResponse(mprMap map[string]interface{}) (*response.MPRequestHeaderResponse, error) {
	mprData, ok := mprMap["mp_request_header"].(map[string]interface{})
	if !ok {
		d.Log.Errorf("MPRequestHeader information is missing or invalid")
		return nil, errors.New("MPRequestHeader information is missing or invalid")
	}

	mprID, ok := mprData["id"].(string)
	if !ok {
		d.Log.Errorf("MPRequestHeader ID is missing or invalid")
		return nil, errors.New("MPRequestHeader ID is missing or invalid")
	}

	organizationID, ok := mprData["organization_id"].(string)
	if !ok {
		d.Log.Errorf("Organization ID is missing or invalid")
		return nil, errors.New("Organization ID is missing or invalid")
	}

	// gradeID, ok := mprData["grade_id"].(string)
	// if !ok {
	// 	d.Log.Errorf("Grade ID is missing or invalid")
	// 	return nil, errors.New("Grade ID is missing or invalid")
	// }

	// var parsedGradeID *uuid.UUID
	// if gradeID != "" {
	// 	gradeUUID, err := uuid.Parse(gradeID)
	// 	if err != nil {
	// 		d.Log.Errorf("Invalid grade ID format: %v", err)
	// 		parsedGradeID = nil
	// 	} else {
	// 		parsedGradeID = &gradeUUID
	// 	}
	// } else {
	// 	parsedGradeID = nil
	// }

	organizationLocationID, ok := mprData["organization_location_id"].(string)
	if !ok {
		d.Log.Errorf("Organization Location ID is missing or invalid")
		return nil, errors.New("Organization Location ID is missing or invalid")
	}
	forOrganizationID, ok := mprData["for_organization_id"].(string)
	if !ok {
		d.Log.Errorf("For Organization ID is missing or invalid")
		return nil, errors.New("For Organization ID is missing or invalid")
	}

	forOrganizationLocationID, ok := mprData["for_organization_location_id"].(string)
	if !ok {
		d.Log.Errorf("For Organization Location ID is missing or invalid")
		return nil, errors.New("For Organization Location ID is missing or invalid")
	}

	forOrganizationStructureID, ok := mprData["for_organization_structure_id"].(string)
	if !ok {
		d.Log.Errorf("For Organization Structure ID is missing or invalid")
		return nil, errors.New("For Organization Structure ID is missing or invalid")
	}

	jobID, ok := mprData["job_id"].(string)
	if !ok {
		d.Log.Errorf("Job ID is missing or invalid")
		return nil, errors.New("Job ID is missing or invalid")
	}

	requestCategoryID, ok := mprData["request_category_id"].(string)
	if !ok {
		d.Log.Errorf("Request Category ID is missing or invalid")
		return nil, errors.New("Request Category ID is missing or invalid")
	}

	expectedDate, ok := mprData["expected_date"].(string)
	if !ok {
		d.Log.Errorf("Expected Date is missing or invalid")
		return nil, errors.New("Expected Date is missing or invalid")
	}

	experiences, ok := mprData["experiences"].(string)
	if !ok {
		d.Log.Errorf("Experiences is missing or invalid")
		return nil, errors.New("Experiences is missing or invalid")
	}

	documentNumber, ok := mprData["document_number"].(string)
	if !ok {
		d.Log.Errorf("Document Number is missing or invalid")
		return nil, errors.New("Document Number is missing or invalid")
	}

	documentDate, ok := mprData["document_date"].(string)
	if !ok {
		d.Log.Errorf("Document Date is missing or invalid")
		return nil, errors.New("Document Date is missing or invalid")
	}

	maleNeedsFloat, ok := mprData["male_needs"].(float64)
	if !ok {
		d.Log.Errorf("Male needs is missing or invalid")
		return nil, errors.New("Male needs is missing or invalid")
	}
	maleNeeds := int(maleNeedsFloat)

	femaleNeedsFloat, ok := mprData["female_needs"].(float64)
	if !ok {
		d.Log.Errorf("female needs is missing or invalid")
		return nil, errors.New("female needs is missing or invalid")
	}
	femaleNeeds := int(femaleNeedsFloat)

	minimumAgeFloat, ok := mprData["minimum_age"].(float64)
	if !ok {
		d.Log.Errorf("Minimum Age is missing or invalid")
		return nil, errors.New("Minimum Age is missing or invalid")
	}
	minimumAge := int(minimumAgeFloat)

	maximumAgeFloat, ok := mprData["maximum_age"].(float64)
	if !ok {
		d.Log.Errorf("Maximum Age is missing or invalid")
		return nil, errors.New("Maximum Age is missing or invalid")
	}
	maximumAge := int(maximumAgeFloat)

	minimumExperienceFloat, ok := mprData["minimum_experience"].(float64)
	if !ok {
		d.Log.Errorf("Minimum Experience is missing or invalid")
		return nil, errors.New("Minimum Experience is missing or invalid")
	}
	minimumExperienceInt := int(minimumExperienceFloat)

	maritalStatus, ok := mprData["marital_status"].(string)
	if !ok {
		d.Log.Errorf("Marital Status is missing or invalid")
		return nil, errors.New("Marital Status is missing or invalid")
	}

	minimumEducation, ok := mprData["minimum_education"].(string)
	if !ok {
		d.Log.Errorf("Minimum Education is missing or invalid")
		return nil, errors.New("Minimum Education is missing or invalid")
	}

	requiredQualification, ok := mprData["required_qualification"].(string)
	if !ok {
		d.Log.Errorf("Required Qualification is missing or invalid")
		return nil, errors.New("Required Qualification is missing or invalid")
	}

	certificate, ok := mprData["certificate"].(string)
	if !ok {
		d.Log.Errorf("Certificate is missing or invalid")
		return nil, errors.New("Certificate is missing or invalid")
	}

	computerSkill, ok := mprData["computer_skill"].(string)
	if !ok {
		d.Log.Errorf("Computer Skill is missing or invalid")
		return nil, errors.New("Computer Skill is missing or invalid")
	}

	languageSkill, ok := mprData["language_skill"].(string)
	if !ok {
		d.Log.Errorf("Language Skill is missing or invalid")
		return nil, errors.New("Language Skill is missing or invalid")
	}

	otherSkill, ok := mprData["other_skill"].(string)
	if !ok {
		d.Log.Errorf("Other Skill is missing or invalid")
		return nil, errors.New("Other Skill is missing or invalid")
	}

	jobdesc, ok := mprData["jobdesc"].(string)
	if !ok {
		d.Log.Errorf("Jobdesc is missing or invalid")
		return nil, errors.New("Jobdesc is missing or invalid")
	}

	salaryMin, ok := mprData["salary_min"].(string)
	if !ok {
		d.Log.Errorf("Salary Min is missing or invalid")
		return nil, errors.New("Salary Min is missing or invalid")
	}

	salaryMax, ok := mprData["salary_max"].(string)
	if !ok {
		d.Log.Errorf("Salary Max is missing or invalid")
		return nil, errors.New("Salary Max is missing or invalid")
	}

	requestorID, ok := mprData["requestor_id"].(string)
	if !ok {
		d.Log.Errorf("Requestor ID is missing or invalid")
		return nil, errors.New("Requestor ID is missing or invalid")
	}

	var departmentHead *string
	departmentHeadVal, ok := mprData["department_head"].(string)
	if !ok {
		d.Log.Errorf("Department Head is missing or invalid")
		// return nil, errors.New("Department Head is missing or invalid")
	} else {
		departmentHead = &departmentHeadVal
	}

	var vpGmDirector *string
	vpGmDirectorVal, ok := mprData["vp_gm_director"].(string)
	if !ok {
		d.Log.Errorf("VP GM Director is missing or invalid")
		// return nil, errors.New("VP GM Director is missing or invalid")
	} else {
		vpGmDirector = &vpGmDirectorVal
	}

	var ceo *string
	if ceoVal, ok := mprData["ceo"].(string); ok {
		ceo = &ceoVal
	} else {
		d.Log.Warn("CEO is missing or invalid")
	}

	var hrdHoUnit *string
	hrdHoUnitVal, ok := mprData["hrd_ho_unit"].(string)
	if !ok {
		d.Log.Errorf("HRD HO Unit is missing or invalid")
		// return nil, errors.New("HRD HO Unit is missing or invalid")
	} else {
		hrdHoUnit = &hrdHoUnitVal
	}

	var mpPlanningHeaderID *string
	mpPlanningHeaderIDVal, ok := mprData["mp_planning_header_id"].(string)
	if !ok {
		d.Log.Errorf("MP Planning Header ID is missing or invalid")
		// return nil, errors.New("MP Planning Header ID is missing or invalid")
	} else {
		mpPlanningHeaderID = &mpPlanningHeaderIDVal
	}

	status, ok := mprData["status"].(string)
	if !ok {
		d.Log.Errorf("Status is missing or invalid")
		return nil, errors.New("Status is missing or invalid")
	}

	mpRequestType, ok := mprData["mp_request_type"].(string)
	if !ok {
		d.Log.Errorf("MP Request Type is missing or invalid")
		return nil, errors.New("MP Request Type is missing or invalid")
	}

	recruitmentType, ok := mprData["recruitment_type"].(string)
	if !ok {
		d.Log.Errorf("Recruitment Type is missing or invalid")
		return nil, errors.New("Recruitment Type is missing or invalid")
	}

	mppPeriodID, ok := mprData["mpp_period_id"].(string)
	if !ok {
		d.Log.Errorf("MPP Period ID is missing or invalid")
		return nil, errors.New("MPP Period ID is missing or invalid")
	}

	empOrganizationID, ok := mprData["emp_organization_id"].(string)
	if !ok {
		d.Log.Errorf("Emp Organization ID is missing or invalid")
		return nil, errors.New("Emp Organization ID is missing or invalid")
	}

	jobLevelID, ok := mprData["job_level_id"].(string)
	if !ok {
		d.Log.Errorf("Job Level ID is missing or invalid")
		return nil, errors.New("Job Level ID is missing or invalid")
	}

	isReplacement, ok := mprData["is_replacement"].(bool)
	if !ok {
		d.Log.Errorf("Is Replacement is missing or invalid")
		return nil, errors.New("Is Replacement is missing or invalid")
	}

	createdAt, ok := mprData["created_at"].(string)
	if !ok {
		d.Log.Errorf("Created At is missing or invalid")
		return nil, errors.New("Created At is missing or invalid")
	}

	updatedAt, ok := mprData["updated_at"].(string)
	if !ok {
		d.Log.Errorf("Updated At is missing or invalid")
		return nil, errors.New("Updated At is missing or invalid")
	}

	requestCategory, ok := mprData["request_category"].(map[string]interface{})
	if !ok {
		d.Log.Errorf("Request Category is missing or invalid")
		return nil, errors.New("Request Category is missing or invalid")
	}

	requestMajorsInterface, ok := mprData["request_majors"].([]interface{})
	if !ok {
		d.Log.Errorf("Request Majors is missing or invalid")
		return nil, errors.New("Request Majors is missing or invalid")
	}
	var requestMajors []map[string]interface{}
	for _, v := range requestMajorsInterface {
		if m, ok := v.(map[string]interface{}); ok {
			requestMajors = append(requestMajors, m)
		} else {
			d.Log.Errorf("Invalid type in Request Majors")
			return nil, errors.New("Invalid type in Request Majors")
		}
	}

	gradeName, ok := mprData["grade_name"].(string)
	if !ok {
		d.Log.Errorf("Grade Name is missing or invalid")
		return nil, errors.New("Grade Name is missing or invalid")
	}

	organizationName, ok := mprData["organization_name"].(string)
	if !ok {
		d.Log.Errorf("Organization Name is missing or invalid")
		return nil, errors.New("Organization Name is missing or invalid")
	}

	organizationCategory, ok := mprData["organization_category"].(string)
	if !ok {
		d.Log.Errorf("Organization Category Name is missing or invalid")
		return nil, errors.New("Organization Category Name is missing or invalid")
	}

	organizationLocationName, ok := mprData["organization_location_name"].(string)
	if !ok {
		d.Log.Errorf("Organization Location Name is missing or invalid")
		return nil, errors.New("Organization Location Name is missing or invalid")
	}

	forOrganizationName, ok := mprData["for_organization_name"].(string)
	if !ok {
		d.Log.Errorf("For Organization Name is missing or invalid")
		return nil, errors.New("For Organization Name is missing or invalid")
	}

	forOrganizationLocation, ok := mprData["for_organization_location"].(string)
	if !ok {
		d.Log.Errorf("For Organization Location is missing or invalid")
		return nil, errors.New("For Organization Location is missing or invalid")
	}

	forOrganizationStructure, ok := mprData["for_organization_structure"].(string)
	if !ok {
		d.Log.Errorf("For Organization Structure is missing or invalid")
		return nil, errors.New("For Organization Structure is missing or invalid")
	}

	jobName, ok := mprData["job_name"].(string)
	if !ok {
		d.Log.Errorf("Job Name is missing or invalid")
		return nil, errors.New("Job Name is missing or invalid")
	}

	requestorName, ok := mprData["requestor_name"].(string)
	if !ok {
		d.Log.Errorf("Requestor Name is missing or invalid")
		return nil, errors.New("Requestor Name is missing or invalid")
	}

	departmentHeadName, ok := mprData["department_head_name"].(string)
	if !ok {
		d.Log.Errorf("Department Head Name is missing or invalid")
		return nil, errors.New("Department Head Name is missing or invalid")
	}

	hrdHoUnitName, ok := mprData["hrd_ho_unit_name"].(string)
	if !ok {
		d.Log.Errorf("HRD HO Unit Name is missing or invalid")
		return nil, errors.New("HRD HO Unit Name is missing or invalid")
	}

	vpGmDirectorName, ok := mprData["vp_gm_director_name"].(string)
	if !ok {
		d.Log.Errorf("VP GM Director Name is missing or invalid")
		return nil, errors.New("VP GM Director Name is missing or invalid")
	}

	ceoName, ok := mprData["ceo_name"].(string)
	if !ok {
		d.Log.Errorf("CEO Name is missing or invalid")
		return nil, errors.New("CEO Name is missing or invalid")
	}

	empOrganizationName, ok := mprData["emp_organization_name"].(string)
	if !ok {
		d.Log.Errorf("Emp Organization Name is missing or invalid")
		return nil, errors.New("Emp Organization Name is missing or invalid")
	}

	jobLevelName, ok := mprData["job_level_name"].(string)
	if !ok {
		d.Log.Errorf("Job Level Name is missing or invalid")
		return nil, errors.New("Job Level Name is missing or invalid")
	}

	jobLevelFloat, ok := mprData["job_level"].(float64)
	if !ok {
		d.Log.Errorf("Job Level is missing or invalid")
		return nil, errors.New("Job Level is missing or invalid")
	}
	jobLevel := int(jobLevelFloat)

	approvedByDepartmentHead, ok := mprData["approved_by_department_head"].(bool)
	if !ok {
		d.Log.Errorf("Approved By Department Head is missing or invalid")
		return nil, errors.New("Approved By Department Head is missing or invalid")
	}

	approvedByVpGmDirector, ok := mprData["approved_by_vp_gm_director"].(bool)
	if !ok {
		d.Log.Errorf("Approved By VP GM Director is missing or invalid")
		return nil, errors.New("Approved By VP GM Director is missing or invalid")
	}

	approvedByCEO, ok := mprData["approved_by_ceo"].(bool)
	if !ok {
		d.Log.Errorf("Approved By CEO is missing or invalid")
		return nil, errors.New("Approved By CEO is missing or invalid")
	}

	approvedByHrdHoUnit, ok := mprData["approved_by_hrd_ho_unit"].(bool)
	if !ok {
		d.Log.Errorf("Approved By HRD HO Unit is missing or invalid")
		return nil, errors.New("Approved By HRD HO Unit is missing or invalid")
	}

	mprCloneID := uuid.MustParse(mprID)
	parsedRequestorID := uuid.MustParse(requestorID)
	var departmentHeadUUID *uuid.UUID
	if departmentHead != nil {
		uuidVal := uuid.MustParse(*departmentHead)
		departmentHeadUUID = &uuidVal
	} else {
		departmentHeadUUID = nil
	}
	var vpGmDirectorUUID *uuid.UUID
	if vpGmDirector != nil {
		uuidVal := uuid.MustParse(*vpGmDirector)
		vpGmDirectorUUID = &uuidVal
	} else {
		vpGmDirectorUUID = nil
	}
	var ceoUUID *uuid.UUID
	if ceo != nil {
		uuidVal := uuid.MustParse(*ceo)
		ceoUUID = &uuidVal
	} else {
		ceoUUID = nil
	}
	var hrdHoUnitUUID *uuid.UUID
	if hrdHoUnit != nil {
		uuidVal := uuid.MustParse(*hrdHoUnit)
		hrdHoUnitUUID = &uuidVal
	} else {
		hrdHoUnitUUID = nil
	}

	var mpPlanningHeaderUUID *uuid.UUID
	if mpPlanningHeaderID != nil {
		uuidVal := uuid.MustParse(*mpPlanningHeaderID)
		mpPlanningHeaderUUID = &uuidVal
	} else {
		mpPlanningHeaderUUID = nil
	}
	parsedMppPeriodID := uuid.MustParse(mppPeriodID)
	parsedEmpOrganizationID := uuid.MustParse(empOrganizationID)
	parsedJobLevelID := uuid.MustParse(jobLevelID)

	parsedExpectedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", expectedDate)
	if err != nil {
		d.Log.Errorf("Invalid expected date format: %v", err)
		return nil, errors.New("Invalid expected date format")
	}
	parsedDocumentDate, err := time.Parse("2006-01-02T15:04:05Z07:00", documentDate)
	if err != nil {
		d.Log.Errorf("Invalid document date format: %v", err)
		return nil, errors.New("Invalid document date format")
	}
	parsedCreatedAt, err := time.Parse("2006-01-02T15:04:05Z07:00", createdAt)
	if err != nil {
		d.Log.Errorf("Invalid created at format: %v", err)
		return nil, errors.New("Invalid created at format")
	}
	parsedUpdatedAt, err := time.Parse("2006-01-02T15:04:05Z07:00", updatedAt)
	if err != nil {
		d.Log.Errorf("Invalid updated at format: %v", err)
		return nil, errors.New("Invalid updated at format")
	}

	return &response.MPRequestHeaderResponse{
		ID:                         uuid.MustParse(mprID),
		MPRCloneID:                 &mprCloneID,
		OrganizationID:             uuid.MustParse(organizationID),
		OrganizationLocationID:     uuid.MustParse(organizationLocationID),
		ForOrganizationID:          uuid.MustParse(forOrganizationID),
		ForOrganizationLocationID:  uuid.MustParse(forOrganizationLocationID),
		ForOrganizationStructureID: uuid.MustParse(forOrganizationStructureID),
		JobID:                      uuid.MustParse(jobID),
		RequestCategoryID:          uuid.MustParse(requestCategoryID),
		// GradeID:                    parsedGradeID,
		ExpectedDate:             parsedExpectedDate,
		Experiences:              experiences,
		DocumentNumber:           documentNumber,
		DocumentDate:             parsedDocumentDate,
		MaleNeeds:                maleNeeds,
		FemaleNeeds:              femaleNeeds,
		MinimumAge:               minimumAge,
		MaximumAge:               maximumAge,
		MinimumExperience:        minimumExperienceInt,
		MaritalStatus:            maritalStatus,
		MinimumEducation:         minimumEducation,
		RequiredQualification:    requiredQualification,
		Certificate:              certificate,
		ComputerSkill:            computerSkill,
		LanguageSkill:            languageSkill,
		OtherSkill:               otherSkill,
		Jobdesc:                  jobdesc,
		SalaryMin:                salaryMin,
		SalaryMax:                salaryMax,
		RequestorID:              &parsedRequestorID,
		DepartmentHead:           departmentHeadUUID,
		VpGmDirector:             vpGmDirectorUUID,
		CEO:                      ceoUUID,
		HrdHoUnit:                hrdHoUnitUUID,
		MPPlanningHeaderID:       mpPlanningHeaderUUID,
		Status:                   status,
		MPRequestType:            mpRequestType,
		RecruitmentType:          recruitmentType,
		MPPPeriodID:              &parsedMppPeriodID,
		EmpOrganizationID:        &parsedEmpOrganizationID,
		JobLevelID:               &parsedJobLevelID,
		IsReplacement:            isReplacement,
		CreatedAt:                parsedCreatedAt,
		UpdatedAt:                parsedUpdatedAt,
		RequestCategory:          requestCategory,
		RequestMajors:            requestMajors,
		GradeName:                gradeName,
		OrganizationName:         organizationName,
		OrganizationCategory:     organizationCategory,
		OrganizationLocationName: organizationLocationName,
		ForOrganizationName:      forOrganizationName,
		ForOrganizationLocation:  forOrganizationLocation,
		ForOrganizationStructure: forOrganizationStructure,
		JobName:                  jobName,
		RequestorName:            requestorName,
		DepartmentHeadName:       departmentHeadName,
		HrdHoUnitName:            hrdHoUnitName,
		VpGmDirectorName:         vpGmDirectorName,
		CeoName:                  ceoName,
		EmpOrganizationName:      empOrganizationName,
		JobLevelName:             jobLevelName,
		JobLevel:                 jobLevel,
		ApprovedByDepartmentHead: approvedByDepartmentHead,
		ApprovedByVpGmDirector:   approvedByVpGmDirector,
		ApprovedByCEO:            approvedByCEO,
		ApprovedByHrdHoUnit:      approvedByHrdHoUnit,
	}, nil
}

func (d *MPRequestHelper) ConvertMapInterfaceToResponseMinimal(mprMap map[string]interface{}) (*response.MPRequestHeaderResponse, error) {
	mprData, ok := mprMap["mp_request_header"].(map[string]interface{})
	if !ok {
		d.Log.Errorf("MPRequestHeader information is missing or invalid")
		return nil, errors.New("MPRequestHeader information is missing or invalid")
	}

	mprID, ok := mprData["id"].(string)
	if !ok {
		d.Log.Errorf("MPRequestHeader ID is missing or invalid")
		return nil, errors.New("MPRequestHeader ID is missing or invalid")
	}

	// gradeID, ok := mprData["grade_id"].(string)
	// if !ok {
	// 	d.Log.Errorf("Grade ID is missing or invalid")
	// 	return nil, errors.New("Grade ID is missing or invalid")
	// }

	// var parsedGradeID *uuid.UUID
	// if gradeID != "" {
	// 	gradeUUID, err := uuid.Parse(gradeID)
	// 	if err != nil {
	// 		d.Log.Errorf("Invalid grade ID format: %v", err)
	// 		parsedGradeID = nil
	// 	} else {
	// 		parsedGradeID = &gradeUUID
	// 	}
	// } else {
	// 	parsedGradeID = nil
	// }

	organizationID, ok := mprData["organization_id"].(string)
	if !ok {
		d.Log.Errorf("Organization ID is missing or invalid")
		return nil, errors.New("Organization ID is missing or invalid")
	}

	organizationLocationID, ok := mprData["organization_location_id"].(string)
	if !ok {
		d.Log.Errorf("Organization Location ID is missing or invalid")
		return nil, errors.New("Organization Location ID is missing or invalid")
	}
	forOrganizationID, ok := mprData["for_organization_id"].(string)
	if !ok {
		d.Log.Errorf("For Organization ID is missing or invalid")
		return nil, errors.New("For Organization ID is missing or invalid")
	}

	forOrganizationLocationID, ok := mprData["for_organization_location_id"].(string)
	if !ok {
		d.Log.Errorf("For Organization Location ID is missing or invalid")
		return nil, errors.New("For Organization Location ID is missing or invalid")
	}

	forOrganizationStructureID, ok := mprData["for_organization_structure_id"].(string)
	if !ok {
		d.Log.Errorf("For Organization Structure ID is missing or invalid")
		return nil, errors.New("For Organization Structure ID is missing or invalid")
	}

	jobID, ok := mprData["job_id"].(string)
	if !ok {
		d.Log.Errorf("Job ID is missing or invalid")
		return nil, errors.New("Job ID is missing or invalid")
	}

	requestCategoryID, ok := mprData["request_category_id"].(string)
	if !ok {
		d.Log.Errorf("Request Category ID is missing or invalid")
		return nil, errors.New("Request Category ID is missing or invalid")
	}

	expectedDate, ok := mprData["expected_date"].(string)
	if !ok {
		d.Log.Errorf("Expected Date is missing or invalid")
		return nil, errors.New("Expected Date is missing or invalid")
	}

	experiences, ok := mprData["experiences"].(string)
	if !ok {
		d.Log.Errorf("Experiences is missing or invalid")
		return nil, errors.New("Experiences is missing or invalid")
	}

	documentNumber, ok := mprData["document_number"].(string)
	if !ok {
		d.Log.Errorf("Document Number is missing or invalid")
		return nil, errors.New("Document Number is missing or invalid")
	}

	documentDate, ok := mprData["document_date"].(string)
	if !ok {
		d.Log.Errorf("Document Date is missing or invalid")
		return nil, errors.New("Document Date is missing or invalid")
	}

	maleNeedsFloat, ok := mprData["male_needs"].(float64)
	if !ok {
		d.Log.Errorf("Male needs is missing or invalid")
		return nil, errors.New("Male needs is missing or invalid")
	}
	maleNeeds := int(maleNeedsFloat)

	femaleNeedsFloat, ok := mprData["female_needs"].(float64)
	if !ok {
		d.Log.Errorf("female needs is missing or invalid")
		return nil, errors.New("female needs is missing or invalid")
	}
	femaleNeeds := int(femaleNeedsFloat)

	minimumAgeFloat, ok := mprData["minimum_age"].(float64)
	if !ok {
		d.Log.Errorf("Minimum Age is missing or invalid")
		return nil, errors.New("Minimum Age is missing or invalid")
	}
	minimumAge := int(minimumAgeFloat)

	maximumAgeFloat, ok := mprData["maximum_age"].(float64)
	if !ok {
		d.Log.Errorf("Maximum Age is missing or invalid")
		return nil, errors.New("Maximum Age is missing or invalid")
	}
	maximumAge := int(maximumAgeFloat)

	minimumExperienceFloat, ok := mprData["minimum_experience"].(float64)
	if !ok {
		d.Log.Errorf("Minimum Experience is missing or invalid")
		return nil, errors.New("Minimum Experience is missing or invalid")
	}
	minimumExperienceInt := int(minimumExperienceFloat)

	maritalStatus, ok := mprData["marital_status"].(string)
	if !ok {
		d.Log.Errorf("Marital Status is missing or invalid")
		return nil, errors.New("Marital Status is missing or invalid")
	}

	minimumEducation, ok := mprData["minimum_education"].(string)
	if !ok {
		d.Log.Errorf("Minimum Education is missing or invalid")
		return nil, errors.New("Minimum Education is missing or invalid")
	}

	requiredQualification, ok := mprData["required_qualification"].(string)
	if !ok {
		d.Log.Errorf("Required Qualification is missing or invalid")
		return nil, errors.New("Required Qualification is missing or invalid")
	}

	certificate, ok := mprData["certificate"].(string)
	if !ok {
		d.Log.Errorf("Certificate is missing or invalid")
		return nil, errors.New("Certificate is missing or invalid")
	}

	computerSkill, ok := mprData["computer_skill"].(string)
	if !ok {
		d.Log.Errorf("Computer Skill is missing or invalid")
		return nil, errors.New("Computer Skill is missing or invalid")
	}

	languageSkill, ok := mprData["language_skill"].(string)
	if !ok {
		d.Log.Errorf("Language Skill is missing or invalid")
		return nil, errors.New("Language Skill is missing or invalid")
	}

	otherSkill, ok := mprData["other_skill"].(string)
	if !ok {
		d.Log.Errorf("Other Skill is missing or invalid")
		return nil, errors.New("Other Skill is missing or invalid")
	}

	jobdesc, ok := mprData["jobdesc"].(string)
	if !ok {
		d.Log.Errorf("Jobdesc is missing or invalid")
		return nil, errors.New("Jobdesc is missing or invalid")
	}

	salaryMin, ok := mprData["salary_min"].(string)
	if !ok {
		d.Log.Errorf("Salary Min is missing or invalid")
		return nil, errors.New("Salary Min is missing or invalid")
	}

	salaryMax, ok := mprData["salary_max"].(string)
	if !ok {
		d.Log.Errorf("Salary Max is missing or invalid")
		return nil, errors.New("Salary Max is missing or invalid")
	}

	status, ok := mprData["status"].(string)
	if !ok {
		d.Log.Errorf("Status is missing or invalid")
		return nil, errors.New("Status is missing or invalid")
	}

	mpRequestType, ok := mprData["mp_request_type"].(string)
	if !ok {
		d.Log.Errorf("MP Request Type is missing or invalid")
		return nil, errors.New("MP Request Type is missing or invalid")
	}

	recruitmentType, ok := mprData["recruitment_type"].(string)
	if !ok {
		d.Log.Errorf("Recruitment Type is missing or invalid")
		return nil, errors.New("Recruitment Type is missing or invalid")
	}

	mppPeriodID, ok := mprData["mpp_period_id"].(string)
	if !ok {
		d.Log.Errorf("MPP Period ID is missing or invalid")
		return nil, errors.New("MPP Period ID is missing or invalid")
	}

	empOrganizationID, ok := mprData["emp_organization_id"].(string)
	if !ok {
		d.Log.Errorf("Emp Organization ID is missing or invalid")
		return nil, errors.New("Emp Organization ID is missing or invalid")
	}

	jobLevelID, ok := mprData["job_level_id"].(string)
	if !ok {
		d.Log.Errorf("Job Level ID is missing or invalid")
		return nil, errors.New("Job Level ID is missing or invalid")
	}

	isReplacement, ok := mprData["is_replacement"].(bool)
	if !ok {
		d.Log.Errorf("Is Replacement is missing or invalid")
		return nil, errors.New("Is Replacement is missing or invalid")
	}

	createdAt, ok := mprData["created_at"].(string)
	if !ok {
		d.Log.Errorf("Created At is missing or invalid")
		return nil, errors.New("Created At is missing or invalid")
	}

	updatedAt, ok := mprData["updated_at"].(string)
	if !ok {
		d.Log.Errorf("Updated At is missing or invalid")
		return nil, errors.New("Updated At is missing or invalid")
	}

	mprCloneID := uuid.MustParse(mprID)
	parsedMppPeriodID := uuid.MustParse(mppPeriodID)
	parsedEmpOrganizationID := uuid.MustParse(empOrganizationID)
	parsedJobLevelID := uuid.MustParse(jobLevelID)

	parsedExpectedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", expectedDate)
	if err != nil {
		d.Log.Errorf("Invalid expected date format: %v", err)
		return nil, errors.New("Invalid expected date format")
	}
	parsedDocumentDate, err := time.Parse("2006-01-02T15:04:05Z07:00", documentDate)
	if err != nil {
		d.Log.Errorf("Invalid document date format: %v", err)
		return nil, errors.New("Invalid document date format")
	}
	parsedCreatedAt, err := time.Parse("2006-01-02T15:04:05Z07:00", createdAt)
	if err != nil {
		d.Log.Errorf("Invalid created at format: %v", err)
		return nil, errors.New("Invalid created at format")
	}
	parsedUpdatedAt, err := time.Parse("2006-01-02T15:04:05Z07:00", updatedAt)
	if err != nil {
		d.Log.Errorf("Invalid updated at format: %v", err)
		return nil, errors.New("Invalid updated at format")
	}

	return &response.MPRequestHeaderResponse{
		ID:                         uuid.MustParse(mprID),
		MPRCloneID:                 &mprCloneID,
		OrganizationID:             uuid.MustParse(organizationID),
		OrganizationLocationID:     uuid.MustParse(organizationLocationID),
		ForOrganizationID:          uuid.MustParse(forOrganizationID),
		ForOrganizationLocationID:  uuid.MustParse(forOrganizationLocationID),
		ForOrganizationStructureID: uuid.MustParse(forOrganizationStructureID),
		JobID:                      uuid.MustParse(jobID),
		// GradeID:                    parsedGradeID,
		RequestCategoryID:     uuid.MustParse(requestCategoryID),
		ExpectedDate:          parsedExpectedDate,
		Experiences:           experiences,
		DocumentNumber:        documentNumber,
		DocumentDate:          parsedDocumentDate,
		MaleNeeds:             maleNeeds,
		FemaleNeeds:           femaleNeeds,
		MinimumAge:            minimumAge,
		MaximumAge:            maximumAge,
		MinimumExperience:     minimumExperienceInt,
		MaritalStatus:         maritalStatus,
		MinimumEducation:      minimumEducation,
		RequiredQualification: requiredQualification,
		Certificate:           certificate,
		ComputerSkill:         computerSkill,
		LanguageSkill:         languageSkill,
		OtherSkill:            otherSkill,
		Jobdesc:               jobdesc,
		SalaryMin:             salaryMin,
		SalaryMax:             salaryMax,
		Status:                status,
		MPRequestType:         mpRequestType,
		RecruitmentType:       recruitmentType,
		MPPPeriodID:           &parsedMppPeriodID,
		EmpOrganizationID:     &parsedEmpOrganizationID,
		JobLevelID:            &parsedJobLevelID,
		IsReplacement:         isReplacement,
		CreatedAt:             parsedCreatedAt,
		UpdatedAt:             parsedUpdatedAt,
	}, nil
}
