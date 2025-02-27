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

type IInterviewResultUseCase interface {
	FillInterviewResult(req *request.FillInterviewResultRequest) (*response.InterviewResultResponse, error)
	FindByInterviewApplicantAndAssessorID(interviewApplicantID, interviewAssessorID uuid.UUID) (*response.InterviewResultResponse, error)
}

type InterviewResultUseCase struct {
	Log                          *logrus.Logger
	Repository                   repository.IInterviewResultRepository
	InterviewApplicantRepository repository.IInterviewApplicantRepository
	InterviewAssessorRepository  repository.IInterviewAssessorRepository
	Viper                        *viper.Viper
	DTO                          dto.IInterviewResultDTO
}

func NewInterviewResultUseCase(
	log *logrus.Logger,
	repo repository.IInterviewResultRepository,
	interviewApplicantRepo repository.IInterviewApplicantRepository,
	interviewAssessorRepo repository.IInterviewAssessorRepository,
	viper *viper.Viper,
	dto dto.IInterviewResultDTO,
) IInterviewResultUseCase {
	return &InterviewResultUseCase{
		Log:                          log,
		Repository:                   repo,
		InterviewApplicantRepository: interviewApplicantRepo,
		InterviewAssessorRepository:  interviewAssessorRepo,
		Viper:                        viper,
		DTO:                          dto,
	}
}

func InterviewResultUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IInterviewResultUseCase {
	repo := repository.InterviewResultRepositoryFactory(log)
	interviewApplicantRepo := repository.InterviewApplicantRepositoryFactory(log)
	interviewAssessorRepo := repository.InterviewAssessorRepositoryFactory(log)
	dto := dto.InterviewResultDTOFactory(log)
	return NewInterviewResultUseCase(log, repo, interviewApplicantRepo, interviewAssessorRepo, viper, dto)
}

func (uc *InterviewResultUseCase) FillInterviewResult(req *request.FillInterviewResultRequest) (*response.InterviewResultResponse, error) {
	parsedInterviewApplicantID, err := uuid.Parse(req.InterviewApplicantID)
	if err != nil {
		uc.Log.Errorf("[InterviewResultUseCase.FillInterviewResult] " + err.Error())
		return nil, err
	}

	// interviewApplicant, err := uc.InterviewApplicantRepository.FindByID(parsedInterviewApplicantID)
	interviewApplicant, err := uc.InterviewApplicantRepository.FindByKeys(map[string]interface{}{
		"applicant_id": parsedInterviewApplicantID,
	})
	if err != nil {
		uc.Log.Errorf("[InterviewResultUseCase.FillInterviewResult] " + err.Error())
		return nil, err
	}
	if interviewApplicant == nil {
		return nil, errors.New("Interview applicant not found")
	}

	parsedInterviewAssessorID, err := uuid.Parse(req.InterviewAssessorID)
	if err != nil {
		uc.Log.Errorf("[InterviewResultUseCase.FillInterviewResult] " + err.Error())
		return nil, err
	}

	interviewAssessor, err := uc.InterviewAssessorRepository.FindByID(parsedInterviewAssessorID)
	if err != nil {
		uc.Log.Errorf("[InterviewResultUseCase.FillInterviewResult] " + err.Error())
		return nil, err
	}
	if interviewAssessor == nil {
		return nil, errors.New("Interview assessor not found")
	}

	interviewResult, err := uc.Repository.FindByKeys(map[string]interface{}{
		"interview_applicant_id": req.InterviewApplicantID,
		"interview_assessor_id":  req.InterviewAssessorID,
	})

	if err != nil {
		uc.Log.Errorf("[InterviewResultUseCase.FillInterviewResult] " + err.Error())
		return nil, err
	}

	if interviewResult != nil {
		return nil, errors.New("Interview result already exists")
	}

	createdInterviewResult, err := uc.Repository.CreateInterviewResult(&entity.InterviewResult{
		InterviewApplicantID: interviewApplicant.ID,
		InterviewAssessorID:  parsedInterviewAssessorID,
		Status:               entity.InterviewResultStatus(req.Status),
	})
	if err != nil {
		uc.Log.Errorf("[InterviewResultUseCase.FillInterviewResult] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(createdInterviewResult), nil
}

func (uc *InterviewResultUseCase) FindByInterviewApplicantAndAssessorID(interviewApplicantID, interviewAssessorID uuid.UUID) (*response.InterviewResultResponse, error) {
	interviewResult, err := uc.Repository.FindByKeys(map[string]interface{}{
		"interview_applicant_id": interviewApplicantID,
		"interview_assessor_id":  interviewAssessorID,
	})
	if err != nil {
		uc.Log.Errorf("[InterviewResultUseCase.FindByInterviewApplicantAndAssessorID] " + err.Error())
		return nil, err
	}
	if interviewResult == nil {
		return nil, nil
	}

	return uc.DTO.ConvertEntityToResponse(interviewResult), nil
}
