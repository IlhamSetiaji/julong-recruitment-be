package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IFgdResultDTO interface {
	ConvertEntityToResponse(ent *entity.FgdResult) *response.FgdResultResponse
}

type FgdResultDTO struct {
	Log            *logrus.Logger
	FgdAssessorDTO IFgdAssessorDTO
}

func NewFgdResultDTO(
	log *logrus.Logger,
	faDTO IFgdAssessorDTO,
) IFgdResultDTO {
	return &FgdResultDTO{
		Log:            log,
		FgdAssessorDTO: faDTO,
	}
}

func FgdResultDTOFactory(log *logrus.Logger) IFgdResultDTO {
	faDTO := FgdAssessorDTOFactory(log)
	return NewFgdResultDTO(log, faDTO)
}

func (dto *FgdResultDTO) ConvertEntityToResponse(ent *entity.FgdResult) *response.FgdResultResponse {
	return &response.FgdResultResponse{
		ID:             ent.ID,
		FgdApplicantID: ent.FgdApplicantID,
		FgdAssessorID:  ent.FgdAssessorID,
		Status:         ent.Status,
		CreatedAt:      ent.CreatedAt,
		UpdatedAt:      ent.UpdatedAt,
		FgdAssessor: func() *response.FgdAssessorResponse {
			if ent.FgdAssessor == nil {
				return nil
			}
			return dto.FgdAssessorDTO.ConvertEntityToResponse(ent.FgdAssessor)
		}(),
	}
}
