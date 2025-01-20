package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IEducationDTO interface {
	ConvertEntityToResponse(ent *entity.Education) *response.EducationResponse
}

type EducationDTO struct {
	Log *logrus.Logger
}

func NewEducationDTO(log *logrus.Logger) IEducationDTO {
	return &EducationDTO{
		Log: log,
	}
}

func EducationDTOFactory(log *logrus.Logger) IEducationDTO {
	return NewEducationDTO(log)
}

func (dto *EducationDTO) ConvertEntityToResponse(ent *entity.Education) *response.EducationResponse {
	return &response.EducationResponse{
		ID:             ent.ID,
		UserProfileID:  ent.UserProfileID,
		EducationLevel: ent.EducationLevel,
		Major:          ent.Major,
		SchoolName:     ent.SchoolName,
		GraduateYear:   ent.GraduateYear,
		EndDate:        ent.EndDate,
		Certificate:    ent.Certificate,
		Gpa:            ent.Gpa,
		CreatedAt:      ent.CreatedAt,
		UpdatedAt:      ent.UpdatedAt,
	}
}
