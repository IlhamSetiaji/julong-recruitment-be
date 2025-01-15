package request

type CreateOrUpdateQuestions struct {
	TemplateQuestionID string `json:"template_question_id" validate:"required,uuid"`
	Questions          []struct {
		ID              string `json:"id" validate:"omitempty,uuid"`
		AnswerTypeID    string `json:"answer_type_id" validate:"required,uuid"`
		Name            string `json:"name" validate:"required"`
		QuestionOptions []struct {
			OptionText string `json:"option_text" validate:"required"`
		} `json:"question_options" validate:"omitempty,dive"`
	} `json:"questions" validate:"required,dive"`
	DeletedQuestionIDs []string `json:"deleted_question_ids" validate:"omitempty,dive,uuid"`
}
