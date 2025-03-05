package usecase

import (
	"errors"
	"sort"
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
	GetAllByKeys(keys map[string]interface{}) ([]*response.ProjectRecruitmentLineResponse, error)
	GetAllByKeysWithoutPic(keys map[string]interface{}) ([]*response.ProjectRecruitmentLineResponse, error)
	FindAllByFormType(formType entity.TemplateQuestionFormType) ([]*response.ProjectRecruitmentLineResponse, error)
	FindAllByProjectRecruitmentHeaderIDAndEmployeeID(projectRecruitmentHeaderID, employeeID uuid.UUID) ([]*response.ProjectRecruitmentLineResponse, error)
	FindByIDForAnswer(id, jobPostingID, userProfileID uuid.UUID) (*response.ProjectRecruitmentLineResponse, error)
	FindByIDForAnswerInterview(id, jobPostingID, userProfileID, interviewAssessorID uuid.UUID) (*response.ProjectRecruitmentLineResponse, error)
	FindByIDForAnswerFgd(id, jobPostingID, userProfileID, fgdAssessorID uuid.UUID) (*response.ProjectRecruitmentLineResponse, error)
	FindAllByHeaderID(headerID uuid.UUID) (*[]response.ProjectRecruitmentLineResponse, error)
	FindAllByHeaderIDAndFormType(headerID uuid.UUID, formType entity.TemplateQuestionFormType) ([]*response.ProjectRecruitmentLineResponse, error)
	FindAllByMonthAndYear(month, year int, employeeID uuid.UUID) ([]*response.ProjectRecruitmentLineResponse, error)
}

type ProjectRecruitmentLineUseCase struct {
	Log                                *logrus.Logger
	Repository                         repository.IProjectRecruitmentLineRepository
	DTO                                dto.IProjectRecruitmentLineDTO
	TemplateActivityLineRepository     repository.ITemplateActivityLineRepository
	ProjectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository
	ProjectPicRepository               repository.IProjectPicRepository
	ProjectRecruitmentHeaderDTO        dto.IProjectRecruitmentHeaderDTO
	TemplateQuestionRepository         repository.ITemplateQuestionRepository
}

func NewProjectRecruitmentLineUseCase(
	log *logrus.Logger,
	repo repository.IProjectRecruitmentLineRepository,
	dto dto.IProjectRecruitmentLineDTO,
	talRepo repository.ITemplateActivityLineRepository,
	prhRepo repository.IProjectRecruitmentHeaderRepository,
	picRepo repository.IProjectPicRepository,
	prhDTO dto.IProjectRecruitmentHeaderDTO,
	tqRepo repository.ITemplateQuestionRepository,
) IProjectRecruitmentLineUseCase {
	return &ProjectRecruitmentLineUseCase{
		Log:                                log,
		Repository:                         repo,
		DTO:                                dto,
		TemplateActivityLineRepository:     talRepo,
		ProjectRecruitmentHeaderRepository: prhRepo,
		ProjectPicRepository:               picRepo,
		ProjectRecruitmentHeaderDTO:        prhDTO,
		TemplateQuestionRepository:         tqRepo,
	}
}

func ProjectRecruitmentLineUseCaseFactory(log *logrus.Logger) IProjectRecruitmentLineUseCase {
	repo := repository.ProjectRecruitmentLineRepositoryFactory(log)
	prlDTO := dto.ProjectRecruitmentLineDTOFactory(log)
	talRepo := repository.TemplateActivityLineRepositoryFactory(log)
	prhRepo := repository.ProjectRecruitmentHeaderRepositoryFactory(log)
	picRepo := repository.ProjectPicRepositoryFactory(log)
	prhDTO := dto.ProjectRecruitmentHeaderDTOFactory(log)
	tqRepo := repository.TemplateQuestionRepositoryFactory(log)
	return NewProjectRecruitmentLineUseCase(log, repo, prlDTO, talRepo, prhRepo, picRepo, prhDTO, tqRepo)
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

			if parsedStartDate.After(parsedEndDate) {
				uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] start date is after end date")
				return nil, errors.New("project recruitment line error: start date is after end date")
			}

			if parsedStartDate.Before(prh.StartDate) || parsedEndDate.After(prh.EndDate) {
				uc.Log.Errorf("[ProjectRecruitmentLineUseCase.CreateOrUpdateProjectRecruitmentLines] start date or end date is out of range")
				return nil, errors.New("project recruitment line error: start date or end date is out of range of project recruitment header")
			}

			if exist == nil {
				createdData, err := uc.Repository.CreateProjectRecruitmentLine(&entity.ProjectRecruitmentLine{
					TemplateActivityLineID:     uuid.MustParse(line.TemplateActivityLineID),
					ProjectRecruitmentHeaderID: prh.ID,
					StartDate:                  parsedStartDate,
					EndDate:                    parsedEndDate,
					Order:                      line.Order,
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
					Order:                      line.Order,
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
				Order:                      line.Order,
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

func (uc *ProjectRecruitmentLineUseCase) GetAllByKeys(keys map[string]interface{}) ([]*response.ProjectRecruitmentLineResponse, error) {
	data, err := uc.Repository.GetAllByKeys(keys)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.GetAllByKeys] error when getting project recruitment lines: %s", err.Error())
		return nil, err
	}

	responses := make([]*response.ProjectRecruitmentLineResponse, 0)
	for _, d := range data {
		responses = append(responses, uc.DTO.ConvertEntityToResponse(&d))
	}

	return responses, nil
}

func (uc *ProjectRecruitmentLineUseCase) GetAllByKeysWithoutPic(keys map[string]interface{}) ([]*response.ProjectRecruitmentLineResponse, error) {
	data, err := uc.Repository.GetAllByKeysWithoutProjectPic(keys)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.GetAllByKeys] error when getting project recruitment lines: %s", err.Error())
		return nil, err
	}

	responses := make([]*response.ProjectRecruitmentLineResponse, 0)
	for _, d := range data {
		responses = append(responses, uc.DTO.ConvertEntityToResponse(&d))
	}

	return responses, nil
}

func (uc *ProjectRecruitmentLineUseCase) generateOrder(prhID uuid.UUID) (int, error) {
	data, err := uc.Repository.GetAllByKeys(map[string]interface{}{
		"project_recruitment_header_id": prhID,
	})
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.generateOrder] error when getting project recruitment lines: %s", err.Error())
		return 0, err
	}

	if len(data) == 0 {
		return 1, nil
	}

	maxOrder := 0
	for _, d := range data {
		if d.Order > maxOrder {
			maxOrder = d.Order
		}
	}

	return maxOrder + 1, nil
}

func (uc *ProjectRecruitmentLineUseCase) FindAllByFormType(formType entity.TemplateQuestionFormType) ([]*response.ProjectRecruitmentLineResponse, error) {
	tQuestions, err := uc.TemplateQuestionRepository.FindAllByFormType(formType)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindAllByFormType] error when finding template questions by form type: %s", err.Error())
		return nil, err
	}

	questionIDs := make([]uuid.UUID, 0)
	for _, tq := range *tQuestions {
		questionIDs = append(questionIDs, tq.ID)
	}

	tActivityLines, err := uc.TemplateActivityLineRepository.FindAllByTemplateQuestionIDs(questionIDs)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindAllByFormType] error when finding template activity lines by template question ids: %s", err.Error())
		return nil, err
	}

	tActivityLineIDs := make([]uuid.UUID, 0)
	for _, tal := range *tActivityLines {
		tActivityLineIDs = append(tActivityLineIDs, tal.ID)
	}

	data, err := uc.Repository.FindAllByTemplateActivityLineIDs(tActivityLineIDs)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindAllByFormType] error when finding project recruitment lines by template activity line ids: %s", err.Error())
		return nil, err
	}

	responses := make([]*response.ProjectRecruitmentLineResponse, 0)
	for _, d := range *data {
		responses = append(responses, uc.DTO.ConvertEntityToResponse(&d))
	}

	return responses, nil
}

func (uc *ProjectRecruitmentLineUseCase) FindAllByProjectRecruitmentHeaderIDAndEmployeeID(projectRecruitmentHeaderID, employeeID uuid.UUID) ([]*response.ProjectRecruitmentLineResponse, error) {
	pics, err := uc.ProjectPicRepository.FindAllByEmployeeID(employeeID)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindAllByProjectRecruitmentHeaderIDAndEmployeeID] error when finding project pics by employee id: %s", err.Error())
		return nil, err
	}

	projectRecruitmentHeader, err := uc.ProjectRecruitmentHeaderRepository.FindByID(projectRecruitmentHeaderID)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindAllByProjectRecruitmentHeaderIDAndEmployeeID] error when finding project recruitment header by id: %s", err.Error())
		return nil, err
	}

	if projectRecruitmentHeader == nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindAllByProjectRecruitmentHeaderIDAndEmployeeID] project recruitment header with id %s not found", projectRecruitmentHeaderID)
		return nil, errors.New("project recruitment header not found")
	}

	projectRecruitmentLineIDs := make([]uuid.UUID, 0)
	for _, pic := range pics {
		projectRecruitmentLineIDs = append(projectRecruitmentLineIDs, pic.ProjectRecruitmentLineID)
	}

	data, err := uc.Repository.FindAllByProjectRecruitmentHeaderIdAndIds(projectRecruitmentHeaderID, projectRecruitmentLineIDs)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindAllByProjectRecruitmentHeaderIDAndEmployeeID] error when finding project recruitment lines by project recruitment header id and ids: %s", err.Error())
		return nil, err
	}

	responses := make([]*response.ProjectRecruitmentLineResponse, 0)
	for _, d := range *data {
		responses = append(responses, uc.DTO.ConvertEntityToResponse(&d))
	}

	return responses, nil
}

func (uc *ProjectRecruitmentLineUseCase) FindByIDForAnswer(id, jobPostingID, userProfileID uuid.UUID) (*response.ProjectRecruitmentLineResponse, error) {
	data, err := uc.Repository.FindByIDForAnswer(id, jobPostingID, userProfileID)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindByIDForAnswer] error when finding project recruitment line by id for answer: %s", err.Error())
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	return uc.DTO.ConvertEntityToResponse(data), nil
}

func (uc *ProjectRecruitmentLineUseCase) FindByIDForAnswerInterview(id, jobPostingID, userProfileID, interviewAssessorID uuid.UUID) (*response.ProjectRecruitmentLineResponse, error) {
	data, err := uc.Repository.FindByIDForAnswerInterview(id, jobPostingID, userProfileID, interviewAssessorID)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindByIDForAnswerInterview] error when finding project recruitment line by id for answer interview: %s", err.Error())
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	return uc.DTO.ConvertEntityToResponse(data), nil
}

func (uc *ProjectRecruitmentLineUseCase) FindByIDForAnswerFgd(id, jobPostingID, userProfileID, fgdAssessorID uuid.UUID) (*response.ProjectRecruitmentLineResponse, error) {
	data, err := uc.Repository.FindByIDForAnswerFgd(id, jobPostingID, userProfileID, fgdAssessorID)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindByIDForAnswerFgd] error when finding project recruitment line by id for answer fgd: %s", err.Error())
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	return uc.DTO.ConvertEntityToResponse(data), nil
}

func (uc *ProjectRecruitmentLineUseCase) FindAllByHeaderID(headerID uuid.UUID) (*[]response.ProjectRecruitmentLineResponse, error) {
	data, err := uc.Repository.GetAllByKeys(map[string]interface{}{
		"project_recruitment_header_id": headerID,
	})
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindAllByHeaderID] error when finding project recruitment lines by header id: %s", err.Error())
		return nil, err
	}

	responses := make([]response.ProjectRecruitmentLineResponse, 0)
	for _, d := range data {
		responses = append(responses, *uc.DTO.ConvertEntityToResponse(&d))
	}

	return &responses, nil
}

func (uc *ProjectRecruitmentLineUseCase) FindAllByHeaderIDAndFormType(headerID uuid.UUID, formType entity.TemplateQuestionFormType) ([]*response.ProjectRecruitmentLineResponse, error) {
	templateQuestions, err := uc.TemplateQuestionRepository.FindAllByFormType(formType)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindAllByHeaderIDAndFormType] error when finding template questions by form type: %s", err.Error())
		return nil, err
	}

	var projectRecruitmentLines []*entity.ProjectRecruitmentLine
	for _, tq := range *templateQuestions {
		for _, tal := range tq.TemplateActivityLines {
			for _, prl := range tal.ProjectRecruitmentLines {
				projectRecruitmentLines = append(projectRecruitmentLines, &prl)
			}
		}
	}

	// Sort the projectRecruitmentLines by Order field in ascending order
	sort.Slice(projectRecruitmentLines, func(i, j int) bool {
		return projectRecruitmentLines[i].Order < projectRecruitmentLines[j].Order
	})

	responses := make([]*response.ProjectRecruitmentLineResponse, 0)
	for _, prl := range projectRecruitmentLines {
		if prl.ProjectRecruitmentHeaderID == headerID {
			responses = append(responses, uc.DTO.ConvertEntityToResponse(prl))
		}
	}

	return responses, nil
}

func (uc *ProjectRecruitmentLineUseCase) fetchDataForDay(date time.Time, employeeID uuid.UUID) ([]*response.ProjectRecruitmentLineResponse, error) {
	projectRecruitmentLines, err := uc.Repository.FindAllByStartDate(date, employeeID)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentLineUseCase.fetchDataForDay] error when finding project recruitment lines by start date: %s", err.Error())
		return nil, err
	}
	var responses []*response.ProjectRecruitmentLineResponse
	for _, prl := range *projectRecruitmentLines {
		responses = append(responses, uc.DTO.ConvertEntityToResponse(&prl))
	}
	return responses, nil
}

func (uc *ProjectRecruitmentLineUseCase) FindAllByMonthAndYear(month, year int, employeeID uuid.UUID) ([]*response.ProjectRecruitmentLineResponse, error) {
	firstDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDate := firstDate.AddDate(0, 1, -1)
	var responses []*response.ProjectRecruitmentLineResponse

	for current := firstDate; !current.After(lastDate); current = current.AddDate(0, 0, 1) {
		resp, err := uc.fetchDataForDay(current, employeeID)
		if err != nil {
			uc.Log.Errorf("[ProjectRecruitmentLineUseCase.FindAllByMonthAndYear] error when fetching data for day: %s", err.Error())
			return nil, err
		}
		responses = append(responses, resp...)
	}

	return responses, nil
}
