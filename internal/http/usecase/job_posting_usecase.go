package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/service"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IJobPostingUseCase interface {
	CreateJobPosting(req *request.CreateJobPostingRequest) (*response.JobPostingResponse, error)
	FindAllPaginatedWithoutUserID(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.JobPostingResponse, int64, error)
	FindAllPaginatedWithoutUserIDShowOnly(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.JobPostingResponse, int64, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}, userID uuid.UUID) (*[]response.JobPostingResponse, int64, error)
	FindAllPaginatedShowOnly(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}, userID uuid.UUID, userMajors []string) (*[]response.JobPostingResponse, int64, error)
	FindByID(id uuid.UUID, userID uuid.UUID) (*response.JobPostingResponse, error)
	UpdateJobPosting(req *request.UpdateJobPostingRequest) (*response.JobPostingResponse, error)
	DeleteJobPosting(id uuid.UUID) error
	GenerateDocumentNumber(dateNow time.Time) (string, error)
	UpdateJobPostingOrganizationLogoToNull(id uuid.UUID) error
	UpdateJobPostingPosterToNull(id uuid.UUID) error
	FindAllAppliedJobPostingByUserID(userID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}) (*[]response.JobPostingResponse, int64, error)
	InsertSavedJob(userID, jobPostingID uuid.UUID) error
	FindAllSavedJobsByUserID(userID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.JobPostingResponse, int64, error)
	DeleteSavedJob(userID, jobPostingID uuid.UUID) error
	FindAllJobPostingsByEmployeeID(employeeID uuid.UUID) (*[]response.JobPostingResponse, error)
	FindAllByProjectRecruitmentHeaderID(prhID uuid.UUID) (*[]response.JobPostingResponse, error)
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
	ProjectRecruitmentLineRepository   repository.IProjectRecruitmentLineRepository
	AdministrativeSelectionRepository  repository.IAdministrativeSelectionRepository
	AdministrativeResultRepository     repository.IAdministrativeResultRepository
	ProjectPicRepository               repository.IProjectPicRepository
	MPRequestMessage                   messaging.IMPRequestMessage
	MPRequestService                   service.IMPRequestService
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
	prlRepository repository.IProjectRecruitmentLineRepository,
	asRepo repository.IAdministrativeSelectionRepository,
	arRepo repository.IAdministrativeResultRepository,
	projectPicRepo repository.IProjectPicRepository,
	mpRequestMessage messaging.IMPRequestMessage,
	mpRequestService service.IMPRequestService,
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
		ProjectRecruitmentLineRepository:   prlRepository,
		AdministrativeSelectionRepository:  asRepo,
		AdministrativeResultRepository:     arRepo,
		ProjectPicRepository:               projectPicRepo,
		MPRequestMessage:                   mpRequestMessage,
		MPRequestService:                   mpRequestService,
	}
}

func JobPostingUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IJobPostingUseCase {
	repo := repository.JobPostingRepositoryFactory(log)
	dto := dto.JobPostingDTOFactory(log, viper)
	prhRepository := repository.ProjectRecruitmentHeaderRepositoryFactory(log)
	mpRequestRepository := repository.MPRequestRepositoryFactory(log)
	applicantRepository := repository.ApplicantRepositoryFactory(log)
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	prlRepository := repository.ProjectRecruitmentLineRepositoryFactory(log)
	asRepo := repository.AdministrativeSelectionRepositoryFactory(log)
	arRepo := repository.AdministrativeResultRepositoryFactory(log)
	projectPicRepo := repository.ProjectPicRepositoryFactory(log)
	mpRequestMessage := messaging.MPRequestMessageFactory(log)
	mpRequestService := service.MPRequestServiceFactory(log)
	return NewJobPostingUseCase(
		log,
		repo,
		dto,
		prhRepository,
		mpRequestRepository,
		viper,
		applicantRepository,
		userProfileRepository,
		prlRepository,
		asRepo,
		arRepo,
		projectPicRepo,
		mpRequestMessage,
		mpRequestService,
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

	if data["startDate"].(time.Time).Before(prh.StartDate) && data["endDate"].(time.Time).After(prh.EndDate) {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + "Start Date or End Date is not in the range of Project Recruitment Header")
		return nil, fmt.Errorf("Start Date or End Date is not in the range of Project Recruitment Header[Start Date: %v, End Date: %v]", prh.StartDate, prh.EndDate)
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
		MinimumWorkExperience:      req.MinimumWorkExperience,
		Name:                       req.Name,
		IsShow:                     req.IsShow,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
		return nil, err
	}

	_, err = uc.MPRequestRepository.Update(&entity.MPRequest{
		ID:     mpRequest.ID,
		Status: entity.MPR_STATUS_ON_GOING,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
		return nil, err
	}

	// get project recruitment lines
	prls, err := uc.ProjectRecruitmentLineRepository.GetAllByKeys(map[string]interface{}{
		"project_recruitment_header_id": prhID,
		"order":                         1,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
		return nil, err
	}

	for _, prl := range prls {
		for _, pic := range prl.ProjectPics {
			_, err = uc.AdministrativeSelectionRepository.CreateAdministrativeSelection(&entity.AdministrativeSelection{
				JobPostingID:    jobPosting.ID,
				ProjectPicID:    pic.ID,
				Status:          entity.ADMINISTRATIVE_SELECTION_STATUS_IN_PROGRESS,
				DocumentDate:    jobPosting.DocumentDate,
				DocumentNumber:  jobPosting.DocumentNumber,
				TotalApplicants: 0,
			})
			if err != nil {
				uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
				return nil, err
			}
		}
	}

	return uc.DTO.ConvertEntityToResponse(jobPosting), nil
}

func (uc *JobPostingUseCase) FindAllPaginatedWithoutUserID(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.JobPostingResponse, int64, error) {
	jobPostings, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllPaginatedWithoutUserID] " + err.Error())
		return nil, 0, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	for _, jobPosting := range *jobPostings {
		totalApplicant := len(jobPosting.Applicants)
		jobPosting.TotalApplicant = totalApplicant
		jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jobPosting))
	}

	return &jobPostingResponses, total, nil
}

func (uc *JobPostingUseCase) FindAllPaginatedWithoutUserIDShowOnly(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.JobPostingResponse, int64, error) {
	jobPostings, total, err := uc.Repository.FindAllPaginatedShowOnly(page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllPaginatedWithoutUserID] " + err.Error())
		return nil, 0, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	for _, jobPosting := range *jobPostings {
		totalApplicant := len(jobPosting.Applicants)
		jobPosting.TotalApplicant = totalApplicant
		jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jobPosting))
	}

	return &jobPostingResponses, total, nil
}

func (uc *JobPostingUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}, userID uuid.UUID) (*[]response.JobPostingResponse, int64, error) {
	jobPostings, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	userProfile, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	if userProfile != nil {
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

			savedJob, err := uc.Repository.FindSavedJob(userProfile.ID, jobPosting.ID)
			if err != nil {
				uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
				return nil, 0, err
			}

			isSaved := false
			if savedJob != nil {
				isSaved = true
			}

			uc.Log.Info("IsSaved: ", isSaved)

			totalApplicant := len(jobPosting.Applicants)

			jobPosting.IsApplied = isApplied
			jobPosting.IsSaved = isSaved
			jobPosting.TotalApplicant = totalApplicant

			jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jobPosting))
		}
	} else {
		for _, jobPosting := range *jobPostings {
			totalApplicant := len(jobPosting.Applicants)
			jobPosting.TotalApplicant = totalApplicant
			jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jobPosting))
		}
	}

	return &jobPostingResponses, total, nil
}

func (uc *JobPostingUseCase) FindAllPaginatedShowOnly(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}, userID uuid.UUID, userMajors []string) (*[]response.JobPostingResponse, int64, error) {
	userProfile, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if userProfile == nil {
		uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + "User Profile not found")
		return nil, 0, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	var jumlah int64
	if len(userMajors) > 0 {
		mprCloneIds, err := uc.MPRequestMessage.SendFindIdsByMajorsMessage(userMajors)
		if err != nil {
			uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
			return nil, 0, err
		}

		mprs, err := uc.MPRequestRepository.FindAllByMPRCloneIDs(mprCloneIds)
		if err != nil {
			uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
			return nil, 0, err
		}

		mprIds := make([]uuid.UUID, 0)
		for _, mpr := range *mprs {
			mprIds = append(mprIds, mpr.ID)
		}

		jobPostings, total, err := uc.Repository.FindAlPaginatedByMPRequestIDs(page, pageSize, search, sort, filter, mprIds)
		if err != nil {
			uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
			return nil, 0, err
		}
		jumlah = total

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

			savedJob, err := uc.Repository.FindSavedJob(userProfile.ID, jobPosting.ID)
			if err != nil {
				uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
				return nil, 0, err
			}

			isSaved := false
			if savedJob != nil {
				isSaved = true
			}

			uc.Log.Info("IsSaved: ", isSaved)

			totalApplicant := len(jobPosting.Applicants)

			jobPosting.IsApplied = isApplied
			jobPosting.IsSaved = isSaved
			jobPosting.TotalApplicant = totalApplicant

			jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jobPosting))
		}
	} else {
		jobPostings, total, err := uc.Repository.FindAllPaginatedShowOnly(page, pageSize, search, sort, filter)
		if err != nil {
			uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
			return nil, 0, err
		}

		jumlah = total

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

			savedJob, err := uc.Repository.FindSavedJob(userProfile.ID, jobPosting.ID)
			if err != nil {
				uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
				return nil, 0, err
			}

			isSaved := false
			if savedJob != nil {
				isSaved = true
			}

			uc.Log.Info("IsSaved: ", isSaved)

			totalApplicant := len(jobPosting.Applicants)

			jobPosting.IsApplied = isApplied
			jobPosting.IsSaved = isSaved
			jobPosting.TotalApplicant = totalApplicant

			jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jobPosting))
		}
	}

	return &jobPostingResponses, jumlah, nil
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

	if userProfile != nil {
		if userProfile == nil {
			uc.Log.Error("[JobPostingUseCase.FindByID] " + "User Profile not found")
			return nil, err
		}

		applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
			"job_posting_id":  jobPosting.ID,
			"user_profile_id": userProfile.ID,
		})
		if err != nil {
			uc.Log.Error("[JobPostingUseCase.FindByID] " + err.Error())
			return nil, err
		}

		isApplied := false
		if applicant != nil {
			isApplied = true
		}

		jobPosting.IsApplied = isApplied
	}

	jobResp := uc.DTO.ConvertEntityToResponse(jobPosting)
	resp, err := uc.MPRequestMessage.SendFindByIdMessage(jobPosting.MPRequest.MPRCloneID.String())
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when send find by id message: %v", err)
		jobResp.MPRequest = nil
	} else {
		convertedData, err := uc.MPRequestService.CheckPortalData(resp)
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when check portal data: %v", err)
			return nil, err
		}
		convertedData.Status = string(jobPosting.MPRequest.Status)
		convertedData.ID = jobPosting.MPRequest.ID
		jobResp.MPRequest = convertedData
	}

	return jobResp, nil
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

	exist, err := uc.Repository.FindByID(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + err.Error())
		return nil, err
	}

	if exist == nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + "Job Posting not found")
		return nil, errors.New("Job Posting not found")
	}

	_, err = uc.MPRequestRepository.Update(&entity.MPRequest{
		ID:     *exist.MPRequestID,
		Status: entity.MPR_STATUS_OPEN,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
		return nil, err
	}

	// Parse data
	data, err := uc.parseData(req.ForOrganizationID, req.ForOrganizationLocationID, req.JobID, req.DocumentDate, req.StartDate, req.EndDate)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + err.Error())
		return nil, err
	}

	if data["startDate"].(time.Time).Before(prh.StartDate) && data["endDate"].(time.Time).After(prh.EndDate) {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + "Start Date or End Date is not in the range of Project Recruitment Header")
		return nil, fmt.Errorf("Start Date or End Date is not in the range of Project Recruitment Header [Start Date: %v, End Date: %v]", prh.StartDate, prh.EndDate)
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
		MinimumWorkExperience:      req.MinimumWorkExperience,
		Name:                       req.Name,
		IsShow:                     req.IsShow,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.UpdateJobPosting] " + err.Error())
		return nil, err
	}

	_, err = uc.MPRequestRepository.Update(&entity.MPRequest{
		ID:     mpRequest.ID,
		Status: entity.MPR_STATUS_ON_GOING,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.CreateJobPosting] " + err.Error())
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

func (uc *JobPostingUseCase) FindAllAppliedJobPostingByUserID(userID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}) (*[]response.JobPostingResponse, int64, error) {
	userProfile, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllAppliedJobPostingByUserID] " + err.Error())
		return nil, 0, err
	}

	if userProfile == nil {
		uc.Log.Error("[JobPostingUseCase.FindAllAppliedJobPostingByUserID] " + "User Profile not found")
		return nil, 0, err
	}

	applicants, err := uc.ApplicantRepository.GetAllByKeys(map[string]interface{}{
		"user_profile_id": userProfile.ID,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllAppliedJobPostingByUserID] " + err.Error())
		return nil, 0, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	for _, applicant := range applicants {
		jobPosting, err := uc.Repository.FindByID(applicant.JobPostingID)
		if err != nil {
			uc.Log.Error("[JobPostingUseCase.FindAllAppliedJobPostingByUserID] " + err.Error())
			return nil, 0, err
		}

		if jobPosting == nil {
			uc.Log.Error("[JobPostingUseCase.FindAllAppliedJobPostingByUserID] " + "Job Posting not found")
			return nil, 0, err
		}

		jobPosting.AppliedDate = applicant.AppliedDate
		jobPosting.ApplicantStatus = applicant.Status
		jobPosting.ApplicantProcessStatus = applicant.ProcessStatus

		jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(jobPosting))
	}

	// Implement pagination
	total := int64(len(jobPostingResponses))
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > int(total) {
		start = int(total)
	}
	if end > int(total) {
		end = int(total)
	}

	paginatedResponses := jobPostingResponses[start:end]

	return &paginatedResponses, total, nil
}

func (uc *JobPostingUseCase) InsertSavedJob(userID, jobPostingID uuid.UUID) error {
	jobPostingExist, err := uc.Repository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.InsertSavedJob] " + err.Error())
		return err
	}
	if jobPostingExist == nil {
		uc.Log.Error("[JobPostingUseCase.InsertSavedJob] " + "Job Posting not found")
		return err
	}

	userProfileExist, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.InsertSavedJob] " + err.Error())
		return err
	}
	if userProfileExist == nil {
		uc.Log.Error("[JobPostingUseCase.InsertSavedJob] " + "User Profile not found")
		return err
	}

	savedJob, err := uc.Repository.FindSavedJob(userProfileExist.ID, jobPostingID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.InsertSavedJob] " + err.Error())
		return err
	}
	if savedJob != nil {
		err = uc.Repository.DeleteSavedJob(userProfileExist.ID, jobPostingID)
		if err != nil {
			uc.Log.Error("[JobPostingUseCase.InsertSavedJob] " + err.Error())
			return err
		}
	} else {
		err = uc.Repository.InsertSavedJob(userProfileExist.ID, jobPostingID)
		if err != nil {
			uc.Log.Error("[JobPostingUseCase.InsertSavedJob] " + err.Error())
			return err
		}
	}

	return nil
}

func (uc *JobPostingUseCase) FindAllSavedJobsByUserID(userID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.JobPostingResponse, int64, error) {
	userProfileExist, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllSavedJobsByUserID] " + err.Error())
		return nil, 0, err
	}
	if userProfileExist == nil {
		uc.Log.Error("[JobPostingUseCase.FindAllSavedJobsByUserID] " + "User Profile not found")
		return nil, 0, err
	}

	jobPostings, total, err := uc.Repository.FindAllSavedJobsByUserProfileID(userProfileExist.ID, page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllSavedJobsByUserID] " + err.Error())
		return nil, 0, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	for _, jobPosting := range *jobPostings {
		jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jobPosting))
	}

	return &jobPostingResponses, total, nil
}

func (uc *JobPostingUseCase) DeleteSavedJob(userID, jobPostingID uuid.UUID) error {
	userProfileExist, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.DeleteSavedJob] " + err.Error())
		return err
	}
	if userProfileExist == nil {
		uc.Log.Error("[JobPostingUseCase.DeleteSavedJob] " + "User Profile not found")
		return err
	}

	err = uc.Repository.DeleteSavedJob(userProfileExist.ID, jobPostingID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.DeleteSavedJob] " + err.Error())
		return err
	}

	return nil
}

func (uc *JobPostingUseCase) FindAllJobPostingsByEmployeeID(employeeID uuid.UUID) (*[]response.JobPostingResponse, error) {
	pics, err := uc.ProjectPicRepository.FindAllByEmployeeID(employeeID)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllJobPostingsByEmployeeID] " + err.Error())
		return nil, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	jobPostingIDs := make(map[uuid.UUID]bool)

	for _, pic := range pics {
		for _, jp := range pic.ProjectRecruitmentLine.ProjectRecruitmentHeader.JobPostings {
			if _, exists := jobPostingIDs[jp.ID]; !exists {
				jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jp))
				jobPostingIDs[jp.ID] = true
			}
		}
	}

	return &jobPostingResponses, nil
}

func (uc *JobPostingUseCase) FindAllByProjectRecruitmentHeaderID(prhID uuid.UUID) (*[]response.JobPostingResponse, error) {
	jobPostings, err := uc.Repository.GetAllByKeys(map[string]interface{}{
		"project_recruitment_header_id": prhID,
	})
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllByProjectRecruitmentHeaderID] " + err.Error())
		return nil, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	for _, jobPosting := range *jobPostings {
		jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jobPosting))
	}

	return &jobPostingResponses, nil
}
