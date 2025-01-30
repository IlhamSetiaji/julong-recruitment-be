package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/service"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IMPRequestUseCase interface {
	CreateMPRequest(req *request.CreateMPRequest) (*entity.MPRequest, error)
	FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) (*response.MPRequestPaginatedResponse, error)
}

type MPRequestUseCase struct {
	Log        *logrus.Logger
	Repository repository.IMPRequestRepository
	Message    messaging.IMPRequestMessage
	Service    service.IMPRequestService
}

func NewMPRequestUseCase(
	log *logrus.Logger,
	repo repository.IMPRequestRepository,
	message messaging.IMPRequestMessage,
	mprService service.IMPRequestService,
) IMPRequestUseCase {
	return &MPRequestUseCase{
		Log:        log,
		Repository: repo,
		Message:    message,
		Service:    mprService,
	}
}

func MPRequestUseCaseFactory(log *logrus.Logger) IMPRequestUseCase {
	repo := repository.MPRequestRepositoryFactory(log)
	message := messaging.MPRequestMessageFactory(log)
	mprService := service.MPRequestServiceFactory(log)
	return NewMPRequestUseCase(log, repo, message, mprService)
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
		resp, err := uc.Message.SendFindByIdMessage(mpRequest.MPRCloneID.String())
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when send find by id message: %v", err)
			return nil, err
		}

		convertedData, err := uc.Service.CheckPortalData(resp)
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when check portal data: %v", err)
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
