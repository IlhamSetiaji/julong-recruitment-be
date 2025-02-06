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
	"github.com/spf13/viper"
)

type IQuestionResponseUseCase interface {
	CreateOrUpdateQuestionResponses(req *request.QuestionResponseRequest) (*response.QuestionResponse, error)
	AnswerInterviewQuestionResponses(req *request.InterviewQuestionResponseRequest) (*response.TemplateQuestionResponse, error)
}

type QuestionResponseUseCase struct {
	Log                         *logrus.Logger
	Repository                  repository.IQuestionResponseRepository
	JobPostingRepository        repository.IJobPostingRepository
	UserProfileRepository       repository.IUserProfileRepository
	QuestionRepository          repository.IQuestionRepository
	QuestionDTO                 dto.IQuestionDTO
	Viper                       *viper.Viper
	TemplateQuestionRepository  repository.ITemplateQuestionRepository
	TemplateQuestionDTO         dto.ITemplateQuestionDTO
	InterviewAssessorRepository repository.IInterviewAssessorRepository
}

func NewQuestionResponseUseCase(
	log *logrus.Logger,
	repo repository.IQuestionResponseRepository,
	jpRepo repository.IJobPostingRepository,
	upRepo repository.IUserProfileRepository,
	qRepo repository.IQuestionRepository,
	qDTO dto.IQuestionDTO,
	viper *viper.Viper,
	tqRepo repository.ITemplateQuestionRepository,
	tqDTO dto.ITemplateQuestionDTO,
	iaRepo repository.IInterviewAssessorRepository,
) IQuestionResponseUseCase {
	return &QuestionResponseUseCase{
		Log:                         log,
		Repository:                  repo,
		JobPostingRepository:        jpRepo,
		UserProfileRepository:       upRepo,
		QuestionRepository:          qRepo,
		QuestionDTO:                 qDTO,
		Viper:                       viper,
		TemplateQuestionRepository:  tqRepo,
		TemplateQuestionDTO:         tqDTO,
		InterviewAssessorRepository: iaRepo,
	}
}

func QuestionResponseUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IQuestionResponseUseCase {
	repo := repository.QuestionResponseRepositoryFactory(log)
	jpRepo := repository.JobPostingRepositoryFactory(log)
	upRepo := repository.UserProfileRepositoryFactory(log)
	qRepo := repository.QuestionRepositoryFactory(log)
	qDTO := dto.QuestionDTOFactory(log)
	tqRepo := repository.TemplateQuestionRepositoryFactory(log)
	tqDTO := dto.TemplateQuestionDTOFactory(log)
	iaRepo := repository.InterviewAssessorRepositoryFactory(log)
	return NewQuestionResponseUseCase(log, repo, jpRepo, upRepo, qRepo, qDTO, viper, tqRepo, tqDTO, iaRepo)
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

	var userProfileUUID uuid.UUID

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
			return nil, errors.New("job posting not found")
		}

		parsedUserProfileID, err := uuid.Parse(ans.UserProfileID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when parsing user profile id: %s", err.Error())
			return nil, err
		}
		userProfileUUID = parsedUserProfileID
		up, err := uc.UserProfileRepository.FindByID(parsedUserProfileID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when finding user profile by id: %s", err.Error())
			return nil, err
		}
		if up == nil {
			uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] user profile with id %s not found", ans.UserProfileID)
			return nil, errors.New("user profile not found")
		}
		uc.Log.Info("Halooo")

		// check if answer is exist
		if ans.ID != nil {
			parsedAnswerID, err := uuid.Parse(*ans.ID)
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
				uc.Log.Infof("kontol: %+v", ans)
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
				uc.Log.Infof("memek: %+v", ans)
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
			uc.Log.Infof("cok: %+v", ans)
			hasil, err := uc.Repository.CreateQuestionResponse(&entity.QuestionResponse{
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

			uc.Log.Infof("hasil: %+v", hasil)
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

	rQuestion, err := uc.QuestionRepository.FindQuestionWithResponsesByIDAndUserProfileID(question.ID, userProfileUUID)
	if err != nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] error when finding question by id: %s", err.Error())
		return nil, err
	}
	if rQuestion == nil {
		uc.Log.Errorf("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] question with id %s not found", req.QuestionID)
		return nil, err
	}

	// embed url to answer file
	for _, qr := range rQuestion.QuestionResponses {
		if qr.AnswerFile != "" {
			qr.AnswerFile = uc.Viper.GetString("app.url") + "/" + qr.AnswerFile
			uc.Log.Infof("[QuestionResponseUseCase.CreateOrUpdateQuestionResponses] answer file url: %s", qr.AnswerFile)
		}
	}

	return uc.QuestionDTO.ConvertEntityToResponse(rQuestion), nil
}

func (uc *QuestionResponseUseCase) AnswerInterviewQuestionResponses(req *request.InterviewQuestionResponseRequest) (*response.TemplateQuestionResponse, error) {
	parsedTemplateQuestionID, err := uuid.Parse(req.TemplateQuestionID)
	if err != nil {
		uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when parsing template question id: %s", err.Error())
		return nil, err
	}

	tq, err := uc.TemplateQuestionRepository.FindByID(parsedTemplateQuestionID)
	if err != nil {
		uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when finding template question by id: %s", err.Error())
		return nil, err
	}

	if tq == nil {
		uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] template question with id %s not found", req.TemplateQuestionID)
		return nil, errors.New("template question not found")
	}

	var jobPostingID uuid.UUID
	var userProfileID uuid.UUID

	// create or update answers
	for _, questionPayload := range req.Questions {
		parsedQuestionID, err := uuid.Parse(questionPayload.ID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when parsing question id: %s", err.Error())
			return nil, err
		}

		question, err := uc.QuestionRepository.FindByID(parsedQuestionID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when finding question by id: %s", err.Error())
			return nil, err
		}

		if question == nil {
			uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] question with id %s not found", questionPayload.ID)
			return nil, errors.New("question not found")
		}

		userProfileID, err = uuid.Parse(questionPayload.Answers[0].UserProfileID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when parsing user profile id: %s", err.Error())
			return nil, err
		}

		jobPostingID, err = uuid.Parse(questionPayload.Answers[0].JobPostingID)
		if err != nil {
			uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when parsing job posting id: %s", err.Error())
			return nil, err
		}

		for _, ans := range questionPayload.Answers {
			parsedJobPostingID, err := uuid.Parse(ans.JobPostingID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when parsing job posting id: %s", err.Error())
				return nil, err
			}
			jp, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when finding job posting by id: %s", err.Error())
				return nil, err
			}

			if jp == nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] job posting with id %s not found", ans.JobPostingID)
				return nil, errors.New("job posting not found")
			}

			parsedUserProfileID, err := uuid.Parse(ans.UserProfileID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when parsing user profile id: %s", err.Error())
				return nil, err
			}
			up, err := uc.UserProfileRepository.FindByID(parsedUserProfileID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when finding user profile by id: %s", err.Error())
				return nil, err
			}

			if up == nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] user profile with id %s not found", ans.UserProfileID)
				return nil, errors.New("user profile not found")
			}

			parsedAssessorInterviewID, err := uuid.Parse(ans.InterviewAssessorID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when parsing interview assessor id: %s", err.Error())
				return nil, err
			}

			ia, err := uc.InterviewAssessorRepository.FindByID(parsedAssessorInterviewID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when finding interview assessor by id: %s", err.Error())
				return nil, err
			}

			if ia == nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] interview assessor with id %s not found", ans.InterviewAssessorID)
				return nil, errors.New("interview assessor not found")
			}

			if ans.ID != "" && ans.ID != uuid.Nil.String() {
				parsedAnswerID, err := uuid.Parse(ans.ID)
				if err != nil {
					uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when parsing answer id: %s", err.Error())
					return nil, err
				}

				exist, err := uc.Repository.FindByID(parsedAnswerID)
				if err != nil {
					uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when finding answer by id: %s", err.Error())
					return nil, err
				}

				if exist == nil {
					_, err := uc.Repository.CreateQuestionResponse(&entity.QuestionResponse{
						QuestionID:          question.ID,
						JobPostingID:        jp.ID,
						UserProfileID:       up.ID,
						InterviewAssessorID: ia.ID,
						Answer:              ans.Answer,
					})
					if err != nil {
						uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when creating answer: %s", err.Error())
						return nil, err
					}
				} else {
					_, err := uc.Repository.UpdateQuestionResponse(&entity.QuestionResponse{
						ID:                  exist.ID,
						QuestionID:          question.ID,
						JobPostingID:        jp.ID,
						UserProfileID:       up.ID,
						InterviewAssessorID: ia.ID,
						Answer:              ans.Answer,
					})
					if err != nil {
						uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when updating answer: %s", err.Error())
						return nil, err
					}
				}
			} else {
				_, err := uc.Repository.CreateQuestionResponse(&entity.QuestionResponse{
					QuestionID:          question.ID,
					JobPostingID:        jp.ID,
					UserProfileID:       up.ID,
					InterviewAssessorID: ia.ID,
					Answer:              ans.Answer,
				})
				if err != nil {
					uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when creating answer: %s", err.Error())
					return nil, err
				}
			}
		}
	}

	// delete answers
	if len(req.DeletedAnswerIDs) > 0 {
		for _, id := range req.DeletedAnswerIDs {
			parsedAnswerID, err := uuid.Parse(id)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when parsing answer id: %s", err.Error())
				return nil, err
			}
			err = uc.Repository.DeleteQuestionResponse(parsedAnswerID)
			if err != nil {
				uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when deleting answer: %s", err.Error())
				return nil, err
			}
		}
	}

	tqRes, err := uc.TemplateQuestionRepository.FindByIDForInterviewAnswer(parsedTemplateQuestionID, userProfileID, jobPostingID)
	if err != nil {
		uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] error when finding template question by id: %s", err.Error())
		return nil, err
	}

	if tqRes == nil {
		uc.Log.Errorf("[QuestionResponseUseCase.AnswerInterviewQuestionResponses] template question with id %s not found", req.TemplateQuestionID)
		return nil, errors.New("template question not found")
	}

	return uc.TemplateQuestionDTO.ConvertEntityToResponse(tqRes), nil
}
