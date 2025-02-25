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
	TemplateActivityDTO ITemplateActivityDTO
	TemplateQuestionDTO ITemplateQuestionDTO
}

func NewTemplateActivityLineDTO(
	log *logrus.Logger,
	TemplateActivityDTO ITemplateActivityDTO,
	templateQuestionDTO ITemplateQuestionDTO,
) ITemplateActivityLineDTO {
	return &TemplateActivityLineDTO{
		Log:                 log,
		TemplateActivityDTO: TemplateActivityDTO,
		TemplateQuestionDTO: templateQuestionDTO,
	}
}

func TemplateActivityLineDTOFactory(log *logrus.Logger) ITemplateActivityLineDTO {
	TemplateActivityDTO := TemplateActivityDTOFactory(log)
	templateQuestionDTO := TemplateQuestionDTOFactory(log)
	return NewTemplateActivityLineDTO(log, TemplateActivityDTO, templateQuestionDTO)
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
