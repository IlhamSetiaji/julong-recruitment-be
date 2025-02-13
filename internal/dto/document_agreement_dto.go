package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentAgreementDTO interface {
	ConvertEntityToResponse(ent *entity.DocumentAgreement) *response.DocumentAgreementResponse
}

type DocumentAgreementDTOIDocumentAgreementDTO struct {
	Log                *logrus.Logger
	DocumentSendingDTO IDocumentSendingDTO
	ApplicantDTO       IApplicantDTO
	Viper              *viper.Viper
}

func NewDocumentAgreementDTOIDocumentAgreementDTO(log *logrus.Logger, documentSendingDTO IDocumentSendingDTO, applicantDTO IApplicantDTO, viper *viper.Viper) IDocumentAgreementDTO {
	return &DocumentAgreementDTOIDocumentAgreementDTO{
		Log:                log,
		DocumentSendingDTO: documentSendingDTO,
		ApplicantDTO:       applicantDTO,
		Viper:              viper,
	}
}

func DocumentAgreementDTOIDocumentAgreementDTOFactory(log *logrus.Logger, viper *viper.Viper) IDocumentAgreementDTO {
	documentSendingDTO := DocumentSendingDTOFactory(log, viper)
	applicantDTO := ApplicantDTOFactory(log, viper)
	return NewDocumentAgreementDTOIDocumentAgreementDTO(log, documentSendingDTO, applicantDTO, viper)
}

func (dto *DocumentAgreementDTOIDocumentAgreementDTO) ConvertEntityToResponse(ent *entity.DocumentAgreement) *response.DocumentAgreementResponse {
	return &response.DocumentAgreementResponse{
		ID:                ent.ID,
		DocumentSendingID: ent.DocumentSendingID,
		ApplicantID:       ent.ApplicantID,
		Status:            ent.Status,
		Path:              dto.Viper.GetString("app.url") + ent.Path,
		CreatedAt:         ent.CreatedAt,
		UpdatedAt:         ent.UpdatedAt,
		DocumentSending: func() *response.DocumentSendingResponse {
			if ent.DocumentSending != nil {
				return dto.DocumentSendingDTO.ConvertEntityToResponse(ent.DocumentSending)
			} else {
				return nil
			}
		}(),
		Applicant: func() *response.ApplicantResponse {
			if ent.Applicant != nil {
				res, err := dto.ApplicantDTO.ConvertEntityToResponse(ent.Applicant)
				if err != nil {
					dto.Log.Error(err)
					return nil
				}
				return res
			} else {
				return nil
			}
		}(),
	}
}
