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

type IDocumentVerificationHeaderUseCase interface {
	CreateDocumentVerificationHeader(req *request.CreateDocumentVerificationHeaderRequest) (*response.DocumentVerificationHeaderResponse, error)
	UpdateDocumentVerificationHeader(req *request.UpdateDocumentVerificationHeaderRequest) (*response.DocumentVerificationHeaderResponse, error)
	FindByID(id string) (*response.DocumentVerificationHeaderResponse, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.DocumentVerificationHeaderResponse, int64, error)
	DeleteDocumentVerificationHeader(id string) error
	FindByJobPostingAndApplicant(jobPostingID, applicantID uuid.UUID) (*response.DocumentVerificationHeaderResponse, error)
}

type DocumentVerificationHeaderUseCase struct {
	Log                              *logrus.Logger
	Repository                       repository.IDocumentVerificationHeaderRepository
	DTO                              dto.IDocumentVerificationHeaderDTO
	ProjectRecruitmentLineRepository repository.IProjectRecruitmentLineRepository
	ApplicantRepository              repository.IApplicantRepository
	JobPostingRepository             repository.IJobPostingRepository
	Viper                            *viper.Viper
}

func NewDocumentVerificationHeaderUseCase(
	log *logrus.Logger,
	repo repository.IDocumentVerificationHeaderRepository,
	dto dto.IDocumentVerificationHeaderDTO,
	prlRepository repository.IProjectRecruitmentLineRepository,
	applicantRepository repository.IApplicantRepository,
	jobPostingRepository repository.IJobPostingRepository,
	viper *viper.Viper,
) IDocumentVerificationHeaderUseCase {
	return &DocumentVerificationHeaderUseCase{
		Log:                              log,
		Repository:                       repo,
		DTO:                              dto,
		ProjectRecruitmentLineRepository: prlRepository,
		ApplicantRepository:              applicantRepository,
		JobPostingRepository:             jobPostingRepository,
		Viper:                            viper,
	}
}

func DocumentVerificationHeaderUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IDocumentVerificationHeaderUseCase {
	repo := repository.DocumentVerificationHeaderRepositoryFactory(log)
	dto := dto.DocumentVerificationHeaderDTOFactory(log, viper)
	prlRepository := repository.ProjectRecruitmentLineRepositoryFactory(log)
	applicantRepository := repository.ApplicantRepositoryFactory(log)
	jobPostingRepository := repository.JobPostingRepositoryFactory(log)
	return NewDocumentVerificationHeaderUseCase(log, repo, dto, prlRepository, applicantRepository, jobPostingRepository, viper)
}

func (uc *DocumentVerificationHeaderUseCase) CreateDocumentVerificationHeader(req *request.CreateDocumentVerificationHeaderRequest) (*response.DocumentVerificationHeaderResponse, error) {
	parsedPrlID, err := uuid.Parse(req.ProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.CreateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	prl, err := uc.ProjectRecruitmentLineRepository.FindByID(parsedPrlID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.CreateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	if prl == nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.CreateDocumentVerificationHeader] " + "Project Recruitment Line not found")
		return nil, errors.New("Project Recruitment Line not found")
	}

	parseApplicantID, err := uuid.Parse(req.ApplicantID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.CreateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{"id": parseApplicantID})
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.CreateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.CreateDocumentVerificationHeader] " + "Applicant not found")
		return nil, errors.New("Applicant not found")
	}

	parseJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.CreateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parseJobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.CreateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.CreateDocumentVerificationHeader] " + "Job Posting not found")
		return nil, errors.New("Job Posting not found")
	}

	ent, err := uc.Repository.CreateDocumentVerificationHeader(&entity.DocumentVerificationHeader{
		ProjectRecruitmentLineID: parsedPrlID,
		ApplicantID:              parseApplicantID,
		JobPostingID:             parseJobPostingID,
		Status:                   entity.DocumentVerificationHeaderStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.CreateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(ent), nil
}

func (uc *DocumentVerificationHeaderUseCase) UpdateDocumentVerificationHeader(req *request.UpdateDocumentVerificationHeaderRequest) (*response.DocumentVerificationHeaderResponse, error) {
	ent, err := uc.Repository.FindByID(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	if ent == nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + "Document Verification Header not found")
		return nil, errors.New("Document Verification Header not found")
	}

	parsedPrlID, err := uuid.Parse(req.ProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	prl, err := uc.ProjectRecruitmentLineRepository.FindByID(parsedPrlID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	if prl == nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + "Project Recruitment Line not found")
		return nil, errors.New("Project Recruitment Line not found")
	}

	parseApplicantID, err := uuid.Parse(req.ApplicantID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{"id": parseApplicantID})
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + "Applicant not found")
		return nil, errors.New("Applicant not found")
	}

	parseJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parseJobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + "Job Posting not found")
		return nil, errors.New("Job Posting not found")
	}

	var verifiedBy *uuid.UUID
	uc.Log.Info("Verified by ", req.VerifiedBy)
	if req.VerifiedBy != nil {
		parsedVerifiedBy, err := uuid.Parse(*req.VerifiedBy)
		verifiedBy = &parsedVerifiedBy
		if err != nil {
			uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader verified by] " + err.Error())
			return nil, err
		}
	} else {
		verifiedBy = nil
	}

	updatedData, err := uc.Repository.UpdateDocumentVerificationHeader(&entity.DocumentVerificationHeader{
		ID:                       uuid.MustParse(req.ID),
		ProjectRecruitmentLineID: parsedPrlID,
		ApplicantID:              parseApplicantID,
		JobPostingID:             parseJobPostingID,
		VerifiedBy:               verifiedBy,
		Status:                   entity.DocumentVerificationHeaderStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(updatedData), nil
}

func (uc *DocumentVerificationHeaderUseCase) FindByID(id string) (*response.DocumentVerificationHeaderResponse, error) {
	ent, err := uc.Repository.FindByID(uuid.MustParse(id))
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(ent), nil
}

func (uc *DocumentVerificationHeaderUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.DocumentVerificationHeaderResponse, int64, error) {
	documentVerifications, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	documentVerificationResponses := make([]response.DocumentVerificationHeaderResponse, 0)
	for _, documentVerification := range *documentVerifications {
		documentVerificationResponses = append(documentVerificationResponses, *uc.DTO.ConvertEntityToResponse(&documentVerification))
	}

	return &documentVerificationResponses, total, nil
}

func (uc *DocumentVerificationHeaderUseCase) DeleteDocumentVerificationHeader(id string) error {
	return uc.Repository.DeleteDocumentVerificationHeader(uuid.MustParse(id))
}

func (uc *DocumentVerificationHeaderUseCase) FindByJobPostingAndApplicant(jobPostingID, applicantID uuid.UUID) (*response.DocumentVerificationHeaderResponse, error) {
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.FindByJobPostingAndApplicant] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.FindByJobPostingAndApplicant] " + "Job Posting not found")
		return nil, errors.New("Job Posting not found")
	}

	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{"id": applicantID})
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.FindByJobPostingAndApplicant] " + err.Error())
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.FindByJobPostingAndApplicant] " + "Applicant not found")
		return nil, errors.New("Applicant not found")
	}

	ent, err := uc.Repository.FindByKeys(map[string]interface{}{
		"job_posting_id": jobPostingID,
		"applicant_id":   applicantID,
	})
	if err != nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.FindByJobPostingAndApplicant] " + err.Error())
		return nil, err
	}
	if ent == nil {
		uc.Log.Error("[DocumentVerificationHeaderUseCase.FindByJobPostingAndApplicant] " + "Document Verification Header not found")
		return nil, errors.New("Document Verification Header not found")
	}

	return uc.DTO.ConvertEntityToResponse(ent), nil
}
