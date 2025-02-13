package usecase

// type IDocumentAgreementUseCase interface {
// 	CreateDocumentAgreement(req *request.CreateDocumentAgreementRequest) (*response.DocumentAgreementResponse, error)
// 	UpdateDocumentAgreement(req *request.UpdateDocumentAgreementRequest) (*response.DocumentAgreementResponse, error)
// 	FindByDocumentSendingIDAndApplicantID(documentSendingID string, applicantID string) (*response.DocumentAgreementResponse, error)
// 	UpdateStatusDocumentAgreement(req *request.UpdateStatusDocumentAgreementRequest) (*response.DocumentAgreementResponse, error)
// }

// type DocumentAgreementUseCase struct {
// 	Log                       *logrus.Logger
// 	Repository                repository.IDocumentAgreementRepository
// 	DocumentSendingRepository repository.IDocumentSendingRepository
// 	DTO                       dto.IDocumentAgreementDTO
// 	ApplicantRepository       repository.IApplicantRepository
// 	Viper                     *viper.Viper
// }

// func NewDocumentAgreementUseCase(log *logrus.Logger, repository repository.IDocumentAgreementRepository, documentSendingRepository repository.IDocumentSendingRepository, dto dto.IDocumentAgreementDTO, applicantRepository repository.IApplicantRepository, viper *viper.Viper) IDocumentAgreementUseCase {
// 	return &DocumentAgreementUseCase{
// 		Log:                       log,
// 		Repository:                repository,
// 		DocumentSendingRepository: documentSendingRepository,
// 		DTO:                       dto,
// 		ApplicantRepository:       applicantRepository,
// 		Viper:                     viper,
// 	}
// }

// func DocumentAgreementUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IDocumentAgreementUseCase {
// 	daRepository := repository.DocumentAgreementRepositoryFactory(log)
// 	documentSendingRepository := repository.DocumentSendingRepositoryFactory(log)
// 	applicantRepository := repository.ApplicantRepositoryFactory(log)
// 	dto := dto.DocumentAgreementDTOIDocumentAgreementDTOFactory(log, viper)
// 	return NewDocumentAgreementUseCase(log, daRepository, documentSendingRepository, dto, applicantRepository, viper)
// }

// func (uc *DocumentAgreementUseCase) CreateDocumentAgreement(req *request.CreateDocumentAgreementRequest) (*response.DocumentAgreementResponse, error) {
// 	parsedDocumentSendingID, err := uuid.Parse(req.DocumentSendingID)
// 	if err != nil {
// 		uc.Log.Error(err)
// 		return nil, err
// 	}

// 	documentSending, err := uc.DocumentSendingRepository.FindByID(req.DocumentSendingID)
// 	if err != nil {
// 		uc.Log.Error(err)
// 		return nil, err
// 	}

// 	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{"id": req.ApplicantID})
// 	if err != nil {
// 		uc.Log.Error(err)
// 		return nil, err
// 	}

// 	result, err := uc.Repository.CreateDocumentAgreement()
// 	if err != nil {
// 		uc.Log.Error(err)
// 		return nil, err
// 	}
// }
