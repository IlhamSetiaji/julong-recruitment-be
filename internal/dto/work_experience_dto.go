package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IWorkExperienceDTO interface {
	ConvertEntityToResponse(ent *entity.WorkExperience) *response.WorkExperienceResponse
}

type WorkExperienceDTO struct {
	Log *logrus.Logger
}

func NewWorkExperienceDTO(log *logrus.Logger) IWorkExperienceDTO {
	return &WorkExperienceDTO{
		Log: log,
	}
}

func WorkExperienceDTOFactory(log *logrus.Logger) IWorkExperienceDTO {
	return NewWorkExperienceDTO(log)
}

func (dto *WorkExperienceDTO) ConvertEntityToResponse(ent *entity.WorkExperience) *response.WorkExperienceResponse {
	return &response.WorkExperienceResponse{
		ID:             ent.ID,
		UserProfileID:  ent.UserProfileID,
		Name:           ent.Name,
		CompanyName:    ent.CompanyName,
		YearExperience: ent.YearExperience,
		JobDescription: ent.JobDescription,
		Certificate:    ent.Certificate,
		CreatedAt:      ent.CreatedAt,
		UpdatedAt:      ent.UpdatedAt,
	}
}
