package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IMPRequestUseCase interface {
	CreateMPRequest(req *request.CreateMPRequest) (*entity.MPRequest, error)
}

type MPRequestUseCase struct {
	Log        *logrus.Logger
	Repository repository.IMPRequestRepository
}

func NewMPRequestUseCase(log *logrus.Logger, repo repository.IMPRequestRepository) IMPRequestUseCase {
	return &MPRequestUseCase{
		Log:        log,
		Repository: repo,
	}
}

func MPRequestUseCaseFactory(log *logrus.Logger) IMPRequestUseCase {
	repo := repository.MPRequestRepositoryFactory(log)
	return NewMPRequestUseCase(log, repo)
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
