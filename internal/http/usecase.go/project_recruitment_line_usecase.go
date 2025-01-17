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
)

type IProjectRecruitmentLineUseCase interface {
	CreateOrUpdateProjectRecruitmentLines(req *request.CreateOrUpdateProjectRecruitmentLinesRequest) (*response.ProjectRecruitmentHeaderResponse, error)
}

type ProjectRecruitmentLineUseCase struct {
	Log                                *logrus.Logger
	Repository                         repository.IProjectRecruitmentLineRepository
	DTO                                dto.IProjectRecruitmentLineDTO
	TemplateActivityLineRepository     repository.ITemplateActivityLineRepository
	ProjectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository
	ProjectPicRepository               repository.IProjectPicRepository
	ProjectRecruitmentHeaderDTO        dto.IProjectRecruitmentHeaderDTO
}

func NewProjectRecruitmentLineUseCase(
	log *logrus.Logger,
	repo repository.IProjectRecruitmentLineRepository,
	dto dto.IProjectRecruitmentLineDTO,
	talRepo repository.ITemplateActivityLineRepository,
	prhRepo repository.IProjectRecruitmentHeaderRepository,
	picRepo repository.IProjectPicRepository,
	prhDTO dto.IProjectRecruitmentHeaderDTO,
) IProjectRecruitmentLineUseCase {
	return &ProjectRecruitmentLineUseCase{
		Log:                                log,
		Repository:                         repo,
		DTO:                                dto,
		TemplateActivityLineRepository:     talRepo,
		ProjectRecruitmentHeaderRepository: prhRepo,
		ProjectPicRepository:               picRepo,
		ProjectRecruitmentHeaderDTO:        prhDTO,
	}
}

func ProjectRecruitmentLineUseCaseFactory(log *logrus.Logger) IProjectRecruitmentLineUseCase {
	repo := repository.ProjectRecruitmentLineRepositoryFactory(log)
	prlDTO := dto.ProjectRecruitmentLineDTOFactory(log)
	talRepo := repository.TemplateActivityLineRepositoryFactory(log)
	prhRepo := repository.ProjectRecruitmentHeaderRepositoryFactory(log)
	picRepo := repository.ProjectPicRepositoryFactory(log)
	prhDTO := dto.ProjectRecruitmentHeaderDTOFactory(log)
	return NewProjectRecruitmentLineUseCase(log, repo, prlDTO, talRepo, prhRepo, picRepo, prhDTO)
}

func (uc *ProjectRecruitmentLineUseCase) CreateOrUpdateProjectRecruitmentLines(req *request.CreateOrUpdateProjectRecruitmentLinesRequest) (*response.ProjectRecruitmentHeaderResponse, error) {
	// check if project recruitment header exist
	prh, err := uc.ProjectRecruitmentHeaderRepository.FindByID(uuid.MustParse(req.ProjectRecruitmentHeaderID))
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when finding project recruitment header by id: %s", err.Error())
		return nil, err
	}

	if prh == nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] project recruitment header with id %s not found", req.ProjectRecruitmentHeaderID)
		return nil, errors.New("project recruitment header not found")
	}

	// create or update project recruitment lines
	for _, line := range req.ProjectRecruitmentLines {
		if line.ID != "" && line.ID != uuid.Nil.String() {
			exist, err := uc.Repository.FindByID(uuid.MustParse(line.ID))
			if err != nil {
				uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when finding project recruitment line by id: %s", err.Error())
				return nil, err
			}

			parsedStartDate, err := time.Parse("2006-01-02", line.StartDate)
			if err != nil {
				uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when parsing start date: %s", err.Error())
				return nil, err
			}
			parsedEndDate, err := time.Parse("2006-01-02", line.EndDate)
			if err != nil {
				uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when parsing end date: %s", err.Error())
				return nil, err
			}

			if exist == nil {
				createdData, err := uc.Repository.CreateProjectRecruitmentLine(&entity.ProjectRecruitmentLine{
					TemplateActivityLineID:     uuid.MustParse(line.TemplateActivityLineID),
					ProjectRecruitmentHeaderID: prh.ID,
					StartDate:                  parsedStartDate,
					EndDate:                    parsedEndDate,
				})
				if err != nil {
					uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when creating project recruitment line: %s", err.Error())
					return nil, err
				}

				if len(line.ProjectPics) > 0 {
					for _, pic := range line.ProjectPics {
						var empID uuid.UUID
						if pic.EmployeeID != "" {
							empID = uuid.MustParse(pic.EmployeeID)
						}
						_, err := uc.ProjectPicRepository.CreateProjectPic(&entity.ProjectPic{
							EmployeeID:               &empID,
							AdministrativeTotal:      pic.AdministrativeTotal,
							ProjectRecruitmentLineID: createdData.ID,
						})
						if err != nil {
							uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when creating project pic: %s", err.Error())
							return nil, err
						}
					}
				}
			} else {
				updatedData, err := uc.Repository.UpdateProjectRecruitmentLine(&entity.ProjectRecruitmentLine{
					ID:                         exist.ID,
					TemplateActivityLineID:     uuid.MustParse(line.TemplateActivityLineID),
					ProjectRecruitmentHeaderID: prh.ID,
					StartDate:                  parsedStartDate,
					EndDate:                    parsedEndDate,
				})
				if err != nil {
					uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when updating project recruitment line: %s", err.Error())
					return nil, err
				}

				// delete project pics
				err = uc.ProjectPicRepository.DeleteProjectPicByProjectRecruitmentLineID(updatedData.ID)
				if err != nil {
					uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when deleting project pics: %s", err.Error())
					return nil, err
				}

				// create project pics
				if len(line.ProjectPics) > 0 {
					for _, pic := range line.ProjectPics {
						var empID uuid.UUID
						if pic.EmployeeID != "" {
							empID = uuid.MustParse(pic.EmployeeID)
						}
						_, err := uc.ProjectPicRepository.CreateProjectPic(&entity.ProjectPic{
							EmployeeID:               &empID,
							AdministrativeTotal:      pic.AdministrativeTotal,
							ProjectRecruitmentLineID: updatedData.ID,
						})
						if err != nil {
							uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when creating project pic: %s", err.Error())
							return nil, err
						}
					}
				}
			}
		} else {
			parsedStartDate, err := time.Parse("2006-01-02", line.StartDate)
			if err != nil {
				uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when parsing start date: %s", err.Error())
				return nil, err
			}
			parsedEndDate, err := time.Parse("2006-01-02", line.EndDate)
			if err != nil {
				uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when parsing end date: %s", err.Error())
				return nil, err
			}
			createdData, err := uc.Repository.CreateProjectRecruitmentLine(&entity.ProjectRecruitmentLine{
				TemplateActivityLineID:     uuid.MustParse(line.TemplateActivityLineID),
				ProjectRecruitmentHeaderID: prh.ID,
				StartDate:                  parsedStartDate,
				EndDate:                    parsedEndDate,
			})
			if err != nil {
				uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when creating project recruitment line: %s", err.Error())
				return nil, err
			}

			if len(line.ProjectPics) > 0 {
				for _, pic := range line.ProjectPics {
					var empID uuid.UUID
					if pic.EmployeeID != "" {
						empID = uuid.MustParse(pic.EmployeeID)
					}
					_, err := uc.ProjectPicRepository.CreateProjectPic(&entity.ProjectPic{
						EmployeeID:               &empID,
						AdministrativeTotal:      pic.AdministrativeTotal,
						ProjectRecruitmentLineID: createdData.ID,
					})
					if err != nil {
						uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when creating project pic: %s", err.Error())
						return nil, err
					}
				}
			}
		}
	}

	// delete project recruitment lines
	if len(req.DeletedProjectRecruitmentLineIDs) > 0 {
		for _, id := range req.DeletedProjectRecruitmentLineIDs {
			err := uc.Repository.DeleteProjectRecruitmentLine(uuid.MustParse(id))
			if err != nil {
				uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when deleting project recruitment line: %s", err.Error())
				return nil, err
			}
		}
	}

	// get project recruitment header response
	prhExist, err := uc.ProjectRecruitmentHeaderRepository.FindByID(prh.ID)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] error when finding project recruitment header by id: %s", err.Error())
		return nil, err
	}

	if prhExist == nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] project recruitment header with id %s not found", req.ProjectRecruitmentHeaderID)
		return nil, errors.New("project recruitment header not found")
	}

	return uc.ProjectRecruitmentHeaderDTO.ConvertEntityToResponse(prhExist), nil
}
