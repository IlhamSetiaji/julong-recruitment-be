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

type IDocumentVerificationLineUsecase interface {
	CreateOrUpdateDocumentVerificationLine(req *request.CreateOrUpdateDocumentVerificationLine) (*response.DocumentVerificationHeaderResponse, error)
	FindByID(id string) (*response.DocumentVerificationLineResponse, error)
	FindAllByDocumentVerificationHeaderID(documentVerificationHeaderID string) (*[]response.DocumentVerificationLineResponse, error)
}

type DocumentVerificationLineUsecase struct {
	Log                                  *logrus.Logger
	Repository                           repository.IDocumentVerificationLineRepository
	DTO                                  dto.IDocumentVerificationLineDTO
	DocumentVerificationHeaderDTO        dto.IDocumentVerificationHeaderDTO
	DocumentVerificationHeaderRepository repository.IDocumentVerificationHeaderRepository
	DocumentVerificationRepository       repository.IDocumentVerificationRepository
	Viper                                *viper.Viper
}

func NewDocumentVerificationLineUsecase(log *logrus.Logger, repository repository.IDocumentVerificationLineRepository, dto dto.IDocumentVerificationLineDTO, documentVerificationHeaderDTO dto.IDocumentVerificationHeaderDTO, documentVerificationHeaderRepository repository.IDocumentVerificationHeaderRepository, documentVerificationRepository repository.IDocumentVerificationRepository, viper *viper.Viper) IDocumentVerificationLineUsecase {
	return &DocumentVerificationLineUsecase{
		Log:                                  log,
		Repository:                           repository,
		DTO:                                  dto,
		DocumentVerificationHeaderDTO:        documentVerificationHeaderDTO,
		DocumentVerificationHeaderRepository: documentVerificationHeaderRepository,
		DocumentVerificationRepository:       documentVerificationRepository,
		Viper:                                viper,
	}
}

func DocumentVerificationLineFactory(log *logrus.Logger, viper *viper.Viper) IDocumentVerificationLineUsecase {
	dvlRepo := repository.DocumentVerificationLineRepositoryFactory(log)
	dvlDTO := dto.DocumentVerificationLineDTOFactory(log, viper)
	documentVerificationHeaderDTO := dto.DocumentVerificationHeaderDTOFactory(log, viper)
	documentVerificationHeaderRepository := repository.DocumentVerificationHeaderRepositoryFactory(log)
	documentVerificationRepository := repository.DocumentVerificationRepositoryFactory(log)
	return NewDocumentVerificationLineUsecase(log, dvlRepo, dvlDTO, documentVerificationHeaderDTO, documentVerificationHeaderRepository, documentVerificationRepository, viper)
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
