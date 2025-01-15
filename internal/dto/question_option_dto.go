package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IQuestionOptionDTO interface {
	ConvertEntityToResponse(ent *entity.QuestionOption) *response.QuestionOptionResponse
}

type QuestionOptionDTO struct {
	Log *logrus.Logger
}

func NewQuestionOptionDTO(log *logrus.Logger) IQuestionOptionDTO {
	return &QuestionOptionDTO{
		Log: log,
	}
}

func QuestionOptionDTOFactory(log *logrus.Logger) IQuestionOptionDTO {
	return NewQuestionOptionDTO(log)
}

func (dto *QuestionOptionDTO) ConvertEntityToResponse(ent *entity.QuestionOption) *response.QuestionOptionResponse {
	return &response.QuestionOptionResponse{
		ID:         ent.ID,
		QuestionID: ent.QuestionID,
		OptionText: ent.OptionText,
		CreatedAt:  ent.CreatedAt,
		UpdatedAt:  ent.UpdatedAt,
	}
}
