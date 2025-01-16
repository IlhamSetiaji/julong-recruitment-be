package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type ITemplateActivityDTO interface {
	ConvertEntityToResponse(ent *entity.TemplateActivity) *response.TemplateActivityResponse
}

type TemplateActivityDTO struct {
	Log                     *logrus.Logger
	TemplateActivityLineDTO ITemplateActivityLineDTO
}

func NewTemplateActivityDTO(
	log *logrus.Logger,
	templateActivityLineDTO ITemplateActivityLineDTO,
) ITemplateActivityDTO {
	return &TemplateActivityDTO{
		Log:                     log,
		TemplateActivityLineDTO: templateActivityLineDTO,
	}
}

func TemplateActivityDTOFactory(log *logrus.Logger) ITemplateActivityDTO {
	templateActivityLineDTO := TemplateActivityLineDTOFactory(log)
	return NewTemplateActivityDTO(log, templateActivityLineDTO)
}

func (dto *TemplateActivityDTO) ConvertEntityToResponse(ent *entity.TemplateActivity) *response.TemplateActivityResponse {
	return &response.TemplateActivityResponse{
		ID:              ent.ID,
		Name:            ent.Name,
		Description:     ent.Description,
		RecruitmentType: ent.RecruitmentType,
		Status:          ent.Status,
		CreatedAt:       ent.CreatedAt,
		UpdatedAt:       ent.UpdatedAt,
		TemplateActivityLines: func() *[]response.TemplateActivityLineResponse {
			var templateActivityLineResponses []response.TemplateActivityLineResponse
			if ent.TemplateActivityLines == nil {
				return nil
			}
			for _, templateActivityLine := range ent.TemplateActivityLines {
				templateActivityLineResponses = append(templateActivityLineResponses, *dto.TemplateActivityLineDTO.ConvertEntityToResponse(&templateActivityLine))
			}
			return &templateActivityLineResponses
		}(),
	}
}
