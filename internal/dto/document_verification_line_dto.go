package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentVerificationLineDTO interface {
	ConvertEntityToResponse(ent *entity.DocumentVerificationLine) *response.DocumentVerificationLineResponse
}

type DocumentVerificationLineDTO struct {
	Log                     *logrus.Logger
	DocumentVerificationDTO IDocumentVerificationDTO
	Viper                   *viper.Viper
}

func NewDocumentVerificationLineDTO(
	log *logrus.Logger,
	documentVerificationDTO IDocumentVerificationDTO,
	viper *viper.Viper,
) IDocumentVerificationLineDTO {
	return &DocumentVerificationLineDTO{
		Log:                     log,
		DocumentVerificationDTO: documentVerificationDTO,
		Viper:                   viper,
	}
}

func DocumentVerificationLineDTOFactory(log *logrus.Logger, viper *viper.Viper) IDocumentVerificationLineDTO {
	documentVerificationDTO := DocumentVerificationDTOFactory(log)
	return NewDocumentVerificationLineDTO(log, documentVerificationDTO, viper)
}

func (dto *DocumentVerificationLineDTO) ConvertEntityToResponse(ent *entity.DocumentVerificationLine) *response.DocumentVerificationLineResponse {
	return &response.DocumentVerificationLineResponse{
		ID:                           ent.ID,
		DocumentVerificationHeaderID: ent.DocumentVerificationHeaderID,
		DocumentVerificationID:       ent.DocumentVerificationID,
		Path: func() string {
			if ent.Path != "" {
				return dto.Viper.GetString("app.url") + ent.Path
			}
			return ""
		}(),
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
		DocumentVerification: func() *response.DocumentVerificationResponse {
			if ent.DocumentVerification != nil {
				return dto.DocumentVerificationDTO.ConvertEntityToResponse(ent.DocumentVerification)
			}
			return nil
		}(),
	}
}
