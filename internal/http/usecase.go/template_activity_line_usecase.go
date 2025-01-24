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
)

type ITemplateActivityLineUseCase interface {
	CreateOrUpdateTemplateActivityLine(req *request.CreateOrUpdateTemplateActivityLineRequest) (*response.TemplateActivityResponse, error)
	FindByTemplateActivityID(templateActivityID string) (*[]response.TemplateActivityLineResponse, error)
	FindByID(id uuid.UUID) (*response.TemplateActivityLineResponse, error)
}

type TemplateActivityLineUseCase struct {
	Log                        *logrus.Logger
	Repository                 repository.ITemplateActivityLineRepository
	DTO                        dto.ITemplateActivityLineDTO
	TemplateQuestionRepository repository.ITemplateQuestionRepository
	TemplateActivityRepository repository.ITemplateActivityRepository
	TemplateActivityDTO        dto.ITemplateActivityDTO
}

func NewTemplateActivityLineUseCase(
	log *logrus.Logger,
	repo repository.ITemplateActivityLineRepository,
	talDTO dto.ITemplateActivityLineDTO,
	tqRepository repository.ITemplateQuestionRepository,
	taRepository repository.ITemplateActivityRepository,
	taDTO dto.ITemplateActivityDTO,
) ITemplateActivityLineUseCase {
	return &TemplateActivityLineUseCase{
		Log:                        log,
		Repository:                 repo,
		DTO:                        talDTO,
		TemplateQuestionRepository: tqRepository,
		TemplateActivityRepository: taRepository,
		TemplateActivityDTO:        taDTO,
	}
}

func TemplateActivityLineUseCaseFactory(log *logrus.Logger) ITemplateActivityLineUseCase {
	repo := repository.TemplateActivityLineRepositoryFactory(log)
	talDTO := dto.TemplateActivityLineDTOFactory(log)
	tqRepository := repository.TemplateQuestionRepositoryFactory(log)
	taRepository := repository.TemplateActivityRepositoryFactory(log)
	taDTO := dto.TemplateActivityDTOFactory(log)
	return NewTemplateActivityLineUseCase(log, repo, talDTO, tqRepository, taRepository, taDTO)
}

func (uc *TemplateActivityLineUseCase) CreateOrUpdateTemplateActivityLine(req *request.CreateOrUpdateTemplateActivityLineRequest) (*response.TemplateActivityResponse, error) {
	// check if template activity exist
	parsedTemplateActivityID, err := uuid.Parse(req.TemplateActivityID)
	if err != nil {
		uc.Log.Errorf("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when parsing template activity id: %s", err.Error())
		return nil, err
	}
	ta, err := uc.TemplateActivityRepository.FindByID(parsedTemplateActivityID)
	if err != nil {
		uc.Log.Errorf("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when finding template activity by id: %s", err.Error())
		return nil, err
	}

	if ta == nil {
		uc.Log.Errorf("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] template activity not found")
		return nil, errors.New("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] template activity not found")
	}

	// create of update template activity line
	if len(req.TemplateActivityLines) > 0 {
		for _, templateActivity := range req.TemplateActivityLines {
			tq, err := uc.TemplateQuestionRepository.FindByID(uuid.MustParse(templateActivity.TemplateQuestionID))
			if err != nil {
				uc.Log.Errorf("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when finding template question by id: %s", err.Error())
				return nil, err
			}

			if tq == nil {
				uc.Log.Errorf("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] template question with id %s not found", templateActivity.TemplateQuestionID)
				return nil, errors.New("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] template question with id " + templateActivity.TemplateQuestionID + " not found")
			}
			if templateActivity.ID != "" && templateActivity.ID != uuid.Nil.String() {
				exist, err := uc.Repository.FindByID(uuid.MustParse(templateActivity.ID))
				if err != nil {
					uc.Log.Errorf("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when finding template activity line by id: %s", err.Error())
					return nil, err
				}

				if exist == nil {
					_, err := uc.Repository.CreateTemplateActivityLine(&entity.TemplateActivityLine{
						TemplateActivityID: ta.ID,
						Name:               templateActivity.Name,
						Description:        templateActivity.Description,
						Status:             entity.TemplateActivityLineStatus(templateActivity.Status),
						QuestionTemplateID: tq.ID,
						ColorHexCode:       templateActivity.ColorHexCode,
					})
					if err != nil {
						uc.Log.Errorf("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when creating template activity line: %s", err.Error())
						return nil, errors.New("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when creating template activity line: " + err.Error())
					}
				} else {
					_, err := uc.Repository.UpdateTemplateActivityLine(&entity.TemplateActivityLine{
						ID:                 exist.ID,
						TemplateActivityID: ta.ID,
						Name:               templateActivity.Name,
						Description:        templateActivity.Description,
						Status:             entity.TemplateActivityLineStatus(templateActivity.Status),
						QuestionTemplateID: tq.ID,
						ColorHexCode:       templateActivity.ColorHexCode,
					})
					if err != nil {
						uc.Log.Errorf("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when updating template activity line: %s", err.Error())
						return nil, errors.New("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when updating template activity line: " + err.Error())
					}
				}
			} else {
				_, err := uc.Repository.CreateTemplateActivityLine(&entity.TemplateActivityLine{
					TemplateActivityID: ta.ID,
					Name:               templateActivity.Name,
					Description:        templateActivity.Description,
					Status:             entity.TemplateActivityLineStatus(templateActivity.Status),
					QuestionTemplateID: tq.ID,
					ColorHexCode:       templateActivity.ColorHexCode,
				})
				if err != nil {
					uc.Log.Errorf("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when creating template activity line: %s", err.Error())
					return nil, errors.New("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when creating template activity line: " + err.Error())
				}
			}
		}
	}
	// delete template activity line
	if len(req.DeletedTemplateActivityLineIDs) > 0 {
		for _, id := range req.DeletedTemplateActivityLineIDs {
			err := uc.Repository.DeleteTemplateActivityLine(uuid.MustParse(id))
			if err != nil {
				uc.Log.Errorf("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when deleting template activity line: %s", err.Error())
				return nil, errors.New("[TemplateActivityLineUseCase.CreateOrUpdateTemplateActivityLine] error when deleting template activity line: " + err.Error())
			}
		}
	}

	findTa, err := uc.TemplateActivityRepository.FindByID(parsedTemplateActivityID)

	return uc.TemplateActivityDTO.ConvertEntityToResponse(findTa), nil
}

func (uc *TemplateActivityLineUseCase) FindByTemplateActivityID(templateActivityID string) (*[]response.TemplateActivityLineResponse, error) {
	parsedTemplateActivityID, err := uuid.Parse(templateActivityID)
	if err != nil {
		uc.Log.Errorf("[TemplateActivityLineUseCase.FindByTemplateActivityID] error when parsing template activity id: %s", err.Error())
		return nil, err
	}

	templateActivityLines, err := uc.Repository.FindByTemplateActivityID(parsedTemplateActivityID)
	if err != nil {
		uc.Log.Errorf("[TemplateActivityLineUseCase.FindByTemplateActivityID] error when finding template activity line by template activity id: %s", err.Error())
		return nil, err
	}

	var templateActivityLineResponses []response.TemplateActivityLineResponse
	for _, templateActivityLine := range *templateActivityLines {
		templateActivityLineResponses = append(templateActivityLineResponses, *uc.DTO.ConvertEntityToResponse(&templateActivityLine))
	}

	return &templateActivityLineResponses, nil
}

func (uc *TemplateActivityLineUseCase) FindByID(id uuid.UUID) (*response.TemplateActivityLineResponse, error) {
	templateActivityLine, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Errorf("[TemplateActivityLineUseCase.FindByID] error when finding template activity line by id: %s", err.Error())
		return nil, err
	}

	if templateActivityLine == nil {
		return nil, errors.New("[TemplateActivityLineUseCase.FindByID] template activity line not found")
	}

	return uc.DTO.ConvertEntityToResponse(templateActivityLine), nil
}
