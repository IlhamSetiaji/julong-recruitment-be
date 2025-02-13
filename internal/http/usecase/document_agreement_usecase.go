package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentAgreementUseCase interface {
	CreateDocumentAgreement(req *request.CreateDocumentAgreementRequest) (*response.DocumentAgreementResponse, error)
	UpdateDocumentAgreement(req *request.UpdateDocumentAgreementRequest) (*response.DocumentAgreementResponse, error)
	FindByDocumentSendingIDAndApplicantID(documentSendingID string, applicantID string) (*response.DocumentAgreementResponse, error)
	UpdateStatusDocumentAgreement(req *request.UpdateStatusDocumentAgreementRequest) (*response.DocumentAgreementResponse, error)
}

type DocumentAgreementUseCase struct {
	Log                       *logrus.Logger
	Repository                repository.IDocumentAgreementRepository
	DocumentSendingRepository repository.IDocumentSendingRepository
	DTO                       dto.IDocumentAgreementDTO
	ApplicantRepository       repository.IApplicantRepository
	Viper                     *viper.Viper
}

func NewDocumentAgreementUseCase(log *logrus.Logger, repository repository.IDocumentAgreementRepository, documentSendingRepository repository.IDocumentSendingRepository, dto dto.IDocumentAgreementDTO, applicantRepository repository.IApplicantRepository, viper *viper.Viper) IDocumentAgreementUseCase {
	return &DocumentAgreementUseCase{
		Log:                       log,
		Repository:                repository,
		DocumentSendingRepository: documentSendingRepository,
		DTO:                       dto,
		ApplicantRepository:       applicantRepository,
		Viper:                     viper,
	}
}

func DocumentAgreementUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IDocumentAgreementUseCase {
	daRepository := repository.DocumentAgreementRepositoryFactory(log)
	documentSendingRepository := repository.DocumentSendingRepositoryFactory(log)
	applicantRepository := repository.ApplicantRepositoryFactory(log)
	dto := dto.DocumentAgreementDTOIDocumentAgreementDTOFactory(log, viper)
	return NewDocumentAgreementUseCase(log, daRepository, documentSendingRepository, dto, applicantRepository, viper)
}

func (uc *DocumentAgreementUseCase) CreateDocumentAgreement(req *request.CreateDocumentAgreementRequest) (*response.DocumentAgreementResponse, error) {
	parsedDocumentSendingID, err := uuid.Parse(req.DocumentSendingID)
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}
	documentSending, err := uc.DocumentSendingRepository.FindByID(parsedDocumentSendingID)
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}
	if documentSending == nil {
		uc.Log.Error("document sending not found")
		return nil, errors.New("document sending not found")
	}

	parsedApplicantID, err := uuid.Parse(req.ApplicantID)
	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{"id": req.ApplicantID})
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("applicant not found")
		return nil, errors.New("applicant not found")
	}

	result, err := uc.Repository.CreateDocumentAgreement(&entity.DocumentAgreement{
		DocumentSendingID: parsedDocumentSendingID,
		ApplicantID:       parsedApplicantID,
		Path:              req.Path,
	})
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(result), nil
}

func (uc *DocumentAgreementUseCase) UpdateDocumentAgreement(req *request.UpdateDocumentAgreementRequest) (*response.DocumentAgreementResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}

	documentAgreement, err := uc.Repository.FindByKeys(map[string]interface{}{"id": parsedID})
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}
	if documentAgreement == nil {
		uc.Log.Error("document agreement not found")
		return nil, errors.New("document agreement not found")
	}

	documentAgreement.Path = req.Path
	result, err := uc.Repository.UpdateDocumentAgreement(documentAgreement)
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(result), nil
}

func (uc *DocumentAgreementUseCase) FindByDocumentSendingIDAndApplicantID(documentSendingID string, applicantID string) (*response.DocumentAgreementResponse, error) {
	parsedDocumentSendingID, err := uuid.Parse(documentSendingID)
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}
	documentSending, err := uc.DocumentSendingRepository.FindByID(parsedDocumentSendingID)
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}
	if documentSending == nil {
		uc.Log.Error("document sending not found")
		return nil, errors.New("document sending not found")
	}

	parsedApplicantID, err := uuid.Parse(applicantID)
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}
	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{"id": applicantID})
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("applicant not found")
		return nil, errors.New("applicant not found")
	}

	documentAgreement, err := uc.Repository.FindByKeys(map[string]interface{}{
		"document_sending_id": parsedDocumentSendingID,
		"applicant_id":        parsedApplicantID,
		"status":              entity.DOCUMENT_AGREEMENT_STATUS_SUBMITTED,
	})
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}
	if documentAgreement == nil {
		uc.Log.Error("document agreement not found")
		return nil, errors.New("document agreement not found")
	}

	return uc.DTO.ConvertEntityToResponse(documentAgreement), nil
}

func (uc *DocumentAgreementUseCase) UpdateStatusDocumentAgreement(req *request.UpdateStatusDocumentAgreementRequest) (*response.DocumentAgreementResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}

	documentAgreement, err := uc.Repository.FindByKeys(map[string]interface{}{"id": parsedID})
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}
	if documentAgreement == nil {
		uc.Log.Error("document agreement not found")
		return nil, errors.New("document agreement not found")
	}

	documentAgreement.Status = entity.DocumentAgreementStatus(req.Status)
	result, err := uc.Repository.UpdateDocumentAgreement(documentAgreement)
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(result), nil
}
