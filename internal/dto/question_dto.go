package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IQuestionDTO interface {
	ConvertEntityToResponse(ent *entity.Question) *response.QuestionResponse
}

type QuestionDTO struct {
	Log           *logrus.Logger
	AnswerTypeDTO IAnswerTypeDTO
}

func NewQuestionDTO(log *logrus.Logger, answerTypeDTO IAnswerTypeDTO) IQuestionDTO {
	return &QuestionDTO{
		Log:           log,
		AnswerTypeDTO: answerTypeDTO,
	}
}

func QuestionDTOFactory(log *logrus.Logger) IQuestionDTO {
	answerTypeDTO := AnswerTypeDTOFactory(log)
	return NewQuestionDTO(log, answerTypeDTO)
}

func (dto *QuestionDTO) ConvertEntityToResponse(ent *entity.Question) *response.QuestionResponse {
	return &response.QuestionResponse{
		ID:                 ent.ID,
		TemplateQuestionID: ent.TemplateQuestionID,
		AnswerTypeID:       ent.AnswerTypeID,
		Name:               ent.Name,
		AnswerTypeResponse: func() *response.AnswerTypeResponse {
			if ent.AnswerType == nil {
				return nil
			}
			return dto.AnswerTypeDTO.ConvertEntityToResponse(ent.AnswerType)
		}(),
	}
}
