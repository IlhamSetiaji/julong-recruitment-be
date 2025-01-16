package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IDocumentVerificationDTO interface {
	ConvertEntityToResponse(ent *entity.DocumentVerification) *response.DocumentVerificationResponse
}

type DocumentVerificationDTO struct {
	Log *logrus.Logger
}

func NewDocumentVerificationDTO(log *logrus.Logger) IDocumentVerificationDTO {
	return &DocumentVerificationDTO{
		Log: log,
	}
}

func DocumentVerificationDTOFactory(log *logrus.Logger) IDocumentVerificationDTO {
	return NewDocumentVerificationDTO(log)
}

func (dto *DocumentVerificationDTO) ConvertEntityToResponse(ent *entity.DocumentVerification) *response.DocumentVerificationResponse {
	return &response.DocumentVerificationResponse{
		ID:                 ent.ID,
		TemplateQuestionID: ent.TemplateQuestionID,
		Name:               ent.Name,
		Format:             ent.Format,
		TemplateQuestion: func() *response.TemplateQuestionResponse {
			if ent.TemplateQuestion == nil {
				return nil
			}
			return &response.TemplateQuestionResponse{
				ID: ent.TemplateQuestion.ID,
				DocumentSetupID: func() *uuid.UUID {
					if ent.TemplateQuestion.DocumentSetupID == nil {
						return nil
					}
					return ent.TemplateQuestion.DocumentSetupID
				}(),
				Name:        ent.TemplateQuestion.Name,
				FormType:    entity.TemplateQuestionFormType(ent.TemplateQuestion.FormType),
				Description: ent.TemplateQuestion.Description,
				Duration:    ent.TemplateQuestion.Duration,
				Status:      ent.TemplateQuestion.Status,
			}
		}(),
	}
}
