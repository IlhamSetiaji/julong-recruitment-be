package response

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type TemplateQuestionResponse struct {
	ID              uuid.UUID                       `json:"id"`
	DocumentSetupID *uuid.UUID                      `json:"document_setup_id"`
	Name            string                          `json:"name"`
	FormType        entity.TemplateQuestionFormType `json:"form_type"`
	Description     string                          `json:"description"`
	Duration        *int                            `json:"duration"`
	Status          entity.TemplateQuestionStatus   `json:"status"`

	Questions             *[]QuestionResponse             `json:"questions"`
	DocumentSetup         *DocumentSetupResponse          `json:"document_setup"`
	DocumentVerifications *[]DocumentVerificationResponse `json:"document_verifications"`
}

type FormTypeResponse struct {
	Value string `json:"value"`
}
