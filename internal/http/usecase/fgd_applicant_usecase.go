package usecase

import (
	"context"
	"errors"
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

type IFgdApplicantUseCase interface {
	CreateOrUpdateFgdApplicants(req *request.CreateOrUpdateFgdApplicantsRequest) (*response.FgdScheduleResponse, error)
	UpdateStatusFgdApplicant(req *request.UpdateStatusFgdApplicantRequest) error
	UpdateFinalResultStatusFgdApplicant(req *request.UpdateFinalResultFgdApplicantRequest) error
	FindByUserProfileIDAndFgdID(userProfileID, fgdID string) (*response.FgdApplicantResponse, error)
	FindAllByFgdIDPaginated(fgdID string, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]response.FgdApplicantResponse, int64, error)
	UpdateFinalResultStatusTestApplicant(ctx context.Context, id uuid.UUID, finalResult entity.FinalResultStatus) (*response.FgdApplicantResponse, error)
}

type FgdApplicantUseCase struct {
	Log                   *logrus.Logger
	Repository            repository.IFgdApplicantRepository
	DTO                   dto.IFgdApplicantDTO
	FgdRepository         repository.IFgdScheduleRepository
	FgdDTO                dto.IFgdDTO
	UserProfileRepository repository.IUserProfileRepository
	Viper                 *viper.Viper
	ApplicantRepository   repository.IApplicantRepository
	JobPostingRepository  repository.IJobPostingRepository
}

func NewFgdApplicantUseCase(
	log *logrus.Logger,
	repository repository.IFgdApplicantRepository,
	iaDto dto.IFgdApplicantDTO,
	fgdRepository repository.IFgdScheduleRepository,
	fgdDTO dto.IFgdDTO,
	userProfileRepository repository.IUserProfileRepository,
	viper *viper.Viper,
	applicantRepo repository.IApplicantRepository,
	jobPostingRepo repository.IJobPostingRepository,
) IFgdApplicantUseCase {
	return &FgdApplicantUseCase{
		Log:                   log,
		Repository:            repository,
		DTO:                   iaDto,
		FgdRepository:         fgdRepository,
		FgdDTO:                fgdDTO,
		UserProfileRepository: userProfileRepository,
		Viper:                 viper,
		ApplicantRepository:   applicantRepo,
		JobPostingRepository:  jobPostingRepo,
	}
}

func FgdApplicantUseCaseFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IFgdApplicantUseCase {
	iaRepo := repository.FgdApplicantRepositoryFactory(log)
	iaDto := dto.FgdApplicantDTOFactory(log, viper)
	fgdRepo := repository.FgdScheduleRepositoryFactory(log)
	FgdDTO := dto.FgdDTOFactory(log, viper)
	userProfileRepo := repository.UserProfileRepositoryFactory(log)
	applicantRepo := repository.ApplicantRepositoryFactory(log)
	jobPostingRepo := repository.JobPostingRepositoryFactory(log)
	return NewFgdApplicantUseCase(
		log,
		iaRepo,
		iaDto,
		fgdRepo,
		FgdDTO,
		userProfileRepo,
		viper,
		applicantRepo,
		jobPostingRepo,
	)
}

func (uc *FgdApplicantUseCase) CreateOrUpdateFgdApplicants(req *request.CreateOrUpdateFgdApplicantsRequest) (*response.FgdScheduleResponse, error) {
	// check if Fgd exist
	parsedFgdID, err := uuid.Parse(req.FgdID)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error parsing Fgd id: %v", err)
		return nil, err
	}
	Fgd, err := uc.FgdRepository.FindByID(parsedFgdID)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error finding Fgd by id: %v", err)
		return nil, err
	}
	if Fgd == nil {
		return nil, errors.New("Fgd not found")
	}

	// create or update Fgd applicants
	for _, ia := range req.FgdApplicants {
		parsedUserProfileID, err := uuid.Parse(ia.UserProfileID)
		if err != nil {
			uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error parsing user profile id: %v", err)
			return nil, err
		}
		up, err := uc.UserProfileRepository.FindByID(parsedUserProfileID)
		if err != nil {
			uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error finding user profile by id: %v", err)
			return nil, err
		}
		if up == nil {
			return nil, errors.New("user profile not found")
		}

		parsedApplicantID, err := uuid.Parse(ia.ApplicantID)
		if err != nil {
			uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error parsing applicant id: %v", err)
			return nil, err
		}
		applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
			"id": parsedApplicantID,
		})
		if err != nil {
			uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error finding applicant by id: %v", err)
			return nil, err
		}
		if applicant == nil {
			return nil, errors.New("applicant not found")
		}

		parsedStartTime, err := time.Parse("2006-01-02 15:04:05", ia.StartTime)
		if err != nil {
			uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error parsing start time: %v", err)
			return nil, err
		}

		parsedEndTime, err := time.Parse("2006-01-02 15:04:05", ia.EndTime)
		if err != nil {
			uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error parsing end time: %v", err)
			return nil, err
		}

		if ia.ID != "" && ia.ID != uuid.Nil.String() {
			parsedIaID, err := uuid.Parse(ia.ID)
			if err != nil {
				uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error parsing Fgd applicant id: %v", err)
				return nil, err
			}
			exist, err := uc.Repository.FindByID(parsedIaID)
			if err != nil {
				uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error finding Fgd applicant by id: %v", err)
				return nil, err
			}
			if exist == nil {
				_, err := uc.Repository.CreateFgdApplicant(&entity.FgdApplicant{
					FgdScheduleID: parsedFgdID,
					ApplicantID:   parsedApplicantID,
					UserProfileID: parsedUserProfileID,
					StartTime:     parsedStartTime,
					EndTime:       parsedEndTime,
					FinalResult:   entity.FinalResultStatus(ia.FinalResult),
				})
				if err != nil {
					uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error creating Fgd applicant: %v", err)
					return nil, err
				}
			} else {
				_, err := uc.Repository.UpdateFgdApplicant(&entity.FgdApplicant{
					ID:            parsedIaID,
					FgdScheduleID: parsedFgdID,
					ApplicantID:   parsedApplicantID,
					UserProfileID: parsedUserProfileID,
					StartTime:     parsedStartTime,
					EndTime:       parsedEndTime,
					FinalResult:   entity.FinalResultStatus(ia.FinalResult),
				})
				if err != nil {
					uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error updating Fgd applicant: %v", err)
					return nil, err
				}

				findUpdatedFgdApplicant, err := uc.Repository.FindByID(parsedIaID)
				if err != nil {
					uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error finding updated Fgd applicant: %v", err)
					return nil, err
				}
				if findUpdatedFgdApplicant == nil {
					return nil, errors.New("Fgd applicant not found")
				}
				if findUpdatedFgdApplicant.FinalResult == entity.FINAL_RESULT_STATUS_ACCEPTED || findUpdatedFgdApplicant.FinalResult == entity.FINAL_RESULT_STATUS_REJECTED || findUpdatedFgdApplicant.FinalResult == entity.FINAL_RESULT_STATUS_SHORTLISTED {
					continue
				}

				jpExist, err := uc.JobPostingRepository.FindByID(Fgd.JobPostingID)
				if err != nil {
					uc.Log.Error("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] " + err.Error())
					return nil, err
				}

				if jpExist == nil {
					uc.Log.Error("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] " + "Job Posting not found")
					return nil, errors.New("job posting not found")
				}

				applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
					"id": parsedApplicantID,
				})
				if err != nil {
					uc.Log.Error("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] " + err.Error())
					return nil, err
				}

				if applicant == nil {
					uc.Log.Error("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] " + "Applicant not found")
					return nil, errors.New("applicant not found")
				}

				if entity.FinalResultStatus(ia.FinalResult) == entity.FINAL_RESULT_STATUS_ACCEPTED {
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
						Order:              applicant.Order + 1,
						TemplateQuestionID: *TemplateQuestionID,
					})
					if err != nil {
						uc.Log.Error("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] " + err.Error())
						return nil, err
					}
				} else if entity.FinalResultStatus(ia.FinalResult) == entity.FINAL_RESULT_STATUS_REJECTED {
					_, err = uc.ApplicantRepository.UpdateApplicantWhenRejected(&entity.Applicant{
						ID: applicant.ID,
					})
					if err != nil {
						uc.Log.Error("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] " + err.Error())
						return nil, err
					}
				} else if entity.FinalResultStatus(ia.FinalResult) == entity.FINAL_RESULT_STATUS_SHORTLISTED {
					_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
						ID:     applicant.ID,
						Status: entity.APPLICANT_STATUS_SHORTLIST,
					})
					if err != nil {
						uc.Log.Error("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] " + err.Error())
						return nil, err
					}
				}
			}
		} else {
			_, err := uc.Repository.CreateFgdApplicant(&entity.FgdApplicant{
				FgdScheduleID: parsedFgdID,
				ApplicantID:   parsedApplicantID,
				UserProfileID: parsedUserProfileID,
				StartTime:     parsedStartTime,
				EndTime:       parsedEndTime,
				FinalResult:   entity.FinalResultStatus(ia.FinalResult),
			})
			if err != nil {
				uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error creating Fgd applicant: %v", err)
				return nil, err
			}
		}
	}

	// delete Fgd applicants
	if len(req.DeletedFgdApplicantIDs) > 0 {
		for _, id := range req.DeletedFgdApplicantIDs {
			parsedID, err := uuid.Parse(id)
			if err != nil {
				uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error parsing Fgd applicant id: %v", err)
				return nil, err
			}
			err = uc.Repository.DeleteFgdApplicant(parsedID)
			if err != nil {
				uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error deleting Fgd applicant: %v", err)
				return nil, err
			}
		}
	}

	// get Fgd
	FgdRes, err := uc.FgdRepository.FindByID(parsedFgdID)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error finding Fgd by id: %v", err)
		return nil, err
	}
	if FgdRes == nil {
		return nil, errors.New("Fgd not found")
	}

	// convert Fgd to response
	resp, err := uc.FgdDTO.ConvertEntityToResponse(FgdRes)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.CreateOrUpdateFgdApplicants] error converting FgdSchedule to response: %v", err)
		return nil, err
	}

	return resp, nil
}

func (uc *FgdApplicantUseCase) UpdateStatusFgdApplicant(req *request.UpdateStatusFgdApplicantRequest) error {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.UpdateStatusFgdApplicant] error parsing Fgd applicant id: %v", err.Error())
		return err
	}

	ia, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.UpdateStatusFgdApplicant] error finding Fgd applicant by id: %v", err.Error())
		return err
	}

	if ia == nil {
		return errors.New("Fgd applicant not found")
	}

	_, err = uc.Repository.UpdateFgdApplicant(&entity.FgdApplicant{
		ID:               parsedID,
		AssessmentStatus: entity.AssessmentStatus(req.Status),
	})
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.UpdateStatusFgdApplicant] error updating Fgd applicant status: %v", err.Error())
		return err
	}

	return nil
}

func (uc *FgdApplicantUseCase) UpdateFinalResultStatusFgdApplicant(req *request.UpdateFinalResultFgdApplicantRequest) error {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.UpdateFinalResultStatusFgdApplicant] error parsing Fgd applicant id: %v", err.Error())
		return err
	}

	ia, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.UpdateFinalResultStatusFgdApplicant] error finding Fgd applicant by id: %v", err.Error())
		return err
	}

	if ia == nil {
		return errors.New("Fgd applicant not found")
	}

	if ia.FinalResult == entity.FINAL_RESULT_STATUS_ACCEPTED || ia.FinalResult == entity.FINAL_RESULT_STATUS_REJECTED || ia.FinalResult == entity.FINAL_RESULT_STATUS_SHORTLISTED {
		return errors.New("final result status already set")
	}

	_, err = uc.Repository.UpdateFgdApplicant(&entity.FgdApplicant{
		ID:          parsedID,
		FinalResult: entity.FinalResultStatus(req.FinalResult),
	})
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.UpdateFinalResultStatusFgdApplicant] error updating Fgd applicant status: %v", err.Error())
		return err
	}

	jpExist, err := uc.JobPostingRepository.FindByID(ia.FgdSchedule.JobPostingID)
	if err != nil {
		uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusFgdApplicant] " + err.Error())
		return err
	}

	if jpExist == nil {
		uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusFgdApplicant] " + "Job Posting not found")
		return errors.New("job posting not found")
	}

	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": ia.ApplicantID,
	})
	if err != nil {
		uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusFgdApplicant] " + err.Error())
		return err
	}

	if applicant == nil {
		uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusFgdApplicant] " + "Applicant not found")
		return errors.New("applicant not found")
	}

	if entity.FinalResultStatus(req.FinalResult) == entity.FINAL_RESULT_STATUS_ACCEPTED {
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
			Order:              applicant.Order + 1,
			TemplateQuestionID: *TemplateQuestionID,
		})
		if err != nil {
			uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusFgdApplicant] " + err.Error())
			return err
		}
	} else if entity.FinalResultStatus(ia.FinalResult) == entity.FINAL_RESULT_STATUS_REJECTED {
		_, err = uc.ApplicantRepository.UpdateApplicantWhenRejected(&entity.Applicant{
			ID: applicant.ID,
		})
		if err != nil {
			uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusFgdApplicant] " + err.Error())
			return err
		}
	} else if entity.FinalResultStatus(ia.FinalResult) == entity.FINAL_RESULT_STATUS_SHORTLISTED {
		_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
			ID:     applicant.ID,
			Status: entity.APPLICANT_STATUS_SHORTLIST,
		})
		if err != nil {
			uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusFgdApplicant] " + err.Error())
			return err
		}
	}

	return nil
}

func (uc *FgdApplicantUseCase) FindByUserProfileIDAndFgdID(userProfileID, FgdID string) (*response.FgdApplicantResponse, error) {
	parsedUserProfileID, err := uuid.Parse(userProfileID)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.FindByUserProfileIDAndFgdID] error parsing user profile id: %v", err)
		return nil, err
	}

	parsedFgdID, err := uuid.Parse(FgdID)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.FindByUserProfileIDAndFgdID] error parsing Fgd id: %v", err)
		return nil, err
	}

	exist, err := uc.Repository.FindByKeys(map[string]interface{}{
		"user_profile_id": parsedUserProfileID,
		"Fgd_id":          parsedFgdID,
	})
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.FindByUserProfileIDAndFgdID] error finding Fgd applicant by user profile id and Fgd id: %v", err)
		return nil, err
	}
	if exist == nil {
		return nil, nil
	}

	resp, err := uc.DTO.ConvertEntityToResponse(exist)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.FindByUserProfileIDAndFgdID] error converting Fgd applicant to response: %v", err)
		return nil, err
	}

	return resp, nil
}

func (uc *FgdApplicantUseCase) FindAllByFgdIDPaginated(FgdID string, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]response.FgdApplicantResponse, int64, error) {
	parsedFgdID, err := uuid.Parse(FgdID)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.FindAllByFgdIDPaginated] error parsing Fgd id: %v", err)
		return nil, 0, err
	}

	exist, count, err := uc.Repository.FindAllByFgdIDPaginated(parsedFgdID, page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.FindAllByFgdIDPaginated] error finding Fgd applicants by Fgd id: %v", err)
		return nil, 0, err
	}

	var res []response.FgdApplicantResponse
	for _, ia := range exist {
		resp, err := uc.DTO.ConvertEntityToResponse(&ia)
		if err != nil {
			uc.Log.Errorf("[FgdApplicantUseCase.FindAllByFgdIDPaginated] error converting Fgd applicant to response: %v", err)
			return nil, 0, err
		}
		res = append(res, *resp)
	}

	return res, count, nil
}

func (uc *FgdApplicantUseCase) UpdateFinalResultStatusTestApplicant(ctx context.Context, id uuid.UUID, finalResult entity.FinalResultStatus) (*response.FgdApplicantResponse, error) {
	ia, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when finding test applicant by id: %s", err.Error())
		return nil, errors.New("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when finding test applicant by id: " + err.Error())
	}
	if ia == nil {
		uc.Log.Errorf("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] test applicant with id %s not found", id.String())
		return nil, errors.New("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] test applicant with id " + id.String() + " not found")
	}

	if ia.FinalResult == entity.FINAL_RESULT_STATUS_ACCEPTED || ia.FinalResult == entity.FINAL_RESULT_STATUS_REJECTED || ia.FinalResult == entity.FINAL_RESULT_STATUS_SHORTLISTED {
		return nil, errors.New("final result status already set")
	}

	jpExist, err := uc.JobPostingRepository.FindByID(ia.FgdSchedule.JobPostingID)
	if err != nil {
		uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] " + err.Error())
		return nil, err
	}

	if jpExist == nil {
		uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": ia.ApplicantID,
	})
	if err != nil {
		uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] " + err.Error())
		return nil, err
	}

	if applicant == nil {
		uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] " + "Applicant not found")
		return nil, errors.New("applicant not found")
	}

	if finalResult == entity.FINAL_RESULT_STATUS_ACCEPTED {
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
			Order:              applicant.Order + 1,
			TemplateQuestionID: *TemplateQuestionID,
		})
		if err != nil {
			uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] " + err.Error())
			return nil, err
		}
	} else if finalResult == entity.FINAL_RESULT_STATUS_REJECTED {
		// zero, err := strconv.Atoi("0")
		if err != nil {
			uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] " + err.Error())
			return nil, err
		}
		if applicant != nil {
			_, err = uc.ApplicantRepository.UpdateApplicantWhenRejected(&entity.Applicant{
				ID: applicant.ID,
			})
			if err != nil {
				uc.Log.Error("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] " + err.Error())
				return nil, err
			}
		}
	}

	_, err = uc.Repository.UpdateFgdApplicant(&entity.FgdApplicant{
		ID:          ia.ID,
		FinalResult: finalResult,
	})
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when updating test applicant: %s", err.Error())
		// tx.Rollback()
		return nil, errors.New("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when updating test applicant: " + err.Error())
	}

	resp, err := uc.DTO.ConvertEntityToResponse(ia)
	if err != nil {
		uc.Log.Errorf("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when converting test applicant entity to response: %s", err.Error())
		// tx.Rollback()
		return nil, errors.New("[FgdApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when converting test applicant entity to response: " + err.Error())
	}

	return resp, nil
}
