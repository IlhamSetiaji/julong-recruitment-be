package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IRecruitmentTypeUseCase interface {
	FindAll() ([]*response.RecruitmentTypeResponse, error)
}

type RecruitmentTypeUseCase struct {
	Log *logrus.Logger
}

func NewRecruitmentTypeUseCase(log *logrus.Logger) IRecruitmentTypeUseCase {
	return &RecruitmentTypeUseCase{
		Log: log,
	}
}

func RecruitmentTypeUseCaseFactory(log *logrus.Logger) IRecruitmentTypeUseCase {
	return NewRecruitmentTypeUseCase(log)
}

func (uc *RecruitmentTypeUseCase) FindAll() ([]*response.RecruitmentTypeResponse, error) {
	recruitmentTypes := []*response.RecruitmentTypeResponse{
		{Key: "PROJECT_RECRUITMENT_TYPE_MT", Value: string(entity.PROJECT_RECRUITMENT_TYPE_MT)},
		{Key: "PROJECT_RECRUITMENT_TYPE_PH", Value: string(entity.PROJECT_RECRUITMENT_TYPE_PH)},
		{Key: "PROJECT_RECRUITMENT_TYPE_NS", Value: string(entity.PROJECT_RECRUITMENT_TYPE_NS)},
	}

	return recruitmentTypes, nil
}
