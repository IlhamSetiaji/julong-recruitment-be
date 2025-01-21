package usecase

import (
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

type IJobPostingUseCase interface {
	CreateJobPosting(req *request.CreateJobPostingRequest) (*response.JobPostingResponse, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}, userID uuid.UUID) (*[]response.JobPostingResponse, int64, error)
	FindByID(id uuid.UUID, userID uuid.UUID) (*response.JobPostingResponse, error)
	UpdateJobPosting(req *request.UpdateJobPostingRequest) (*response.JobPostingResponse, error)
	DeleteJobPosting(id uuid.UUID) error
	GenerateDocumentNumber(dateNow time.Time) (string, error)
	UpdateJobPostingOrganizationLogoToNull(id uuid.UUID) error
	UpdateJobPostingPosterToNull(id uuid.UUID) error
	FindAllAppliedJobPostingByUserID(userID uuid.UUID) (*[]response.JobPostingResponse, error)
}

type JobPostingUseCase struct {
	Log                                *logrus.Logger
	Repository                         repository.IJobPostingRepository
	DTO                                dto.IJobPostingDTO
	ProjectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository
	MPRequestRepository                repository.IMPRequestRepository
	Viper                              *viper.Viper
	ApplicantRepository                repository.IApplicantRepository
	UserProfileRepository              repository.IUserProfileRepository
}

func NewJobPostingUseCase(
	log *logrus.Logger,
	repo repository.IJobPostingRepository,
	dto dto.IJobPostingDTO,
	prhRepository repository.IProjectRecruitmentHeaderRepository,
	mpRequestRepository repository.IMPRequestRepository,
	viper *viper.Viper,
	applicantRepository repository.IApplicantRepository,
	userProfileRepository repository.IUserProfileRepository,
) IJobPostingUseCase {
	return &JobPostingUseCase{
		Log:                                log,
		Repository:                         repo,
		DTO:                                dto,
		ProjectRecruitmentHeaderRepository: prhRepository,
		MPRequestRepository:                mpRequestRepository,
		Viper:                              viper,
		ApplicantRepository:                applicantRepository,
		UserProfileRepository:              userProfileRepository,
	}
}

func JobPostingUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IJobPostingUseCase {
	repo := repository.JobPostingRepositoryFactory(log)
	dto := dto.JobPostingDTOFactory(log, viper)
	prhRepository := repository.ProjectRecruitmentHeaderRepositoryFactory(log)
	mpRequestRepository := repository.MPRequestRepositoryFactory(log)
	applicantRepository := repository.ApplicantRepositoryFactory(log)
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	return NewJobPostingUseCase(
		log,
		repo,
		dto,
		prhRepository,
		mpRequestRepository,
		viper,
		applicantRepository,
		userProfileRepository,
	)
}

func (uc *JobPostingUseCase) CreateJobPosting(req *request.CreateJobPostingRequest) (*response.JobPostingResponse, error) {
	// Get project recruitment header
	prhID, err := uuid.Parse(req.ProjectRecruitmentHeaderID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
		return nil, err
	}

	prh, err := uc.ProjectRecruitmentHeaderRepository.FindByID(prhID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
		return nil, err
	}

	if prh == nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + "Project Recruitment Header not found")
		return nil, err
	}

	// Get MP Request
	mpRequestID, err := uuid.Parse(req.MPRequestID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
		return nil, err
	}

	mpRequest, err := uc.MPRequestRepository.FindByID(mpRequestID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
		return nil, err
	}

	if mpRequest == nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + "MP Request not found")
		return nil, err
	}

	// Parse data
	data, err := uc.parseData(req.ForOrganizationID, req.ForOrganizationLocationID, req.JobID, req.DocumentDate, req.StartDate, req.EndDate)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
		return nil, err
	}

	jobPosting, err := uc.Repository.CreateJobPosting(&entity.JobPosting{
		ProjectRecruitmentHeaderID: prhID,
		MPRequestID:                &mpRequestID,
		ForOrganizationID:          data["forOrgID"].(*uuid.UUID),
		ForOrganizationLocationID:  data["forOrgLocID"].(*uuid.UUID),
		JobID:                      data["jobID"].(*uuid.UUID),
		DocumentNumber:             req.DocumentNumber,
		DocumentDate:               data["documentDate"].(time.Time),
		RecruitmentType:            entity.ProjectRecruitmentType(req.RecruitmentType),
		StartDate:                  data["startDate"].(time.Time),
		EndDate:                    data["endDate"].(time.Time),
		Status:                     entity.JobPostingStatus(req.Status),
		SalaryMin:                  req.SalaryMin,
		SalaryMax:                  req.SalaryMax,
		ContentDescription:         req.ContentDescription,
		OrganizationLogo:           req.OrganizationLogoPath,
		Poster:                     req.PosterPath,
		Link:                       req.Link,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(jobPosting), nil
}

func (uc *JobPostingUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}, userID uuid.UUID) (*[]response.JobPostingResponse, int64, error) {
	userProfile, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if userProfile == nil {
		uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + "User Profile not found")
		return nil, 0, err
	}

	jobPostings, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	for _, jobPosting := range *jobPostings {
		applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
			"job_posting_id":  jobPosting.ID,
			"user_profile_id": userProfile.ID,
		})
		if err != nil {
			uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
			return nil, 0, err
		}

		isApplied := false
		if applicant != nil {
			isApplied = true
		}

		jobPosting.IsApplied = isApplied

		jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jobPosting))
	}

	return &jobPostingResponses, total, nil
}

func (uc *JobPostingUseCase) FindByID(id uuid.UUID, userID uuid.UUID) (*response.JobPostingResponse, error) {
	jobPosting, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindByID] " + err.Error())
		return nil, err
	}

	if jobPosting == nil {
		uc.Log.Error("[JobPostingUseCase.FindByID] " + "Job Posting not found")
		return nil, err
	}

	userProfile, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindByID] " + err.Error())
		return nil, err
	}

	if userProfile == nil {
		uc.Log.Error("[JobPostingUseCase.FindByID] " + "User Profile not found")
		return nil, err
	}

	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"job_posting_id":  jobPosting.ID,
		"user_profile_id": userProfile.ID,
	})

	isApplied := false
	if applicant != nil {
		isApplied = true
	}

	jobPosting.IsApplied = isApplied

	return uc.DTO.ConvertEntityToResponse(jobPosting), nil
}

func (uc *JobPostingUseCase) UpdateJobPosting(req *request.UpdateJobPostingRequest) (*response.JobPostingResponse, error) {
	// Get project recruitment header
	prhID, err := uuid.Parse(req.ProjectRecruitmentHeaderID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + err.Error())
		return nil, err
	}

	prh, err := uc.ProjectRecruitmentHeaderRepository.FindByID(prhID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + err.Error())
		return nil, err
	}

	if prh == nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + "Project Recruitment Header not found")
		return nil, err
	}

	// Get MP Request
	mpRequestID, err := uuid.Parse(req.MPRequestID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + err.Error())
		return nil, err
	}

	mpRequest, err := uc.MPRequestRepository.FindByID(mpRequestID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + err.Error())
		return nil, err
	}

	if mpRequest == nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + "MP Request not found")
		return nil, err
	}

	// Parse data
	data, err := uc.parseData(req.ForOrganizationID, req.ForOrganizationLocationID, req.JobID, req.DocumentDate, req.StartDate, req.EndDate)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + err.Error())
		return nil, err
	}

	jobPosting, err := uc.Repository.UpdateJobPosting(&entity.JobPosting{
		ID:                         uuid.MustParse(req.ID),
		ProjectRecruitmentHeaderID: prhID,
		MPRequestID:                &mpRequestID,
		ForOrganizationID:          data["forOrgID"].(*uuid.UUID),
		ForOrganizationLocationID:  data["forOrgLocID"].(*uuid.UUID),
		JobID:                      data["jobID"].(*uuid.UUID),
		DocumentNumber:             req.DocumentNumber,
		DocumentDate:               data["documentDate"].(time.Time),
		RecruitmentType:            entity.ProjectRecruitmentType(req.RecruitmentType),
		StartDate:                  data["startDate"].(time.Time),
		EndDate:                    data["endDate"].(time.Time),
		Status:                     entity.JobPostingStatus(req.Status),
		SalaryMin:                  req.SalaryMin,
		SalaryMax:                  req.SalaryMax,
		ContentDescription:         req.ContentDescription,
		OrganizationLogo:           req.OrganizationLogoPath,
		Poster:                     req.PosterPath,
		Link:                       req.Link,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(jobPosting), nil
}

func (uc *JobPostingUseCase) DeleteJobPosting(id uuid.UUID) error {
	return uc.Repository.DeleteJobPosting(id)
}

func (uc *JobPostingUseCase) parseData(forOrgID, forOrgLocID, forJobID, docDate, stDate, edDate string) (map[string]interface{}, error) {
	orgID, err := uuid.Parse(forOrgID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.parseData] " + err.Error())
		return nil, err
	}
	orgLocID, err := uuid.Parse(forOrgLocID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.parseData] " + err.Error())
		return nil, err
	}
	jobID, err := uuid.Parse(forJobID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.parseData] " + err.Error())
		return nil, err
	}
	documentDate, err := time.Parse("2006-01-02", docDate)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.parseData] " + err.Error())
		return nil, err
	}
	startDate, err := time.Parse("2006-01-02", stDate)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.parseData] " + err.Error())
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", edDate)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.parseData] " + err.Error())
		return nil, err
	}

	orgIDPtr := &orgID
	orgLocIDPtr := &orgLocID
	jobIDPtr := &jobID

	return map[string]interface{}{
		"forOrgID":     orgIDPtr,
		"forOrgLocID":  orgLocIDPtr,
		"jobID":        jobIDPtr,
		"documentDate": documentDate,
		"startDate":    startDate,
		"endDate":      endDate,
	}, nil
}

func (uc *JobPostingUseCase) GenerateDocumentNumber(dateNow time.Time) (string, error) {
	dateStr := dateNow.Format("2006-01-02")
	highestNumber, err := uc.Repository.GetHighestDocumentNumberByDate(dateStr)
	if err != nil {
		uc.Log.Errorf("[JobPostingUseCase.GenerateDocumentNumber] " + err.Error())
		return "", err
	}

	newNumber := highestNumber + 1
	documentNumber := fmt.Sprintf("JP/%s/%03d", dateNow.Format("20060102"), newNumber)
	return documentNumber, nil
}

func (uc *JobPostingUseCase) UpdateJobPostingOrganizationLogoToNull(id uuid.UUID) error {
	return uc.Repository.UpdateJobPostingOrganizationLogoToNull(id)
}

func (uc *JobPostingUseCase) UpdateJobPostingPosterToNull(id uuid.UUID) error {
	return uc.Repository.UpdateJobPostingPosterToNull(id)
}

func (uc *JobPostingUseCase) FindAllAppliedJobPostingByUserID(userID uuid.UUID) (*[]response.JobPostingResponse, error) {
	userProfile, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllAppliedJobPostingByUserID] " + err.Error())
		return nil, err
	}

	if userProfile == nil {
		uc.Log.Error("[JobPostingUseCase.FindAllAppliedJobPostingByUserID] " + "User Profile not found")
		return nil, err
	}

	applicants, err := uc.ApplicantRepository.GetAllByKeys(map[string]interface{}{
		"user_profile_id": userProfile.ID,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllAppliedJobPostingByUserID] " + err.Error())
		return nil, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	for _, applicant := range applicants {
		jobPosting, err := uc.Repository.FindByID(applicant.JobPostingID)
		if err != nil {
			uc.Log.Error("[JobPostingUseCase.FindAllAppliedJobPostingByUserID] " + err.Error())
			return nil, err
		}

		if jobPosting == nil {
			uc.Log.Error("[JobPostingUseCase.FindAllAppliedJobPostingByUserID] " + "Job Posting not found")
			return nil, err
		}

		jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(jobPosting))
	}

	return &jobPostingResponses, nil
}
