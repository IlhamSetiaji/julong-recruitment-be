package usecase

import (
	"errors"
	"fmt"
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

type IFgdScheduleUseCase interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.FgdScheduleResponse, int64, error)
	CreateFgdSchedule(req *request.CreateFgdScheduleRequest) (*response.FgdScheduleResponse, error)
	UpdateFgdSchedule(req *request.UpdateFgdScheduleRequest) (*response.FgdScheduleResponse, error)
	FindByID(id uuid.UUID) (*response.FgdScheduleResponse, error)
	DeleteByID(id uuid.UUID) error
	GenerateDocumentNumber(dateNow time.Time) (string, error)
	UpdateStatusFgdSchedule(req *request.UpdateStatusFgdScheduleRequest) (*response.FgdScheduleResponse, error)
	FindMySchedule(userID, projectRecruitmentLineID, jobPostingID uuid.UUID) (*response.FgdScheduleMyselfResponse, error)
	FindMyScheduleForAssessor(employeeID, projectRecruitmentLineID, jobPostingID uuid.UUID) (*[]response.FgdScheduleMyselfForAssessorResponse, error)
	FindScheduleForApplicant(applicantID, projectRecruitmentLineID, jobPostingID, employeeID uuid.UUID) (*response.FgdScheduleMyselfResponse, error)
	FindByIDForAnswer(id, jobPostingID uuid.UUID) (*response.FgdScheduleResponse, error)
}

type FgdScheduleUseCase struct {
	Log                                *logrus.Logger
	Repository                         repository.IFgdScheduleRepository
	DTO                                dto.IFgdDTO
	Viper                              *viper.Viper
	JobPostingRepository               repository.IJobPostingRepository
	ProjectPicRepository               repository.IProjectPicRepository
	ProjectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository
	ProjectRecruitmentLineRepository   repository.IProjectRecruitmentLineRepository
	UserProfileRepository              repository.IUserProfileRepository
	FgdAssessorRepository              repository.IFgdAssessorRepository
	FgdApplicantRepository             repository.IFgdApplicantRepository
	ApplicantRepository                repository.IApplicantRepository
}

func NewFgdUseCase(
	log *logrus.Logger,
	repository repository.IFgdScheduleRepository,
	dto dto.IFgdDTO,
	viper *viper.Viper,
	jobPostingRepository repository.IJobPostingRepository,
	projectPicRepository repository.IProjectPicRepository,
	projectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository,
	projectRecruitmentLineRepository repository.IProjectRecruitmentLineRepository,
	userProfileRepository repository.IUserProfileRepository,
	fgdAssessorRepository repository.IFgdAssessorRepository,
	fgdApplicantRepository repository.IFgdApplicantRepository,
	applicantRepository repository.IApplicantRepository,
) IFgdScheduleUseCase {
	return &FgdScheduleUseCase{
		Log:                                log,
		Repository:                         repository,
		DTO:                                dto,
		Viper:                              viper,
		JobPostingRepository:               jobPostingRepository,
		ProjectPicRepository:               projectPicRepository,
		ProjectRecruitmentHeaderRepository: projectRecruitmentHeaderRepository,
		ProjectRecruitmentLineRepository:   projectRecruitmentLineRepository,
		UserProfileRepository:              userProfileRepository,
		FgdAssessorRepository:              fgdAssessorRepository,
		FgdApplicantRepository:             fgdApplicantRepository,
		ApplicantRepository:                applicantRepository,
	}
}

func FgdScheduleUseCaseFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IFgdScheduleUseCase {
	FgdRepository := repository.FgdScheduleRepositoryFactory(log)
	dto := dto.FgdDTOFactory(log, viper)
	jobPostingRepository := repository.JobPostingRepositoryFactory(log)
	projectPicRepository := repository.ProjectPicRepositoryFactory(log)
	projectRecruitmentHeaderRepository := repository.ProjectRecruitmentHeaderRepositoryFactory(log)
	projectRecruitmentLineRepository := repository.ProjectRecruitmentLineRepositoryFactory(log)
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	fgdAssessorRepository := repository.FgdAssessorRepositoryFactory(log)
	fgdApplicantRepository := repository.FgdApplicantRepositoryFactory(log)
	applicantRepository := repository.ApplicantRepositoryFactory(log)
	return NewFgdUseCase(
		log,
		FgdRepository,
		dto,
		viper,
		jobPostingRepository,
		projectPicRepository,
		projectRecruitmentHeaderRepository,
		projectRecruitmentLineRepository,
		userProfileRepository,
		fgdAssessorRepository,
		fgdApplicantRepository,
		applicantRepository,
	)
}

func (uc *FgdScheduleUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.FgdScheduleResponse, int64, error) {
	fgdSchedules, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		return nil, 0, err
	}

	fgdScheduleResponses := make([]response.FgdScheduleResponse, 0)
	for _, fgdSchedule := range *fgdSchedules {
		resp, err := uc.DTO.ConvertEntityToResponse(&fgdSchedule)
		if err != nil {
			uc.Log.Error("[FgdScheduleUseCase.FindAllPaginated] " + err.Error())
			return nil, 0, err
		}
		fgdScheduleResponses = append(fgdScheduleResponses, *resp)
	}

	return &fgdScheduleResponses, total, nil
}

func (uc *FgdScheduleUseCase) CreateFgdSchedule(req *request.CreateFgdScheduleRequest) (*response.FgdScheduleResponse, error) {
	parsedJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] Job Posting not found")
		return nil, err
	}

	parsedProjectPicID, err := uuid.Parse(req.ProjectPicID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdSchedule] " + err.Error())
		return nil, err
	}
	projectPic, err := uc.ProjectPicRepository.FindByID(parsedProjectPicID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdSchedule] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdSchedule] " + err.Error())
		return nil, err
	}
	if projectPic == nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdSchedule] Project PIC not found")
		return nil, err
	}

	combinedDateStart := req.ScheduleDate + " " + req.StartTime
	combinedDateEnd := req.ScheduleDate + " " + req.EndTime

	parsedStartTime, err := time.Parse("2006-01-02 15:04:05", combinedDateStart)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	parsedEndTime, err := time.Parse("2006-01-02 15:04:05", combinedDateEnd)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	if parsedStartTime.After(parsedEndTime) {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + "Start time must be before end time")
		return nil, errors.New("start time must be before end time")
	}

	parsedScheduleDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	if parsedScheduleDate.Before(jobPosting.StartDate) && parsedScheduleDate.After(jobPosting.EndDate) {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + "Schedule date must be between job posting start date and end date")
		return nil, errors.New("schedule date must be between job posting start date and end date [Start Date: " + jobPosting.StartDate.String() + ", End Date: " + jobPosting.EndDate.String() + "]")
	}

	parsedPrhID, err := uuid.Parse(req.ProjectRecruitmentHeaderID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	parsedPrlID, err := uuid.Parse(req.ProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	fgdSchedule, err := uc.Repository.CreateFgdSchedule(&entity.FgdSchedule{
		JobPostingID:               parsedJobPostingID,
		ProjectPicID:               parsedProjectPicID,
		ProjectRecruitmentHeaderID: parsedPrhID,
		ProjectRecruitmentLineID:   parsedPrlID,
		Name:                       req.Name,
		DocumentNumber:             req.DocumentNumber,
		ScheduleDate:               parsedScheduleDate,
		StartTime:                  parsedStartTime,
		EndTime:                    parsedEndTime,
		LocationLink:               req.LocationLink,
		Description:                req.Description,
		RangeDuration:              req.RangeDuration,
		TotalCandidate:             req.TotalCandidate,
		Status:                     entity.FgdScheduleStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	// insert FgdSchedule assessors
	if len(req.FgdScheduleAssessors) > 0 {
		for _, assessor := range req.FgdScheduleAssessors {
			parsedEmployeeID, err := uuid.Parse(assessor.EmployeeID)
			if err != nil {
				uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
				return nil, err
			}

			_, err = uc.FgdAssessorRepository.CreateFgdAssessor(&entity.FgdAssessor{
				FgdScheduleID: fgdSchedule.ID,
				EmployeeID:    &parsedEmployeeID,
			})
			if err != nil {
				uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
				return nil, err
			}
		}
	}

	// get applicants
	applicantsPayload, err := uc.getApplicantIDsByJobPostingID(parsedJobPostingID, parsedPrlID, 1, req.TotalCandidate)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	if applicantsPayload.Total == 0 {
		err := uc.Repository.DeleteFgdSchedule(fgdSchedule.ID)
		if err != nil {
			uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
			return nil, err
		}

		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + "No applicants found")
		return nil, errors.New("no applicants found")
	}

	// insert FgdSchedule applicants
	for i, applicantID := range applicantsPayload.ApplicantIDs {
		_, err = uc.FgdApplicantRepository.CreateFgdApplicant(&entity.FgdApplicant{
			FgdScheduleID:    fgdSchedule.ID,
			ApplicantID:      applicantID,
			UserProfileID:    applicantsPayload.UserProfileIDs[i],
			StartTime:        parsedStartTime,
			EndTime:          parsedEndTime,
			AssessmentStatus: entity.ASSESSMENT_STATUS_DRAFT,
			FinalResult:      entity.FINAL_RESULT_STATUS_DRAFT,
		})
		if err != nil {
			uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
			return nil, err
		}
	}

	if applicantsPayload.Total < req.TotalCandidate {
		_, err = uc.Repository.UpdateFgdSchedule(&entity.FgdSchedule{
			ID:             fgdSchedule.ID,
			TotalCandidate: applicantsPayload.Total,
		})
		if err != nil {
			uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
			return nil, err
		}
	}

	findByID, err := uc.Repository.FindByID(fgdSchedule.ID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}
	if findByID == nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] FgdSchedule not found")
		return nil, errors.New("FgdSchedule not found")
	}

	resp, err := uc.DTO.ConvertEntityToResponse(findByID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	uc.Log.Error("Total Candidate: ", applicantsPayload)

	return resp, nil
}

func (uc *FgdScheduleUseCase) UpdateFgdSchedule(req *request.UpdateFgdScheduleRequest) (*response.FgdScheduleResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	parsedJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] Job Posting not found")
		return nil, err
	}

	parsedProjectPicID, err := uuid.Parse(req.ProjectPicID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdSchedule] " + err.Error())
		return nil, err
	}
	projectPic, err := uc.ProjectPicRepository.FindByID(parsedProjectPicID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdSchedule] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdSchedule] " + err.Error())
		return nil, err
	}
	if projectPic == nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdSchedule] Project PIC not found")
		return nil, err
	}

	parsedStartTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	parsedEndTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	if parsedStartTime.After(parsedEndTime) {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + "Start time must be before end time")
		return nil, errors.New("start time must be before end time")
	}

	parsedScheduleDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	if parsedScheduleDate.Before(jobPosting.StartDate) && parsedScheduleDate.After(jobPosting.EndDate) {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + "Schedule date must be between job posting start date and end date")
		return nil, errors.New("schedule date must be between job posting start date and end date [Start Date: " + jobPosting.StartDate.String() + ", End Date: " + jobPosting.EndDate.String() + "]")
	}

	parsedPrhID, err := uuid.Parse(req.ProjectRecruitmentHeaderID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	parsedPrlID, err := uuid.Parse(req.ProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	fgdSchedule, err := uc.Repository.UpdateFgdSchedule(&entity.FgdSchedule{
		ID:                         parsedID,
		JobPostingID:               parsedJobPostingID,
		ProjectPicID:               parsedProjectPicID,
		ProjectRecruitmentHeaderID: parsedPrhID,
		ProjectRecruitmentLineID:   parsedPrlID,
		Name:                       req.Name,
		DocumentNumber:             req.DocumentNumber,
		ScheduleDate:               parsedScheduleDate,
		StartTime:                  parsedStartTime,
		EndTime:                    parsedEndTime,
		LocationLink:               req.LocationLink,
		Description:                req.Description,
		RangeDuration:              req.RangeDuration,
		Status:                     entity.FgdScheduleStatus(req.Status),
	})

	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	// delete FgdSchedule assessors
	err = uc.FgdAssessorRepository.DeleteFgdAssessorByFgdID(parsedID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	// insert FgdSchedule assessors
	if len(req.FgdScheduleAssessors) > 0 {
		if len(req.FgdScheduleAssessors) > 0 {
			for _, assessor := range req.FgdScheduleAssessors {
				parsedEmployeeID, err := uuid.Parse(assessor.EmployeeID)
				if err != nil {
					uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
					return nil, err
				}

				_, err = uc.FgdAssessorRepository.CreateFgdAssessor(&entity.FgdAssessor{
					FgdScheduleID: fgdSchedule.ID,
					EmployeeID:    &parsedEmployeeID,
				})
				if err != nil {
					uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
					return nil, err
				}
			}
		}
	}

	if req.Status == string(entity.FGD_SCHEDULE_STATUS_DRAFT) {
		// delete FgdSchedule applicants
		err = uc.FgdApplicantRepository.DeleteByFgdScheduleID(parsedID)
		if err != nil {
			uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
			return nil, err
		}
		// get applicants
		applicantsPayload, err := uc.getApplicantIDsByJobPostingID(parsedJobPostingID, parsedPrlID, 1, req.TotalCandidate)
		if err != nil {
			uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
			return nil, err
		}

		if applicantsPayload.Total == 0 {
			err := uc.Repository.DeleteFgdSchedule(fgdSchedule.ID)
			if err != nil {
				uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
				return nil, err
			}

			uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + "No applicants found")
			return nil, errors.New("no applicants found")
		}

		// insert FgdSchedule applicants
		for i, applicantID := range applicantsPayload.ApplicantIDs {
			_, err = uc.FgdApplicantRepository.CreateFgdApplicant(&entity.FgdApplicant{
				FgdScheduleID:    fgdSchedule.ID,
				ApplicantID:      applicantID,
				UserProfileID:    applicantsPayload.UserProfileIDs[i],
				StartTime:        parsedStartTime,
				EndTime:          parsedEndTime,
				AssessmentStatus: entity.ASSESSMENT_STATUS_DRAFT,
				FinalResult:      entity.FINAL_RESULT_STATUS_DRAFT,
			})
			if err != nil {
				uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
				return nil, err
			}
		}

		if applicantsPayload.Total < req.TotalCandidate {
			_, err = uc.Repository.UpdateFgdSchedule(&entity.FgdSchedule{
				ID:             fgdSchedule.ID,
				TotalCandidate: applicantsPayload.Total,
			})
			if err != nil {
				uc.Log.Error("[FgdScheduleUseCase.CreateFgdScheduleRequest] " + err.Error())
				return nil, err
			}
		}
	}

	findByID, err := uc.Repository.FindByID(fgdSchedule.ID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	resp, err := uc.DTO.ConvertEntityToResponse(findByID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (uc *FgdScheduleUseCase) FindByID(id uuid.UUID) (*response.FgdScheduleResponse, error) {
	fgdSchedule, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindByID] " + err.Error())
		return nil, err
	}

	if fgdSchedule == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindByID] fgdSchedule not found")
		return nil, errors.New("fgdSchedule not found")
	}

	resp, err := uc.DTO.ConvertEntityToResponse(fgdSchedule)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (uc *FgdScheduleUseCase) DeleteByID(id uuid.UUID) error {
	fgdSchedule, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.DeleteByID] " + err.Error())
		return err
	}

	if fgdSchedule == nil {
		uc.Log.Error("[FgdScheduleUseCase.DeleteByID] fgdSchedule not found")
		return errors.New("fgdSchedule not found")
	}

	err = uc.Repository.DeleteFgdSchedule(id)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.DeleteByID] " + err.Error())
		return err
	}

	return nil
}

func (uc *FgdScheduleUseCase) GenerateDocumentNumber(dateNow time.Time) (string, error) {
	dateStr := dateNow.Format("2006-01-02")
	highestNumber, err := uc.Repository.GetHighestDocumentNumberByDate(dateStr)
	if err != nil {
		uc.Log.Errorf("[FgdScheduleUseCase.GenerateDocumentNumber] " + err.Error())
		return "", err
	}

	newNumber := highestNumber + 1
	documentNumber := fmt.Sprintf("FGD/%s/%03d", dateNow.Format("20060102"), newNumber)
	return documentNumber, nil
}

func (uc *FgdScheduleUseCase) UpdateStatusFgdSchedule(req *request.UpdateStatusFgdScheduleRequest) (*response.FgdScheduleResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateStatusFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	FgdSchedule, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateStatusFgdScheduleRequest] " + err.Error())
		return nil, err
	}
	if FgdSchedule == nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateStatusFgdScheduleRequest] FgdSchedule not found")
		return nil, errors.New("FgdSchedule not found")
	}

	FgdSchedule, err = uc.Repository.UpdateFgdSchedule(&entity.FgdSchedule{
		ID:     parsedID,
		Status: entity.FgdScheduleStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateStatusFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	resp, err := uc.DTO.ConvertEntityToResponse(FgdSchedule)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.UpdateStatusFgdScheduleRequest] " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (uc *FgdScheduleUseCase) FindMySchedule(userID, projectRecruitmentLineID, jobPostingID uuid.UUID) (*response.FgdScheduleMyselfResponse, error) {
	// find project recruitment line
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByID(projectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "Project Recruitment Line not found")
		return nil, errors.New("project recruitment line not found")
	}

	// find user profile
	userProfile, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if userProfile == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "User Profile not found")
		return nil, errors.New("user profile not found")
	}

	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	// find FgdSchedule
	fgdSchedules, err := uc.Repository.FindAllByKeys(map[string]interface{}{
		"job_posting_id":              jobPostingID,
		"project_recruitment_line_id": projectRecruitmentLineID,
	})
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}

	fgdScheduleIDS := make([]uuid.UUID, 0)
	for _, fgdSchedule := range *fgdSchedules {
		fgdScheduleIDS = append(fgdScheduleIDS, fgdSchedule.ID)
	}

	fgdScheduleApplicant, err := uc.FgdApplicantRepository.FindByUserProfileIDAndFgdIDs(userProfile.ID, fgdScheduleIDS)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if fgdScheduleApplicant == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "FgdSchedule Applicant not found")
		return nil, errors.New("FgdSchedule applicant not found")
	}

	resp, err := uc.Repository.FindByIDForMyself(fgdScheduleApplicant.FgdScheduleID, fgdScheduleApplicant.UserProfileID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if resp == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "FgdSchedule not found")
		return nil, errors.New("FgdSchedule not found")
	}

	convertResp, err := uc.DTO.ConvertEntityToMyselfResponse(resp)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}

	return convertResp, nil
}

func (uc *FgdScheduleUseCase) FindScheduleForApplicant(applicantID, projectRecruitmentLineID, jobPostingID, employeeID uuid.UUID) (*response.FgdScheduleMyselfResponse, error) {
	// find project recruitment line
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByID(projectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "Project Recruitment Line not found")
		return nil, errors.New("project recruitment line not found")
	}

	// find applicant
	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": applicantID,
	})
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "Applicant not found")
		return nil, errors.New("applicant not found")
	}

	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	// find fgdSchedule assessor id
	fgdScheduleAssessor, err := uc.FgdAssessorRepository.FindByKeys(map[string]interface{}{
		"employee_id": employeeID,
	})
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}

	// find FgdSchedule
	fgdSchedules, err := uc.Repository.FindAllByKeys(map[string]interface{}{
		"job_posting_id":              jobPostingID,
		"project_recruitment_line_id": projectRecruitmentLineID,
	})
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}

	fgdScheduleIDS := make([]uuid.UUID, 0)
	for _, fgdSchedule := range *fgdSchedules {
		fgdScheduleIDS = append(fgdScheduleIDS, fgdSchedule.ID)
	}

	fgdScheduleApplicant, err := uc.FgdApplicantRepository.FindByUserProfileIDAndFgdIDs(applicant.UserProfileID, fgdScheduleIDS)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if fgdScheduleApplicant == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "FgdSchedule Applicant not found")
		return nil, errors.New("FgdSchedule applicant not found")
	}

	resp, err := uc.Repository.FindByIDForMyselfAndAssessor(fgdScheduleApplicant.FgdScheduleID, fgdScheduleApplicant.UserProfileID, fgdScheduleAssessor.ID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if resp == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "FgdSchedule not found")
		return nil, errors.New("FgdSchedule not found")
	}

	convertResp, err := uc.DTO.ConvertEntityToMyselfResponse(resp)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}

	return convertResp, nil
}

func (uc *FgdScheduleUseCase) FindMyScheduleForAssessor(employeeID, projectRecruitmentLineID, jobPostingID uuid.UUID) (*[]response.FgdScheduleMyselfForAssessorResponse, error) {
	// find project recruitment line
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByID(projectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMyScheduleForAssessor] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMyScheduleForAssessor] " + "Project Recruitment Line not found")
		return nil, errors.New("project recruitment line not found")
	}

	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	// find fgdSchedule
	fgdSchedules, err := uc.Repository.FindAllByKeys(map[string]interface{}{
		"job_posting_id":              jobPostingID,
		"project_recruitment_line_id": projectRecruitmentLineID,
	})
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}

	fgdScheduleIDS := make([]uuid.UUID, 0)
	for _, fgdSchedule := range *fgdSchedules {
		fgdScheduleIDS = append(fgdScheduleIDS, fgdSchedule.ID)
	}

	fgdScheduleAssessor, err := uc.FgdAssessorRepository.FindByEmployeeIDAndFgdIDs(employeeID, fgdScheduleIDS)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMyScheduleForAssessor] " + err.Error())
		return nil, err
	}
	if fgdScheduleAssessor == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMyScheduleForAssessor] " + "FgdSchedule Assessor not found")
		return nil, errors.New("FgdSchedule assessor not found")
	}

	fgdScheduleResps, err := uc.Repository.FindByIDsForMyselfAssessor(fgdScheduleIDS, fgdScheduleAssessor.ID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindMyScheduleForAssessor] " + err.Error())
		return nil, err
	}

	var convertResp []response.FgdScheduleMyselfForAssessorResponse
	for _, FgdScheduleResp := range *fgdScheduleResps {
		convertedResp, err := uc.DTO.ConvertEntityToMyselfAssessorResponse(&FgdScheduleResp)
		if err != nil {
			uc.Log.Error("[FgdScheduleUseCase.FindMyScheduleForAssessor] " + err.Error())
			return nil, err
		}
		convertResp = append(convertResp, *convertedResp)
	}

	return &convertResp, nil
}

func (uc *FgdScheduleUseCase) getApplicantIDsByJobPostingID(jobPostingID uuid.UUID, projectRecruitmentLineID uuid.UUID, order int, total int) (*response.TestApplicantsPayload, error) {
	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.GetApplicantsByJobPostingID] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[FgdScheduleUseCase.GetApplicantsByJobPostingID] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	var totalResult int = total

	// find project recruitment line that has order
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByKeys(map[string]interface{}{
		"project_recruitment_header_id": jobPosting.ProjectRecruitmentHeaderID,
		"id":                            projectRecruitmentLineID,
		// "order":                         order,
	})
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.GetApplicantsByJobPostingID] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[FgdScheduleUseCase.GetApplicantsByJobPostingID] " + "Project Recruitment Line not found")
		return nil, errors.New("project recruitment line not found")
	}

	applicants, err := uc.ApplicantRepository.GetAllByKeys(map[string]interface{}{
		"job_posting_id": jobPostingID,
		"order":          projectRecruitmentLine.Order,
	})
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.GetApplicantsByJobPostingID] " + err.Error())
		return nil, err
	}

	applicantIDs := []uuid.UUID{}
	for _, applicant := range applicants {
		applicantIDs = append(applicantIDs, applicant.ID)
	}

	resultApplicants := &[]entity.Applicant{}
	*resultApplicants = applicants

	if projectRecruitmentLine.TemplateActivityLine != nil {
		if projectRecruitmentLine.TemplateActivityLine.TemplateQuestion != nil {
			if projectRecruitmentLine.TemplateActivityLine.TemplateQuestion.FormType == string(entity.TQ_FORM_TYPE_FGD) {
				testApplicants, err := uc.FgdApplicantRepository.FindAllByApplicantIDs(applicantIDs)
				if err != nil {
					uc.Log.Error("[FgdScheduleUseCase.GetApplicantsByJobPostingID] " + err.Error())
					return nil, err
				}

				// filter applicants that have not taken the test
				resultApplicants = &[]entity.Applicant{}
				for _, applicant := range applicants {
					var found bool
					for _, testApplicant := range testApplicants {
						if applicant.ID == testApplicant.ApplicantID && applicant.Order != projectRecruitmentLine.Order {
							found = true
							break
						}
					}

					if !found {
						*resultApplicants = append(*resultApplicants, applicant)
					}
				}
			}
		}
	}

	if total > 0 {
		if len(*resultApplicants) > total {
			*resultApplicants = (*resultApplicants)[:total]
		} else {
			totalResult = len(*resultApplicants)
		}
	}

	if len(*resultApplicants) == 0 {
		uc.Log.Warn("[FgdScheduleUseCase.GetApplicantsByJobPostingID] " + "No applicants found")
		totalResult = 0
	}

	resultApplicantIDs := []uuid.UUID{}
	userProfileIDs := []uuid.UUID{}
	for _, applicant := range *resultApplicants {
		resultApplicantIDs = append(resultApplicantIDs, applicant.ID)
		userProfileIDs = append(userProfileIDs, applicant.UserProfileID)
	}

	return &response.TestApplicantsPayload{
		ApplicantIDs:   resultApplicantIDs,
		UserProfileIDs: userProfileIDs,
		Total:          totalResult,
	}, nil
}

func (uc *FgdScheduleUseCase) FindByIDForAnswer(id, jobPostingID uuid.UUID) (*response.FgdScheduleResponse, error) {
	fgdSchedule, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindByIDForAnswer] " + err.Error())
		return nil, err
	}

	if fgdSchedule == nil {
		uc.Log.Error("[FgdScheduleUseCase.FindByIDForAnswer] fgdSchedule not found")
		return nil, errors.New("fgdSchedule not found")
	}

	resp, err := uc.DTO.ConvertEntityToResponse(fgdSchedule)
	if err != nil {
		uc.Log.Error("[FgdScheduleUseCase.FindByIDForAnswer] " + err.Error())
		return nil, err
	}

	return resp, nil
}
