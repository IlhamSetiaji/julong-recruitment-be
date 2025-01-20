package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type ISkillDTO interface {
	ConvertEntityToResponse(ent *entity.Skill) *response.SkillResponse
}

type SkillDTO struct {
	Log *logrus.Logger
}

func NewSkillDTO(log *logrus.Logger) ISkillDTO {
	return &SkillDTO{
		Log: log,
	}
}

func SkillDTOFactory(log *logrus.Logger) ISkillDTO {
	return NewSkillDTO(log)
}

func (dto *SkillDTO) ConvertEntityToResponse(ent *entity.Skill) *response.SkillResponse {
	return &response.SkillResponse{
		ID:            ent.ID,
		UserProfileID: ent.UserProfileID,
		Name:          ent.Name,
		Description:   ent.Description,
		Certificate:   ent.Certificate,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
	}
}
