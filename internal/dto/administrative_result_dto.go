package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IAdministrativeResultDTO interface {
	ConvertEntityToResponse(ent *entity.AdministrativeResult) (*response.AdministrativeResultResponse, error)
}

type AdministrativeResultDTO struct {
	Log          *logrus.Logger
	ApplicantDTO IApplicantDTO
	Viper        *viper.Viper
}

func NewAdministrativeResultDTO(log *logrus.Logger, applicantDTO IApplicantDTO, viper *viper.Viper) IAdministrativeResultDTO {
	return &AdministrativeResultDTO{
		Log:          log,
		ApplicantDTO: applicantDTO,
		Viper:        viper,
	}
}

func AdministrativeResultDTOFactory(log *logrus.Logger, viper *viper.Viper) IAdministrativeResultDTO {
	applicantDTO := ApplicantDTOFactory(log, viper)
	return NewAdministrativeResultDTO(log, applicantDTO, viper)
}

func (dto *AdministrativeResultDTO) ConvertEntityToResponse(ent *entity.AdministrativeResult) (*response.AdministrativeResultResponse, error) {
	var applicantResponse *response.ApplicantResponse

	if ent.Applicant != nil {
		appRes, err := dto.ApplicantDTO.ConvertEntityToResponse(ent.Applicant)
		if err != nil {
			return nil, err
		}
		applicantResponse = appRes
	} else {
		applicantResponse = nil
	}

	return &response.AdministrativeResultResponse{
		ID:                        ent.ID,
		AdministrativeSelectionID: ent.AdministrativeSelectionID,
		ApplicantID:               ent.ApplicantID,
		Status:                    ent.Status,
		CreatedAt:                 ent.CreatedAt.Format(dto.Viper.GetString("time_format")),
		UpdatedAt:                 ent.UpdatedAt.Format(dto.Viper.GetString("time_format")),
		Applicant:                 applicantResponse,
	}, nil
}
