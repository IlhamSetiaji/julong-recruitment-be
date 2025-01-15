package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ITemplateQuestionDTO interface {
	ConvertEntityToResponse(ent *entity.TemplateQuestion) *response.TemplateQuestionResponse
}

type TemplateQuestionDTO struct {
	Log              *logrus.Logger
	QuestionDTO      IQuestionDTO
	DocumentSetupDTO IDocumentSetupDTO
}

func NewTemplateQuestionDTO(log *logrus.Logger, questionDTO IQuestionDTO, documentSetupDTO IDocumentSetupDTO) ITemplateQuestionDTO {
	return &TemplateQuestionDTO{
		Log:              log,
		QuestionDTO:      questionDTO,
		DocumentSetupDTO: documentSetupDTO,
	}
}

func TemplateQuestionDTOFactory(log *logrus.Logger) ITemplateQuestionDTO {
	questionDTO := QuestionDTOFactory(log)
	documentSetupDTO := DocumentSetupDTOFactory(log)
	return NewTemplateQuestionDTO(log, questionDTO, documentSetupDTO)
}

func (dto *TemplateQuestionDTO) ConvertEntityToResponse(ent *entity.TemplateQuestion) *response.TemplateQuestionResponse {
	return &response.TemplateQuestionResponse{
		ID: ent.ID,
		DocumentSetupID: func() *uuid.UUID {
			if ent.DocumentSetupID == nil {
				return nil
			}
			return ent.DocumentSetupID
		}(),
		Name:        ent.Name,
		FormType:    entity.TemplateQuestionFormType(ent.FormType),
		Description: ent.Description,
		Duration:    ent.Duration,
		Status:      ent.Status,
		Questions: func() *[]response.QuestionResponse {
			var questionResponses []response.QuestionResponse
			if len(ent.Questions) == 0 || ent.Questions == nil {
				return nil
			}
			for _, question := range ent.Questions {
				questionResponses = append(questionResponses, *dto.QuestionDTO.ConvertEntityToResponse(&question))
			}
			return &questionResponses
		}(),
		DocumentSetup: func() *response.DocumentSetupResponse {
			if ent.DocumentSetup == nil {
				return nil
			}
			return dto.DocumentSetupDTO.ConvertEntityToResponse(ent.DocumentSetup)
		}(),
	}
}
