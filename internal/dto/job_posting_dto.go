package dto

import (
	"strings"

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

	organizationId := ent.ForOrganizationID.String()
	dto.Log.Infof("[JobPostingDTO.ConvertEntityToResponse] organizationId: %s", organizationId)

	// Validasi untuk organization
	organization, err := dto.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: ent.ForOrganizationID.String(),
	})
	if err != nil {
		dto.Log.Errorf("[JobPostingDTO.ConvertEntityToResponse] Failed to find organization: %s", err.Error())
		organizationName = "Unknown" // Nilai default jika terjadi error
	} else if organization != nil {
		organizationName = organization.Name
	} else {
		organizationName = "Unknown" // Nilai default jika organization nil
	}

	// Validasi untuk organizationLocation
	organizationLocation, err := dto.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
		ID: ent.ForOrganizationLocationID.String(),
	})
	if err != nil {
		dto.Log.Errorf("[JobPostingDTO.ConvertEntityToResponse] Failed to find organization location: %s", err.Error())
		organizationLocationName = "Unknown" // Nilai default jika terjadi error
	} else if organizationLocation != nil {
		organizationLocationName = organizationLocation.Name
	} else {
		organizationLocationName = "Unknown" // Nilai default jika organizationLocation nil
	}

	// Validasi untuk job
	job, err := dto.JobMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		ID: ent.JobID.String(),
	})
	if err != nil {
		dto.Log.Errorf("[JobPostingDTO.ConvertEntityToResponse] Failed to find job: %s", err.Error())
		jobName = "Unknown" // Nilai default jika terjadi error
	} else if job != nil {
		jobName = job.Name
	} else {
		jobName = "Unknown" // Nilai default jika job nil
	}

	// Potong jobName jika ada " - "
	if idx := strings.Index(jobName, " - "); idx != -1 {
		jobName = jobName[:idx]
	}

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
		AppliedDate:                ent.AppliedDate,
		ApplicantStatus:            ent.ApplicantStatus,
		ApplicantProcessStatus:     ent.ApplicantProcessStatus,
		ContentDescription:         ent.ContentDescription,
		MinimumWorkExperience:      ent.MinimumWorkExperience,
		Name:                       ent.Name,
		IsShow:                     ent.IsShow,
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
		TotalApplicant:          ent.TotalApplicant,

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
