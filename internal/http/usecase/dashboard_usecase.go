package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
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

	return &response.DashboardResponse{
		TotalRecruitmentTargetResponse:      *totalRecruitmentTarget,
		TotalRecruitmentRealizationResponse: *totalRecruitmentRealization,
		TotalBilingualResponse:              *totalBilingual,
	}, nil
}

func (uc *DashboardUseCase) getTotalRecruitmentTarget() (*response.TotalRecruitmentTargetResponse, error) {
	mprResponses := make([]response.MPRequestHeaderResponse, 0)
	mpRequests, err := uc.MPRequestRepository.FindAll()
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.getTotalRecruitmentTarget] failed to get all MP requests")
		return nil, err
	}
	for _, mpRequest := range *mpRequests {
		resp, err := uc.MPRequestMessage.SendFindByIdMessage(mpRequest.MPRCloneID.String())
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when send find by id message: %v", err)
			continue
		}

		convertedData, err := uc.MPRequestService.CheckPortalData(resp)
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when check portal data: %v", err)
			return nil, err
		}
		convertedData.Status = string(mpRequest.Status)
		convertedData.ID = mpRequest.ID

		mprResponses = append(mprResponses, *convertedData)
	}
	applicants, err := uc.ApplicantRepository.GetAllByKeys(map[string]interface{}{
		"status": entity.APPLICANT_STATUS_HIRED,
	})
	if err != nil {
		uc.Log.WithError(err).Error("[DashboardUseCase.getTotalRecruitmentTarget] failed to get all hired applicants")
		return nil, err
	}
	totalRecruitmentTarget := len(mprResponses)
	var percentage int
	if len(applicants) > 0 {
		percentage = len(applicants) / totalRecruitmentTarget * 100
	}

	return &response.TotalRecruitmentTargetResponse{
		TotalRecruitmentTarget: totalRecruitmentTarget,
		Percentage:             percentage,
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
	var percentage int
	if len(applicants) > 0 {
		percentage = len(applicantsHired) / len(applicants) * 100
	}

	return &response.TotalRecruitmentRealizationResponse{
		TotalRecruitmentRealization: totalRecruitmentRealization,
		Percentage:                  percentage,
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

	totalBilingual := len(applicants)

	return &response.TotalBilingualResponse{
		TotalBilingual: totalBilingual,
	}, nil
}
