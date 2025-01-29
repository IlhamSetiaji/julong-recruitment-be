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
	Viper                   *viper.Viper
}

func NewAdministrativeSelectionDTO(log *logrus.Logger, administrativeResultDTO IAdministrativeResultDTO, viper *viper.Viper) IAdministrativeSelectionDTO {
	return &AdministrativeSelectionDTO{
		Log:                     log,
		AdministrativeResultDTO: administrativeResultDTO,
		Viper:                   viper,
	}
}

func AdministrativeSelectionDTOFactory(log *logrus.Logger, viper *viper.Viper) IAdministrativeSelectionDTO {
	administrativeResultDTO := AdministrativeResultDTOFactory(log, viper)
	return NewAdministrativeSelectionDTO(log, administrativeResultDTO, viper)
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
		ID:                    ent.ID,
		JobPostingID:          ent.JobPostingID,
		ProjectPicID:          ent.ProjectPicID,
		Status:                ent.Status,
		VerifiedAt:            ent.VerifiedAt,
		VerifiedBy:            &ent.VerifiedBy,
		DocumentDate:          ent.DocumentDate,
		DocumentNumber:        ent.DocumentNumber,
		TotalApplicants:       len(ent.AdministrativeResults),
		CreatedAt:             ent.CreatedAt,
		UpdatedAt:             ent.UpdatedAt,
		JobPosting:            nil,
		ProjectPIC:            nil,
		AdministrativeResults: administrativeResultsResponse,
	}, nil
}
