package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IQuestionUseCase interface {
	CreateOrUpdateQuestions(req *request.CreateOrUpdateQuestions) error
}

type QuestionUseCase struct {
	Log                        *logrus.Logger
	Repository                 repository.IQuestionRepository
	DTO                        dto.IQuestionDTO
	QuestionOptionRepository   repository.IQuestionOptionRepository
	TemplateQuestionRepository repository.ITemplateQuestionRepository
}

func NewQuestionUseCase(
	log *logrus.Logger,
	repo repository.IQuestionRepository,
	qDTO dto.IQuestionDTO,
	qoRepository repository.IQuestionOptionRepository,
	tqRepository repository.ITemplateQuestionRepository,
) IQuestionUseCase {
	return &QuestionUseCase{
		Log:                        log,
		Repository:                 repo,
		DTO:                        qDTO,
		QuestionOptionRepository:   qoRepository,
		TemplateQuestionRepository: tqRepository,
	}
}

func QuestionUseCaseFactory(log *logrus.Logger) IQuestionUseCase {
	repo := repository.QuestionRepositoryFactory(log)
	qDTO := dto.QuestionDTOFactory(log)
	qoRepository := repository.QuestionOptionRepositoryFactory(log)
	tqRepository := repository.TemplateQuestionRepositoryFactory(log)
	return NewQuestionUseCase(log, repo, qDTO, qoRepository, tqRepository)
}

func (uc *QuestionUseCase) CreateOrUpdateQuestions(req *request.CreateOrUpdateQuestions) error {
	// check if template question exist
	tq, err := uc.TemplateQuestionRepository.FindByID(uuid.MustParse(req.TemplateQuestionID))
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when finding template question by id: %s", err.Error())
		return errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when finding template question by id: " + err.Error())
	}

	if tq == nil {
		uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] template question with id %s not found", req.TemplateQuestionID)
		return errors.New("[QuestionUseCase.CreateOrUpdateQuestions] template question with id " + req.TemplateQuestionID + " not found")
	}

	// create or update questions
	for _, question := range req.Questions {
		exist, err := uc.Repository.FindByID(uuid.MustParse(question.ID))
		if err != nil {
			uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when finding question by id: %s", err.Error())
			return errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when finding question by id: " + err.Error())
		}

		if exist == nil {
			createdQuestion, err := uc.Repository.CreateQuestion(&entity.Question{
				TemplateQuestionID: tq.ID,
				AnswerTypeID:       uuid.MustParse(question.AnswerTypeID),
				Name:               question.Name,
			})
			if err != nil {
				uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question: %s", err.Error())
				return errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question: " + err.Error())
			}

			if len(question.QuestionOptions) > 0 {
				for _, questionOption := range question.QuestionOptions {
					_, err := uc.QuestionOptionRepository.CreateQuestionOption(&entity.QuestionOption{
						QuestionID: createdQuestion.ID,
						OptionText: questionOption.OptionText,
					})
					if err != nil {
						uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: %s", err.Error())
						return errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: " + err.Error())
					}
				}
			}
		} else {
			updatedQuestion, err := uc.Repository.UpdateQuestion(&entity.Question{
				ID:                 exist.ID,
				TemplateQuestionID: tq.ID,
				AnswerTypeID:       uuid.MustParse(question.AnswerTypeID),
				Name:               question.Name,
			})
			if err != nil {
				uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when updating question: %s", err.Error())
				return errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when updating question: " + err.Error())
			}

			// delete question options
			err = uc.QuestionOptionRepository.DeleteQuestionOptionsByQuestionID(updatedQuestion.ID)
			if err != nil {
				uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question options: %s", err.Error())
				return errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question options: " + err.Error())
			}

			// create question options
			if len(question.QuestionOptions) > 0 {
				for _, questionOption := range question.QuestionOptions {
					_, err := uc.QuestionOptionRepository.CreateQuestionOption(&entity.QuestionOption{
						QuestionID: updatedQuestion.ID,
						OptionText: questionOption.OptionText,
					})
					if err != nil {
						uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: %s", err.Error())
						return errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: " + err.Error())
					}
				}
			}
		}
	}

	// delete questions
	if len(req.DeletedQuestionIDs) > 0 {
		for _, id := range req.DeletedQuestionIDs {
			err := uc.Repository.DeleteQuestion(uuid.MustParse(id))
			if err != nil {
				uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question: %s", err.Error())
				return errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question: " + err.Error())
			}
		}
	}

	return nil
}
