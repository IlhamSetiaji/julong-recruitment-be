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
	validate.RegisterValidation("template_activity_line_status_validation", request.TemplateActivityLineStatusValidation)
	validate.RegisterValidation("project_recruitment_header_status_validation", request.ProjectRecruitmentHeaderStatusValidation)
	validate.RegisterValidation("job_posting_status_validation", request.JobPostingStatusValidation)
	validate.RegisterValidation("marital_status_validation", request.MaritalStatusValidation)
	validate.RegisterValidation("user_status_validation", request.UserStatusValidation)
	validate.RegisterValidation("user_gender_validation", request.UserGenderValidation)
	validate.RegisterValidation("education_level_validation", request.EducationLevelValidation)
	validate.RegisterValidation("test_type_status_validation", request.TestTypeStatusValidation)
	validate.RegisterValidation("test_schedule_header_status_validation", request.TestScheduleHeaderStatusValidation)
	validate.RegisterValidation("test_applicant_status_validation", request.TestApplicantStatusValidation)
	validate.RegisterValidation("administrative_selection_status_validation", request.AdministrativeSelectionStatusValidation)
	validate.RegisterValidation("administrative_result_status_validation", request.AdministrativeResultStatusValidation)
	validate.RegisterValidation("final_result_status_validation", request.FinalResultStatusValidation)
	validate.RegisterValidation("assessment_status_validation", request.AssessmentStatusValidation)
	validate.RegisterValidation("interview_status_validation", request.InterviewStatusValidation)
	validate.RegisterValidation("interview_result_status_validation", request.InterviewResultStatusValidation)
	validate.RegisterValidation("fgd_result_status_validation", request.FgdResultStatusValidation)
	validate.RegisterValidation("fgd_schedule_status_validation", request.FgdScheduleStatusValidation)
	validate.RegisterValidation("document_sending_status_validation", request.DocumentSendingStatusValidation)
	validate.RegisterValidation("document_agreement_status_validation", request.DocumentAgreementStatusValidation)
	validate.RegisterValidation("document_verification_header_status_validation", request.DocumentVerificationHeaderStatusValidation)
	validate.RegisterValidation("document_sending_hired_status_validation", request.DocumentSendingHiredStatusValidation)
	return validate
}
