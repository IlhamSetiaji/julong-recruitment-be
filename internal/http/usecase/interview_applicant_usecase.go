package usecase

import (
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

type IInterviewApplicantUseCase interface {
	CreateOrUpdateInterviewApplicants(req *request.CreateOrUpdateInterviewApplicantsRequest) (*response.InterviewResponse, error)
	UpdateStatusInterviewApplicant(req *request.UpdateStatusInterviewApplicantRequest) error
	UpdateFinalResultStatusInterviewApplicant(req *request.UpdateFinalResultInterviewApplicantRequest) error
	FindByUserProfileIDAndInterviewID(userProfileID, interviewID string) (*response.InterviewApplicantResponse, error)
	FindAllByInterviewIDPaginated(interviewID string, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]response.InterviewApplicantResponse, int64, error)
}

type InterviewApplicantUseCase struct {
	Log                   *logrus.Logger
	Repository            repository.IInterviewApplicantRepository
	DTO                   dto.IInterviewApplicantDTO
	InterviewRepository   repository.IInterviewRepository
	InterviewDTO          dto.IInterviewDTO
	UserProfileRepository repository.IUserProfileRepository
	Viper                 *viper.Viper
	ApplicantRepository   repository.IApplicantRepository
	JobPostingRepository  repository.IJobPostingRepository
}

func NewInterviewApplicantUseCase(
	log *logrus.Logger,
	repository repository.IInterviewApplicantRepository,
	iaDto dto.IInterviewApplicantDTO,
	interviewRepository repository.IInterviewRepository,
	interviewDTO dto.IInterviewDTO,
	userProfileRepository repository.IUserProfileRepository,
	viper *viper.Viper,
	applicantRepo repository.IApplicantRepository,
	jobPostingRepo repository.IJobPostingRepository,
) IInterviewApplicantUseCase {
	return &InterviewApplicantUseCase{
		Log:                   log,
		Repository:            repository,
		DTO:                   iaDto,
		InterviewRepository:   interviewRepository,
		InterviewDTO:          interviewDTO,
		UserProfileRepository: userProfileRepository,
		Viper:                 viper,
		ApplicantRepository:   applicantRepo,
		JobPostingRepository:  jobPostingRepo,
	}
}

func InterviewApplicantUseCaseFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IInterviewApplicantUseCase {
	iaRepo := repository.InterviewApplicantRepositoryFactory(log)
	iaDto := dto.InterviewApplicantDTOFactory(log, viper)
	interviewRepo := repository.InterviewRepositoryFactory(log)
	interviewDTO := dto.InterviewDTOFactory(log, viper)
	userProfileRepo := repository.UserProfileRepositoryFactory(log)
	applicantRepo := repository.ApplicantRepositoryFactory(log)
	jobPostingRepo := repository.JobPostingRepositoryFactory(log)
	return NewInterviewApplicantUseCase(
		log,
		iaRepo,
		iaDto,
		interviewRepo,
		interviewDTO,
		userProfileRepo,
		viper,
		applicantRepo,
		jobPostingRepo,
	)
}

func (uc *InterviewApplicantUseCase) CreateOrUpdateInterviewApplicants(req *request.CreateOrUpdateInterviewApplicantsRequest) (*response.InterviewResponse, error) {
	// check if interview exist
	parsedInterviewID, err := uuid.Parse(req.InterviewID)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error parsing interview id: %v", err)
		return nil, err
	}
	interview, err := uc.InterviewRepository.FindByID(parsedInterviewID)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error finding interview by id: %v", err)
		return nil, err
	}
	if interview == nil {
		return nil, errors.New("interview not found")
	}

	// create or update interview applicants
	for _, ia := range req.InterviewApplicants {
		parsedUserProfileID, err := uuid.Parse(ia.UserProfileID)
		if err != nil {
			uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error parsing user profile id: %v", err)
			return nil, err
		}
		up, err := uc.UserProfileRepository.FindByID(parsedUserProfileID)
		if err != nil {
			uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error finding user profile by id: %v", err)
			return nil, err
		}
		if up == nil {
			return nil, errors.New("user profile not found")
		}

		parsedApplicantID, err := uuid.Parse(ia.ApplicantID)
		if err != nil {
			uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error parsing applicant id: %v", err)
			return nil, err
		}
		applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
			"id": parsedApplicantID,
		})
		if err != nil {
			uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error finding applicant by id: %v", err)
			return nil, err
		}
		if applicant == nil {
			return nil, errors.New("applicant not found")
		}

		parsedStartTime, err := time.Parse("2006-01-02 15:04:05", ia.StartTime)
		if err != nil {
			uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error parsing start time: %v", err)
			return nil, err
		}

		parsedEndTime, err := time.Parse("2006-01-02 15:04:05", ia.EndTime)
		if err != nil {
			uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error parsing end time: %v", err)
			return nil, err
		}

		if ia.ID != "" && ia.ID != uuid.Nil.String() {
			parsedIaID, err := uuid.Parse(ia.ID)
			if err != nil {
				uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error parsing interview applicant id: %v", err)
				return nil, err
			}
			exist, err := uc.Repository.FindByID(parsedIaID)
			if err != nil {
				uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error finding interview applicant by id: %v", err)
				return nil, err
			}
			if exist == nil {
				_, err := uc.Repository.CreateInterviewApplicant(&entity.InterviewApplicant{
					InterviewID:   parsedInterviewID,
					ApplicantID:   parsedApplicantID,
					UserProfileID: parsedUserProfileID,
					StartTime:     parsedStartTime,
					EndTime:       parsedEndTime,
					FinalResult:   entity.FinalResultStatus(ia.FinalResult),
				})
				if err != nil {
					uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error creating interview applicant: %v", err)
					return nil, err
				}
			} else {
				_, err := uc.Repository.UpdateInterviewApplicant(&entity.InterviewApplicant{
					ID:            parsedIaID,
					InterviewID:   parsedInterviewID,
					ApplicantID:   parsedApplicantID,
					UserProfileID: parsedUserProfileID,
					StartTime:     parsedStartTime,
					EndTime:       parsedEndTime,
					FinalResult:   entity.FinalResultStatus(ia.FinalResult),
				})
				if err != nil {
					uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error updating interview applicant: %v", err)
					return nil, err
				}

				jpExist, err := uc.JobPostingRepository.FindByID(interview.JobPostingID)
				if err != nil {
					uc.Log.Error("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] " + err.Error())
					return nil, err
				}

				if jpExist == nil {
					uc.Log.Error("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] " + "Job Posting not found")
					return nil, errors.New("job posting not found")
				}

				applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
					"id": parsedApplicantID,
				})
				if err != nil {
					uc.Log.Error("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] " + err.Error())
					return nil, err
				}

				if applicant == nil {
					uc.Log.Error("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] " + "Applicant not found")
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
						uc.Log.Error("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] " + err.Error())
						return nil, err
					}
				} else if entity.FinalResultStatus(ia.FinalResult) == entity.FINAL_RESULT_STATUS_REJECTED {
					_, err = uc.ApplicantRepository.UpdateApplicantWhenRejected(&entity.Applicant{
						ID: applicant.ID,
					})
					if err != nil {
						uc.Log.Error("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] " + err.Error())
						return nil, err
					}
				} else if entity.FinalResultStatus(ia.FinalResult) == entity.FINAL_RESULT_STATUS_SHORTLISTED {
					_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
						ID:     applicant.ID,
						Status: entity.APPLICANT_STATUS_SHORTLIST,
					})
					if err != nil {
						uc.Log.Error("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] " + err.Error())
						return nil, err
					}
				}
			}
		} else {
			_, err := uc.Repository.CreateInterviewApplicant(&entity.InterviewApplicant{
				InterviewID:   parsedInterviewID,
				ApplicantID:   parsedApplicantID,
				UserProfileID: parsedUserProfileID,
				StartTime:     parsedStartTime,
				EndTime:       parsedEndTime,
				FinalResult:   entity.FinalResultStatus(ia.FinalResult),
			})
			if err != nil {
				uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error creating interview applicant: %v", err)
				return nil, err
			}
		}
	}

	// delete interview applicants
	if len(req.DeletedInterviewApplicantIDs) > 0 {
		for _, id := range req.DeletedInterviewApplicantIDs {
			parsedID, err := uuid.Parse(id)
			if err != nil {
				uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error parsing interview applicant id: %v", err)
				return nil, err
			}
			err = uc.Repository.DeleteInterviewApplicant(parsedID)
			if err != nil {
				uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error deleting interview applicant: %v", err)
				return nil, err
			}
		}
	}

	// get interview
	interviewRes, err := uc.InterviewRepository.FindByID(parsedInterviewID)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error finding interview by id: %v", err)
		return nil, err
	}
	if interviewRes == nil {
		return nil, errors.New("interview not found")
	}

	// convert interview to response
	resp, err := uc.InterviewDTO.ConvertEntityToResponse(interviewRes)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.CreateOrUpdateInterviewApplicants] error converting interview to response: %v", err)
		return nil, err
	}

	return resp, nil
}

func (uc *InterviewApplicantUseCase) UpdateStatusInterviewApplicant(req *request.UpdateStatusInterviewApplicantRequest) error {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.UpdateStatusInterviewApplicant] error parsing interview applicant id: %v", err.Error())
		return err
	}

	ia, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.UpdateStatusInterviewApplicant] error finding interview applicant by id: %v", err.Error())
		return err
	}

	if ia == nil {
		return errors.New("interview applicant not found")
	}

	_, err = uc.Repository.UpdateInterviewApplicant(&entity.InterviewApplicant{
		ID:               parsedID,
		AssessmentStatus: entity.AssessmentStatus(req.Status),
	})
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.UpdateStatusInterviewApplicant] error updating interview applicant status: %v", err.Error())
		return err
	}

	return nil
}

func (uc *InterviewApplicantUseCase) UpdateFinalResultStatusInterviewApplicant(req *request.UpdateFinalResultInterviewApplicantRequest) error {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] error parsing interview applicant id: %v", err.Error())
		return err
	}

	ia, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] error finding interview applicant by id: %v", err.Error())
		return err
	}

	if ia == nil {
		return errors.New("interview applicant not found")
	}

	if ia.FinalResult == entity.FINAL_RESULT_STATUS_ACCEPTED || ia.FinalResult == entity.FINAL_RESULT_STATUS_REJECTED || ia.FinalResult == entity.FINAL_RESULT_STATUS_SHORTLISTED {
		return errors.New("final result status already set")
	}

	_, err = uc.Repository.UpdateInterviewApplicant(&entity.InterviewApplicant{
		ID:          parsedID,
		FinalResult: entity.FinalResultStatus(req.FinalResult),
	})
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] error updating interview applicant status: %v", err.Error())
		return err
	}

	jpExist, err := uc.JobPostingRepository.FindByID(ia.Interview.JobPostingID)
	if err != nil {
		uc.Log.Error("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] " + err.Error())
		return err
	}

	if jpExist == nil {
		uc.Log.Error("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] " + "Job Posting not found")
		return errors.New("job posting not found")
	}

	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": ia.ApplicantID,
	})
	if err != nil {
		uc.Log.Error("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] " + err.Error())
		return err
	}

	if applicant == nil {
		uc.Log.Error("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] " + "Applicant not found")
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
			uc.Log.Error("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] " + err.Error())
			return err
		}
	} else if entity.FinalResultStatus(ia.FinalResult) == entity.FINAL_RESULT_STATUS_REJECTED {
		_, err = uc.ApplicantRepository.UpdateApplicantWhenRejected(&entity.Applicant{
			ID: applicant.ID,
		})
		if err != nil {
			uc.Log.Error("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] " + err.Error())
			return err
		}
	} else if entity.FinalResultStatus(ia.FinalResult) == entity.FINAL_RESULT_STATUS_SHORTLISTED {
		_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
			ID:     applicant.ID,
			Status: entity.APPLICANT_STATUS_SHORTLIST,
		})
		if err != nil {
			uc.Log.Error("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] " + err.Error())
			return err
		}
	}

	return nil
}

func (uc *InterviewApplicantUseCase) FindByUserProfileIDAndInterviewID(userProfileID, interviewID string) (*response.InterviewApplicantResponse, error) {
	parsedUserProfileID, err := uuid.Parse(userProfileID)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.FindByUserProfileIDAndInterviewID] error parsing user profile id: %v", err)
		return nil, err
	}

	parsedInterviewID, err := uuid.Parse(interviewID)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.FindByUserProfileIDAndInterviewID] error parsing interview id: %v", err)
		return nil, err
	}

	exist, err := uc.Repository.FindByKeys(map[string]interface{}{
		"user_profile_id": parsedUserProfileID,
		"interview_id":    parsedInterviewID,
	})
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.FindByUserProfileIDAndInterviewID] error finding interview applicant by user profile id and interview id: %v", err)
		return nil, err
	}
	if exist == nil {
		return nil, nil
	}

	resp, err := uc.DTO.ConvertEntityToResponse(exist)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.FindByUserProfileIDAndInterviewID] error converting interview applicant to response: %v", err)
		return nil, err
	}

	return resp, nil
}

func (uc *InterviewApplicantUseCase) FindAllByInterviewIDPaginated(interviewID string, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]response.InterviewApplicantResponse, int64, error) {
	parsedInterviewID, err := uuid.Parse(interviewID)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.FindAllByInterviewIDPaginated] error parsing interview id: %v", err)
		return nil, 0, err
	}

	exist, count, err := uc.Repository.FindAllByInterviewIDPaginated(parsedInterviewID, page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Errorf("[InterviewApplicantUseCase.FindAllByInterviewIDPaginated] error finding interview applicants by interview id: %v", err)
		return nil, 0, err
	}

	var res []response.InterviewApplicantResponse
	for _, ia := range exist {
		resp, err := uc.DTO.ConvertEntityToResponse(&ia)
		if err != nil {
			uc.Log.Errorf("[InterviewApplicantUseCase.FindAllByInterviewIDPaginated] error converting interview applicant to response: %v", err)
			return nil, 0, err
		}
		res = append(res, *resp)
	}

	return res, count, nil
}
