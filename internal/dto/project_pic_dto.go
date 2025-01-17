package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IProjectPicDTO interface {
	ConvertEntityToResponse(ent *entity.ProjectPic) *response.ProjectPicResponse
}

type ProjectPicDTO struct {
	Log *logrus.Logger
}

func NewProjectPicDTO(log *logrus.Logger) IProjectPicDTO {
	return &ProjectPicDTO{
		Log: log,
	}
}

func ProjectPicDTOFactory(log *logrus.Logger) IProjectPicDTO {
	return NewProjectPicDTO(log)
}

func (dto *ProjectPicDTO) ConvertEntityToResponse(ent *entity.ProjectPic) *response.ProjectPicResponse {
	return &response.ProjectPicResponse{
		ID:                       ent.ID,
		ProjectRecruitmentLineID: ent.ProjectRecruitmentLineID,
		EmployeeID:               ent.EmployeeID,
		EmployeeName:             "Halo",
		AdministrativeTotal:      ent.AdministrativeTotal,
		CreatedAt:                ent.CreatedAt,
		UpdatedAt:                ent.UpdatedAt,
	}
}
