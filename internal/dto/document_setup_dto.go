package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IDocumentSetupDTO interface {
	ConvertEntityToResponse(ent *entity.DocumentSetup) *response.DocumentSetupResponse
}

type DocumentSetupDTO struct {
	Log             *logrus.Logger
	DocumentTypeDTO IDocumentTypeDTO
}

func NewDocumentSetupDTO(
	log *logrus.Logger,
	documentTypeDTO IDocumentTypeDTO,
) IDocumentSetupDTO {
	return &DocumentSetupDTO{
		Log:             log,
		DocumentTypeDTO: documentTypeDTO,
	}
}

func DocumentSetupDTOFactory(log *logrus.Logger) IDocumentSetupDTO {
	documentTypeDTO := DocumentTypeDTOFactory(log)
	return NewDocumentSetupDTO(log, documentTypeDTO)
}

func (dto *DocumentSetupDTO) ConvertEntityToResponse(ent *entity.DocumentSetup) *response.DocumentSetupResponse {
	return &response.DocumentSetupResponse{
		ID:              ent.ID,
		DocumentTypeID:  ent.DocumentTypeID,
		RecruitmentType: ent.RecruitmentType,
		Title:           ent.Title,
		Header:          ent.Header,
		Body:            ent.Body,
		Footer:          ent.Footer,
		CreatedAt:       ent.CreatedAt,
		UpdatedAt:       ent.UpdatedAt,
		DocumentType: func() *response.DocumentTypeResponse {
			if ent.DocumentType == nil {
				return nil
			}
			return dto.DocumentTypeDTO.ConvertEntityToResponse(ent.DocumentType)
		}(),
	}
}
