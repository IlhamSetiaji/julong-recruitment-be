package usecase

import (
	"strconv"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/service"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDashboardUseCase interface {
	GetDashboard() (*response.DashboardResponse, error)
}

type DashboardUseCase struct {
	Log                                *logrus.Logger
	Viper                              *viper.Viper
	JobPostingRepository               repository.IJobPostingRepository
	ProjectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository
	MPRequestRepository                repository.IMPRequestRepository
	ProjectRecruitmentLineRepository   repository.IProjectRecruitmentLineRepository
	ApplicantRepository                repository.IApplicantRepository
	DocumentSendingRepository          repository.IDocumentSendingRepository
	MPRequestMessage                   messaging.IMPRequestMessage
	MPRequestService                   service.IMPRequestService
	UserHelper                         helper.IUserHelper
	MPRequestHelper                    helper.IMPRequestHelper
	EmployeeMessage                    messaging.IEmployeeMessage
	JobPlafonMessage                   messaging.IJobPlafonMessage
}

func NewDashboardUseCase(
	log *logrus.Logger,
	viper *viper.Viper,
	jobPostingRepository repository.IJobPostingRepository,
	projectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository,
	MPRequestRepository repository.IMPRequestRepository,
	projectRecruitmentLineRepository repository.IProjectRecruitmentLineRepository,
	applicantRepository repository.IApplicantRepository,
	documentSendingRepository repository.IDocumentSendingRepository,
	MPRequestMessage messaging.IMPRequestMessage,
	MPRequestService service.IMPRequestService,
	userHelper helper.IUserHelper,
	MPRequestHelper helper.IMPRequestHelper,
	employeeMessage messaging.IEmployeeMessage,
	jobPlafonMessage messaging.IJobPlafonMessage,
) IDashboardUseCase {
	return &DashboardUseCase{
		Log:                                log,
		Viper:                              viper,
		JobPostingRepository:               jobPostingRepository,
		ProjectRecruitmentHeaderRepository: projectRecruitmentHeaderRepository,
		MPRequestRepository:                MPRequestRepository,
		ProjectRecruitmentLineRepository:   projectRecruitmentLineRepository,
		ApplicantRepository:                applicantRepository,
		DocumentSendingRepository:          documentSendingRepository,
		MPRequestMessage:                   MPRequestMessage,
		MPRequestService:                   MPRequestService,
		UserHelper:                         userHelper,
		MPRequestHelper:                    MPRequestHelper,
		EmployeeMessage:                    employeeMessage,
		JobPlafonMessage:                   jobPlafonMessage,
	}
}

func DashboardUseCaseFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IDashboardUseCase {
	jobPostingRepository := repository.JobPostingRepositoryFactory(log)
	projectRecruitmentHeaderRepository := repository.ProjectRecruitmentHeaderRepositoryFactory(log)
	MPRequestRepository := repository.MPRequestRepositoryFactory(log)
	projectRecruitmentLineRepository := repository.ProjectRecruitmentLineRepositoryFactory(log)
	applicantRepository := repository.ApplicantRepositoryFactory(log)
	documentSendingRepository := repository.DocumentSendingRepositoryFactory(log)
	MPRequestMessage := messaging.MPRequestMessageFactory(log)
	MPRequestService := service.MPRequestServiceFactory(log)
	userHelper := helper.UserHelperFactory(log)
	MPRequestHelper := helper.MPRequestHelperFactory(log)
	employeeMessage := messaging.EmployeeMessageFactory(log)
	jobPlafonMessage := messaging.JobPlafonMessageFactory(log)
	return NewDashboardUseCase(
		log,
		viper,
		jobPostingRepository,
		projectRecruitmentHeaderRepository,
		MPRequestRepository,
		projectRecruitmentLineRepository,
		applicantRepository,
		documentSendingRepository,
		MPRequestMessage,
		MPRequestService,
		userHelper,
		MPRequestHelper,
		employeeMessage,
		jobPlafonMessage,
	)
}

func (uc *DashboardUseCase) GetDashboard() (*response.DashboardResponse, error) {
	// get total recruitment target
	totalRecruitmentTarget, err := uc.getTotalRecruitmentTarget()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.GetDashboard] failed to get total recruitment target")
		return nil, err
	}

	// get total recruitment realization
	totalRecruitmentRealization, err := uc.getTotalRecruitmentRealization()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.GetDashboard] failed to get total recruitment realization")
		return nil, err
	}

	// get total bilingual
	totalBilingual, err := uc.getTotalBilingual()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.GetDashboard] failed to get total bilingual")
		return nil, err
	}

	// get chart duration recruitment
	chartDurationRecruitment, err := uc.getChartDurationRecruitment()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.GetDashboard] failed to get chart duration recruitment")
		return nil, err
	}

	// get chart job level
	chartJobLevel, err := uc.getChartJobLevel()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.GetDashboard] failed to get chart job level")
		return nil, err
	}

	// get chart department
	chartDepartment, err := uc.getChartDepartment()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.GetDashboard] failed to get chart department")
		return nil, err
	}

	// get avg time to hire
	avgTimeToHire, err := uc.getAvgTimeToHire()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.GetDashboard] failed to get avg time to hire")
		return nil, err
	}

	return &response.DashboardResponse{
		TotalRecruitmentTargetResponse:      *totalRecruitmentTarget,
		TotalRecruitmentRealizationResponse: *totalRecruitmentRealization,
		TotalBilingualResponse:              *totalBilingual,
		ChartDurationRecruitmentResponse:    *chartDurationRecruitment,
		ChartJobLevelResponse:               *chartJobLevel,
		ChartDepartmentResponse:             *chartDepartment,
		AvgTimeToHireResponse:               *avgTimeToHire,
	}, nil
}

func (uc *DashboardUseCase) getTotalRecruitmentTarget() (*response.TotalRecruitmentTargetResponse, error) {
	countMpr := 0
	mpRequests, err := uc.MPRequestRepository.FindAll()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.getTotalRecruitmentTarget] failed to get all MP requests")
		return nil, err
	}
	for _, mpRequest := range *mpRequests {
		_, err := uc.MPRequestMessage.SendFindByIdMessage(mpRequest.MPRCloneID.String())
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when send find by id message: %v", err)
			continue
		}
		countMpr++
	}
	applicantsHired, err := uc.ApplicantRepository.GetAllByKeys(map[string]interface{}{
		"status": entity.APPLICANT_STATUS_HIRED,
	})
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.getTotalRecruitmentTarget] failed to get all hired applicants")
		return nil, err
	}
	totalRecruitmentTarget := countMpr
	var percentage float64
	if totalRecruitmentTarget > 0 {
		percentage = (float64(len(applicantsHired)) / float64(totalRecruitmentTarget)) * 100
	}

	return &response.TotalRecruitmentTargetResponse{
		TotalRecruitmentTarget: totalRecruitmentTarget,
		Percentage:             int(percentage),
	}, nil
}

func (uc *DashboardUseCase) getTotalRecruitmentRealization() (*response.TotalRecruitmentRealizationResponse, error) {
	applicants, err := uc.ApplicantRepository.GetAllByKeys(map[string]interface{}{})
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.getTotalRecruitmentRealization] failed to get all applicants")
		return nil, err
	}

	applicantsHired, err := uc.ApplicantRepository.GetAllByKeys(map[string]interface{}{
		"status": entity.APPLICANT_STATUS_HIRED,
	})
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.getTotalRecruitmentRealization] failed to get all hired applicants")
		return nil, err
	}

	totalRecruitmentRealization := len(applicantsHired)
	var percentage float64
	if len(applicants) > 0 {
		percentage = (float64(len(applicantsHired)) / float64(len(applicants))) * 100
	}

	return &response.TotalRecruitmentRealizationResponse{
		TotalRecruitmentRealization: totalRecruitmentRealization,
		Percentage:                  int(percentage),
	}, nil
}

func (uc *DashboardUseCase) getTotalBilingual() (*response.TotalBilingualResponse, error) {
	applicants, err := uc.ApplicantRepository.GetAllByKeys(map[string]interface{}{
		"status": entity.APPLICANT_STATUS_HIRED,
	})
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.getTotalBilingual] failed to get all applicants")
		return nil, err
	}
	var totalBilingual int
	var TotalNonBilingual int
	for _, applicant := range applicants {
		if applicant.UserProfile.Bilingual == "yes" {
			totalBilingual++
		} else {
			TotalNonBilingual++
		}
	}

	return &response.TotalBilingualResponse{
		TotalBilingual:    totalBilingual,
		TotalNonBilingual: TotalNonBilingual,
	}, nil
}

func (uc *DashboardUseCase) getChartDurationRecruitment() (*response.ChartDurationRecruitmentResponse, error) {
	labels := []string{
		"> 30 Hari",
		"21 - 30 Hari",
		"11 - 20 Hari",
		"1 - 10 Hari",
	}
	datasets := make([]int, 0)
	for _, label := range labels {
		count, err := uc.ProjectRecruitmentHeaderRepository.CountDaysToHireByTotalDays(label)
		if err != nil {
			uc.Log.WithError(err).Error("[DashboardUseCase.getChartDurationRecruitment] failed to count days to hire by total days")
			return nil, err
		}
		datasets = append(datasets, count)
	}

	return &response.ChartDurationRecruitmentResponse{
		Labels:   labels,
		Datasets: datasets,
	}, nil
}

func (uc *DashboardUseCase) getChartJobLevel() (*response.ChartJobLevelResponse, error) {
	jobLevelIDs, err := uc.DocumentSendingRepository.GetJobLevelIdDistinct()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.getChartJobLevel] failed to get job level id distinct")
		return nil, err
	}

	var labels []string
	var datasets []int

	for _, jobLevelID := range jobLevelIDs {
		jobLevelResp, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: jobLevelID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPPlanningUseCase.FindAllLinesByHeaderIdPaginated Message] " + err.Error())
			// return nil, err
			continue
		} else {
			name := strconv.Itoa(int(jobLevelResp.Level)) + " - " + jobLevelResp.Name
			labels = append(labels, name)

			count, err := uc.DocumentSendingRepository.CountByJobLevelID(jobLevelID)
			if err != nil {
				uc.Log.WithError(err).Error("[DashboardUseCase.getChartJobLevel] failed to count by job level id")
				return nil, err
			}

			datasets = append(datasets, count)
		}
	}

	return &response.ChartJobLevelResponse{
		Labels:   labels,
		Datasets: datasets,
	}, nil
}

func (uc *DashboardUseCase) getChartDepartment() (*response.ChartDepartmentResponse, error) {
	chartResp, err := uc.EmployeeMessage.SendGetChartEmployeeOrganizationStructureMessage()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.getChartDepartment] failed to get chart employee organization structure message")
		return nil, err
	}

	return chartResp, nil
}

func (uc *DashboardUseCase) getAvgTimeToHire() (*response.AvgTimeToHireResponse, error) {
	avg, err := uc.ProjectRecruitmentHeaderRepository.CountAverageDaysToHireAll()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.getAvgTimeToHire] failed to count average days to hire all")
		return nil, err
	}

	return &response.AvgTimeToHireResponse{
		AvgTimeToHire: int(avg),
	}, nil
}
