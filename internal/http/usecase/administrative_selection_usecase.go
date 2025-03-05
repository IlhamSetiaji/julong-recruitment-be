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

type IAdministrativeSelectionUsecase interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.AdministrativeSelectionResponse, int64, error)
	FindAllPaginatedPic(employeeID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.AdministrativeSelectionResponse, int64, error)
	CreateAdministrativeSelection(req *request.CreateAdministrativeSelectionRequest) (*response.AdministrativeSelectionResponse, error)
	FindByID(id string) (*response.AdministrativeSelectionResponse, error)
	UpdateAdministrativeSelection(req *request.UpdateAdministrativeSelectionRequest) (*response.AdministrativeSelectionResponse, error)
	DeleteAdministrativeSelection(id string) error
	VerifyAdministrativeSelection(id, verifiedBy string) error
	GenerateDocumentNumber(dateNow time.Time) (string, error)
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

func (uc *AdministrativeSelectionUsecase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.AdministrativeSelectionResponse, int64, error) {
	entities, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort, filter)
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

func (uc *AdministrativeSelectionUsecase) FindAllPaginatedPic(employeeID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.AdministrativeSelectionResponse, int64, error) {
	projectPics, err := uc.ProjectPicRepository.FindAllByEmployeeID(employeeID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.FindAllPaginatedPic] " + err.Error())
		return nil, 0, err
	}

	if len(projectPics) == 0 {
		return &[]response.AdministrativeSelectionResponse{}, 0, nil
	}

	projectPicIDs := make([]uuid.UUID, 0)
	for _, projectPic := range projectPics {
		projectPicIDs = append(projectPicIDs, projectPic.ID)
	}

	entities, total, err := uc.Repository.FindAllPaginatedPic(projectPicIDs, page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.FindAllPaginatedPic] " + err.Error())
		return nil, 0, err
	}

	responses := make([]response.AdministrativeSelectionResponse, 0)
	for _, entity := range *entities {
		res, err := uc.DTO.ConvertEntityToResponse(&entity)
		if err != nil {
			uc.Log.Error("[AdministrativeSelectionUsecase.FindAllPaginatedPic] " + err.Error())
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

func (uc *AdministrativeSelectionUsecase) VerifyAdministrativeSelection(id, verifiedBy string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.VerifyAdministrativeSelection - parsed ID] " + err.Error())
		return err
	}

	parsedVerifiedBy, err := uuid.Parse(verifiedBy)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.VerifyAdministrativeSelection - parsed verified by] " + err.Error())
		return err
	}

	entity, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[AdministrativeSelectionUsecase.VerifyAdministrativeSelection] " + err.Error())
		return err
	}

	if entity == nil {
		return errors.New("Administrative Selection not found")
	}

	return uc.Repository.VerifyAdministrativeSelection(parsedID, parsedVerifiedBy)
}

func (uc *AdministrativeSelectionUsecase) GenerateDocumentNumber(dateNow time.Time) (string, error) {
	dateStr := dateNow.Format("2006-01-02")
	highestNumber, err := uc.Repository.GetHighestDocumentNumberByDate(dateStr)
	if err != nil {
		uc.Log.Errorf("[AdministrativeSelectionUsecase.GenerateDocumentNumber] " + err.Error())
		return "", err
	}

	newNumber := highestNumber + 1
	documentNumber := fmt.Sprintf("JP/%s/%03d", dateNow.Format("20060102"), newNumber)
	return documentNumber, nil
}
