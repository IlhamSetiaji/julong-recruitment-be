package usecase

import (
	"errors"
	"strconv"

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
	FindAllByAdministrativeSelectionID(administrativeSelectionID string, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.AdministrativeResultResponse, int64, error)
	FindByID(id uuid.UUID) (*response.AdministrativeResultResponse, error)
	UpdateStatusAdministrativeResult(id uuid.UUID, status entity.AdministrativeResultStatus) (*response.AdministrativeResultResponse, error)
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
						uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + "Job Posting not found")
						return nil, errors.New("job posting not found")
					}

					applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
						"user_profile_id": parsedUserProfileID,
						"job_posting_id":  as.JobPostingID,
					})
					if err != nil {
						uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
						return nil, err
					}
					if applicant != nil {
						applicantOrder := applicant.Order
						var TemplateQuestionID *uuid.UUID
						for i := range jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines {
							if jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines[i].Order == applicantOrder+1 {
								projectRecruitmentLine := &jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines[i]
								TemplateQuestionID = &projectRecruitmentLine.TemplateActivityLine.QuestionTemplateID
								break
							} else {
								TemplateQuestionID = &applicant.TemplateQuestionID
							}
						}
						_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
							ID:                 applicant.ID,
							UserProfileID:      parsedUserProfileID,
							JobPostingID:       as.JobPostingID,
							Order:              applicant.Order + 1,
							TemplateQuestionID: *TemplateQuestionID,
						})
						if err != nil {
							uc.Log.Error("[AdministrativeResultUseCase.ApplyJobPosting] " + err.Error())
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
						uc.Log.Infof("Masuk sini bos")
						_, err = uc.ApplicantRepository.UpdateApplicantWhenRejected(&entity.Applicant{
							ID:            applicant.ID,
							UserProfileID: parsedUserProfileID,
							JobPostingID:  as.JobPostingID,
							Order:         zero,
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

func (uc *AdministrativeResultUseCase) FindAllByAdministrativeSelectionID(administrativeSelectionID string, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.AdministrativeResultResponse, int64, error) {
	parsedAdministrativeSelectionID, err := uuid.Parse(administrativeSelectionID)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.FindAllByAdministrativeSelectionID] " + err.Error())
		return nil, 0, err
	}

	entities, total, err := uc.Repository.FindAllByAdministrativeSelectionID(parsedAdministrativeSelectionID, page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.FindAllByAdministrativeSelectionID] " + err.Error())
		return nil, 0, err
	}

	var res []response.AdministrativeResultResponse
	for _, entity := range *entities {
		response, err := uc.DTO.ConvertEntityToResponse(&entity)
		if err != nil {
			uc.Log.Error("[AdministrativeResultUseCase.FindAllByAdministrativeSelectionID] " + err.Error())
			return nil, 0, err
		}

		res = append(res, *response)
	}

	return &res, total, nil
}

func (uc *AdministrativeResultUseCase) FindByID(id uuid.UUID) (*response.AdministrativeResultResponse, error) {
	entity, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.FindByID] " + err.Error())
		return nil, err
	}

	if entity == nil {
		return nil, nil
	}

	res, err := uc.DTO.ConvertEntityToResponse(entity)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return res, nil
}

func (uc *AdministrativeResultUseCase) UpdateStatusAdministrativeResult(id uuid.UUID, status entity.AdministrativeResultStatus) (*response.AdministrativeResultResponse, error) {
	exist, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.UpdateStatusAdministrativeResult] " + err.Error())
		return nil, err
	}

	if exist == nil {
		return nil, nil
	}

	admResult, err := uc.Repository.UpdateAdministrativeResult(&entity.AdministrativeResult{
		ID:     id,
		Status: status,
	})
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.UpdateStatusAdministrativeResult] " + err.Error())
		return nil, err
	}

	if status == entity.ADMINISTRATIVE_RESULT_STATUS_ACCEPTED {
		jpExist, err := uc.JobPostingRepository.FindByID(exist.AdministrativeSelection.JobPostingID)
		if err != nil {
			uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
			return nil, err
		}

		if jpExist == nil {
			uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + "Job Posting not found")
			return nil, errors.New("job posting not found")
		}

		applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
			"user_profile_id": exist.UserProfileID,
			"job_posting_id":  exist.AdministrativeSelection.JobPostingID,
		})
		if err != nil {
			uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
			return nil, err
		}
		if applicant != nil {
			applicantOrder := applicant.Order
			var TemplateQuestionID *uuid.UUID
			for i := range jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines {
				if jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines[i].Order == applicantOrder+1 {
					projectRecruitmentLine := &jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines[i]
					TemplateQuestionID = &projectRecruitmentLine.TemplateActivityLine.QuestionTemplateID
					break
				} else {
					TemplateQuestionID = &applicant.TemplateQuestionID
				}
			}
			_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
				ID:                 applicant.ID,
				UserProfileID:      exist.UserProfileID,
				JobPostingID:       exist.AdministrativeSelection.JobPostingID,
				Status:             entity.APPLICANT_STATUS_APPLIED,
				Order:              applicant.Order + 1,
				TemplateQuestionID: *TemplateQuestionID,
			})
			if err != nil {
				uc.Log.Error("[AdministrativeResultUseCase.ApplyJobPosting] " + err.Error())
				return nil, err
			}
		}
	} else if status == entity.ADMINISTRATIVE_RESULT_STATUS_REJECTED {
		jpExist, err := uc.JobPostingRepository.FindByID(exist.AdministrativeSelection.JobPostingID)
		if err != nil {
			uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
			return nil, err
		}

		if jpExist == nil {
			uc.Log.Error("[ApplicantUseCase.CreateOrUpdateAdministrativeResults] " + "Job Posting not found")
			return nil, errors.New("job posting not found")
		}

		applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
			"user_profile_id": exist.UserProfileID,
			"job_posting_id":  exist.AdministrativeSelection.JobPostingID,
		})
		if err != nil {
			uc.Log.Error("[ApplicantUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
			return nil, err
		}
		// zero, err := strconv.Atoi("0")
		// if err != nil {
		// 	uc.Log.Error("[ApplicantUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
		// 	return nil, err
		// }
		if applicant != nil {
			_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
				ID:                 applicant.ID,
				UserProfileID:      exist.UserProfileID,
				JobPostingID:       exist.AdministrativeSelection.JobPostingID,
				TemplateQuestionID: uuid.Nil,
				Status:             entity.APPLICANT_STATUS_REJECTED,
				ProcessStatus:      entity.APPLICANT_PROCESS_STATUS_REJECTED,
			})
			if err != nil {
				uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
				return nil, err
			}
		}
	}

	res, err := uc.DTO.ConvertEntityToResponse(admResult)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.UpdateStatusAdministrativeResult] " + err.Error())
		return nil, err
	}

	return res, nil
}
