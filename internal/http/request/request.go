package request

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/go-playground/validator/v10"
)

func FormTypeValidation(fl validator.FieldLevel) bool {
	formType := fl.Field().String()
	if formType == "" {
		return true
	}
	switch entity.TemplateQuestionFormType(formType) {
	case entity.TQ_FORM_TYPE_ADMINISTRATIVE_SELECTION,
		entity.TQ_FORM_TYPE_TEST,
		entity.TQ_FORM_TYPE_INTERVIEW,
		entity.TQ_FORM_TYPE_FGD,
		entity.TQ_FORM_TYPE_FINAL_INTERVIEW,
		entity.TQ_FORM_TYPE_OFFERING_LETTER,
		entity.TQ_FORM_TYPE_CONTRACT_DOCUMENT,
		entity.TQ_FORM_TYPE_DOCUMENT_CHECKING:
		return true
	default:
		return false
	}
}

func TemplateQuestionStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.TemplateQuestionStatus(status) {
	case entity.TEMPLATE_QUESTION_STATUS_ACTIVE,
		entity.TEMPLATE_QUESTION_STATUS_INACTIVE:
		return true
	default:
		return false
	}
}

func RecruitmentTypeValidation(fl validator.FieldLevel) bool {
	recruitmentType := fl.Field().String()
	if recruitmentType == "" {
		return true
	}
	switch entity.ProjectRecruitmentType(recruitmentType) {
	case entity.PROJECT_RECRUITMENT_TYPE_MT,
		entity.PROJECT_RECRUITMENT_TYPE_PH,
		entity.PROJECT_RECRUITMENT_TYPE_NS:
		return true
	default:
		return false
	}
}

func TemplateActivityStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.TemplateActivityStatus(status) {
	case entity.TEMPLATE_ACTIVITY_STATUS_ACTIVE,
		entity.TEMPLATE_ACTIVITY_STATUS_INACTIVE:
		return true
	default:
		return false
	}
}
