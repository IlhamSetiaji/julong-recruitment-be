package usecase

import (
	"errors"
	"strconv"
	"time"

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
	JobPostingRepository              repository.IJobPostingRepository
	ApplicantRepository               repository.IApplicantRepository
}

func NewAdministrativeResultUseCase(
	log *logrus.Logger,
	repo repository.IAdministrativeResultRepository,
	arDTO dto.IAdministrativeResultDTO,
	asDTO dto.IAdministrativeSelectionDTO,
	viper *viper.Viper,
	userProfileRepository repository.IUserProfileRepository,
	asRepository repository.IAdministrativeSelectionRepository,
	jpRepo repository.IJobPostingRepository,
	applicantRepo repository.IApplicantRepository,
) IAdministrativeResultUseCase {
	return &AdministrativeResultUseCase{
		Log:                               log,
		Repository:                        repo,
		DTO:                               arDTO,
		asDTO:                             asDTO,
		Viper:                             viper,
		UserProfileRepository:             userProfileRepository,
		AdministrativeSelectionRepository: asRepository,
		JobPostingRepository:              jpRepo,
		ApplicantRepository:               applicantRepo,
	}
}

func AdministrativeResultUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IAdministrativeResultUseCase {
	repo := repository.AdministrativeResultRepositoryFactory(log)
	arDTO := dto.AdministrativeResultDTOFactory(log, viper)
	asDTO := dto.AdministrativeSelectionDTOFactory(log, viper)
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	asRepository := repository.AdministrativeSelectionRepositoryFactory(log)
	jpRepo := repository.JobPostingRepositoryFactory(log)
	applicantRepo := repository.ApplicantRepositoryFactory(log)
	return NewAdministrativeResultUseCase(log, repo, arDTO, asDTO, viper, userProfileRepository, asRepository, jpRepo, applicantRepo)
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

				if entity.AdministrativeResultStatus(administrativeResult.Status) == entity.ADMINISTRATIVE_RESULT_STATUS_ACCEPTED {
					jpExist, err := uc.JobPostingRepository.FindByID(as.JobPostingID)
					if err != nil {
						uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
						return nil, err
					}

					if jpExist == nil {
						uc.Log.Error("[ApplicantUseCase.CreateOrUpdateAdministrativeResults] " + "Job Posting not found")
						return nil, errors.New("job posting not found")
					}

					applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
						"user_profile_id": parsedUserProfileID,
						"job_posting_id":  as.JobPostingID,
					})
					if err != nil {
						uc.Log.Error("[ApplicantUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
						return nil, err
					}
					if applicant != nil {
						var TemplateQuestionID *uuid.UUID
						for i := range jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines {
							if jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines[i].Order == applicant.Order+1 {
								projectRecruitmentLine := &jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines[i]
								TemplateQuestionID = &projectRecruitmentLine.TemplateActivityLine.QuestionTemplateID
								break
							} else {
								TemplateQuestionID = &applicant.TemplateQuestionID
							}
						}
						_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
							UserProfileID:      parsedUserProfileID,
							JobPostingID:       as.JobPostingID,
							Status:             entity.APPLICANT_STATUS_APPLIED,
							AppliedDate:        time.Now(),
							Order:              applicant.Order + 1,
							TemplateQuestionID: *TemplateQuestionID,
						})
						if err != nil {
							uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
							return nil, err
						}
					}
				} else if entity.AdministrativeResultStatus(administrativeResult.Status) == entity.ADMINISTRATIVE_RESULT_STATUS_REJECTED {
					jpExist, err := uc.JobPostingRepository.FindByID(as.JobPostingID)
					if err != nil {
						uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
						return nil, err
					}

					if jpExist == nil {
						uc.Log.Error("[ApplicantUseCase.CreateOrUpdateAdministrativeResults] " + "Job Posting not found")
						return nil, errors.New("job posting not found")
					}

					applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
						"user_profile_id": parsedUserProfileID,
						"job_posting_id":  as.JobPostingID,
					})
					if err != nil {
						uc.Log.Error("[ApplicantUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
						return nil, err
					}
					zero, err := strconv.Atoi("0")
					if err != nil {
						uc.Log.Error("[ApplicantUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
						return nil, err
					}
					if applicant != nil {
						_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
							UserProfileID:      parsedUserProfileID,
							JobPostingID:       as.JobPostingID,
							Status:             entity.APPLICANT_STATUS_APPLIED,
							AppliedDate:        time.Now(),
							Order:              zero,
							TemplateQuestionID: uuid.Nil,
						})
						if err != nil {
							uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
							return nil, err
						}
					}
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
		uc.Log.Error("[AdministrativeResultUseCase.FindByAdministrativeSelectionID - parsed id]" + err.Error())
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
