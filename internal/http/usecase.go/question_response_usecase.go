package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IQuestionResponseUseCase interface {
	CreateOrUpdateQuestionResponses(req *request.QuestionResponseRequest) (*response.QuestionResponse, error)
}

type QuestionResponseUseCase struct {
	Log                   *logrus.Logger
	Repository            repository.IQuestionResponseRepository
	JobPostingRepository  repository.IJobPostingRepository
	UserProfileRepository repository.IUserProfileRepository
	QuestionRepository    repository.IQuestionRepository
	QuestionDTO           dto.IQuestionDTO
	Viper                 *viper.Viper
}

func NewQuestionResponseUseCase(
	log *logrus.Logger,
	repo repository.IQuestionResponseRepository,
	jpRepo repository.IJobPostingRepository,
	upRepo repository.IUserProfileRepository,
	qRepo repository.IQuestionRepository,
	qDTO dto.IQuestionDTO,
	viper *viper.Viper,
) IQuestionResponseUseCase {
	return &QuestionResponseUseCase{
		Log:                   log,
		Repository:            repo,
		JobPostingRepository:  jpRepo,
		UserProfileRepository: upRepo,
		QuestionRepository:    qRepo,
		QuestionDTO:           qDTO,
		Viper:                 viper,
	}
}

func QuestionResponseUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IQuestionResponseUseCase {
	repo := repository.QuestionResponseRepositoryFactory(log)
	jpRepo := repository.JobPostingRepositoryFactory(log)
	upRepo := repository.UserProfileRepositoryFactory(log)
	qRepo := repository.QuestionRepositoryFactory(log)
	qDTO := dto.QuestionDTOFactory(log)
	return NewQuestionResponseUseCase(log, repo, jpRepo, upRepo, qRepo, qDTO, viper)
}

func (uc *QuestionResponseUseCase) CreateOrUpdateQuestionResponses(req *request.QuestionResponseRequest) (*response.QuestionResponse, error) {
	// check if question is exist
	parsedQuestionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when parsing question id: %s", err.Error())
		return nil, err
	}

	question, err := uc.QuestionRepository.FindByID(parsedQuestionID)
	if err != nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when finding question by id: %s", err.Error())
		return nil, err
	}

	if question == nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] question with id %s not found", req.QuestionID)
		return nil, err
	}

	// create or update answers
	for _, ans := range req.Answers {
		parsedJobPostingID, err := uuid.Parse(ans.JobPostingID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when parsing job posting id: %s", err.Error())
			return nil, err
		}
		jp, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when finding job posting by id: %s", err.Error())
			return nil, err
		}
		if jp == nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] job posting with id %s not found", ans.JobPostingID)
			return nil, err
		}

		parsedUserProfileID, err := uuid.Parse(ans.UserProfileID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when parsing user profile id: %s", err.Error())
			return nil, err
		}
		up, err := uc.UserProfileRepository.FindByID(parsedUserProfileID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when finding user profile by id: %s", err.Error())
			return nil, err
		}
		if up == nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] user profile with id %s not found", ans.UserProfileID)
			return nil, err
		}

		// check if answer is exist
		if ans.ID != "" && ans.ID != uuid.Nil.String() {
			parsedAnswerID, err := uuid.Parse(ans.ID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when parsing answer id: %s", err.Error())
				return nil, err
			}
			exist, err := uc.Repository.FindByID(parsedAnswerID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when finding answer by id: %s", err.Error())
				return nil, err
			}

			if exist == nil {
				_, err := uc.Repository.CreateQuestionResponse(&entity.QuestionResponse{
					QuestionID:    question.ID,
					JobPostingID:  jp.ID,
					UserProfileID: up.ID,
					Answer:        ans.Answer,
					AnswerFile:    ans.AnswerPath,
				})
				if err != nil {
					uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when creating answer: %s", err.Error())
					return nil, err
				}
			} else {
				_, err := uc.Repository.UpdateQuestionResponse(&entity.QuestionResponse{
					ID:            exist.ID,
					QuestionID:    question.ID,
					JobPostingID:  jp.ID,
					UserProfileID: up.ID,
					Answer:        ans.Answer,
					AnswerFile:    ans.AnswerPath,
				})

				if err != nil {
					uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when updating answer: %s", err.Error())
					return nil, err
				}
			}
		} else {
			_, err := uc.Repository.CreateQuestionResponse(&entity.QuestionResponse{
				QuestionID:    question.ID,
				JobPostingID:  jp.ID,
				UserProfileID: up.ID,
				Answer:        ans.Answer,
				AnswerFile:    ans.AnswerPath,
			})
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when creating answer: %s", err.Error())
				return nil, err
			}
		}
	}

	// delete answers
	if len(req.DeletedAnswerIDs) > 0 {
		for _, id := range req.DeletedAnswerIDs {
			parsedAnswerID, err := uuid.Parse(id)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when parsing answer id: %s", err.Error())
				return nil, err
			}
			err = uc.Repository.DeleteQuestionResponse(parsedAnswerID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when deleting answer: %s", err.Error())
				return nil, err
			}
		}
	}

	rQuestion, err := uc.QuestionRepository.FindByID(question.ID)
	if err != nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when finding question by id: %s", err.Error())
		return nil, err
	}
	if rQuestion == nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] question with id %s not found", req.QuestionID)
		return nil, err
	}

	return uc.QuestionDTO.ConvertEntityToResponse(rQuestion), nil
}
