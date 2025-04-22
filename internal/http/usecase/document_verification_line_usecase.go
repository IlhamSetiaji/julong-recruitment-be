package usecase

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/service"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentVerificationLineUsecase interface {
	CreateOrUpdateDocumentVerificationLine(req *request.CreateOrUpdateDocumentVerificationLine) (*response.DocumentVerificationHeaderResponse, error)
	FindByID(id string) (*response.DocumentVerificationLineResponse, error)
	FindAllByDocumentVerificationHeaderID(documentVerificationHeaderID string) (*[]response.DocumentVerificationLineResponse, error)
	UploadDocumentVerificationLine(req *request.UploadDocumentVerificationLine) (*response.DocumentVerificationLineResponse, error)
	UpdateAnswer(id uuid.UUID, payload *request.UpdateAnswer) (*response.DocumentVerificationLineResponse, error)
}

type DocumentVerificationLineUsecase struct {
	Log                                  *logrus.Logger
	Repository                           repository.IDocumentVerificationLineRepository
	DTO                                  dto.IDocumentVerificationLineDTO
	DocumentVerificationHeaderDTO        dto.IDocumentVerificationHeaderDTO
	DocumentVerificationHeaderRepository repository.IDocumentVerificationHeaderRepository
	DocumentVerificationRepository       repository.IDocumentVerificationRepository
	Viper                                *viper.Viper
	UserMessage                          messaging.IUserMessage
	ApplicantRepository                  repository.IApplicantRepository
	UserProfileRepository                repository.IUserProfileRepository
	EmployeeMessage                      messaging.IEmployeeMessage
	MidsuitService                       service.IMidsuitService
	UserHelper                           helper.IUserHelper
}

func NewDocumentVerificationLineUsecase(
	log *logrus.Logger,
	repository repository.IDocumentVerificationLineRepository,
	dto dto.IDocumentVerificationLineDTO,
	documentVerificationHeaderDTO dto.IDocumentVerificationHeaderDTO,
	documentVerificationHeaderRepository repository.IDocumentVerificationHeaderRepository,
	documentVerificationRepository repository.IDocumentVerificationRepository,
	viper *viper.Viper,
	userMessage messaging.IUserMessage,
	applicantRepository repository.IApplicantRepository,
	userProfileRepository repository.IUserProfileRepository,
	employeeMessage messaging.IEmployeeMessage,
	midsuitService service.IMidsuitService,
	userHelper helper.IUserHelper,
) IDocumentVerificationLineUsecase {
	return &DocumentVerificationLineUsecase{
		Log:                                  log,
		Repository:                           repository,
		DTO:                                  dto,
		DocumentVerificationHeaderDTO:        documentVerificationHeaderDTO,
		DocumentVerificationHeaderRepository: documentVerificationHeaderRepository,
		DocumentVerificationRepository:       documentVerificationRepository,
		Viper:                                viper,
		UserMessage:                          userMessage,
		ApplicantRepository:                  applicantRepository,
		UserProfileRepository:                userProfileRepository,
		EmployeeMessage:                      employeeMessage,
		MidsuitService:                       midsuitService,
		UserHelper:                           userHelper,
	}
}

func DocumentVerificationLineFactory(log *logrus.Logger, viper *viper.Viper) IDocumentVerificationLineUsecase {
	dvlRepo := repository.DocumentVerificationLineRepositoryFactory(log)
	dvlDTO := dto.DocumentVerificationLineDTOFactory(log, viper)
	documentVerificationHeaderDTO := dto.DocumentVerificationHeaderDTOFactory(log, viper)
	documentVerificationHeaderRepository := repository.DocumentVerificationHeaderRepositoryFactory(log)
	documentVerificationRepository := repository.DocumentVerificationRepositoryFactory(log)
	userMessage := messaging.UserMessageFactory(log)
	applicantRepository := repository.ApplicantRepositoryFactory(log)
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	employeeMessage := messaging.EmployeeMessageFactory(log)
	midsuitService := service.MidsuitServiceFactory(viper, log)
	userHelper := helper.UserHelperFactory(log)
	return NewDocumentVerificationLineUsecase(
		log,
		dvlRepo,
		dvlDTO,
		documentVerificationHeaderDTO,
		documentVerificationHeaderRepository,
		documentVerificationRepository,
		viper,
		userMessage,
		applicantRepository,
		userProfileRepository,
		employeeMessage,
		midsuitService,
		userHelper,
	)
}

func (uc *DocumentVerificationLineUsecase) UpdateAnswer(id uuid.UUID, payload *request.UpdateAnswer) (*response.DocumentVerificationLineResponse, error) {
	exist, err := uc.Repository.FindByIDPreload(id)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, errors.New("Document Verification Line not found")
	}

	if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
		umResponse, err := uc.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
			ID: exist.DocumentVerificationHeader.Applicant.UserProfile.UserID.String(),
		})
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		if umResponse.User == nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] user not found")
			return nil, errors.New("user not found")
		}
		employeeID, err := uc.UserHelper.GetEmployeeId(umResponse.User)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		empResp, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: employeeID.String(),
		})
		if err != nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find employee by id message: ", err)
			return nil, err
		}
		if empResp == nil {
			uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] employee not found in midsuit")
			return nil, errors.New("employee not found in midsuit")
		}

		empMidsuitIdInt, err := strconv.Atoi(empResp.MidsuitID)
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}

		authResp, err := uc.MidsuitService.AuthOneStep()
		if err != nil {
			uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
			return nil, err
		}
		switch exist.DocumentVerification.Name {
		case "Nomor KTP":
			if len(payload.Answer) < 16 || len(payload.Answer) > 16 {
				return nil, errors.New("Nomor KTP tidak valid")
			}
			uc.Log.Info("[DocumentSendingUseCase.UpdateDocumentSending] KTP - " + payload.Answer)
			midsuitPayload := request.SyncUpdateEmployeeNationalDataMidsuitRequest{
				HcNationalID1: payload.Answer,
			}
			_, err = uc.MidsuitService.SyncUpdateEmployeeNationalDataMidsuit(empMidsuitIdInt, midsuitPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending - KTP] " + err.Error())
				return nil, err
			}
		case "Nomor NPWP":
			if len(payload.Answer) < 15 || len(payload.Answer) > 16 {
				return nil, errors.New("Nomor NPWP tidak valid")
			}
			uc.Log.Info("[DocumentSendingUseCase.UpdateDocumentSending] NPWP - " + payload.Answer)
			midsuitPayload := request.SyncUpdateEmployeeNationalDataMidsuitRequest{
				HcNationalID3: payload.Answer,
			}
			_, err = uc.MidsuitService.SyncUpdateEmployeeNationalDataMidsuit(empMidsuitIdInt, midsuitPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending - NPWP] " + err.Error())
				return nil, err
			}
		case "Nomor Kartu BPJS TK":
			if len(payload.Answer) < 11 || len(payload.Answer) > 11 {
				return nil, errors.New("Nomor Kartu BPJS TK tidak valid")
			}
			uc.Log.Info("[DocumentSendingUseCase.UpdateDocumentSending] BPJS TK - " + payload.Answer)
			midsuitPayload := request.SyncUpdateEmployeeNationalDataMidsuitRequest{
				HcNationalID4: payload.Answer,
			}
			_, err = uc.MidsuitService.SyncUpdateEmployeeNationalDataMidsuit(empMidsuitIdInt, midsuitPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending - BPJS TK] " + err.Error())
				return nil, err
			}
		case "Nomor Kartu BPJS KS":
			if len(payload.Answer) < 13 || len(payload.Answer) > 13 {
				return nil, errors.New("Nomor Kartu BPJS KS tidak valid")
			}
			uc.Log.Info("[DocumentSendingUseCase.UpdateDocumentSending] BPJS KS - " + payload.Answer)
			midsuitPayload := request.SyncUpdateEmployeeNationalDataMidsuitRequest{
				HcNationalID5: payload.Answer,
			}

			_, err = uc.MidsuitService.SyncUpdateEmployeeNationalDataMidsuit(empMidsuitIdInt, midsuitPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending - BPJS KS] " + err.Error())
				return nil, err
			}
		default:
			uc.Log.Info("[DocumentSendingUseCase.UpdateDocumentSending] " + exist.DocumentVerification.Name + " - " + payload.Answer)
		}
	}

	err = uc.Repository.UpdateAnswer(id, payload.Answer)
	if err != nil {
		return nil, err
	}

	return &response.DocumentVerificationLineResponse{
		ID:     id,
		Answer: payload.Answer,
	}, nil
}

func (uc *DocumentVerificationLineUsecase) CreateOrUpdateDocumentVerificationLine(req *request.CreateOrUpdateDocumentVerificationLine) (*response.DocumentVerificationHeaderResponse, error) {
	parsedDocumentVerificationHeaderID, err := uuid.Parse(req.DocumentVerificationHeaderID)
	if err != nil {
		return nil, err
	}
	documentVerificationHeader, err := uc.DocumentVerificationHeaderRepository.FindByID(parsedDocumentVerificationHeaderID)
	if err != nil {
		return nil, err
	}
	if documentVerificationHeader == nil {
		return nil, errors.New("Document Verification Header not found")
	}

	// create or update document verification line
	for _, documentVerificationLine := range req.DocumentVerificationLines {
		parsedDocumentVerificationID, err := uuid.Parse(documentVerificationLine.DocumentVerificationID)
		if err != nil {
			return nil, err
		}
		documentVerification, err := uc.DocumentVerificationRepository.FindByID(parsedDocumentVerificationID)
		if err != nil {
			return nil, err
		}
		if documentVerification == nil {
			return nil, errors.New("Document Verification not found")
		}
		if documentVerificationLine.ID != "" && documentVerificationLine.ID != uuid.Nil.String() {
			parsedId, err := uuid.Parse(documentVerificationLine.ID)
			if err != nil {
				return nil, errors.New("ID not valid")
			}
			exist, err := uc.Repository.FindByID(parsedId)
			if err != nil {
				return nil, err
			}
			if exist == nil {
				_, err := uc.Repository.CreateDocumentVerificationLine(&entity.DocumentVerificationLine{
					DocumentVerificationHeaderID: parsedDocumentVerificationHeaderID,
					DocumentVerificationID:       parsedDocumentVerificationID,
				})
				if err != nil {
					return nil, err
				}
			} else {
				_, err := uc.Repository.UpdateDocumentVerificationLine(&entity.DocumentVerificationLine{
					ID:                           parsedId,
					DocumentVerificationHeaderID: parsedDocumentVerificationHeaderID,
					DocumentVerificationID:       parsedDocumentVerificationID,
				})
				if err != nil {
					return nil, err
				}
			}
		} else {
			_, err := uc.Repository.CreateDocumentVerificationLine(&entity.DocumentVerificationLine{
				DocumentVerificationHeaderID: parsedDocumentVerificationHeaderID,
				DocumentVerificationID:       parsedDocumentVerificationID,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	// delete document verification lines
	if len(req.DeletedDocumentVerificationLineIDs) > 0 {
		for _, deletedDocumentVerificationLineID := range req.DeletedDocumentVerificationLineIDs {
			parsedDeletedDocumentVerificationLineID, err := uuid.Parse(deletedDocumentVerificationLineID)
			if err != nil {
				return nil, err
			}
			err = uc.Repository.DeleteDocumentVerificationLine(parsedDeletedDocumentVerificationLineID)
			if err != nil {
				return nil, err
			}
		}
	}

	// find document verification header
	dvh, err := uc.DocumentVerificationHeaderRepository.FindByID(parsedDocumentVerificationHeaderID)
	if err != nil {
		return nil, err
	}
	if dvh == nil {
		return nil, errors.New("Document Verification Header not found")
	}

	res := uc.DocumentVerificationHeaderDTO.ConvertEntityToResponse(dvh)
	return res, nil
}

func (uc *DocumentVerificationLineUsecase) FindByID(id string) (*response.DocumentVerificationLineResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	documentVerificationLine, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		return nil, err
	}
	if documentVerificationLine == nil {
		return nil, errors.New("Document Verification Line not found")
	}
	res := uc.DTO.ConvertEntityToResponse(documentVerificationLine)
	return res, nil
}

func (uc *DocumentVerificationLineUsecase) FindAllByDocumentVerificationHeaderID(documentVerificationHeaderID string) (*[]response.DocumentVerificationLineResponse, error) {
	parsedDocumentVerificationHeaderID, err := uuid.Parse(documentVerificationHeaderID)
	if err != nil {
		return nil, err
	}
	documentVerificationLines, err := uc.Repository.FindAllByDocumentVerificationHeaderID(parsedDocumentVerificationHeaderID)
	if err != nil {
		return nil, err
	}

	response := []response.DocumentVerificationLineResponse{}
	for _, documentVerificationLine := range *documentVerificationLines {
		res := uc.DTO.ConvertEntityToResponse(&documentVerificationLine)
		response = append(response, *res)
	}

	return &response, nil
}

func (uc *DocumentVerificationLineUsecase) UploadDocumentVerificationLine(req *request.UploadDocumentVerificationLine) (*response.DocumentVerificationLineResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}
	documentVerificationLine, err := uc.Repository.FindByIDPreload(parsedID)
	if err != nil {
		return nil, err
	}
	if documentVerificationLine == nil {
		return nil, errors.New("Document Verification Line not found")
	}

	if uc.Viper.GetString("midsuit.sync") == "ACTIVE" {
		if documentVerificationLine.DocumentVerification.Name == "Foto Formal" {
			fileContent, err := os.ReadFile(req.Path)
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.CreateEmployeeTaskUseCase] error reading file: ", err)
				return nil, err
			}

			// Extract the file name from the path
			fileName := filepath.Base(req.Path)

			// Encode the file content to base64
			encodedData := base64.StdEncoding.EncodeToString(fileContent)

			umResponse, err := uc.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
				ID: documentVerificationLine.DocumentVerificationHeader.Applicant.UserProfile.UserID.String(),
			})
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}
			if umResponse.User == nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] user not found")
				return nil, errors.New("user not found")
			}
			employeeID, err := uc.UserHelper.GetEmployeeId(umResponse.User)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}
			empResp, err := uc.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
				ID: employeeID.String(),
			})
			if err != nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] error sending find employee by id message: ", err)
				return nil, err
			}
			if empResp == nil {
				uc.Log.Error("[EmployeeTaskUseCase.UpdateEmployeeTaskUseCase] employee not found in midsuit")
				return nil, errors.New("employee not found in midsuit")
			}

			empMidsuitIdInt, err := strconv.Atoi(empResp.MidsuitID)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}

			authResp, err := uc.MidsuitService.AuthOneStep()
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}

			midsuitImagePayload := request.SyncEmployeeImageMidsuitRequest{
				ADClientID: request.ADClientID{
					ID:         1000000,
					Identifier: "Julong Group Indonesia",
				},
				AdOrgId: request.AdOrgId{
					ID:         0,
					Identifier: "*",
				},
				Name:       fileName,
				BinaryData: encodedData,
				ImageURL:   fileName,
				EntityType: request.EntityType{
					ID:         "U",
					Identifier: "User maintained",
				},
			}

			midsuitResp, err := uc.MidsuitService.SyncEmployeeImageMidsuit(midsuitImagePayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}
			if midsuitResp == nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] midsuit response is nil")
				return nil, errors.New("midsuit response is nil")
			}

			empMidsuitPayload := request.SyncUpdateEmployeeImageMidsuitRequest{
				LogoID: request.LogoID{
					Data:     encodedData,
					FileName: fileName,
				},
			}

			empMidsuitResp, err := uc.MidsuitService.SyncUpdateEmployeeImageMidsuit(empMidsuitIdInt, empMidsuitPayload, authResp.Token)
			if err != nil {
				uc.Log.Error("[DocumentSendingUseCase.UpdateDocumentSending] " + err.Error())
				return nil, err
			}

			uc.Log.Info("[DocumentSendingUseCase.UpdateDocumentSending] midsuit response: ", *midsuitResp)
			uc.Log.Info("[DocumentSendingUseCase.UpdateDocumentSending] empMidsuit response: ", *empMidsuitResp)
		}
	}

	// update document verification line
	_, err = uc.Repository.UpdateDocumentVerificationLine(&entity.DocumentVerificationLine{
		ID:   parsedID,
		Path: req.Path,
	})
	if err != nil {
		return nil, err
	}

	ent, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		return nil, err
	}

	res := uc.DTO.ConvertEntityToResponse(ent)
	return res, nil
}
