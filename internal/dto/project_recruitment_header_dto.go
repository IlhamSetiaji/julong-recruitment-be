package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IProjectRecruitmentHeaderDTO interface {
	ConvertEntityToResponse(ent *entity.ProjectRecruitmentHeader) *response.ProjectRecruitmentHeaderResponse
}

type ProjectRecruitmentHeaderDTO struct {
	Log                 *logrus.Logger
	TemplateActivityDTO ITemplateActivityDTO
}

func NewProjectRecruitmentHeaderDTO(
	log *logrus.Logger,
	templateActivityDTO ITemplateActivityDTO,
) IProjectRecruitmentHeaderDTO {
	return &ProjectRecruitmentHeaderDTO{
		Log:                 log,
		TemplateActivityDTO: templateActivityDTO,
	}
}

func ProjectRecruitmentHeaderDTOFactory(log *logrus.Logger) IProjectRecruitmentHeaderDTO {
	templateActivityDTO := TemplateActivityDTOFactory(log)
	return NewProjectRecruitmentHeaderDTO(log, templateActivityDTO)
}

func (dto *ProjectRecruitmentHeaderDTO) ConvertEntityToResponse(ent *entity.ProjectRecruitmentHeader) *response.ProjectRecruitmentHeaderResponse {
	return &response.ProjectRecruitmentHeaderResponse{
		ID:                 ent.ID,
		TemplateActivityID: ent.TemplateActivityID,
		Name:               ent.Name,
		Description:        ent.Description,
		DocumentDate:       ent.DocumentDate,
		DocumentNumber:     ent.DocumentNumber,
		RecruitmentType:    entity.ProjectRecruitmentType(ent.RecruitmentType),
		StartDate:          ent.StartDate,
		EndDate:            ent.EndDate,
		Status:             entity.ProjectRecruitmentHeaderStatus(ent.Status),
		TemplateActivity: func() *response.TemplateActivityResponse {
			if ent.TemplateActivity == nil {
				return nil
			}
			return dto.TemplateActivityDTO.ConvertEntityToResponse(ent.TemplateActivity)
		}(),
	}
}
