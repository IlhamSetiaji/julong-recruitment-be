package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IMailTemplateDTO interface {
	ConvertEntityToResponse(ent *entity.MailTemplate) *response.MailTemplateResponse
}

type MailTemplateDTO struct {
	Log             *logrus.Logger
	DocumentTypeDTO IDocumentTypeDTO
}

func NewMailTemplateDTO(
	log *logrus.Logger,
	documentTypeDTO IDocumentTypeDTO,
) IMailTemplateDTO {
	return &MailTemplateDTO{
		Log:             log,
		DocumentTypeDTO: documentTypeDTO,
	}
}

func MailTemplateDTOFactory(log *logrus.Logger) IMailTemplateDTO {
	documentTypeDTO := DocumentTypeDTOFactory(log)
	return NewMailTemplateDTO(log, documentTypeDTO)
}

func (dto *MailTemplateDTO) ConvertEntityToResponse(ent *entity.MailTemplate) *response.MailTemplateResponse {
	return &response.MailTemplateResponse{
		ID:             ent.ID,
		DocumentTypeID: ent.DocumentTypeID,
		Name:           ent.Name,
		Subject:        ent.Subject,
		Body:           ent.Body,
		CreatedAt:      ent.CreatedAt,
		UpdatedAt:      ent.UpdatedAt,
		DocumentType: func() *response.DocumentTypeResponse {
			if ent.DocumentType == nil {
				return nil
			}
			return dto.DocumentTypeDTO.ConvertEntityToResponse(ent.DocumentType)
		}(),
	}
}
