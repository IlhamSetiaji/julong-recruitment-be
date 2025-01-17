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
	Log                       *logrus.Logger
	TemplateActivityDTO       ITemplateActivityDTO
	ProjectRecruitmentLineDTO IProjectRecruitmentLineDTO
}

func NewProjectRecruitmentHeaderDTO(
	log *logrus.Logger,
	templateActivityDTO ITemplateActivityDTO,
	prl IProjectRecruitmentLineDTO,
) IProjectRecruitmentHeaderDTO {
	return &ProjectRecruitmentHeaderDTO{
		Log:                       log,
		TemplateActivityDTO:       templateActivityDTO,
		ProjectRecruitmentLineDTO: prl,
	}
}

func ProjectRecruitmentHeaderDTOFactory(log *logrus.Logger) IProjectRecruitmentHeaderDTO {
	templateActivityDTO := TemplateActivityDTOFactory(log)
	prlDTO := ProjectRecruitmentLineDTOFactory(log)
	return NewProjectRecruitmentHeaderDTO(log, templateActivityDTO, prlDTO)
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
		ProjectRecruitmentLines: func() []response.ProjectRecruitmentLineResponse {
			var projectRecruitmentLineResponses []response.ProjectRecruitmentLineResponse
			if ent.ProjectRecruitmentLines == nil || len(ent.ProjectRecruitmentLines) == 0 {
				return nil
			}
			for _, projectRecruitmentLine := range ent.ProjectRecruitmentLines {
				projectRecruitmentLineResponses = append(projectRecruitmentLineResponses, *dto.ProjectRecruitmentLineDTO.ConvertEntityToResponse(&projectRecruitmentLine))
			}
			return projectRecruitmentLineResponses
		}(),
	}
}
