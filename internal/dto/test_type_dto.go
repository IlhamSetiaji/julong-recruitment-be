package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type ITestTypeDTO interface {
	ConvertEntityToResponse(ent *entity.TestType) *response.TestTypeResponse
}

type TestTypeDTO struct {
	Log *logrus.Logger
}

func NewTestTypeDTO(log *logrus.Logger) ITestTypeDTO {
	return &TestTypeDTO{
		Log: log,
	}
}

func TestTypeDTOFactory(log *logrus.Logger) ITestTypeDTO {
	return NewTestTypeDTO(log)
}

func (dto *TestTypeDTO) ConvertEntityToResponse(ent *entity.TestType) *response.TestTypeResponse {
	return &response.TestTypeResponse{
		ID:              ent.ID,
		Name:            ent.Name,
		RecruitmentType: ent.RecruitmentType,
		Status:          ent.Status,
		CreatedAt:       ent.CreatedAt,
		UpdatedAt:       ent.UpdatedAt,
	}
}
