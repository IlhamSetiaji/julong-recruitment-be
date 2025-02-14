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

type IFgdResultUseCase interface {
	FillFgdResult(req *request.FillFgdResultRequest) (*response.FgdResultResponse, error)
	FindByFgdApplicantAndAssessorID(fgdApplicantID, fgdAssessorID uuid.UUID) (*response.FgdResultResponse, error)
}

type FgdResultUseCase struct {
	Log                    *logrus.Logger
	Repository             repository.IFgdResultRepository
	FgdApplicantRepository repository.IFgdApplicantRepository
	FgdAssessorRepository  repository.IFgdAssessorRepository
	Viper                  *viper.Viper
	DTO                    dto.IFgdResultDTO
}

func NewFgdResultUseCase(
	log *logrus.Logger,
	repo repository.IFgdResultRepository,
	fgdApplicantRepo repository.IFgdApplicantRepository,
	fgdAssessorRepo repository.IFgdAssessorRepository,
	viper *viper.Viper,
	dto dto.IFgdResultDTO,
) IFgdResultUseCase {
	return &FgdResultUseCase{
		Log:                    log,
		Repository:             repo,
		FgdApplicantRepository: fgdApplicantRepo,
		FgdAssessorRepository:  fgdAssessorRepo,
		Viper:                  viper,
		DTO:                    dto,
	}
}

func FgdResultUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IFgdResultUseCase {
	repo := repository.FgdResultRepositoryFactory(log)
	fgdApplicantRepo := repository.FgdApplicantRepositoryFactory(log)
	fgdAssessorRepo := repository.FgdAssessorRepositoryFactory(log)
	dto := dto.FgdResultDTOFactory(log)
	return NewFgdResultUseCase(log, repo, fgdApplicantRepo, fgdAssessorRepo, viper, dto)
}

func (uc *FgdResultUseCase) FillFgdResult(req *request.FillFgdResultRequest) (*response.FgdResultResponse, error) {
	parsedFgdApplicantID, err := uuid.Parse(req.FgdApplicantID)
	if err != nil {
		uc.Log.Errorf("[FgdResultUseCase.FillFgdResult] " + err.Error())
		return nil, err
	}

	FgdApplicant, err := uc.FgdApplicantRepository.FindByKeys(map[string]interface{}{
		"applicant_id": parsedFgdApplicantID,
	})
	if err != nil {
		uc.Log.Errorf("[FgdResultUseCase.FillFgdResult] " + err.Error())
		return nil, err
	}
	if FgdApplicant == nil {
		return nil, errors.New("Fgd applicant not found")
	}

	parsedFgdAssessorID, err := uuid.Parse(req.FgdAssessorID)
	if err != nil {
		uc.Log.Errorf("[FgdResultUseCase.FillFgdResult] " + err.Error())
		return nil, err
	}

	FgdAssessor, err := uc.FgdAssessorRepository.FindByID(parsedFgdAssessorID)
	if err != nil {
		uc.Log.Errorf("[FgdResultUseCase.FillFgdResult] " + err.Error())
		return nil, err
	}
	if FgdAssessor == nil {
		return nil, errors.New("Fgd assessor not found")
	}

	FgdResult, err := uc.Repository.FindByKeys(map[string]interface{}{
		"fgd_applicant_id": req.FgdApplicantID,
		"fgd_assessor_id":  req.FgdAssessorID,
	})

	if err != nil {
		uc.Log.Errorf("[FgdResultUseCase.FillFgdResult] " + err.Error())
		return nil, err
	}

	if FgdResult != nil {
		return nil, errors.New("Fgd result already exists")
	}

	createdFgdResult, err := uc.Repository.CreateFgdResult(&entity.FgdResult{
		FgdApplicantID: FgdApplicant.ID,
		FgdAssessorID:  parsedFgdAssessorID,
		Status:         entity.FgdResultStatus(req.Status),
	})
	if err != nil {
		uc.Log.Errorf("[FgdResultUseCase.FillFgdResult] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(createdFgdResult), nil
}

func (uc *FgdResultUseCase) FindByFgdApplicantAndAssessorID(FgdApplicantID, FgdAssessorID uuid.UUID) (*response.FgdResultResponse, error) {
	FgdResult, err := uc.Repository.FindByKeys(map[string]interface{}{
		"fgd_applicant_id": FgdApplicantID,
		"fgd_assessor_id":  FgdAssessorID,
	})
	if err != nil {
		uc.Log.Errorf("[FgdResultUseCase.FindByFgdApplicantAndAssessorID] " + err.Error())
		return nil, err
	}
	if FgdResult == nil {
		return nil, nil
	}

	return uc.DTO.ConvertEntityToResponse(FgdResult), nil
}
