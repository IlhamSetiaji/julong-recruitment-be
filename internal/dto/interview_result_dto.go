package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IInterviewResultDTO interface {
	ConvertEntityToResponse(ent *entity.InterviewResult) *response.InterviewResultResponse
}

type InterviewResultDTO struct {
	Log                  *logrus.Logger
	InterviewAssessorDTO IInterviewAssessorDTO
}

func NewInterviewResultDTO(
	log *logrus.Logger,
	iaDTO IInterviewAssessorDTO,
) IInterviewResultDTO {
	return &InterviewResultDTO{
		Log:                  log,
		InterviewAssessorDTO: iaDTO,
	}
}

func InterviewResultDTOFactory(log *logrus.Logger) IInterviewResultDTO {
	iaDTO := InterviewAssessorDTOFactory(log)
	return NewInterviewResultDTO(log, iaDTO)
}

func (dto *InterviewResultDTO) ConvertEntityToResponse(ent *entity.InterviewResult) *response.InterviewResultResponse {
	return &response.InterviewResultResponse{
		ID:                   ent.ID,
		InterviewApplicantID: ent.InterviewApplicantID,
		InterviewAssessorID:  ent.InterviewAssessorID,
		Status:               ent.Status,
		CreatedAt:            ent.CreatedAt,
		UpdatedAt:            ent.UpdatedAt,
		InterviewAssessor: func() *response.InterviewAssessorResponse {
			if ent.InterviewAssessor == nil {
				return nil
			}
			return dto.InterviewAssessorDTO.ConvertEntityToResponse(ent.InterviewAssessor)
		}(),
	}
}
