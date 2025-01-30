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

type IAdministrativeSelectionUsecase interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.AdministrativeSelectionResponse, int64, error)
	CreateAdministrativeSelection(req *request.CreateAdministrativeSelectionRequest) (*response.AdministrativeSelectionResponse, error)
	FindByID(id string) (*response.AdministrativeSelectionResponse, error)
	UpdateAdministrativeSelection(req *request.UpdateAdministrativeSelectionRequest) (*response.AdministrativeSelectionResponse, error)
	DeleteAdministrativeSelection(id string) error
}

type AdministrativeSelectionUsecase struct {
	Log                  *logrus.Logger
	Repository           repository.IAdministrativeSelectionRepository
	DTO                  dto.IAdministrativeSelectionDTO
	Viper                *viper.Viper
	JobPostingRepository repository.IJobPostingRepository
	ProjectPicRepository repository.IProjectPicRepository
}

func NewAdministrativeSelectionUsecase(
	log *logrus.Logger,
	repo repository.IAdministrativeSelectionRepository,
	asDto dto.IAdministrativeSelectionDTO,
	viper *viper.Viper,
	jpRepo repository.IJobPostingRepository,
	ppRepo repository.IProjectPicRepository,
) IAdministrativeSelectionUsecase {
	return &AdministrativeSelectionUsecase{
		Log:                  log,
		Repository:           repo,
		DTO:                  asDto,
		Viper:                viper,
		JobPostingRepository: jpRepo,
		ProjectPicRepository: ppRepo,
	}
}

func AdministrativeSelectionUsecaseFactory(log *logrus.Logger, viper *viper.Viper) IAdministrativeSelectionUsecase {
	repo := repository.AdministrativeSelectionRepositoryFactory(log)
	asDto := dto.AdministrativeSelectionDTOFactory(log, viper)
	jpRepo := repository.JobPostingRepositoryFactory(log)
	ppRepo := repository.ProjectPicRepositoryFactory(log)
	return NewAdministrativeSelectionUsecase(log, repo, asDto, viper, jpRepo, ppRepo)
}

func (uc *AdministrativeSelectionUsecase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.AdministrativeSelectionResponse, int64, error) {
	entities, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	responses := make([]response.AdministrativeSelectionResponse, 0)
	for _, entity := range *entities {
		res, err := uc.DTO.ConvertEntityToResponse(&entity)
		if err != nil {
			uc.Log.Error("[AdministrativeSelectionUsecase.FindAllPaginated] " + err.Error())
			return nil, 0, err
		}

		responses = append(responses, *res)
	}

	return &responses, total, nil
}

func (uc *AdministrativeSelectionUsecase) CreateAdministrativeSelection(req *request.CreateAdministrativeSelectionRequest) (*response.AdministrativeSelectionResponse, error) {
	parsedJpID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.CreateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJpID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.CreateAdministrativeSelection] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.CreateAdministrativeSelection] " + "Job Posting not found")
		return nil, err
	}

	parsedPpID, err := uuid.Parse(req.ProjectPicID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.CreateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	projectPic, err := uc.ProjectPicRepository.FindByID(parsedPpID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.CreateAdministrativeSelection] " + err.Error())
		return nil, err
	}
	if projectPic == nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.CreateAdministrativeSelection] " + "Project PIC not found")
		return nil, err
	}

	parsedDocDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.CreateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	entity, err := uc.Repository.CreateAdministrativeSelection(&entity.AdministrativeSelection{
		JobPostingID:   jobPosting.ID,
		ProjectPicID:   projectPic.ID,
		Status:         entity.AdministrativeSelectionStatus(req.Status),
		DocumentDate:   parsedDocDate,
		DocumentNumber: req.DocumentNumber,
	})
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.CreateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(entity)
}

func (uc *AdministrativeSelectionUsecase) FindByID(id string) (*response.AdministrativeSelectionResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.FindByID] " + err.Error())
		return nil, err
	}

	entity, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.FindByID] " + err.Error())
		return nil, err
	}

	if entity == nil {
		return nil, nil
	}

	return uc.DTO.ConvertEntityToResponse(entity)
}

func (uc *AdministrativeSelectionUsecase) UpdateAdministrativeSelection(req *request.UpdateAdministrativeSelectionRequest) (*response.AdministrativeSelectionResponse, error) {
	parsedJpID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.UpdateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJpID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.UpdateAdministrativeSelection] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.UpdateAdministrativeSelection] " + "Job Posting not found")
		return nil, err
	}

	parsedPpID, err := uuid.Parse(req.ProjectPicID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.UpdateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	projectPic, err := uc.ProjectPicRepository.FindByID(parsedPpID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.UpdateAdministrativeSelection] " + err.Error())
		return nil, err
	}
	if projectPic == nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.UpdateAdministrativeSelection] " + "Project PIC not found")
		return nil, err
	}

	parsedDocDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.UpdateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	entity, err := uc.Repository.UpdateAdministrativeSelection(&entity.AdministrativeSelection{
		ID:             uuid.MustParse(req.ID),
		JobPostingID:   jobPosting.ID,
		ProjectPicID:   projectPic.ID,
		Status:         entity.AdministrativeSelectionStatus(req.Status),
		DocumentDate:   parsedDocDate,
		DocumentNumber: req.DocumentNumber,
	})
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.UpdateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(entity)
}

func (uc *AdministrativeSelectionUsecase) DeleteAdministrativeSelection(id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.DeleteAdministrativeSelection] " + err.Error())
		return err
	}

	entity, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.DeleteAdministrativeSelection] " + err.Error())
		return err
	}

	if entity == nil {
		return errors.New("Administrative Selection not found")
	}

	return uc.Repository.DeleteAdministrativeSelection(parsedID)
}
