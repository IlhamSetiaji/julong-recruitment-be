package usecase

import (
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

type IMPRequestUseCase interface {
	CreateMPRequest(req *request.CreateMPRequest) (*entity.MPRequest, error)
	FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) (*response.MPRequestPaginatedResponse, error)
	FindAllPaginatedWhereDoesntHaveJobPosting(jobPostingID string, page int, pageSize int, search string, filter map[string]interface{}) (*response.MPRequestPaginatedResponse, error)
}

type MPRequestUseCase struct {
	Log           *logrus.Logger
	Repository    repository.IMPRequestRepository
	Message       messaging.IMPRequestMessage
	Service       service.IMPRequestService
	JobPostingDTO dto.IJobPostingDTO
	Viper         *viper.Viper
}

func NewMPRequestUseCase(
	log *logrus.Logger,
	repo repository.IMPRequestRepository,
	message messaging.IMPRequestMessage,
	mprService service.IMPRequestService,
	jobPostingDTO dto.IJobPostingDTO,
	viper *viper.Viper,
) IMPRequestUseCase {
	return &MPRequestUseCase{
		Log:           log,
		Repository:    repo,
		Message:       message,
		Service:       mprService,
		JobPostingDTO: jobPostingDTO,
		Viper:         viper,
	}
}

func MPRequestUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IMPRequestUseCase {
	repo := repository.MPRequestRepositoryFactory(log)
	message := messaging.MPRequestMessageFactory(log)
	mprService := service.MPRequestServiceFactory(log)
	jobPostingDTO := dto.JobPostingDTOFactory(log, viper)
	return NewMPRequestUseCase(log, repo, message, mprService, jobPostingDTO, viper)
}

func (uc *MPRequestUseCase) CreateMPRequest(req *request.CreateMPRequest) (*entity.MPRequest, error) {
	mprCloneID := uuid.MustParse(req.MPRCloneID)
	mpRequest := &entity.MPRequest{
		MPRCloneID: &mprCloneID,
		Status:     entity.MPR_STATUS_OPEN,
	}

	createdMPRequest, err := uc.Repository.Create(mpRequest)
	if err != nil {
		return nil, err
	}

	return createdMPRequest, nil
}

func (uc *MPRequestUseCase) FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) (*response.MPRequestPaginatedResponse, error) {
	mprResponses := make([]response.MPRequestHeaderResponse, 0)
	mpRequests, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, filter)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when find all paginated mp request headers: %v", err)
		return nil, err
	}

	// loop and send message to julong_manpower
	for _, mpRequest := range *mpRequests {
		resp, err := uc.Message.SendFindByIdTidakLengkapMessage(mpRequest.MPRCloneID.String())
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when send find by id message: %v", err)
			// return nil, err
			total--
			continue
		}

		convertedData, err := uc.Service.CheckPortalDataMinimal(resp)
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when check portal data: %v", err)
			return nil, err
		}
		convertedData.Status = string(mpRequest.Status)
		convertedData.ID = mpRequest.ID
		if mpRequest.JobPosting != nil {
			convertedData.JobPostingID = &mpRequest.JobPosting.ID
			convertedData.JobPostingDocumentNumber = mpRequest.JobPosting.DocumentNumber
		}

		mprResponses = append(mprResponses, *convertedData)
	}

	return &response.MPRequestPaginatedResponse{
		MPRequestHeader: mprResponses,
		Total:           total,
	}, nil
}

func (uc *MPRequestUseCase) FindAllPaginatedWhereDoesntHaveJobPosting(jobPostingID string, page int, pageSize int, search string, filter map[string]interface{}) (*response.MPRequestPaginatedResponse, error) {
	mprResponses := make([]response.MPRequestHeaderResponse, 0)
	mpRequests, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, filter)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.FindAllPaginatedWhereDoesntHaveJobPosting] error when find all paginated mp request headers where doesn't have job posting: %v", err)
		return nil, err
	}

	// loop and send message to julong_manpower
	for _, mpRequest := range *mpRequests {
		if mpRequest.JobPosting != nil {
			if mpRequest.JobPosting.ID.String() == jobPostingID {
				continue
			}
		}
		resp, err := uc.Message.SendFindByIdMessage(mpRequest.MPRCloneID.String())
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginatedWhereDoesntHaveJobPosting] error when send find by id message: %v", err)
			// return nil, err
			total--
			continue
		}

		convertedData, err := uc.Service.CheckPortalData(resp)
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginatedWhereDoesntHaveJobPosting] error when check portal data: %v", err)
			return nil, err
		}
		convertedData.Status = string(mpRequest.Status)
		convertedData.ID = mpRequest.ID

		mprResponses = append(mprResponses, *convertedData)
	}

	return &response.MPRequestPaginatedResponse{
		MPRequestHeader: mprResponses,
		Total:           total,
	}, nil
}
