package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentSendingDTO interface {
	ConvertEntityToResponse(ent *entity.DocumentSending) *response.DocumentSendingResponse
}

type DocumentSendingDTO struct {
	Log                       *logrus.Logger
	OrganizationMessage       messaging.IOrganizationMessage
	JobMessage                messaging.IJobPlafonMessage
	ProjectRecruitmentLineDTO IProjectRecruitmentLineDTO
	Viper                     *viper.Viper
	ApplicantDTO              IApplicantDTO
	JobPostingDTO             IJobPostingDTO
	DocumentSetupDTO          IDocumentSetupDTO
	GradeMessage              messaging.IGradeMessage
	EmployeeMessage           messaging.IEmployeeMessage
}

func NewDocumentSendingDTO(
	log *logrus.Logger,
	orgMessage messaging.IOrganizationMessage,
	jobMessage messaging.IJobPlafonMessage,
	prlDTO IProjectRecruitmentLineDTO,
	viper *viper.Viper,
	applicantDTO IApplicantDTO,
	jobPostingDTO IJobPostingDTO,
	documentSetupDTO IDocumentSetupDTO,
	gradeMessage messaging.IGradeMessage,
	employeeMessage messaging.IEmployeeMessage,
) IDocumentSendingDTO {
	return &DocumentSendingDTO{
		Log:                       log,
		OrganizationMessage:       orgMessage,
		JobMessage:                jobMessage,
		ProjectRecruitmentLineDTO: prlDTO,
		Viper:                     viper,
		ApplicantDTO:              applicantDTO,
		JobPostingDTO:             jobPostingDTO,
		DocumentSetupDTO:          documentSetupDTO,
		GradeMessage:              gradeMessage,
		EmployeeMessage:           employeeMessage,
	}
}

func DocumentSendingDTOFactory(log *logrus.Logger, viper *viper.Viper) IDocumentSendingDTO {
	orgMessage := messaging.OrganizationMessageFactory(log)
	jobMessage := messaging.JobPlafonMessageFactory(log)
	prlDTO := ProjectRecruitmentLineDTOFactory(log)
	applicantDTO := ApplicantDTOFactory(log, viper)
	jobPostingDTO := JobPostingDTOFactory(log, viper)
	documentSetupDTO := DocumentSetupDTOFactory(log)
	gradeNessage := messaging.GradeMessageFactory(log)
	employeeMessage := messaging.EmployeeMessageFactory(log)
	return NewDocumentSendingDTO(log, orgMessage, jobMessage, prlDTO, viper, applicantDTO, jobPostingDTO, documentSetupDTO, gradeNessage, employeeMessage)
}

func (dto *DocumentSendingDTO) ConvertEntityToResponse(ent *entity.DocumentSending) *response.DocumentSendingResponse {
	jobLevel := &response.SendFindJobLevelByIDMessageResponse{}
	job := &response.SendFindJobByIDMessageResponse{}
	var organizationName string
	var organizationLocationName string
	var gradeName string
	var allowanceApprovalName string
	var err error

	if ent.JobLevelID != nil {
		jobLevel, err = dto.JobMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: ent.JobLevelID.String(),
		})
		if err != nil {
			dto.Log.Errorf("[DocumentSendingDTO.ConvertEntityToResponse] " + err.Error())
			jobLevel = &response.SendFindJobLevelByIDMessageResponse{}
		}
	}

	if ent.JobID != nil {
		job, err = dto.JobMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: ent.JobID.String(),
		})
		if err != nil {
			dto.Log.Errorf("[DocumentSendingDTO.ConvertEntityToResponse] " + err.Error())
			job = &response.SendFindJobByIDMessageResponse{}
		}
	}

	if ent.ForOrganizationID != nil {
		organizationData, err := dto.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: ent.ForOrganizationID.String(),
		})
		if err != nil {
			dto.Log.Errorf("[DocumentSendingDTO.ConvertEntityToResponse] " + err.Error())
			organizationName = ""
		}
		organizationName = organizationData.Name
	}

	if ent.OrganizationLocationID != nil {
		organizationLocationData, err := dto.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
			ID: ent.OrganizationLocationID.String(),
		})
		if err != nil {
			dto.Log.Errorf("[DocumentSendingDTO.ConvertEntityToResponse] " + err.Error())
			organizationLocationName = ""
		}
		if organizationLocationData != nil {
			organizationLocationName = organizationLocationData.Name
		}
	}

	if ent.GradeID != nil {
		gradeData, err := dto.GradeMessage.SendFindByIDMessage(ent.GradeID.String())
		if err != nil {
			dto.Log.Errorf("[DocumentSendingDTO.ConvertEntityToResponse] " + err.Error())
			gradeName = ""
		}
		if gradeData != nil {
			gradeName = gradeData.Name
		}
	}

	if ent.AllowanceApproval != nil {
		allowanceApprovalData, err := dto.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: ent.AllowanceApproval.String(),
		})
		if err != nil {
			dto.Log.Errorf("[DocumentSendingDTO.ConvertEntityToResponse] " + err.Error())
			allowanceApprovalName = ""
		}
		if allowanceApprovalData != nil {
			allowanceApprovalName = allowanceApprovalData.Name
		}
	}

	return &response.DocumentSendingResponse{
		ID:                       ent.ID,
		ProjectRecruitmentLineID: ent.ProjectRecruitmentLineID,
		ApplicantID:              ent.ApplicantID,
		DocumentSetupID:          ent.DocumentSetupID,
		OrganizationLocationID:   ent.OrganizationLocationID,
		GradeID:                  ent.GradeID,
		AllowanceApproval:        ent.AllowanceApproval,
		DocumentDate:             ent.DocumentDate,
		DocumentNumber:           ent.DocumentNumber,
		JoinedDate:               ent.JoinedDate,
		Status:                   ent.Status,
		RecruitmentType:          ent.RecruitmentType,
		BasicWage:                ent.BasicWage,
		PositionalAllowance:      ent.PositionalAllowance,
		OperationalAllowance:     ent.OperationalAllowance,
		MealAllowance:            ent.MealAllowance,
		HouseAllowance:						ent.HouseAllowance,
		JobLocation:              ent.JobLocation,
		HometripTicket:           ent.HometripTicket,
		PeriodAgreement:          ent.PeriodAgreement,
		HomeLocation:             ent.HomeLocation,
		JobLevelID:               ent.JobLevelID,
		JobID:                    ent.JobID,
		JobPostingID:             ent.JobPostingID,
		ForOrganizationID:        ent.ForOrganizationID,
		HiredStatus:              ent.HiredStatus,
		DetailContent:            ent.DetailContent,
		Path: func() string {
			if ent.Path != "" {
				return dto.Viper.GetString("app.url") + ent.Path
			}
			return ""
		}(),
		SyncMidsuit: ent.SyncMidsuit,
		CreatedAt:   ent.CreatedAt,
		UpdatedAt:   ent.UpdatedAt,
		ProjectRecruitmentLine: func() *response.ProjectRecruitmentLineResponse {
			if ent.ProjectRecruitmentLine != nil {
				return dto.ProjectRecruitmentLineDTO.ConvertEntityToResponse(ent.ProjectRecruitmentLine)
			}
			return nil
		}(),
		Applicant: func() *response.ApplicantResponse {
			if ent.Applicant != nil {
				res, err := dto.ApplicantDTO.ConvertEntityToResponse(ent.Applicant)
				if err != nil {
					dto.Log.Errorf("[DocumentSendingDTO.ConvertEntityToResponse] " + err.Error())
					return nil
				}
				return res
			}
			return nil
		}(),
		DocumentSetup: func() *response.DocumentSetupResponse {
			if ent.DocumentSetup != nil {
				return dto.DocumentSetupDTO.ConvertEntityToResponse(ent.DocumentSetup)
			}
			return nil
		}(),
		JobPosting: func() *response.JobPostingResponse {
			if ent.JobPosting != nil {
				return dto.JobPostingDTO.ConvertEntityToResponse(ent.JobPosting)
			}
			return nil
		}(),
		JobLevel:                 jobLevel,
		Job:                      job,
		ForOrganizationName:      &organizationName,
		OrganizationLocationName: &organizationLocationName,
		GradeName:                &gradeName,
		AllowanceApprovalName:    &allowanceApprovalName,
	}
}
