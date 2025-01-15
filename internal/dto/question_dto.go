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
	Log               *logrus.Logger
	AnswerTypeDTO     IAnswerTypeDTO
	QuestionOptionDTO IQuestionOptionDTO
}

func NewQuestionDTO(
	log *logrus.Logger,
	answerTypeDTO IAnswerTypeDTO,
	questionOptionDTO IQuestionOptionDTO,
) IQuestionDTO {
	return &QuestionDTO{
		Log:               log,
		AnswerTypeDTO:     answerTypeDTO,
		QuestionOptionDTO: questionOptionDTO,
	}
}

func QuestionDTOFactory(log *logrus.Logger) IQuestionDTO {
	answerTypeDTO := AnswerTypeDTOFactory(log)
	questionOptionDTO := QuestionOptionDTOFactory(log)
	return NewQuestionDTO(log, answerTypeDTO, questionOptionDTO)
}

func (dto *QuestionDTO) ConvertEntityToResponse(ent *entity.Question) *response.QuestionResponse {
	return &response.QuestionResponse{
		ID:                 ent.ID,
		TemplateQuestionID: ent.TemplateQuestionID,
		AnswerTypeID:       ent.AnswerTypeID,
		Name:               ent.Name,
		CreatedAt:          ent.CreatedAt,
		UpdatedAt:          ent.UpdatedAt,
		AnswerTypeResponse: func() *response.AnswerTypeResponse {
			if ent.AnswerType == nil {
				return nil
			}
			return dto.AnswerTypeDTO.ConvertEntityToResponse(ent.AnswerType)
		}(),
		QuestionOptions: func() *[]response.QuestionOptionResponse {
			var questionOptionResponses []response.QuestionOptionResponse
			if ent.QuestionOptions == nil {
				return nil
			}
			for _, questionOption := range ent.QuestionOptions {
				questionOptionResponses = append(questionOptionResponses, *dto.QuestionOptionDTO.ConvertEntityToResponse(&questionOption))
			}
			return &questionOptionResponses
		}(),
	}
}
