package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IAnswerTypeDTO interface {
	ConvertEntityToResponse(ent *entity.AnswerType) *response.AnswerTypeResponse
}

type AnswerTypeDTO struct {
	Log *logrus.Logger
}

func NewAnswerTypeDTO(log *logrus.Logger) IAnswerTypeDTO {
	return &AnswerTypeDTO{
		Log: log,
	}
}

func AnswerTypeDTOFactory(log *logrus.Logger) IAnswerTypeDTO {
	return NewAnswerTypeDTO(log)
}

func (dto *AnswerTypeDTO) ConvertEntityToResponse(ent *entity.AnswerType) *response.AnswerTypeResponse {
	return &response.AnswerTypeResponse{
		ID:   ent.ID.String(),
		Name: ent.Name,
	}
}
