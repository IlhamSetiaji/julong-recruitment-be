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

type IAdministrativeResultUseCase interface {
	CreateOrUpdateAdministrativeResults(req *request.CreateOrUpdateAdministrativeResults) (*response.AdministrativeSelectionResponse, error)
	FindAllByAdministrativeSelectionID(administrativeSelectionID string) (*[]response.AdministrativeResultResponse, error)
}

type AdministrativeResultUseCase struct {
	Log                               *logrus.Logger
	Repository                        repository.IAdministrativeResultRepository
	DTO                               dto.IAdministrativeResultDTO
	asDTO                             dto.IAdministrativeSelectionDTO
	Viper                             *viper.Viper
	UserProfileRepository             repository.IUserProfileRepository
	AdministrativeSelectionRepository repository.IAdministrativeSelectionRepository
}

func NewAdministrativeResultUseCase(
	log *logrus.Logger,
	repo repository.IAdministrativeResultRepository,
	arDTO dto.IAdministrativeResultDTO,
	asDTO dto.IAdministrativeSelectionDTO,
	viper *viper.Viper,
	userProfileRepository repository.IUserProfileRepository,
	asRepository repository.IAdministrativeSelectionRepository,
) IAdministrativeResultUseCase {
	return &AdministrativeResultUseCase{
		Log:                               log,
		Repository:                        repo,
		DTO:                               arDTO,
		asDTO:                             asDTO,
		Viper:                             viper,
		UserProfileRepository:             userProfileRepository,
		AdministrativeSelectionRepository: asRepository,
	}
}

func AdministrativeResultUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IAdministrativeResultUseCase {
	repo := repository.AdministrativeResultRepositoryFactory(log)
	arDTO := dto.AdministrativeResultDTOFactory(log, viper)
	asDTO := dto.AdministrativeSelectionDTOFactory(log, viper)
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	asRepository := repository.AdministrativeSelectionRepositoryFactory(log)
	return NewAdministrativeResultUseCase(log, repo, arDTO, asDTO, viper, userProfileRepository, asRepository)
}

func (uc *AdministrativeResultUseCase) CreateOrUpdateAdministrativeResults(req *request.CreateOrUpdateAdministrativeResults) (*response.AdministrativeSelectionResponse, error) {
	// Check if administrative selection exists
	parsedAdministrativeSelectionID, err := uuid.Parse(req.AdministrativeSelectionID)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
		return nil, err
	}

	as, err := uc.AdministrativeSelectionRepository.FindByID(parsedAdministrativeSelectionID)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
		return nil, err
	}

	if as == nil {
		uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] administrative selection not found")
		return nil, err
	}

	// create or update administrative results
	for _, administrativeResult := range req.AdministrativeResults {
		parsedUserProfileID, err := uuid.Parse(administrativeResult.UserProfileID)
		if err != nil {
			uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
			return nil, err
		}

		userProfile, err := uc.UserProfileRepository.FindByID(parsedUserProfileID)
		if err != nil {
			uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
			return nil, err
		}

		if userProfile == nil {
			uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] user profile not found")
			return nil, err
		}

		if administrativeResult.ID != "" && administrativeResult.ID != uuid.Nil.String() {
			parsedID, err := uuid.Parse(administrativeResult.ID)
			if err != nil {
				uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
				return nil, err
			}
			exist, err := uc.Repository.FindByID(parsedID)
			if err != nil {
				uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
				return nil, err
			}

			if exist == nil {
				_, err := uc.Repository.CreateAdministrativeResult(&entity.AdministrativeResult{
					AdministrativeSelectionID: parsedAdministrativeSelectionID,
					UserProfileID:             parsedUserProfileID,
					Status:                    entity.AdministrativeResultStatus(administrativeResult.Status),
				})
				if err != nil {
					uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
					return nil, err
				}
			} else {
				_, err := uc.Repository.UpdateAdministrativeResult(&entity.AdministrativeResult{
					ID:                        parsedID,
					AdministrativeSelectionID: parsedAdministrativeSelectionID,
					UserProfileID:             parsedUserProfileID,
					Status:                    entity.AdministrativeResultStatus(administrativeResult.Status),
				})
				if err != nil {
					uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
					return nil, err
				}
			}
		} else {
			_, err := uc.Repository.CreateAdministrativeResult(&entity.AdministrativeResult{
				AdministrativeSelectionID: parsedAdministrativeSelectionID,
				UserProfileID:             parsedUserProfileID,
				Status:                    entity.AdministrativeResultStatus(administrativeResult.Status),
			})
			if err != nil {
				uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
				return nil, err
			}
		}
	}

	// delete administrative results
	if len(req.DeletedAdministrativeResultIDs) > 0 {
		for _, id := range req.DeletedAdministrativeResultIDs {
			parsedID, err := uuid.Parse(id)
			if err != nil {
				uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
				return nil, err
			}

			err = uc.Repository.DeleteAdministrativeResult(parsedID)
			if err != nil {
				uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
				return nil, err
			}
		}
	}

	// get administrative results
	administrativeResults, err := uc.AdministrativeSelectionRepository.FindByID(parsedAdministrativeSelectionID)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
		return nil, err
	}

	res, err := uc.asDTO.ConvertEntityToResponse(administrativeResults)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
		return nil, err
	}

	return res, nil
}

func (uc *AdministrativeResultUseCase) FindAllByAdministrativeSelectionID(administrativeSelectionID string) (*[]response.AdministrativeResultResponse, error) {
	parsedAdministrativeSelectionID, err := uuid.Parse(administrativeSelectionID)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.FindByAdministrativeSelectionID] " + err.Error())
		return nil, err
	}

	administrativeResults, err := uc.Repository.FindAllByAdministrativeSelectionID(parsedAdministrativeSelectionID)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.FindByAdministrativeSelectionID] " + err.Error())
		return nil, err
	}

	res := make([]response.AdministrativeResultResponse, 0)
	for _, administrativeResult := range *administrativeResults {
		r, err := uc.DTO.ConvertEntityToResponse(&administrativeResult)
		if err != nil {
			uc.Log.Error("[AdministrativeResultUseCase.FindByAdministrativeSelectionID] " + err.Error())
			return nil, err
		}
		res = append(res, *r)
	}

	return &res, nil
}
