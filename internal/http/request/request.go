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

func TemplateActivityLineStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.TemplateActivityLineStatus(status) {
	case entity.TEMPLATE_ACTIVITY_LINE_STATUS_ACTIVE,
		entity.TEMPLATE_ACTIVITY_LINE_STATUS_INACTIVE:
		return true
	default:
		return false
	}
}

func ProjectRecruitmentHeaderStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.ProjectRecruitmentHeaderStatus(status) {
	case entity.PROJECT_RECRUITMENT_HEADER_STATUS_DRAFT,
		entity.PROJECT_RECRUITMENT_HEADER_STATUS_SUBMITTED,
		entity.PROJECT_RECRUITMENT_HEADER_STATUS_APPROVED,
		entity.PROJECT_RECRUITMENT_HEADER_STATUS_REJECTED,
		entity.PROJECT_RECRUITMENT_HEADER_STATUS_CLOSE,
		entity.PROJECT_RECRUITMENT_HEADER_STATUS_IN_PROGRESS,
		entity.PROJECT_RECRUITMENT_HEADER_STATUS_COMPLETED:
		return true
	default:
		return false
	}
}

func JobPostingStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true
	}
	switch entity.JobPostingStatus(status) {
	case entity.JOB_POSTING_STATUS_DRAFT,
		entity.JOB_POSTING_STATUS_SUBMITTED,
		entity.JOB_POSTING_STATUS_APPROVED,
		entity.JOB_POSTING_STATUS_REJECTED,
		entity.JOB_POSTING_STATUS_CLOSE,
		entity.JOB_POSTING_STATUS_IN_PROGRESS,
		entity.JOB_POSTING_STATUS_COMPLETED:
		return true
	default:
		return false
	}
}
