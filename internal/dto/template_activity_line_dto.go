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
	TemplateActivityDTO ITemplateActivityDTO
}

func NewTemplateActivityLineDTO(
	log *logrus.Logger,
	templateQuestionDTO ITemplateQuestionDTO,
	TemplateActivityDTO ITemplateActivityDTO,
) ITemplateActivityLineDTO {
	return &TemplateActivityLineDTO{
		Log:                 log,
		TemplateQuestionDTO: templateQuestionDTO,
		TemplateActivityDTO: TemplateActivityDTO,
	}
}

func TemplateActivityLineDTOFactory(log *logrus.Logger) ITemplateActivityLineDTO {
	templateQuestionDTO := TemplateQuestionDTOFactory(log)
	TemplateActivityDTO := TemplateActivityDTOFactory(log)
	return NewTemplateActivityLineDTO(log, templateQuestionDTO, TemplateActivityDTO)
}

func (dto *TemplateActivityLineDTO) ConvertEntityToResponse(ent *entity.TemplateActivityLine) *response.TemplateActivityLineResponse {
	return &response.TemplateActivityLineResponse{
		ID:                 ent.ID,
		TemplateActivityID: ent.TemplateActivityID,
		TemplateQuestionID: ent.QuestionTemplateID,
		Name:               ent.Name,
		Description:        ent.Description,
		Status:             ent.Status,
		ColorHexCode:       ent.ColorHexCode,
		CreatedAt:          ent.CreatedAt,
		UpdatedAt:          ent.UpdatedAt,
		TemplateActivity: func() *response.TemplateActivityResponse {
			if ent.TemplateActivity == nil {
				return nil
			}
			return dto.TemplateActivityDTO.ConvertEntityToResponse(ent.TemplateActivity)
		}(),
		TemplateQuestion: func() *response.TemplateQuestionResponse {
			if ent.TemplateQuestion == nil {
				return nil
			}
			return dto.TemplateQuestionDTO.ConvertEntityToResponse(ent.TemplateQuestion)
		}(),
	}
}
