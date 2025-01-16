package config

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("form_type_validation", request.FormTypeValidation)
	validate.RegisterValidation("template_question_status_validation", request.TemplateQuestionStatusValidation)
	validate.RegisterValidation("recruitment_type_validation", request.RecruitmentTypeValidation)
	validate.RegisterValidation("template_activity_status_validation", request.TemplateActivityStatusValidation)
	return validate
}
