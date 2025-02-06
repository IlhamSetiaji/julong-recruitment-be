package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IQuestionUseCase interface {
	CreateOrUpdateQuestions(req *request.CreateOrUpdateQuestions) (*response.TemplateQuestionResponse, error)
	FindByIDAndUserID(questionID, userID string) (*response.QuestionResponse, error)
	FindAllByProjectRecruitmentLineIDAndJobPostingID(projectRecruitmentLineID uuid.UUID, jobPostingID uuid.UUID) (*[]response.QuestionResponse, error)
	FindByID(questionID string) (*entity.Question, error)
}

type QuestionUseCase struct {
	Log                              *logrus.Logger
	Repository                       repository.IQuestionRepository
	DTO                              dto.IQuestionDTO
	QuestionOptionRepository         repository.IQuestionOptionRepository
	TemplateQuestionRepository       repository.ITemplateQuestionRepository
	TemplateQuestionDTO              dto.ITemplateQuestionDTO
	UserProfileRepository            repository.IUserProfileRepository
	ProjectRecruitmentLineRepository repository.IProjectRecruitmentLineRepository
	TemplateActivityLineRepository   repository.ITemplateActivityLineRepository
}

func NewQuestionUseCase(
	log *logrus.Logger,
	repo repository.IQuestionRepository,
	qDTO dto.IQuestionDTO,
	qoRepository repository.IQuestionOptionRepository,
	tqRepository repository.ITemplateQuestionRepository,
	tqDTO dto.ITemplateQuestionDTO,
	userProfileRepository repository.IUserProfileRepository,
	projectRecruitmentLineRepository repository.IProjectRecruitmentLineRepository,
	talRepo repository.ITemplateActivityLineRepository,
) IQuestionUseCase {
	return &QuestionUseCase{
		Log:                              log,
		Repository:                       repo,
		DTO:                              qDTO,
		QuestionOptionRepository:         qoRepository,
		TemplateQuestionRepository:       tqRepository,
		TemplateQuestionDTO:              tqDTO,
		UserProfileRepository:            userProfileRepository,
		ProjectRecruitmentLineRepository: projectRecruitmentLineRepository,
		TemplateActivityLineRepository:   talRepo,
	}
}

func QuestionUseCaseFactory(log *logrus.Logger) IQuestionUseCase {
	repo := repository.QuestionRepositoryFactory(log)
	qDTO := dto.QuestionDTOFactory(log)
	qoRepository := repository.QuestionOptionRepositoryFactory(log)
	tqRepository := repository.TemplateQuestionRepositoryFactory(log)
	tqDTO := dto.TemplateQuestionDTOFactory(log)
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	projectRecruitmentLineRepository := repository.ProjectRecruitmentLineRepositoryFactory(log)
	talRepo := repository.TemplateActivityLineRepositoryFactory(log)
	return NewQuestionUseCase(log, repo, qDTO, qoRepository, tqRepository, tqDTO, userProfileRepository, projectRecruitmentLineRepository, talRepo)
}

func (uc *QuestionUseCase) CreateOrUpdateQuestions(req *request.CreateOrUpdateQuestions) (*response.TemplateQuestionResponse, error) {
	// check if template question exist
	tq, err := uc.TemplateQuestionRepository.FindByID(uuid.MustParse(req.TemplateQuestionID))
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when finding template question by id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when finding template question by id: " + err.Error())
	}

	if tq == nil {
		uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] template question with id %s not found", req.TemplateQuestionID)
		return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] template question with id " + req.TemplateQuestionID + " not found")
	}

	// create or update questions
	for i, question := range req.Questions {
		if question.ID != "" && question.ID != uuid.Nil.String() {
			exist, err := uc.Repository.FindByID(uuid.MustParse(question.ID))
			if err != nil {
				uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when finding question by id: %s", err.Error())
				return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when finding question by id: " + err.Error())
			}

			if exist == nil {
				createdQuestion, err := uc.Repository.CreateQuestion(&entity.Question{
					TemplateQuestionID: tq.ID,
					AnswerTypeID:       uuid.MustParse(question.AnswerTypeID),
					Name:               question.Name,
					Number:             i + 1,
				})
				if err != nil {
					uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question: %s", err.Error())
					return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question: " + err.Error())
				}

				if len(question.QuestionOptions) > 0 {
					for _, questionOption := range question.QuestionOptions {
						_, err := uc.QuestionOptionRepository.CreateQuestionOption(&entity.QuestionOption{
							QuestionID: createdQuestion.ID,
							OptionText: questionOption.OptionText,
						})
						if err != nil {
							uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: %s", err.Error())
							return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: " + err.Error())
						}
					}
				}
			} else {
				updatedQuestion, err := uc.Repository.UpdateQuestion(&entity.Question{
					ID:                 exist.ID,
					TemplateQuestionID: tq.ID,
					AnswerTypeID:       uuid.MustParse(question.AnswerTypeID),
					Name:               question.Name,
					Number:             i + 1,
				})
				if err != nil {
					uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when updating question: %s", err.Error())
					return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when updating question: " + err.Error())
				}

				// delete question options
				err = uc.QuestionOptionRepository.DeleteQuestionOptionsByQuestionID(updatedQuestion.ID)
				if err != nil {
					uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question options: %s", err.Error())
					return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question options: " + err.Error())
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
							return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: " + err.Error())
						}
					}
				}
			}
		} else {
			uc.Log.Info("Payloadku: ", question.Name)
			createdQuestion, err := uc.Repository.CreateQuestion(&entity.Question{
				TemplateQuestionID: tq.ID,
				AnswerTypeID:       uuid.MustParse(question.AnswerTypeID),
				Name:               question.Name,
				Number:             i + 1,
			})
			if err != nil {
				uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question: %s", err.Error())
				return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question: " + err.Error())
			}

			if len(question.QuestionOptions) > 0 {
				for _, questionOption := range question.QuestionOptions {
					_, err := uc.QuestionOptionRepository.CreateQuestionOption(&entity.QuestionOption{
						QuestionID: createdQuestion.ID,
						OptionText: questionOption.OptionText,
					})
					if err != nil {
						uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: %s", err.Error())
						return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when creating question option: " + err.Error())
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
				return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when deleting question: " + err.Error())
			}
		}
	}

	tQuestion, err := uc.TemplateQuestionRepository.FindByID(tq.ID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.CreateOrUpdateQuestions] error when finding template question by id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.CreateOrUpdateQuestions] error when finding template question by id: " + err.Error())
	}

	return uc.TemplateQuestionDTO.ConvertEntityToResponse(tQuestion), nil
}

func (uc *QuestionUseCase) FindByIDAndUserID(questionID, userID string) (*response.QuestionResponse, error) {
	q, err := uc.Repository.FindByID(uuid.MustParse(questionID))
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] error when finding question by id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] error when finding question by id: " + err.Error())
	}

	if q == nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] question with id %s not found", questionID)
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] question with id " + questionID + " not found")
	}

	up, err := uc.UserProfileRepository.FindByUserID(uuid.MustParse(userID))
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] error when finding user profile by user id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] error when finding user profile by user id: " + err.Error())
	}

	if up == nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] user profile with user id %s not found", userID)
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] user profile with user id " + userID + " not found")
	}

	qr, err := uc.Repository.FindQuestionWithResponsesByIDAndUserProfileID(q.ID, up.ID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] error when finding question with responses by id and user profile id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] error when finding question with responses by id and user profile id: " + err.Error())
	}

	if qr == nil {
		uc.Log.Errorf("[QuestionUseCase.FindByIDAndUserID] question response not found")
		return nil, errors.New("[QuestionUseCase.FindByIDAndUserID] question response not found")
	}

	return uc.DTO.ConvertEntityToResponse(qr), nil
}

func (uc *QuestionUseCase) FindAllByProjectRecruitmentLineIDAndJobPostingID(projectRecruitmentLineID uuid.UUID, jobPostingID uuid.UUID) (*[]response.QuestionResponse, error) {
	prl, err := uc.ProjectRecruitmentLineRepository.FindByID(projectRecruitmentLineID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] error when finding project recruitment line by id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] error when finding project recruitment line by id: " + err.Error())
	}

	if prl == nil {
		uc.Log.Errorf("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] project recruitment line with id %s not found", projectRecruitmentLineID)
		return nil, errors.New("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] project recruitment line with id " + projectRecruitmentLineID.String() + " not found")
	}

	tal, err := uc.TemplateActivityLineRepository.FindByID(prl.TemplateActivityLineID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] error when finding template activity line by id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] error when finding template activity line by id: " + err.Error())
	}

	if tal == nil {
		uc.Log.Errorf("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] template activity line with id %s not found", prl.TemplateActivityLineID)
		return nil, errors.New("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] template activity line with id " + prl.TemplateActivityLineID.String() + " not found")
	}

	tq, err := uc.TemplateQuestionRepository.FindByID(tal.QuestionTemplateID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] error when finding template question by id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] error when finding template question by id: " + err.Error())
	}

	if tq == nil {
		uc.Log.Errorf("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] template question with id %s not found", tal.QuestionTemplateID)
		return nil, errors.New("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] template question with id " + tal.QuestionTemplateID.String() + " not found")
	}

	questions, err := uc.Repository.FindAllByTemplateQuestionIDsAndJobPostingID([]uuid.UUID{tq.ID}, jobPostingID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] error when finding questions by template question ids and job posting id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindAllByProjectRecruitmentLineIDAndJobPostingID] error when finding questions by template question ids and job posting id: " + err.Error())
	}

	responses := make([]response.QuestionResponse, 0)
	for _, q := range *questions {
		responses = append(responses, *uc.DTO.ConvertEntityToResponse(&q))
	}

	return &responses, nil
}

func (uc *QuestionUseCase) FindByID(questionID string) (*entity.Question, error) {
	parsedID, err := uuid.Parse(questionID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindByID] error when parsing question id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindByID] error when parsing question id: " + err.Error())
	}

	q, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Errorf("[QuestionUseCase.FindByID] error when finding question by id: %s", err.Error())
		return nil, errors.New("[QuestionUseCase.FindByID] error when finding question by id: " + err.Error())
	}

	return q, nil
}
