package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IJobPostingDTO interface {
	ConvertEntityToResponse(ent *entity.JobPosting) *response.JobPostingResponse
}

type JobPostingDTO struct {
	Log                         *logrus.Logger
	OrganizationMessage         messaging.IOrganizationMessage
	JobMessage                  messaging.IJobPlafonMessage
	ProjectRecruitmentHeaderDTO IProjectRecruitmentHeaderDTO
	Viper                       *viper.Viper
}

func NewJobPostingDTO(
	log *logrus.Logger,
	orgMessage messaging.IOrganizationMessage,
	jobMessage messaging.IJobPlafonMessage,
	prhDTO IProjectRecruitmentHeaderDTO,
	viper *viper.Viper,
) IJobPostingDTO {
	return &JobPostingDTO{
		Log:                         log,
		OrganizationMessage:         orgMessage,
		JobMessage:                  jobMessage,
		ProjectRecruitmentHeaderDTO: prhDTO,
		Viper:                       viper,
	}
}

func JobPostingDTOFactory(log *logrus.Logger, viper *viper.Viper) IJobPostingDTO {
	orgMessage := messaging.OrganizationMessageFactory(log)
	jobMessage := messaging.JobPlafonMessageFactory(log)
	prhDTO := ProjectRecruitmentHeaderDTOFactory(log)
	return NewJobPostingDTO(log, orgMessage, jobMessage, prhDTO, viper)
}

func (dto *JobPostingDTO) ConvertEntityToResponse(ent *entity.JobPosting) *response.JobPostingResponse {
	var organizationName, organizationLocationName, jobName string

	organization, err := dto.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: ent.ForOrganizationID.String(),
	})
	if err != nil {
		dto.Log.Errorf("[JobPostingDTO.ConvertEntityToResponse] " + err.Error())
		organizationName = ""
	}
	organizationName = organization.Name

	organizationLocation, err := dto.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
		ID: ent.ForOrganizationLocationID.String(),
	})
	if err != nil {
		dto.Log.Errorf("[JobPostingDTO.ConvertEntityToResponse] " + err.Error())
		organizationLocationName = ""
	}
	organizationLocationName = organizationLocation.Name

	job, err := dto.JobMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		ID: ent.JobID.String(),
	})
	if err != nil {
		dto.Log.Errorf("[JobPostingDTO.ConvertEntityToResponse] " + err.Error())
		jobName = ""
	}
	jobName = job.Name

	return &response.JobPostingResponse{
		ID:                         ent.ID,
		ProjectRecruitmentHeaderID: ent.ProjectRecruitmentHeaderID,
		MPRequestID:                ent.MPRequestID,
		JobID:                      ent.JobID,
		ForOrganizationID:          ent.ForOrganizationID,
		ForOrganizationLocationID:  ent.ForOrganizationLocationID,
		DocumentNumber:             ent.DocumentNumber,
		DocumentDate:               ent.DocumentDate.Format("2006-01-02"),
		RecruitmentType:            ent.RecruitmentType,
		StartDate:                  ent.StartDate.Format("2006-01-02"),
		EndDate:                    ent.EndDate.Format("2006-01-02"),
		Status:                     ent.Status,
		SalaryMin:                  ent.SalaryMin,
		SalaryMax:                  ent.SalaryMax,
		IsApplied:                  ent.IsApplied,
		IsSaved:                    ent.IsSaved,
		ContentDescription:         ent.ContentDescription,
		OrganizationLogo: func() *string {
			if ent.OrganizationLogo != "" {
				dto.Log.Info("Organization Logo: ", ent.OrganizationLogo)
				organizationLogoURL := dto.Viper.GetString("app.url") + ent.OrganizationLogo
				return &organizationLogoURL
			}
			return nil
		}(),
		Poster: func() *string {
			if ent.Poster != "" {
				posterURL := dto.Viper.GetString("app.url") + ent.Poster
				return &posterURL
			}
			return nil
		}(),
		Link:                    ent.Link,
		ForOrganizationName:     organizationName,
		ForOrganizationLocation: organizationLocationName,
		JobName:                 jobName,

		ProjectRecruitmentHeader: func() *response.ProjectRecruitmentHeaderResponse {
			if ent.ProjectRecruitmentHeader == nil {
				return nil
			}
			return dto.ProjectRecruitmentHeaderDTO.ConvertEntityToResponse(ent.ProjectRecruitmentHeader)
		}(),
		MPRequest: func() *response.MPRequestHeaderResponse {
			if ent.MPRequest == nil {
				return nil
			}
			return &response.MPRequestHeaderResponse{
				ID:         ent.MPRequest.ID,
				MPRCloneID: ent.MPRequest.MPRCloneID,
				Status:     string(ent.MPRequest.Status),
			}
		}(),
	}
}
