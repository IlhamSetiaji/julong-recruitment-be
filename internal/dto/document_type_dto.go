package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IDocumentTypeDTO interface {
	ConvertEntityToResponse(ent *entity.DocumentType) *response.DocumentTypeResponse
}

type DocumentTypeDTO struct {
	Log *logrus.Logger
}

func NewDocumentTypeDTO(log *logrus.Logger) IDocumentTypeDTO {
	return &DocumentTypeDTO{
		Log: log,
	}
}

func DocumentTypeDTOFactory(log *logrus.Logger) IDocumentTypeDTO {
	return NewDocumentTypeDTO(log)
}

func (dto *DocumentTypeDTO) ConvertEntityToResponse(ent *entity.DocumentType) *response.DocumentTypeResponse {
	return &response.DocumentTypeResponse{
		ID:        ent.ID,
		Name:      ent.Name,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}
}
