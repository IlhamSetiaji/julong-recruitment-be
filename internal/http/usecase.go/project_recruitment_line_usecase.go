package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/sirupsen/logrus"
)

type IProjectRecruitmentLineUseCase interface {
	CreateOrUpdateProjectRecruitmentLines(req *request.CreateOrUpdateProjectRecruitmentLinesRequest) (*response.ProjectRecruitmentHeaderResponse, error)
}

type ProjectRecruitmentLineUseCase struct {
	Log        *logrus.Logger
	Repository repository.IProjectRecruitmentLineRepository
	DTO        dto.IProjectRecruitmentLineDTO
}
