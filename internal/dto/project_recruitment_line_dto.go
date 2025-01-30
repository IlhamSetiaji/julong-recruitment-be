package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IProjectRecruitmentLineDTO interface {
	ConvertEntityToResponse(ent *entity.ProjectRecruitmentLine) *response.ProjectRecruitmentLineResponse
}

type ProjectRecruitmentLineDTO struct {
	Log                     *logrus.Logger
	ProjectPicDTO           IProjectPicDTO
	TemplateActivityLineDTO ITemplateActivityLineDTO
}

func NewProjectRecruitmentLineDTO(
	log *logrus.Logger,
	projectPicDTO IProjectPicDTO,
	taDTO ITemplateActivityLineDTO,
) IProjectRecruitmentLineDTO {
	return &ProjectRecruitmentLineDTO{
		Log:                     log,
		ProjectPicDTO:           projectPicDTO,
		TemplateActivityLineDTO: taDTO,
	}
}

func ProjectRecruitmentLineDTOFactory(log *logrus.Logger) IProjectRecruitmentLineDTO {
	projectPicDTO := ProjectPicDTOFactory(log)
	taDTO := TemplateActivityLineDTOFactory(log)
	return NewProjectRecruitmentLineDTO(log, projectPicDTO, taDTO)
}

func (dto *ProjectRecruitmentLineDTO) ConvertEntityToResponse(ent *entity.ProjectRecruitmentLine) *response.ProjectRecruitmentLineResponse {
	return &response.ProjectRecruitmentLineResponse{
		ID:                         ent.ID,
		ProjectRecruitmentHeaderID: ent.ProjectRecruitmentHeaderID,
		TemplateActivityLineID:     ent.TemplateActivityLineID,
		StartDate:                  ent.StartDate,
		EndDate:                    ent.EndDate,
		Order:                      ent.Order,
		CreatedAt:                  ent.CreatedAt,
		UpdatedAt:                  ent.UpdatedAt,
		ProjectPics: func() []response.ProjectPicResponse {
			var projectPicResponses []response.ProjectPicResponse
			if ent.ProjectPics == nil {
				return nil
			}
			for _, projectPic := range ent.ProjectPics {
				projectPicResponses = append(projectPicResponses, *dto.ProjectPicDTO.ConvertEntityToResponse(&projectPic))
			}
			return projectPicResponses
		}(),
		TemplateActivityLine: func() *response.TemplateActivityLineResponse {
			if ent.TemplateActivityLine == nil {
				return nil
			}
			return dto.TemplateActivityLineDTO.ConvertEntityToResponse(ent.TemplateActivityLine)
		}(),
	}
}
