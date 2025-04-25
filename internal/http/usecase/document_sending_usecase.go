package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/service"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
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
	TestSendEmail() error
	TestSendEmailWithAttachment(path string) error
	TestGenerateHTMLPDF(docSendingID uuid.UUID) (*string, error)
	GeneratePdf(documentSending *entity.DocumentSending) (*string, error)
	GeneratePdfBuffer(documentSendingID uuid.UUID, text string) ([]byte, error)
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
	MailMessage                      messaging.IMailMessage
	UserHelper                       helper.IUserHelper
	UserMessage                      messaging.IUserMessage
	EmployeeMessage                  messaging.IEmployeeMessage
	TemplateQuestionRepository       repository.ITemplateQuestionRepository
	MPRequestMessage                 messaging.IMPRequestMessage
	MPRequestService                 service.IMPRequestService
	OrganizationMessage              messaging.IOrganizationMessage
	DocumentSendingHelper            helper.IDocumentSendingHelper
	JobPlafonMessage                 messaging.IJobPlafonMessage
	MidsuitService                   service.IMidsuitService
	GradeMessage                     messaging.IGradeMessage
	UserProfileRepository            repository.IUserProfileRepository
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
	mailMessage messaging.IMailMessage,
	userHelper helper.IUserHelper,
	userMessage messaging.IUserMessage,
	employeeMessage messaging.IEmployeeMessage,
	templateQuestionRepository repository.ITemplateQuestionRepository,
	mpRequestMessage messaging.IMPRequestMessage,
	mpRequestService service.IMPRequestService,
	organizationMessage messaging.IOrganizationMessage,
	documentSendingHelper helper.IDocumentSendingHelper,
	jobPlafonMessage messaging.IJobPlafonMessage,
	midsuitService service.IMidsuitService,
	gradeMessage messaging.IGradeMessage,
	userProfileRepository repository.IUserProfileRepository,
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
		MailMessage:                      mailMessage,
		UserHelper:                       userHelper,
		UserMessage:                      userMessage,
		EmployeeMessage:                  employeeMessage,
		TemplateQuestionRepository:       templateQuestionRepository,
		MPRequestMessage:                 mpRequestMessage,
		MPRequestService:                 mpRequestService,
		OrganizationMessage:              organizationMessage,
		DocumentSendingHelper:            documentSendingHelper,
		JobPlafonMessage:                 jobPlafonMessage,
		MidsuitService:                   midsuitService,
		GradeMessage:                     gradeMessage,
		UserProfileRepository:            userProfileRepository,
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
	mailMessage := messaging.MailMessageFactory(log)
	userHelper := helper.UserHelperFactory(log)
	userMessage := messaging.UserMessageFactory(log)
	employeeMessage := messaging.EmployeeMessageFactory(log)
	templateQuestionRepository := repository.TemplateQuestionRepositoryFactory(log)
	mpRequestMessage := messaging.MPRequestMessageFactory(log)
	mpRequestService := service.MPRequestServiceFactory(log)
	organizationMessage := messaging.OrganizationMessageFactory(log)
	documentSendingHelper := helper.DocumentSendingHelperFactory(log, viper)
	jobPlafonMesage := messaging.JobPlafonMessageFactory(log)
	midsuitService := service.MidsuitServiceFactory(viper, log)
	gradeMessage := messaging.GradeMessageFactory(log)
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	return NewDocumentSendingUseCase(
		log,
		repo,
		dto,
		jobPostingRepository,
		applicantRepository,
		projectRecruitmentLineRepository,
		documentSetupRepository,
		viper,
		documentTypeRepository,
		documentAgreementRepository,
		mailMessage,
		userHelper,
		userMessage,
		employeeMessage,
		templateQuestionRepository,
		mpRequestMessage,
		mpRequestService,
		organizationMessage,
		documentSendingHelper,
		jobPlafonMesage,
		midsuitService,
		gradeMessage,
		userProfileRepository,
	)
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

	var parsedJobID *uuid.UUID
	if req.JobID != "" {
		parsedJobUUID, err := uuid.Parse(req.JobID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedJobID = &parsedJobUUID
	}

	var gradeID *uuid.UUID
	if req.GradeID != "" {
		parsedGradeID, err := uuid.Parse(req.GradeID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}
		gradeID = &parsedGradeID
	}

	var allowanceApproval *uuid.UUID
	if req.AllowanceApproval != "" {
		parsedAllowanceApprovalID, err := uuid.Parse(req.AllowanceApproval)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}
		allowanceApproval = &parsedAllowanceApprovalID
	}

	documentSending, err := uc.Repository.CreateDocumentSending(&entity.DocumentSending{
		DocumentSetupID:          parsedDocumentSetupID,
		ProjectRecruitmentLineID: parsedProjectRecruitmentLineID,
		ApplicantID:              parsedApplicantID,
		GradeID:                  gradeID,
		JobPostingID:             parsedJobPostingID,
		RecruitmentType:          entity.ProjectRecruitmentType(req.RecruitmentType),
		BasicWage:                req.BasicWage,
		PositionalAllowance:      req.PositionalAllowance,
		OperationalAllowance:     req.OperationalAllowance,
		MealAllowance:            req.MealAllowance,
		HouseAllowance:           req.HouseAllowance,
		JobLocation:              req.JobLocation,
		HometripTicket:           req.HometripTicket,
		PeriodAgreement:          req.PeriodAgreement,
		HomeLocation:             req.HomeLocation,
		JobLevelID:               parsedJobLevelID,
		JobID:                    parsedJobID,
		ForOrganizationID:        parsedForOrganizationID,
		OrganizationLocationID:   parsedOrganizationLocationID,
		DocumentDate:             parsedDocumentDate,
		JoinedDate:               parsedJoinedDate,
		DocumentNumber:           req.DocumentNumber,
		Status:                   entity.DocumentSendingStatus(req.Status),
		DetailContent:            req.DetailContent,
		SyncMidsuit:              entity.SyncMidsuitEnum(req.SyncMidsuit),
		HiredStatus:              entity.HiredStatusEnum(req.HiredStatus),
		AllowanceApproval:        allowanceApproval,
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

func convertToUTF8(text string) string {
	reader := transform.NewReader(strings.NewReader(text), charmap.Windows1252.NewDecoder())
	utf8Text, err := io.ReadAll(reader)
	if err != nil {
		return text
	}
	return string(utf8Text)
}

func (uc *DocumentSendingUseCase) GeneratePdf(documentSending *entity.DocumentSending) (*string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Define the CSS styles
	cssStyles := `
<style>
body {
	font-size: 17px;
}
.tiptap h1 {
  font-size: 1.4rem;
}

.tiptap h2 {
  font-size: 1.2rem;
}

.tiptap h3 {
  font-size: 1.1rem;
}

.tiptap {
  ul,
  ol {
    padding: 0 1rem;
    margin: 1.25rem 1rem 1.25rem 0.4rem;
  }
  li p {
    margin-top: 0.25em;
    margin-bottom: 0.25em;
  }
  code {
    background-color: var(--purple-light);
    border-radius: 0.4rem;
    color: var(--black);
    font-size: 0.85rem;
    padding: 0.25em 0.3em;
  }

  pre {
    background: var(--black);
    border-radius: 0.5rem;
    color: var(--white);
    font-family: "JetBrainsMono", monospace;
    margin: 1.5rem 0;
    padding: 0.75rem 1rem;

    code {
      background: none;
      color: inherit;
      font-size: 0.8rem;
      padding: 0;
    }
  }

  blockquote {
    border-left: 3px solid var(--gray-3);
    margin: 1.5rem 0;
    padding-left: 1rem;
  }

  hr {
    border: none;
    border-top: 1px solid var(--gray-2);
    margin: 2rem 0;
  }
}

.tiptap table {
  border-collapse: collapse;
  margin: 0;
  overflow: hidden;
  table-layout: fixed;
  width: 100%;
}

.tiptap td,
.tiptap th {
  border: 1px solid var(--primary);
  box-sizing: border-box;
  min-width: 1em;
  padding: 6px 8px;
  position: relative;
  vertical-align: top;
}

.tiptap th {
  background-color: var(--second);
  font-weight: normal !important;
  text-align: left;
}

.tiptap .selectedCell:after {
  background: var(--selectGray);
  content: "";
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  pointer-events: none;
  position: absolute;
  z-index: 2;
}

.tiptap .column-resize-handle {
  background-color: var(--gray);
  bottom: -2px;
  pointer-events: none;
  position: absolute;
  right: -2px;
  top: 0;
  width: 1px;
}

.tiptap .tableWrapper {
  margin: 1.5rem 0;
  overflow-x: auto;
}

.tiptap.resize-cursor {
  cursor: ew-resize;
  cursor: col-resize;
}

.tiptap-border-none {
  border: 0px solid transparent !important;
  background-color: transparent !important;
}

.tiptap-border-none th {
  border: 0px solid transparent !important;
  background-color: transparent !important;
}

.tiptap-border-none td {
  border: 0px solid transparent !important;
  background-color: transparent !important;
}
</style>
`

	organizationResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: documentSending.ForOrganizationID.String(),
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.GeneratePdf] " + err.Error())
		return nil, err
	}

	// Use the organization logo URL
	logoURL := organizationResp.Logo

	// check if document sending is cover letter
	htmlText, err := uc.checkCorrespondingPlaceHolders(documentSending)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.GeneratePdf] " + err.Error())
		return nil, err
	}

	// Wrap the HTML content with proper HTML structure, UTF-8 meta tag, and CSS styles
	htmlContent := `<!DOCTYPE html><html><head><meta charset="UTF-8">` + cssStyles + `</head><body><div style="display: flex; flex-direction: column;">
      <img src="` + logoURL + `" alt="Kop Surat" style="width: 100%;">
      <div style="width: 100%; border-bottom: 3px solid black; "></div>
      </div><div class="tiptap">` + string(*htmlText) + `</div></body></html>`
	dataURL := "data:text/html," + url.PathEscape(htmlContent)

	var pdfBuffer []byte

	err = chromedp.Run(ctx,
		chromedp.Navigate(dataURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuffer, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0.5).
				WithMarginRight(1.0).
				WithMarginBottom(1.0).
				WithMarginLeft(1.0).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		uc.Log.Errorf("Gagal membuat PDF: %v", err)
		return nil, err
	}

	timestamp := time.Now().UnixNano()
	filePath := fmt.Sprintf("storage/generated_pdf/%s", strconv.FormatInt(timestamp, 10)+"_document.pdf")
	err = os.MkdirAll("storage/generated_pdf", os.ModePerm)
	if err != nil {
		uc.Log.Errorf("Gagal membuat direktori: %v", err)
		return nil, err
	}
	err = ioutil.WriteFile(filePath, pdfBuffer, 0644)
	if err != nil {
		uc.Log.Errorf("Gagal membuat PDF: %v", err)
		return nil, err
	}

	return &filePath, nil
}

func (uc *DocumentSendingUseCase) GeneratePdfBuffer(documentSendingID uuid.UUID, text string) ([]byte, error) {
	documentSending, err := uc.Repository.FindByID(documentSendingID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.GeneratePdfBuffer] " + err.Error())
		return nil, err
	}

	if documentSending == nil {
		uc.Log.Error("[DocumentSendingUseCase.GeneratePdfBuffer] document sending not found")
		return nil, errors.New("document sending not found")
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Define the CSS styles
	cssStyles := `
<style>
body {
	font-size: 17px;
}
.tiptap h1 {
  font-size: 1.4rem;
}

.tiptap h2 {
  font-size: 1.2rem;
}

.tiptap h3 {
  font-size: 1.1rem;
}

.tiptap {
  ul,
  ol {
    padding: 0 1rem;
    margin: 1.25rem 1rem 1.25rem 0.4rem;
  }
  li p {
    margin-top: 0.25em;
    margin-bottom: 0.25em;
  }
  code {
    background-color: var(--purple-light);
    border-radius: 0.4rem;
    color: var(--black);
    font-size: 0.85rem;
    padding: 0.25em 0.3em;
  }

  pre {
    background: var(--black);
    border-radius: 0.5rem;
    color: var(--white);
    font-family: "JetBrainsMono", monospace;
    margin: 1.5rem 0;
    padding: 0.75rem 1rem;

    code {
      background: none;
      color: inherit;
      font-size: 0.8rem;
      padding: 0;
    }
  }

  blockquote {
    border-left: 3px solid var(--gray-3);
    margin: 1.5rem 0;
    padding-left: 1rem;
  }

  hr {
    border: none;
    border-top: 1px solid var(--gray-2);
    margin: 2rem 0;
  }
}

.tiptap table {
  border-collapse: collapse;
  margin: 0;
  overflow: hidden;
  table-layout: fixed;
  width: 100%;
}

.tiptap td,
.tiptap th {
  border: 1px solid var(--primary);
  box-sizing: border-box;
  min-width: 1em;
  padding: 6px 8px;
  position: relative;
  vertical-align: top;
}

.tiptap th {
  background-color: var(--second);
  font-weight: normal !important;
  text-align: left;
}

.tiptap .selectedCell:after {
  background: var(--selectGray);
  content: "";
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  pointer-events: none;
  position: absolute;
  z-index: 2;
}

.tiptap .column-resize-handle {
  background-color: var(--gray);
  bottom: -2px;
  pointer-events: none;
  position: absolute;
  right: -2px;
  top: 0;
  width: 1px;
}

.tiptap .tableWrapper {
  margin: 1.5rem 0;
  overflow-x: auto;
}

.tiptap.resize-cursor {
  cursor: ew-resize;
  cursor: col-resize;
}

.tiptap-border-none {
  border: 0px solid transparent !important;
  background-color: transparent !important;
}

.tiptap-border-none th {
  border: 0px solid transparent !important;
  background-color: transparent !important;
}

.tiptap-border-none td {
  border: 0px solid transparent !important;
  background-color: transparent !important;
}
</style>
`

	organizationResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: documentSending.ForOrganizationID.String(),
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.generatePdfBuffer] " + err.Error())
		return nil, err
	}
	// Use the organization logo URL
	logoURL := organizationResp.Logo

	// check if document sending is cover letter
	htmlText, err := uc.checkCorrespondingPlaceHolders(documentSending)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.GeneratePdfBuffer] " + err.Error())
		return nil, err
	}

	// Wrap the HTML content with proper HTML structure, UTF-8 meta tag, and CSS styles
	htmlContent := `<!DOCTYPE html><html><head><meta charset="UTF-8">` + cssStyles + `</head><body><div style="display: flex; flex-direction: column;">
      <img src="` + logoURL + `" alt="Kop Surat" style="width: 100%;">
      <div style="width: 100%; border-bottom: 3px solid black; "></div>
      </div><div class="tiptap">` + *htmlText + `</div></body></html>`
	dataURL := "data:text/html," + url.PathEscape(htmlContent)

	var pdfBuffer []byte

	err = chromedp.Run(ctx,
		chromedp.Navigate(dataURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuffer, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0.5).
				WithMarginRight(1.0).
				WithMarginBottom(1.0).
				WithMarginLeft(1.0).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		uc.Log.Errorf("Gagal membuat PDF: %v", err)
		return nil, err
	}

	return pdfBuffer, nil
}

func (uc *DocumentSendingUseCase) checkCorrespondingPlaceHolders(documentSending *entity.DocumentSending) (*string, error) {
	var htmlText *string
	var err error
	htmlText = &documentSending.DetailContent

	if documentSending.DocumentSetup.DocumentType.Name == "FINAL_RESULT" {
		htmlText, err = uc.replaceCoverLetter(*documentSending)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.checkCorrespondingPlaceHolders] " + err.Error())
			return nil, err
		}
	}

	if documentSending.DocumentSetup.DocumentType.Name == "CONTRACT_DOCUMENT" {
		htmlText, err = uc.replaceContractDocument(*documentSending)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.checkCorrespondingPlaceHolders] " + err.Error())
			return nil, err
		}
	}

	if documentSending.DocumentSetup.DocumentType.Name == "OFFERING_LETTER" {
		htmlText, err = uc.replaceOfferLetter(*documentSending)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.checkCorrespondingPlaceHolders] " + err.Error())
			return nil, err
		}
	}

	return htmlText, nil
}

func (uc *DocumentSendingUseCase) replaceCoverLetter(documentSending entity.DocumentSending) (*string, error) {
	organizationResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: documentSending.ForOrganizationID.String(),
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceCoverLetter] " + err.Error())
		return nil, err
	}
	company := organizationResp.Name

	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": documentSending.ApplicantID,
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceCoverLetter] " + err.Error())
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceCoverLetter] applicant not found")
		return nil, errors.New("applicant not found")
	}

	userID := applicant.UserProfile.UserID
	userMessageResponse, err := uc.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
		ID: userID.String(),
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceCoverLetter] " + err.Error())
		return nil, err
	}
	if userMessageResponse.User == nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceCoverLetter] user not found")
		return nil, errors.New("user not found")
	}

	name, err := uc.UserHelper.GetUserName(userMessageResponse.User)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceCoverLetter] " + err.Error())
		return nil, err
	}

	gender := applicant.UserProfile.Gender
	birthPlace := applicant.UserProfile.BirthPlace
	birthDate := applicant.UserProfile.BirthDate.Format("2006-01-02")
	maritalStatus := applicant.UserProfile.MaritalStatus
	educationLevel := applicant.UserProfile.Educations[0]
	// degreeName := strings.TrimSpace(strings.SplitN(string(educationLevel.EducationLevel), "-", 2)[1])
	degreeName := string(educationLevel.EducationLevel)
	major := educationLevel.Major

	jobPosting, err := uc.JobPostingRepository.FindByID(documentSending.JobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceCoverLetter] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceCoverLetter] job posting not found")
		return nil, errors.New("job posting not found")
	}
	position := jobPosting.Name

	var jobLevel string
	if documentSending.JobLevelID != nil {
		jobLevelResp, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: documentSending.JobLevelID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.replaceCoverLetter] " + err.Error())
			return nil, err
		}
		jobLevel = jobLevelResp.Name
	}

	joinedDate := documentSending.JoinedDate.Format("2006-01-02")
	hiredStatus := documentSending.HiredStatus
	documentDate := documentSending.DocumentDate.Format("2006-01-02")

	replacedData := helper.DocumentDataCoverLetter{
		Company:        company,
		DocumentDate:   documentDate,
		Name:           name,
		Gender:         string(gender),
		BirthPlace:     birthPlace,
		BirthDate:      birthDate,
		EducationLevel: degreeName,
		Major:          major,
		Position:       position,
		JobLevel:       jobLevel,
		JoinedDate:     joinedDate,
		HiredStatus:    string(hiredStatus),
		MaritalStatus:  string(maritalStatus),
	}

	// safeHtmlContent := template.HTMLEscapeString(documentSending.DetailContent)

	htmlContent, err := uc.DocumentSendingHelper.ReplacePlaceHoldersCoverLetter(documentSending.DetailContent, replacedData)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceCoverLetter] " + err.Error())
		return nil, err
	}

	return htmlContent, nil
}

func (uc *DocumentSendingUseCase) replaceContractDocument(documentSending entity.DocumentSending) (*string, error) {
	organizationResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: documentSending.ForOrganizationID.String(),
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceContractDocument] " + err.Error())
		return nil, err
	}
	company := organizationResp.Name

	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": documentSending.ApplicantID,
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceContractDocument] " + err.Error())
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceContractDocument] applicant not found")
		return nil, errors.New("applicant not found")
	}

	userID := applicant.UserProfile.UserID
	userMessageResponse, err := uc.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
		ID: userID.String(),
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceContractDocument] " + err.Error())
		return nil, err
	}
	if userMessageResponse.User == nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceContractDocument] user not found")
		return nil, errors.New("user not found")
	}

	name, err := uc.UserHelper.GetUserName(userMessageResponse.User)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceContractDocument] " + err.Error())
		return nil, err
	}

	gender := applicant.UserProfile.Gender
	birthPlace := applicant.UserProfile.BirthPlace
	birthDate := applicant.UserProfile.BirthDate.Format("2006-01-02")
	maritalStatus := applicant.UserProfile.MaritalStatus
	religion := applicant.UserProfile.Religion
	address := applicant.UserProfile.Address
	// educationLevel := applicant.UserProfile.Educations[0]

	var degreeName string
	var major string
	if applicant.UserProfile.Educations != nil {
		educationLevel := applicant.UserProfile.Educations[0]
		parts := strings.SplitN(string(educationLevel.EducationLevel), "-", 2)
		if len(parts) > 1 {
			degreeName = strings.TrimSpace(parts[1])
		} else {
			uc.Log.Warn("[DocumentSendingUseCase.replaceContractDocument] EducationLevel does not contain '-' delimiter")
			degreeName = strings.TrimSpace(parts[0]) // Use the first part as a fallback
		}
		major = educationLevel.Major
	}

	var jobLevel string
	if documentSending.JobLevelID != nil {
		jobLevelResp, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: documentSending.JobLevelID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.replaceContractDocument] " + err.Error())
			return nil, err
		}
		jobLevel = jobLevelResp.Name
	}

	var jobName string
	if documentSending.JobID != nil {
		jobResp, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: documentSending.JobID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.replaceContractDocument] " + err.Error())
			return nil, err
		}
		jobName = jobResp.Name
	}

	var orgLocationName string
	if documentSending.OrganizationLocationID != nil {
		orgLocationResp, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
			ID: documentSending.OrganizationLocationID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.replaceContractDocument] " + err.Error())
			return nil, err
		}
		orgLocationName = orgLocationResp.Name
	}

	joinedDate := documentSending.JoinedDate.Format("2006-01-02")
	hiredStatus := documentSending.HiredStatus
	documentDate := documentSending.DocumentDate.Format("2006-01-02")

	replacedData := helper.DocumentDataContract{
		Name:                 name,
		Gender:               string(gender),
		BirthPlace:           birthPlace,
		BirthDate:            birthDate,
		MaritalStatus:        string(maritalStatus),
		EducationLevel:       degreeName,
		Major:                major,
		Religion:             string(religion),
		ApplicantAddress:     address,
		Position:             jobName,
		JobLevel:             jobLevel,
		Company:              company,
		Location:             orgLocationName,
		BasicWage:            int(documentSending.BasicWage),
		PositionalAllowance:  int(documentSending.PositionalAllowance),
		OperationalAllowance: int(documentSending.OperationalAllowance),
		MealAllowance:        int(documentSending.MealAllowance),
		HouseAllowance:       int(documentSending.HouseAllowance),
		HometripTicket:       documentSending.HometripTicket,
		JoinedDate:           joinedDate,
		HiredStatus:          string(hiredStatus),
		ApprovalBy:           name,
		DocumentDate:         documentDate,
	}

	// safeHtmlContent := template.HTMLEscapeString(documentSending.DetailContent)

	htmlContent, err := uc.DocumentSendingHelper.ReplacePlaceHoldersContract(documentSending.DetailContent, replacedData)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceContractDocument] " + err.Error())
		return nil, err
	}

	return htmlContent, nil
}

func (uc *DocumentSendingUseCase) replaceOfferLetter(documentSending entity.DocumentSending) (*string, error) {
	organizationResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: documentSending.ForOrganizationID.String(),
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceOfferLetter] " + err.Error())
		return nil, err
	}
	company := organizationResp.Name

	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": documentSending.ApplicantID,
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceOfferLetter] " + err.Error())
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceOfferLetter] applicant not found")
		return nil, errors.New("applicant not found")
	}

	userID := applicant.UserProfile.UserID
	userMessageResponse, err := uc.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
		ID: userID.String(),
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceOfferLetter] " + err.Error())
		return nil, err
	}
	if userMessageResponse.User == nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceOfferLetter] user not found")
		return nil, errors.New("user not found")
	}

	name, err := uc.UserHelper.GetUserName(userMessageResponse.User)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceOfferLetter] " + err.Error())
		return nil, err
	}

	jobPosting, err := uc.JobPostingRepository.FindByID(documentSending.JobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceOfferLetter] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceOfferLetter] job posting not found")
		return nil, errors.New("job posting not found")
	}
	position := jobPosting.Name

	documentDate := documentSending.DocumentDate.Format("2006-01-02")

	replacedData := helper.DocumentDataOfferLetter{
		Company:      company,
		DocumentDate: documentDate,
		Name:         name,
		Position:     position,
		ApprovalBy:   name,
		BasicWage:    int(documentSending.BasicWage),
	}

	// safeHtmlContent := template.HTMLEscapeString(documentSending.DetailContent)

	htmlContent, err := uc.DocumentSendingHelper.ReplacePlaceHoldersOfferLetter(documentSending.DetailContent, replacedData)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.replaceOfferLetter] " + err.Error())
		return nil, err
	}

	return htmlContent, nil
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
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByID(parsedProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] project recruitment line not found")
		return nil, errors.New("project recruitment line not found")
	}

	parsedApplicantID, err := uuid.Parse(req.ApplicantID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}
	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": parsedApplicantID,
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}
	if applicant == nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] applicant not found")
		return nil, errors.New("applicant not found")
	}

	parsedDocumentSetupID, err := uuid.Parse(req.DocumentSetupID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}
	documentSetup, err := uc.DocumentSetupRepository.FindByID(parsedDocumentSetupID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}
	if documentSetup == nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] document setup not found")
		return nil, errors.New("document setup not found")
	}

	parsedJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] job posting not found")
		return nil, errors.New("job posting not found")
	}

	var parsedJobLevelID *uuid.UUID
	if req.JobLevelID != "" {
		parsedJobLevelUUID, err := uuid.Parse(req.JobLevelID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedJobLevelID = &parsedJobLevelUUID
	}

	var parsedForOrganizationID *uuid.UUID
	if req.ForOrganizationID != "" {
		parsedForOrganizationUUID, err := uuid.Parse(req.ForOrganizationID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}

		parsedForOrganizationID = &parsedForOrganizationUUID
	}

	parsedDocumentDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}

	var parsedJoinedDate *time.Time
	if req.JoinedDate != "" {
		parsedJoinDate, err := time.Parse("2006-01-02", req.JoinedDate)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedJoinedDate = &parsedJoinDate
	}

	var parsedOrganizationLocationID *uuid.UUID
	if req.OrganizationLocationID != "" {
		parsedOrganizationLocationUUID, err := uuid.Parse(req.OrganizationLocationID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedOrganizationLocationID = &parsedOrganizationLocationUUID
	}

	var parsedJobID *uuid.UUID
	if req.JobID != "" {
		parsedJobUUID, err := uuid.Parse(req.JobID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		parsedJobID = &parsedJobUUID
	}

	var gradeID *uuid.UUID
	if req.GradeID != "" {
		parsedGradeID, err := uuid.Parse(req.GradeID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		gradeID = &parsedGradeID
	}

	var allowanceApproval *uuid.UUID
	if req.AllowanceApproval != "" {
		parsedAllowanceApprovalID, err := uuid.Parse(req.AllowanceApproval)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.CreateDocumentSending] " + err.Error())
			return nil, err
		}
		allowanceApproval = &parsedAllowanceApprovalID
	}

	documentSending, err := uc.Repository.UpdateDocumentSending(&entity.DocumentSending{
		ID:                       parsedID,
		DocumentSetupID:          parsedDocumentSetupID,
		ProjectRecruitmentLineID: parsedProjectRecruitmentLineID,
		ApplicantID:              parsedApplicantID,
		JobPostingID:             parsedJobPostingID,
		GradeID:                  gradeID,
		RecruitmentType:          entity.ProjectRecruitmentType(req.RecruitmentType),
		OrganizationLocationID:   parsedOrganizationLocationID,
		BasicWage:                req.BasicWage,
		PositionalAllowance:      req.PositionalAllowance,
		OperationalAllowance:     req.OperationalAllowance,
		MealAllowance:            req.MealAllowance,
		HouseAllowance:           req.HouseAllowance,
		JobLocation:              req.JobLocation,
		HometripTicket:           req.HometripTicket,
		PeriodAgreement:          req.PeriodAgreement,
		HomeLocation:             req.HomeLocation,
		JobLevelID:               parsedJobLevelID,
		JobID:                    parsedJobID,
		ForOrganizationID:        parsedForOrganizationID,
		DocumentDate:             parsedDocumentDate,
		JoinedDate:               parsedJoinedDate,
		DocumentNumber:           req.DocumentNumber,
		Status:                   entity.DocumentSendingStatus(req.Status),
		DetailContent:            req.DetailContent,
		SyncMidsuit:              entity.SyncMidsuitEnum(req.SyncMidsuit),
		HiredStatus:              entity.HiredStatusEnum(req.HiredStatus),
		AllowanceApproval:        allowanceApproval,
	})
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}

	if entity.DocumentSendingStatus(req.Status) == entity.DOCUMENT_SENDING_STATUS_PENDING {
		userID := documentSending.Applicant.UserProfile.UserID
		userMessageResponse, err := uc.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
			ID: userID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		if userMessageResponse.User == nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] user not found")
			return nil, errors.New("user not found")
		}

		userEmail, err := uc.UserHelper.GetUserEmail(userMessageResponse.User)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}

		uc.Log.Printf("Organization ID Cok %v", *docSending.ForOrganizationID)

		filePath, err := uc.GeneratePdf(docSending)
		if err != nil {
			uc.Log.Errorf("Gagal membuat PDF: %v", err)
			return nil, err
		}

		// Construct the full URL to the document
		documentURL := uc.Viper.GetString("app.url") + *filePath

		// Update the email body with a proper message and a button
		emailBody := fmt.Sprintf(`
    <html>
    <head>
        <style>
            .button {
                display: inline-block;
                padding: 10px 20px;
                font-size: 16px;
                color: #ffffff;
                background-color: #007BFF;
                text-decoration: none;
                border-radius: 5px;
                text-align: center;
            }
            .button:hover {
                background-color: #0056b3;
            }
        </style>
    </head>
    <body>
        <p>Halo,</p>
        <p>Berikut merupakan dokumen Anda, silahkan klik tombol di bawah ini untuk membuka dokumen Anda:</p>
        <p><a href="%s" class="button">Buka Dokumen</a></p>
        <p>Terima kasih,<br>Tim Kami</p>
    </body>
    </html>
`, documentURL)

		if _, err := uc.MailMessage.SendMail(&request.MailRequest{
			Email:   userEmail,
			Subject: "Pengiriman Dokumen",
			Body:    emailBody,
			From:    "ilham.signals99@gmail.com",
			To:      userEmail,
		}); err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}

		_, err = uc.Repository.UpdateDocumentSending(&entity.DocumentSending{
			ID:     documentSending.ID,
			Status: entity.DOCUMENT_SENDING_STATUS_SENT,
			Path:   *filePath,
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
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
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}
		}
		if req.SyncMidsuit == "YES" {
			err = uc.employeeHired(*applicant, *&docSending.ProjectRecruitmentLine.TemplateActivityLine.QuestionTemplateID, *jobPosting, *documentSending)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}
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
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}
		}

		if req.SyncMidsuit == "YES" {
			err = uc.employeeHired(*applicant, *&docSending.ProjectRecruitmentLine.TemplateActivityLine.QuestionTemplateID, *jobPosting, *documentSending)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}
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
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}
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
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}
		}
	}

	findDocSending, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
		return nil, err
	}

	if findDocSending == nil {
		uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] document sending not found")
		return nil, errors.New("document sending not found")
	}

	resp := uc.DTO.ConvertEntityToResponse(findDocSending)

	return resp, nil
}

func (uc *DocumentSendingUseCase) employeeHired(applicant entity.Applicant, templateQuestionID uuid.UUID, jobPosting entity.JobPosting, documentSending entity.DocumentSending) error {
	tq, err := uc.TemplateQuestionRepository.FindByID(templateQuestionID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.EmployeeHired] " + err.Error())
		return err
	}
	if tq == nil {
		uc.Log.Error("[DocumentSendingUseCase.EmployeeHired] template question not found")
		return errors.New("template question not found")
	}
	uc.Log.Info("Masuk sini cok ane", tq.FormType)
	if tq.FormType == string(entity.TQ_FORM_TYPE_CONTRACT_DOCUMENT) {
		uc.Log.Info("[DocumentSendingUseCase.EmployeeHired] sending message to create employee")
		// send message to create employee
		userID := applicant.UserProfile.UserID
		userMessageResponse, err := uc.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
			ID: userID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return err
		}
		if userMessageResponse.User == nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] user not found")
			return errors.New("user not found")
		}

		userEmail, err := uc.UserHelper.GetUserEmail(userMessageResponse.User)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return err
		}

		resp, err := uc.MPRequestMessage.SendFindByIdMessage(jobPosting.MPRequest.MPRCloneID.String())
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when send find by id message: %v", err)
			return err
		}

		convertedData, err := uc.MPRequestService.CheckPortalData(resp)
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when check portal data: %v", err)
			return err
		}

		_, err = uc.EmployeeMessage.SendCreateEmployeeMessage(request.SendCreateEmployeeMessageRequest{
			UserID:                  applicant.UserProfile.UserID.String(),
			Name:                    applicant.UserProfile.Name,
			Email:                   userEmail,
			JobID:                   jobPosting.JobID.String(),
			JobLevelID:              documentSending.JobLevelID.String(),
			OrganizationID:          convertedData.OrganizationID.String(),
			OrganizationLocationID:  convertedData.OrganizationLocationID.String(),
			OrganizationStructureID: convertedData.ForOrganizationStructureID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return err
		}

		umResponse, err := uc.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
			ID: userID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return err
		}
		if umResponse.User == nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] user not found")
			return errors.New("user not found")
		}
		employeeID, err := uc.UserHelper.GetEmployeeId(umResponse.User)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return err
		}

		organizationID, err := uc.UserHelper.GetOrganizationID(umResponse.User)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return err
		}

		organizationResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: organizationID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return err
		}

		empResp, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: employeeID.String(),
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find employee by id message: ", err)
			return err
		}
		if empResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] employee not found in midsuit")
			return errors.New("employee not found in midsuit")
		}

		orgResp, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: empResp.OrganizationID.String(),
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find organization by id message: ", err)
			return err
		}
		if orgResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] organization not found in midsuit")
			return errors.New("organization not found in midsuit")
		}

		orgStructureId := empResp.EmployeeJob["organization_structure_id"].(string)
		orgStructureResp, err := uc.OrganizationMessage.SendFindOrganizationStructureByIDMessage(request.SendFindOrganizationStructureByIDMessageRequest{
			ID: orgStructureId,
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find organization structure by id message: ", err)
			return err
		}
		if orgStructureResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] organization structure not found in midsuit")
			return errors.New("organization structure not found in midsuit")
		}

		_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
			ID:            applicant.ID,
			Status:        entity.APPLICANT_STATUS_HIRED,
			ProcessStatus: entity.APPLICANT_PROCESS_STATUS_COMPLETED,
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return err
		}

		jobResp, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: jobPosting.JobID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return err
		}

		jobLevelResp, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: documentSending.JobLevelID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return err
		}

		// sync to midsuit
		if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
			authResp, err := uc.MidsuitService.AuthOneStep()
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}

			var filter string
			if documentSending.RecruitmentType == entity.PROJECT_RECRUITMENT_TYPE_MT {
				filter = "MT"
			} else if documentSending.RecruitmentType == entity.PROJECT_RECRUITMENT_TYPE_PH {
				filter = "PH"
			} else if documentSending.RecruitmentType == entity.PROJECT_RECRUITMENT_TYPE_NS {
				filter = "NS"
			}

			recTypeResp, err := uc.MidsuitService.RecruitmentTypeMidsuitAPIWithoutFilter(authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}

			if recTypeResp == nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] recruitment type not found")
				return errors.New("recruitment type not found")
			}

			var recTypeID int
			for _, recType := range recTypeResp.Records {
				if recType.Value == filter {
					recTypeID = recType.ID
					break
				}
			}

			uc.Log.Info("Recruitment Type ID cok: ", recTypeID)
			uc.Log.Info("Education Level cok:", strings.TrimSpace(string(applicant.UserProfile.Educations[0].EducationLevel)))

			eduLevel := strings.TrimSpace(string(applicant.UserProfile.Educations[0].EducationLevel))
			if eduLevel != "S3" && eduLevel != "S2" && eduLevel != "S1" && eduLevel != "D4" && eduLevel != "D3" && eduLevel != "D2" && eduLevel != "D1" && eduLevel != "SMA" && eduLevel != "SMP" && eduLevel != "SD" && eduLevel != "TS" {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] education level not found")
				return errors.New("education level not found")
			}

			// sync to midsuit employee
			midsuitPayload := &request.SyncEmployeeMidsuitRequest{
				AdOrgId: request.AdOrgId{
					ID: func() int {
						id, err := strconv.Atoi(organizationResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
					Identifier: organizationResp.Name,
				},
				Name:     applicant.UserProfile.Name,
				Birthday: applicant.UserProfile.BirthDate.Format("2006-01-02"),
				City:     applicant.UserProfile.BirthPlace,
				Email:    userEmail,
				HcGender: func() request.HcGender {
					if applicant.UserProfile.Gender == entity.MALE {
						return request.HcGender{
							ID: "M",
						}
					} else {
						return request.HcGender{
							ID: "F",
						}
					}
				}(),
				HcMaritalStatus: request.HcMaritalStatus{
					Identifier: strings.Title(string(applicant.UserProfile.MaritalStatus)),
				},
				HcNationalID1: applicant.UserProfile.Ktp,
				HcReligionID: request.HcReligionId{
					ID: 1000002,
					// Identifier: string(applicant.UserProfile.Religion),
				},
				HcStatus: request.HcStatus{
					ID: "A",
				},
				HcBasicAcceptance: request.HcBasicAcceptance{
					ID:         strings.TrimSpace(string(applicant.UserProfile.Educations[0].EducationLevel)),
					Identifier: strings.TrimSpace(string(applicant.UserProfile.Educations[0].EducationLevel)),
				},
				HCWorkStartDate: documentSending.JoinedDate.Format("2006-01-02"),
				HcRecruitmentTypeId: request.HcRecruitmentTypeId{
					ID: recTypeID,
				},
			}

			midsuitEmpID, err := uc.MidsuitService.SyncEmployeeMidsuit(*midsuitPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}

			_, err = uc.EmployeeMessage.SendUpdateEmployeeMidsuitMessage(employeeID.String(), *midsuitEmpID)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}
			uc.Log.Info("Midsuit Employee ID: ", *midsuitEmpID)

			orgStructure, err := uc.OrganizationMessage.SendFindOrganizationStructureByIDMessage(request.SendFindOrganizationStructureByIDMessageRequest{
				ID: convertedData.ForOrganizationStructureID.String(),
			})
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}

			forOrganization, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
				ID: documentSending.ForOrganizationID.String(),
			})
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}

			organizationLocation, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
				ID: documentSending.OrganizationLocationID.String(),
			})
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}

			var gradeMidsuitID *string
			if documentSending.GradeID != nil {
				gradeResp, err := uc.GradeMessage.SendFindByIDMessage(documentSending.GradeID.String())
				if err != nil {
					uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
					return err
				}

				if gradeResp == nil {
					uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] grade not found")
					return errors.New("grade not found")
				}

				gradeMidsuitID = &gradeResp.MidsuitID
			}

			// sync to midsuit employee job
			midsuitEmployeeJobPayload := &request.SyncEmployeeJobMidsuitRequest{
				AdOrgId: request.AdOrgId{
					ID: func() int {
						id, err := strconv.Atoi(organizationResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCCompensation1: int(documentSending.BasicWage),
				HCEmployeeID: request.HcEmployeeId{
					ID: func() int {
						id, err := strconv.Atoi(*midsuitEmpID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert midsuitEmpID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCEmployeeCategoryID: func() *request.HcEmployeeCategoryId {
					if documentSending.HiredStatus == "" {
						return nil
					}
					return &request.HcEmployeeCategoryId{
						Identifier: string(documentSending.HiredStatus),
					}
				}(),
				HCJobID: request.HcJobId{
					ID: func() int {
						id, err := strconv.Atoi(jobResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCJobLevelID: request.HcJobLevelId{
					ID: func() int {
						id, err := strconv.Atoi(jobLevelResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCOrgID: request.HcOrgId{
					ID: func() int {
						id, err := strconv.Atoi(orgStructure.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCWorkStartDate: documentSending.JoinedDate.Format("2006-01-02"),
				// HCRecruitmentTypeID: request.HcRecruitmentTypeId{
				// 	// Identifier: string(documentSending.RecruitmentType),
				// 	ID: recTypeID,
				// },
				ADEmploymentOrgID: request.AdOrgId{
					ID: func() int {
						id, err := strconv.Atoi(forOrganization.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCWorkSiteID: request.HcWorkSiteId{
					ID: func() int {
						id, err := strconv.Atoi(organizationLocation.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				IsPrimary: true,
				HCEmployeeGradeID: request.HcEmployeeGradeId{
					ID: func() int {
						if gradeMidsuitID == nil {
							return 0
						}
						id, err := strconv.Atoi(*gradeMidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				ModelName: "hc_employeejob",
			}

			uc.Log.Info("Midsuit Employee Job Payload cok: ", *midsuitEmployeeJobPayload)

			_, err = uc.MidsuitService.SyncEmployeeJobMidsuit(*midsuitEmployeeJobPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}

			// sync to midsuit employee work experiences
			if len(applicant.UserProfile.WorkExperiences) > 0 {
				for _, workExperience := range applicant.UserProfile.WorkExperiences {
					workExperiencePayload := &request.SyncEmployeeWorkExperienceMidsuitRequest{
						AdOrgId: request.AdOrgId{
							ID: func() int {
								id, err := strconv.Atoi(organizationResp.MidsuitID)
								if err != nil {
									uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
									return 0
								}
								return id
							}(),
						},
						HCEmployeeID: request.HcEmployeeId{
							ID: func() int {
								id, err := strconv.Atoi(*midsuitEmpID)
								if err != nil {
									uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert midsuitEmpID to int: " + err.Error())
									return 0
								}
								return id
							}(),
						},
						Name:           workExperience.Name,
						Description:    "Mwehehe",
						YearExperience: string(workExperience.YearExperience),
						ModelName:      "hc_workhistory",
					}

					_, err = uc.MidsuitService.SyncEmployeeWorkExperienceMidsuit(*workExperiencePayload, authResp.Token)
					if err != nil {
						uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
						return err
					}
				}
			}

			// sync to midsuit employee educations
			if len(applicant.UserProfile.Educations) > 0 {
				for i, education := range applicant.UserProfile.Educations {
					educationPayload := &request.SyncEmployeeEducationMidsuitRequest{
						AdOrgId: request.AdOrgId{
							ID: func() int {
								id, err := strconv.Atoi(organizationResp.MidsuitID)
								if err != nil {
									uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
									return 0
								}
								return id
							}(),
						},
						HCEmployeeID: request.HcEmployeeId{
							ID: func() int {
								id, err := strconv.Atoi(*midsuitEmpID)
								if err != nil {
									uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert midsuitEmpID to int: " + err.Error())
									return 0
								}
								return id
							}(),
						},
						BidangPendidikanAkhir: education.Major,
						HcEducationInstitute:  education.SchoolName,
						HcGpaScore:            int(*education.Gpa),
						SeqNo:                 10,
						HCBasicAcceptance: request.HcBasicAcceptance{
							// Identifier: applicant.UserProfile.Educations[i].Major,
							ID: strings.TrimSpace(string(applicant.UserProfile.Educations[i].EducationLevel)),
						},
						ModelName: "hc_employeeeducation",
					}

					_, err = uc.MidsuitService.SyncEmployeeEducationMidsuit(*educationPayload, authResp.Token)
					if err != nil {
						uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
						return err
					}
				}
			}

			// sync to midsuit employee operation
			empRespVerifiedBy, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
				ID: documentSending.AllowanceApproval.String(),
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find employee by id message: ", err)
				return err
			}
			if empRespVerifiedBy == nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] employee not found in midsuit")
				return errors.New("employee not found in midsuit")
			}

			var joinedDate time.Time
			if documentSending.JoinedDate.Day() == 25 {
				joinedDate = documentSending.JoinedDate.AddDate(0, 1, 0)
			} else {
				joinedDate = *documentSending.JoinedDate
			}

			cPeriodID := joinedDate.Format("Jan-06")
			fmt.Println("Formatted cPeriodID:", cPeriodID)
			// cPeriodID := documentSending.JoinedDate.Format("Jan-06")

			allowanceOperationPayload := &request.SyncEmployeeAllowanceMidsuitRequest{
				AdOrgId: request.AdOrgId{
					ID: func() int {
						id, err := strconv.Atoi(organizationResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				CDocTypeID: request.CDocTypeID{
					Identifier: "Allowance Provision",
					ModelName:  "c_doctype",
				},
				CPeriodID: request.CPeriodID{
					// Identifier: "Apr-25",
					Identifier: cPeriodID,
					ModelName:  "c_period",
				},
				DateDoc: documentSending.JoinedDate.Format("2006-01-02"),
				HCEmployee2ID: request.HcEmployeeId{
					ID: func() int {
						id, err := strconv.Atoi(*midsuitEmpID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert midsuitEmpID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCJob2ID: request.HcJobId{
					ID: func() int {
						id, err := strconv.Atoi(jobResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCEmployeeID: request.HcEmployeeId{
					ID: func() int {
						id, err := strconv.Atoi(empRespVerifiedBy.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCOrg2ID: request.HcOrgId{
					ID: func() int {
						id, err := strconv.Atoi(orgStructure.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCAllowanceType: request.HCAllowanceType{
					ID: "OPR",
				},
				Distance: 1,
				Amount:   int(documentSending.OperationalAllowance),
				HCUOM: request.HCUOM{
					ID: "MO",
				},
				IsUseDate: false,
				HCProvisionType: request.HCProvisionType{
					ID: "GRA",
				},
				IsGenerated: true,
				ModelName:   "hc_allowanceprovision",
				DocAction:   "CO",
			}

			_, err = uc.MidsuitService.SyncEmployeeAllowanceMidsuit(*allowanceOperationPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending Allowance Operation] " + err.Error())
				return err
			}

			// sync to midsuit employee allowance meal
			allowanceMealPayload := &request.SyncEmployeeAllowanceMidsuitRequest{
				AdOrgId: request.AdOrgId{
					ID: func() int {
						id, err := strconv.Atoi(organizationResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				CDocTypeID: request.CDocTypeID{
					Identifier: "Allowance Provision",
					ModelName:  "c_doctype",
				},
				CPeriodID: request.CPeriodID{
					// Identifier: "Apr-25",
					Identifier: cPeriodID,
					ModelName:  "c_period",
				},
				DateDoc: documentSending.JoinedDate.Format("2006-01-02"),
				HCEmployee2ID: request.HcEmployeeId{
					ID: func() int {
						id, err := strconv.Atoi(*midsuitEmpID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert midsuitEmpID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCJob2ID: request.HcJobId{
					ID: func() int {
						id, err := strconv.Atoi(jobResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCEmployeeID: request.HcEmployeeId{
					ID: func() int {
						id, err := strconv.Atoi(empRespVerifiedBy.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCOrg2ID: request.HcOrgId{
					ID: func() int {
						id, err := strconv.Atoi(orgStructure.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCAllowanceType: request.HCAllowanceType{
					ID: "MEL",
				},
				Distance: 1,
				Amount:   int(documentSending.MealAllowance),
				HCUOM: request.HCUOM{
					ID: "MO",
				},
				IsUseDate: false,
				HCProvisionType: request.HCProvisionType{
					ID: "GRA",
				},
				IsGenerated: true,
				ModelName:   "hc_allowanceprovision",
				// DocAction:   "CO",
			}

			_, err = uc.MidsuitService.SyncEmployeeAllowanceMidsuit(*allowanceMealPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending Allowance Meal] " + err.Error())
				return err
			}

			// sync to midsuit employee allowance meal
			allowanceHousePayload := &request.SyncEmployeeAllowanceMidsuitRequest{
				AdOrgId: request.AdOrgId{
					ID: func() int {
						id, err := strconv.Atoi(organizationResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				CDocTypeID: request.CDocTypeID{
					Identifier: "Allowance Provision",
					ModelName:  "c_doctype",
				},
				CPeriodID: request.CPeriodID{
					// Identifier: "Apr-25",
					Identifier: cPeriodID,
					ModelName:  "c_period",
				},
				DateDoc: documentSending.JoinedDate.Format("2006-01-02"),
				HCEmployee2ID: request.HcEmployeeId{
					ID: func() int {
						id, err := strconv.Atoi(*midsuitEmpID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert midsuitEmpID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCJob2ID: request.HcJobId{
					ID: func() int {
						id, err := strconv.Atoi(jobResp.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCEmployeeID: request.HcEmployeeId{
					ID: func() int {
						id, err := strconv.Atoi(empRespVerifiedBy.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCOrg2ID: request.HcOrgId{
					ID: func() int {
						id, err := strconv.Atoi(orgStructure.MidsuitID)
						if err != nil {
							uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert MidsuitID to int: " + err.Error())
							return 0
						}
						return id
					}(),
				},
				HCAllowanceType: request.HCAllowanceType{
					ID: "HOS",
				},
				Distance: 1,
				Amount:   int(documentSending.HouseAllowance),
				HCUOM: request.HCUOM{
					ID: "MO",
				},
				IsUseDate: false,
				HCProvisionType: request.HCProvisionType{
					ID: "GRA",
				},
				IsGenerated: true,
				ModelName:   "hc_allowanceprovision",
				DocAction:   "CO",
			}

			_, err = uc.MidsuitService.SyncEmployeeAllowanceMidsuit(*allowanceHousePayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending Allowance House] " + err.Error())
				return err
			}

			midsuitEmpIDInt, err := strconv.Atoi(*midsuitEmpID)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] failed to convert midsuitEmpID to int: " + err.Error())
				return err
			}
			userProfileMidsuitID, err := uc.MidsuitService.SyncGenerateUserMidsuit(midsuitEmpIDInt, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}

			// update user profile
			_, err = uc.UserProfileRepository.UpdateUserProfile(&entity.UserProfile{
				ID:        applicant.UserProfile.ID,
				MidsuitID: userProfileMidsuitID,
			})
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}

			_, err = uc.EmployeeMessage.SendCreateEmployeeTaskMessage(request.SendCreateEmployeeTaskMessageRequest{
				EmployeeID:            employeeID.String(),
				JoinedDate:            documentSending.JoinedDate.String(),
				OrganizationType:      organizationResp.OrganizationType,
				EmployeeMidsuitID:     *midsuitEmpID,
				JobMidsuitID:          jobResp.MidsuitID,
				JobLevelMidsuitID:     jobLevelResp.MidsuitID,
				OrgMidsuitID:          orgResp.MidsuitID,
				OrgStructureMidsuitID: orgStructureResp.MidsuitID,
			})
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return err
			}
		}
	}

	return nil
}

func (uc *DocumentSendingUseCase) DeleteDocumentSending(id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.DeleteDocumentSending] or" + err.Error())
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

func (uc *DocumentSendingUseCase) TestSendEmail() error {
	if _, err := uc.MailMessage.SendMail(&request.MailRequest{
		Email:   "ilham.ahmadz18@gmail.com",
		Subject: "Email Verification",
		Body:    "Halo brother",
		From:    "ilham.signals99@gmail.com",
		To:      "ilham.ahmadz18@gmail.com",
	}); err != nil {
		uc.Log.Error("[UserUseCase.Register] " + err.Error())
		return err
	}

	return nil
}

func (uc *DocumentSendingUseCase) TestSendEmailWithAttachment(attachmentPath string) error {
	// Create the mail request with the attachment
	mailRequest := &request.MailRequest{
		Email:   "ilham.ahmadz18@gmail.com",
		Subject: "Email Verification",
		Body:    "Halo brother",
		From:    "ilham.signals99@gmail.com",
		To:      "ilham.ahmadz18@gmail.com",
		Attach:  attachmentPath,
	}

	// Send the email
	if _, err := uc.MailMessage.SendMail(mailRequest); err != nil {
		uc.Log.Error("[DocumentSendingUseCase.TestSendEmailWithAttachment] " + err.Error())
		return err
	}

	return nil
}

func (uc *DocumentSendingUseCase) TestGenerateHTMLPDF(docSendingID uuid.UUID) (*string, error) {
	docSending, err := uc.Repository.FindByID(docSendingID)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.TestGenerateHTMLPDF] " + err.Error())
		return nil, err
	}
	// Generate the PDF
	filePath, err := uc.GeneratePdf(docSending)
	if err != nil {
		uc.Log.Error("[DocumentSendingUseCase.TestGenerateHTMLPDF] " + err.Error())
		return nil, err
	}

	return filePath, nil
}
