package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IAdministrativeSelectionDTO interface {
	ConvertEntityToResponse(ent *entity.AdministrativeSelection) (*response.AdministrativeSelectionResponse, error)
}

type AdministrativeSelectionDTO struct {
	Log                     *logrus.Logger
	AdministrativeResultDTO IAdministrativeResultDTO
	JobPostingDTO           IJobPostingDTO
	ProjectPICDTO           IProjectPicDTO
	Viper                   *viper.Viper
}

func NewAdministrativeSelectionDTO(
	log *logrus.Logger,
	administrativeResultDTO IAdministrativeResultDTO,
	viper *viper.Viper,
	jobPostingDTO IJobPostingDTO,
	projectPICDTO IProjectPicDTO,
) IAdministrativeSelectionDTO {
	return &AdministrativeSelectionDTO{
		Log:                     log,
		AdministrativeResultDTO: administrativeResultDTO,
		Viper:                   viper,
		JobPostingDTO:           jobPostingDTO,
		ProjectPICDTO:           projectPICDTO,
	}
}

func AdministrativeSelectionDTOFactory(log *logrus.Logger, viper *viper.Viper) IAdministrativeSelectionDTO {
	administrativeResultDTO := AdministrativeResultDTOFactory(log, viper)
	jobPostingDTO := JobPostingDTOFactory(log, viper)
	projectPICDTO := ProjectPicDTOFactory(log)
	return NewAdministrativeSelectionDTO(log, administrativeResultDTO, viper, jobPostingDTO, projectPICDTO)
}

func (dto *AdministrativeSelectionDTO) ConvertEntityToResponse(ent *entity.AdministrativeSelection) (*response.AdministrativeSelectionResponse, error) {
	var administrativeResultsResponse []response.AdministrativeResultResponse

	for _, administrativeResult := range ent.AdministrativeResults {
		administrativeResultResponse, err := dto.AdministrativeResultDTO.ConvertEntityToResponse(&administrativeResult)
		if err != nil {
			return nil, err
		}
		administrativeResultsResponse = append(administrativeResultsResponse, *administrativeResultResponse)
	}

	return &response.AdministrativeSelectionResponse{
		ID:           ent.ID,
		JobPostingID: ent.JobPostingID,
		ProjectPicID: ent.ProjectPicID,
		Status:       ent.Status,
		// VerifiedAt:            ent.VerifiedAt,
		// VerifiedBy:            &ent.VerifiedBy,
		DocumentDate:    ent.DocumentDate,
		DocumentNumber:  ent.DocumentNumber,
		TotalApplicants: ent.TotalApplicants,
		CreatedAt:       ent.CreatedAt,
		UpdatedAt:       ent.UpdatedAt,
		JobPosting: func() *response.JobPostingResponse {
			if ent.JobPosting != nil {
				jobPostingResponse := dto.JobPostingDTO.ConvertEntityToResponse(ent.JobPosting)
				return jobPostingResponse
			}
			return nil
		}(),
		ProjectPIC: func() *response.ProjectPicResponse {
			if ent.ProjectPIC != nil {
				projectPICResponse := dto.ProjectPICDTO.ConvertEntityToResponse(ent.ProjectPIC)
				return projectPICResponse
			}
			return nil
		}(),
		AdministrativeResults: administrativeResultsResponse,
	}, nil
}
