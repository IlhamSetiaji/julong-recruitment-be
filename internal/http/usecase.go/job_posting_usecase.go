package usecase

import (
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
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.JobPostingResponse, int64, error)
	FindByID(id uuid.UUID) (*response.JobPostingResponse, error)
	UpdateJobPosting(req *request.UpdateJobPostingRequest) (*response.JobPostingResponse, error)
	DeleteJobPosting(id uuid.UUID) error
}

type JobPostingUseCase struct {
	Log                                *logrus.Logger
	Repository                         repository.IJobPostingRepository
	DTO                                dto.IJobPostingDTO
	ProjectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository
	MPRequestRepository                repository.IMPRequestRepository
	Viper                              *viper.Viper
}

func NewJobPostingUseCase(
	log *logrus.Logger,
	repo repository.IJobPostingRepository,
	dto dto.IJobPostingDTO,
	prhRepository repository.IProjectRecruitmentHeaderRepository,
	mpRequestRepository repository.IMPRequestRepository,
	viper *viper.Viper,
) IJobPostingUseCase {
	return &JobPostingUseCase{
		Log:                                log,
		Repository:                         repo,
		DTO:                                dto,
		ProjectRecruitmentHeaderRepository: prhRepository,
		MPRequestRepository:                mpRequestRepository,
		Viper:                              viper,
	}
}

func JobPostingUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IJobPostingUseCase {
	repo := repository.JobPostingRepositoryFactory(log)
	dto := dto.JobPostingDTOFactory(log)
	prhRepository := repository.ProjectRecruitmentHeaderRepositoryFactory(log)
	mpRequestRepository := repository.MPRequestRepositoryFactory(log)
	return NewJobPostingUseCase(log, repo, dto, prhRepository, mpRequestRepository, viper)
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

func (uc *JobPostingUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.JobPostingResponse, int64, error) {
	jobPostings, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	jobPostingResponses := make([]response.JobPostingResponse, 0)
	for _, jobPosting := range *jobPostings {
		jobPosting.OrganizationLogo = uc.Viper.GetString("app.url") + jobPosting.OrganizationLogo
		jobPosting.Poster = uc.Viper.GetString("app.url") + jobPosting.Poster
		jobPostingResponses = append(jobPostingResponses, *uc.DTO.ConvertEntityToResponse(&jobPosting))
	}

	return &jobPostingResponses, total, nil
}

func (uc *JobPostingUseCase) FindByID(id uuid.UUID) (*response.JobPostingResponse, error) {
	jobPosting, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[JobPostingUseCase.FindByID] " + err.Error())
		return nil, err
	}

	jobPosting.OrganizationLogo = uc.Viper.GetString("app.url") + jobPosting.OrganizationLogo
	jobPosting.Poster = uc.Viper.GetString("app.url") + jobPosting.Poster

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
