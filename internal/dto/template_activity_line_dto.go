package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type ITemplateActivityLineDTO interface {
	ConvertEntityToResponse(ent *entity.TemplateActivityLine) *response.TemplateActivityLineResponse
}

type TemplateActivityLineDTO struct {
	Log                 *logrus.Logger
	TemplateQuestionDTO ITemplateQuestionDTO
}

func NewTemplateActivityLineDTO(
	log *logrus.Logger,
	templateQuestionDTO ITemplateQuestionDTO,
) ITemplateActivityLineDTO {
	return &TemplateActivityLineDTO{
		Log:                 log,
		TemplateQuestionDTO: templateQuestionDTO,
	}
}

func TemplateActivityLineDTOFactory(log *logrus.Logger) ITemplateActivityLineDTO {
	templateQuestionDTO := TemplateQuestionDTOFactory(log)
	return NewTemplateActivityLineDTO(log, templateQuestionDTO)
}

func (dto *TemplateActivityLineDTO) ConvertEntityToResponse(ent *entity.TemplateActivityLine) *response.TemplateActivityLineResponse {
	return &response.TemplateActivityLineResponse{
		ID:                 ent.ID,
		TemplateActivityID: ent.TemplateActivityID,
		TemplateQuestionID: ent.QuestionTemplateID,
		Description:        ent.Description,
		Status:             ent.Status,
		ColorHexCode:       ent.ColorHexCode,
		CreatedAt:          ent.CreatedAt,
		UpdatedAt:          ent.UpdatedAt,
		TemplateQuestion: func() *response.TemplateQuestionResponse {
			if ent.TemplateQuestion == nil {
				return nil
			}
			return dto.TemplateQuestionDTO.ConvertEntityToResponse(ent.TemplateQuestion)
		}(),
	}
}
