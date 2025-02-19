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

type IDocumentSendingUseCase interface {
	CreateDocumentSending(req *request.CreateDocumentSendingRequest) (*response.DocumentSendingResponse, error)
	FindAllPaginatedByDocumentTypeID(documentTypeID uuid.UUID, page int, pageSize int, search string, sort map[string]interface{}) (*[]response.DocumentSendingResponse, int64, error)
	FindByDocumentTypeIDAndApplicantID(documentTypeID uuid.UUID, applicantID uuid.UUID) (*response.DocumentSendingResponse, error)
	FindByID(id string) (*response.DocumentSendingResponse, error)
	UpdateDocumentSending(req *request.UpdateDocumentSendingRequest) (*response.DocumentSendingResponse, error)
	DeleteDocumentSending(id string) error
	FindAllByDocumentSetupID(documentSetupID string) (*[]response.DocumentSendingResponse, error)
	GenerateDocumentNumber(dateNow time.Time) (string, error)
}

type DocumentSendingUseCase struct {
	Log                              *logrus.Logger
	Repository                       repository.IDocumentSendingRepository
	DTO                              dto.IDocumentSendingDTO
	JobPostingRepository             repository.IJobPostingRepository
	ApplicantRepository              repository.IApplicantRepository
	ProjectRecruitmentLineRepository repository.IProjectRecruitmentLineRepository
	DocumentSetupRepository          repository.IDocumentSetupRepository
	Viper                            *viper.Viper
	DocumentTypeRepository           repository.IDocumentTypeRepository
	DocumentAgreementRepository      repository.IDocumentAgreementRepository
}

func NewDocumentSendingUseCase(
	log *logrus.Logger,
	repo repository.IDocumentSendingRepository,
	dto dto.IDocumentSendingDTO,
	jobPostingRepository repository.IJobPostingRepository,
	applicantRepository repository.IApplicantRepository,
	projectRecruitmentLineRepository repository.IProjectRecruitmentLineRepository,
	documentSetupRepository repository.IDocumentSetupRepository,
	viper *viper.Viper,
	documentTypeRepository repository.IDocumentTypeRepository,
	documentAgreementRepository repository.IDocumentAgreementRepository,
) IDocumentSendingUseCase {
	return &DocumentSendingUseCase{
		Log:                              log,
		Repository:                       repo,
		DTO:                              dto,
		JobPostingRepository:             jobPostingRepository,
		ApplicantRepository:              applicantRepository,
		ProjectRecruitmentLineRepository: projectRecruitmentLineRepository,
		DocumentSetupRepository:          documentSetupRepository,
		Viper:                            viper,
		DocumentTypeRepository:           documentTypeRepository,
		DocumentAgreementRepository:      repository.DocumentAgreementRepositoryFactory(log),
	}
}

func DocumentSendingUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IDocumentSendingUseCase {
	repo := repository.DocumentSendingRepositoryFactory(log)
	dto := dto.DocumentSendingDTOFactory(log, viper)
	jobPostingRepository := repository.JobPostingRepositoryFactory(log)
	applicantRepository := repository.ApplicantRepositoryFactory(log)
	projectRecruitmentLineRepository := repository.ProjectRecruitmentLineRepositoryFactory(log)
	documentSetupRepository := repository.DocumentSetupRepositoryFactory(log)
	documentTypeRepository := repository.DocumentTypeRepositoryFactory(log)
	documentAgreementRepository := repository.DocumentAgreementRepositoryFactory(log)
	return NewDocumentSendingUseCase(log, repo, dto, jobPostingRepository, applicantRepository, projectRecruitmentLineRepository, documentSetupRepository, viper, documentTypeRepository, documentAgreementRepository)
}

func (uc *DocumentSendingUseCase) CreateDocumentSending(req *request.CreateDocumentSendingRequest) (*response.DocumentSendingResponse, error) {
	parsedProjectRecruitmentLineID, err := uuid.Parse(req.ProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByID(parsedProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] project recruitment line not found")
		return nil, errors.New("project recruitment line not found")
	}

	parsedApplicantID, err := uuid.Parse(req.ApplicantID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": parsedApplicantID,
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] applicant not found")
		return nil, errors.New("applicant not found")
	}

	parsedDocumentSetupID, err := uuid.Parse(req.DocumentSetupID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	documentSetup, err := uc.DocumentSetupRepository.FindByID(parsedDocumentSetupID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	if documentSetup == nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] document setup not found")
		return nil, errors.New("document setup not found")
	}

	exist, err := uc.Repository.FindByKeys(map[string]interface{}{
		"document_setup_id": parsedDocumentSetupID,
		"applicant_id":      parsedApplicantID,
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	if exist != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] document sending already exist")
		return nil, errors.New("document sending already exist")
	}

	parsedJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] job posting not found")
		return nil, errors.New("job posting not found")
	}

	var parsedJobLevelID *uuid.UUID
	if req.JobLevelID != "" {
		parsedJobLevelUUID, err := uuid.Parse(req.JobLevelID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedJobLevelID = &parsedJobLevelUUID
	}

	var parsedForOrganizationID *uuid.UUID
	if req.ForOrganizationID != "" {
		parsedForOrganizationUUID, err := uuid.Parse(req.ForOrganizationID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}

		parsedForOrganizationID = &parsedForOrganizationUUID
	}

	parsedDocumentDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}

	var parsedJoinedDate *time.Time
	if req.JoinedDate != "" {
		parsedJoinDate, err := time.Parse("2006-01-02", req.JoinedDate)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedJoinedDate = &parsedJoinDate
	}

	var parsedOrganizationLocationID *uuid.UUID
	if req.OrganizationLocationID != "" {
		parsedOrganizationLocationUUID, err := uuid.Parse(req.OrganizationLocationID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedOrganizationLocationID = &parsedOrganizationLocationUUID
	}

	documentSending, err := uc.Repository.CreateDocumentSending(&entity.DocumentSending{
		DocumentSetupID:          parsedDocumentSetupID,
		ProjectRecruitmentLineID: parsedProjectRecruitmentLineID,
		ApplicantID:              parsedApplicantID,
		JobPostingID:             parsedJobPostingID,
		RecruitmentType:          entity.ProjectRecruitmentType(req.RecruitmentType),
		BasicWage:                req.BasicWage,
		PositionalAllowance:      req.PositionalAllowance,
		OperationalAllowance:     req.OperationalAllowance,
		MealAllowance:            req.MealAllowance,
		JobLocation:              req.JobLocation,
		HometripTicket:           req.HometripTicket,
		PeriodAgreement:          req.PeriodAgreement,
		HomeLocation:             req.HomeLocation,
		JobLevelID:               parsedJobLevelID,
		ForOrganizationID:        parsedForOrganizationID,
		OrganizationLocationID:   parsedOrganizationLocationID,
		DocumentDate:             parsedDocumentDate,
		JoinedDate:               parsedJoinedDate,
		DocumentNumber:           req.DocumentNumber,
		Status:                   entity.DocumentSendingStatus(req.Status),
		DetailContent:            req.DetailContent,
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(documentSending), nil
}

func (uc *DocumentSendingUseCase) FindAllPaginatedByDocumentTypeID(documentTypeID uuid.UUID, page int, pageSize int, search string, sort map[string]interface{}) (*[]response.DocumentSendingResponse, int64, error) {
	docType, err := uc.DocumentTypeRepository.FindByID(documentTypeID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.FindAllPaginatedByDocumentTypeID] " + err.Error())
		return nil, 0, err
	}
	if docType == nil {
		uc.Log.Error("[DocumentSendingUseCase.FindAllPaginatedByDocumentTypeID] document type not found")
		return nil, 0, errors.New("document type not found")
	}

	documentSetups, err := uc.DocumentSetupRepository.FindByDocumentTypeID(documentTypeID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.FindAllPaginatedByDocumentTypeID] " + err.Error())
		return nil, 0, err
	}

	documentSetupIDs := make([]uuid.UUID, 0)

	for _, documentSetup := range documentSetups {
		documentSetupIDs = append(documentSetupIDs, documentSetup.ID)
	}

	documentSendings, total, err := uc.Repository.FindAllPaginatedByDocumentSetupIDs(documentSetupIDs, page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.FindAllPaginatedByDocumentTypeID] " + err.Error())
		return nil, 0, err
	}

	documentSendingResponses := make([]response.DocumentSendingResponse, 0)
	for _, documentSending := range *documentSendings {
		documentSendingResponses = append(documentSendingResponses, *uc.DTO.ConvertEntityToResponse(&documentSending))
	}

	return &documentSendingResponses, total, nil
}

func (uc *DocumentSendingUseCase) FindByDocumentTypeIDAndApplicantID(documentTypeID uuid.UUID, applicantID uuid.UUID) (*response.DocumentSendingResponse, error) {
	docType, err := uc.DocumentTypeRepository.FindByID(documentTypeID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.FindByDocumentTypeIDAndApplicantID] " + err.Error())
		return nil, err
	}
	if docType == nil {
		uc.Log.Error("[DocumentSendingUseCase.FindByDocumentTypeIDAndApplicantID] document type not found")
		return nil, errors.New("document type not found")
	}

	documentSetups, err := uc.DocumentSetupRepository.FindByDocumentTypeID(documentTypeID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.FindByDocumentTypeIDAndApplicantID] " + err.Error())
		return nil, err
	}

	documentSetupIDs := make([]uuid.UUID, 0)

	for _, documentSetup := range documentSetups {
		documentSetupIDs = append(documentSetupIDs, documentSetup.ID)
	}

	documentSending, err := uc.Repository.FindByDocumentSetupIDsAndApplicantID(documentSetupIDs, applicantID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.FindByDocumentTypeIDAndApplicantID] " + err.Error())
		return nil, err
	}

	if documentSending == nil {
		uc.Log.Error("[DocumentSendingUseCase.FindByDocumentTypeIDAndApplicantID] document sending not found")
		return nil, errors.New("document sending not found")
	}

	return uc.DTO.ConvertEntityToResponse(documentSending), nil
}

func (uc *DocumentSendingUseCase) FindByID(id string) (*response.DocumentSendingResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.FindByID] " + err.Error())
		return nil, err
	}

	documentSending, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.FindByID] " + err.Error())
		return nil, err
	}

	if documentSending == nil {
		uc.Log.Error("[DocumentSendingUseCase.FindByID] document sending not found")
		return nil, errors.New("document sending not found")
	}

	return uc.DTO.ConvertEntityToResponse(documentSending), nil
}

func (uc *DocumentSendingUseCase) UpdateDocumentSending(req *request.UpdateDocumentSendingRequest) (*response.DocumentSendingResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}
	docSending, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}
	if docSending == nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] document sending not found")
		return nil, errors.New("document sending not found")
	}

	parsedProjectRecruitmentLineID, err := uuid.Parse(req.ProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByID(parsedProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] project recruitment line not found")
		return nil, errors.New("project recruitment line not found")
	}

	parsedApplicantID, err := uuid.Parse(req.ApplicantID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": parsedApplicantID,
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] applicant not found")
		return nil, errors.New("applicant not found")
	}

	parsedDocumentSetupID, err := uuid.Parse(req.DocumentSetupID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	documentSetup, err := uc.DocumentSetupRepository.FindByID(parsedDocumentSetupID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	if documentSetup == nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] document setup not found")
		return nil, errors.New("document setup not found")
	}

	parsedJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] job posting not found")
		return nil, errors.New("job posting not found")
	}

	var parsedJobLevelID *uuid.UUID
	if req.JobLevelID != "" {
		parsedJobLevelUUID, err := uuid.Parse(req.JobLevelID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedJobLevelID = &parsedJobLevelUUID
	}

	var parsedForOrganizationID *uuid.UUID
	if req.ForOrganizationID != "" {
		parsedForOrganizationUUID, err := uuid.Parse(req.ForOrganizationID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}

		parsedForOrganizationID = &parsedForOrganizationUUID
	}

	parsedDocumentDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
		return nil, err
	}

	var parsedJoinedDate *time.Time
	if req.JoinedDate != "" {
		parsedJoinDate, err := time.Parse("2006-01-02", req.JoinedDate)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedJoinedDate = &parsedJoinDate
	}

	var parsedOrganizationLocationID *uuid.UUID
	if req.OrganizationLocationID != "" {
		parsedOrganizationLocationUUID, err := uuid.Parse(req.OrganizationLocationID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedOrganizationLocationID = &parsedOrganizationLocationUUID
	}

	documentSending, err := uc.Repository.UpdateDocumentSending(&entity.DocumentSending{
		ID:                       parsedID,
		DocumentSetupID:          parsedDocumentSetupID,
		ProjectRecruitmentLineID: parsedProjectRecruitmentLineID,
		ApplicantID:              parsedApplicantID,
		JobPostingID:             parsedJobPostingID,
		RecruitmentType:          entity.ProjectRecruitmentType(req.RecruitmentType),
		OrganizationLocationID:   parsedOrganizationLocationID,
		BasicWage:                req.BasicWage,
		PositionalAllowance:      req.PositionalAllowance,
		OperationalAllowance:     req.OperationalAllowance,
		MealAllowance:            req.MealAllowance,
		JobLocation:              req.JobLocation,
		HometripTicket:           req.HometripTicket,
		PeriodAgreement:          req.PeriodAgreement,
		HomeLocation:             req.HomeLocation,
		JobLevelID:               parsedJobLevelID,
		ForOrganizationID:        parsedForOrganizationID,
		DocumentDate:             parsedDocumentDate,
		JoinedDate:               parsedJoinedDate,
		DocumentNumber:           req.DocumentNumber,
		Status:                   entity.DocumentSendingStatus(req.Status),
		DetailContent:            req.DetailContent,
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}

	if entity.DocumentSendingStatus(req.Status) == entity.DOCUMENT_SENDING_STATUS_APPROVED {
		applicantOrder := applicant.Order
		var TemplateQuestionID *uuid.UUID
		for i := range jobPosting.ProjectRecruitmentHeader.ProjectRecruitmentLines {
			if jobPosting.ProjectRecruitmentHeader.ProjectRecruitmentLines[i].Order == applicantOrder+1 {
				projectRecruitmentLine := &jobPosting.ProjectRecruitmentHeader.ProjectRecruitmentLines[i]
				TemplateQuestionID = &projectRecruitmentLine.TemplateActivityLine.QuestionTemplateID
				break
			} else {
				TemplateQuestionID = &applicant.TemplateQuestionID
			}
		}
		_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
			ID:                 applicant.ID,
			Order:              applicant.Order + 1,
			TemplateQuestionID: *TemplateQuestionID,
		})
		if err != nil {
			uc.Log.Error("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] " + err.Error())
			return nil, err
		}

		documentAgreement, err := uc.DocumentAgreementRepository.FindByKeys(map[string]interface{}{
			"document_sending_id": documentSending.ID,
			"applicant_id":        documentSending.ApplicantID,
			"status":              entity.DOCUMENT_AGREEMENT_STATUS_SUBMITTED,
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		if documentAgreement != nil {
			_, err = uc.DocumentAgreementRepository.UpdateDocumentAgreement(&entity.DocumentAgreement{
				ID:     documentAgreement.ID,
				Status: entity.DOCUMENT_AGREEMENT_STATUS_APPROVED,
			})
		}
	} else if entity.DocumentSendingStatus(req.Status) == entity.DOCUMENT_SENDING_STATUS_COMPLETED {
		applicantOrder := applicant.Order
		var TemplateQuestionID *uuid.UUID
		for i := range jobPosting.ProjectRecruitmentHeader.ProjectRecruitmentLines {
			if jobPosting.ProjectRecruitmentHeader.ProjectRecruitmentLines[i].Order == applicantOrder+1 {
				projectRecruitmentLine := &jobPosting.ProjectRecruitmentHeader.ProjectRecruitmentLines[i]
				TemplateQuestionID = &projectRecruitmentLine.TemplateActivityLine.QuestionTemplateID
				break
			} else {
				TemplateQuestionID = &applicant.TemplateQuestionID
			}
		}
		_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
			ID:                 applicant.ID,
			Order:              applicant.Order + 1,
			TemplateQuestionID: *TemplateQuestionID,
		})
		if err != nil {
			uc.Log.Error("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] " + err.Error())
			return nil, err
		}

		documentAgreement, err := uc.DocumentAgreementRepository.FindByKeys(map[string]interface{}{
			"document_sending_id": documentSending.ID,
			"applicant_id":        documentSending.ApplicantID,
			"status":              entity.DOCUMENT_AGREEMENT_STATUS_SUBMITTED,
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		if documentAgreement != nil {
			_, err = uc.DocumentAgreementRepository.UpdateDocumentAgreement(&entity.DocumentAgreement{
				ID:     documentAgreement.ID,
				Status: entity.DOCUMENT_AGREEMENT_STATUS_COMPLETED,
			})
		}
	} else if entity.DocumentSendingStatus(req.Status) == entity.DOCUMENT_SENDING_STATUS_REJECTED {
		_, err = uc.ApplicantRepository.UpdateApplicantWhenRejected(&entity.Applicant{
			ID: applicant.ID,
		})
		if err != nil {
			uc.Log.Error("[InterviewApplicantUseCase.UpdateFinalResultStatusInterviewApplicant] " + err.Error())
			return nil, err
		}
		documentAgreement, err := uc.DocumentAgreementRepository.FindByKeys(map[string]interface{}{
			"document_sending_id": documentSending.ID,
			"applicant_id":        documentSending.ApplicantID,
			"status":              entity.DOCUMENT_AGREEMENT_STATUS_SUBMITTED,
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		if documentAgreement != nil {
			_, err = uc.DocumentAgreementRepository.UpdateDocumentAgreement(&entity.DocumentAgreement{
				ID:     documentAgreement.ID,
				Status: entity.DOCUMENT_AGREEMENT_STATUS_REJECTED,
			})
		}
	} else {
		documentAgreement, err := uc.DocumentAgreementRepository.FindByKeys(map[string]interface{}{
			"document_sending_id": documentSending.ID,
			"applicant_id":        documentSending.ApplicantID,
			"status":              entity.DOCUMENT_AGREEMENT_STATUS_SUBMITTED,
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		if documentAgreement != nil {
			_, err = uc.DocumentAgreementRepository.UpdateDocumentAgreement(&entity.DocumentAgreement{
				ID:     documentAgreement.ID,
				Status: entity.DOCUMENT_AGREEMENT_STATUS_REVISED,
			})
		}
	}

	return uc.DTO.ConvertEntityToResponse(documentSending), nil
}

func (uc *DocumentSendingUseCase) DeleteDocumentSending(id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.DeleteDocumentSending] " + err.Error())
		return err
	}

	exist, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.DeleteDocumentSending] " + err.Error())
		return err
	}
	if exist == nil {
		uc.Log.Error("[DocumentSendingUseCase.DeleteDocumentSending] document sending not found")
		return errors.New("document sending not found")
	}

	return uc.Repository.DeleteDocumentSending(parsedID)
}

func (uc *DocumentSendingUseCase) FindAllByDocumentSetupID(documentSetupID string) (*[]response.DocumentSendingResponse, error) {
	parsedDocumentSetupID, err := uuid.Parse(documentSetupID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.FindAllByDocumentSetupID] " + err.Error())
		return nil, err
	}

	documentSendings, err := uc.Repository.FindAllByDocumentSetupID(parsedDocumentSetupID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.FindAllByDocumentSetupID] " + err.Error())
		return nil, err
	}

	documentSendingResponses := make([]response.DocumentSendingResponse, 0)
	for _, documentSending := range *documentSendings {
		documentSendingResponses = append(documentSendingResponses, *uc.DTO.ConvertEntityToResponse(&documentSending))
	}

	return &documentSendingResponses, nil
}

func (uc *DocumentSendingUseCase) GenerateDocumentNumber(dateNow time.Time) (string, error) {
	dateStr := dateNow.Format("2006-01-02")
	highestNumber, err := uc.Repository.GetHighestDocumentNumberByDate(dateStr)
	if err != nil {
		uc.Log.Errorf("[DocumentSendingUseCase.GenerateDocumentNumber] " + err.Error())
		return "", err
	}

	newNumber := highestNumber + 1
	documentNumber := fmt.Sprintf("DS/%s/%03d", dateNow.Format("20060102"), newNumber)
	return documentNumber, nil
}
