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

type IUserProfileUseCase interface {
	FillUserProfile(req *request.FillUserProfileRequest, userID uuid.UUID) (*response.UserProfileResponse, error)
}

type UserProfileUseCase struct {
	Log                      *logrus.Logger
	Repository               repository.IUserProfileRepository
	DTO                      dto.IUserProfileDTO
	WorkExperienceRepository repository.IWorkExperienceRepository
	EducationRepository      repository.IEducationRepository
	SkillRepository          repository.ISkillRepository
}

func NewUserProfileUseCase(
	log *logrus.Logger,
	repo repository.IUserProfileRepository,
	uDTO dto.IUserProfileDTO,
	weRepository repository.IWorkExperienceRepository,
	edRepository repository.IEducationRepository,
	sRepository repository.ISkillRepository,
) IUserProfileUseCase {
	return &UserProfileUseCase{
		Log:                      log,
		Repository:               repo,
		DTO:                      uDTO,
		WorkExperienceRepository: weRepository,
		EducationRepository:      edRepository,
		SkillRepository:          sRepository,
	}
}

func UserProfileUseCaseFactory(log *logrus.Logger) IUserProfileUseCase {
	repo := repository.UserProfileRepositoryFactory(log)
	uDTO := dto.UserProfileDTOFactory(log)
	weRepository := repository.WorkExperienceRepositoryFactory(log)
	edRepository := repository.EducationRepositoryFactory(log)
	sRepository := repository.SkillRepositoryFactory(log)
	return NewUserProfileUseCase(log, repo, uDTO, weRepository, edRepository, sRepository)
}

func (uc *UserProfileUseCase) FillUserProfile(req *request.FillUserProfileRequest, userID uuid.UUID) (*response.UserProfileResponse, error) {
	parsedBirthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing birth date: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing birth date: " + err.Error())
	}
	if req.ID == "" || req.ID == uuid.Nil.String() {
		createdProfile, err := uc.Repository.CreateUserProfile(&entity.UserProfile{
			UserID:          &userID,
			MaritalStatus:   entity.MaritalStatusEnum(req.MaritalStatus),
			Gender:          entity.UserGender(req.Gender),
			PhoneNumber:     req.PhoneNumber,
			Age:             req.Age,
			BirthDate:       parsedBirthDate,
			BirthPlace:      req.BirthPlace,
			Ktp:             req.KtpPath,
			CurriculumVitae: req.CvPath,
			Status:          entity.USER_INACTIVE,
		})
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating user profile: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating user profile: " + err.Error())
		}
		if len(req.WorkExperiences) > 0 {
			for _, we := range req.WorkExperiences {
				_, err := uc.WorkExperienceRepository.CreateWorkExperience(&entity.WorkExperience{
					UserProfileID:  createdProfile.ID,
					CompanyName:    we.CompanyName,
					Name:           we.Name,
					YearExperience: we.YearExperience,
					JobDescription: we.JobDescription,
					Certificate:    we.CertificatePath,
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating work experience: %s", err.Error())
					return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating work experience: " + err.Error())
				}
			}
		}
		if len(req.Educations) > 0 {
			for _, ed := range req.Educations {
				_, err := uc.EducationRepository.CreateEducation(&entity.Education{
					UserProfileID: createdProfile.ID,
					SchoolName:    ed.SchoolName,
					Major:         ed.Major,
					GraduateYear:  ed.GraduateYear,
					Certificate:   ed.CertificatePath,
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating education: %s", err.Error())
					return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating education: " + err.Error())
				}
			}
		}
		if len(req.Skills) > 0 {
			for _, s := range req.Skills {
				_, err := uc.SkillRepository.CreateSkill(&entity.Skill{
					UserProfileID: createdProfile.ID,
					Name:          s.Name,
					Description:   s.Description,
					Certificate:   s.CertificatePath,
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating skill: %s", err.Error())
					return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating skill: " + err.Error())
				}
			}
		}
	} else {
		parsedID, err := uuid.Parse(req.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when parsing ID: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when parsing ID: " + err.Error())
		}

		exist, err := uc.Repository.FindByID(parsedID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when finding user profile by ID: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when finding user profile by ID: " + err.Error())
		}
		if exist == nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] user profile not found")
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] user profile not found")
		}

		updatedProfile, err := uc.Repository.UpdateUserProfile(&entity.UserProfile{
			ID:              parsedID,
			UserID:          &userID,
			MaritalStatus:   entity.MaritalStatusEnum(req.MaritalStatus),
			Gender:          entity.UserGender(req.Gender),
			PhoneNumber:     req.PhoneNumber,
			Age:             req.Age,
			BirthDate:       parsedBirthDate,
			BirthPlace:      req.BirthPlace,
			Ktp:             req.KtpPath,
			CurriculumVitae: req.CvPath,
			Status:          entity.USER_INACTIVE,
		})
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when updating user profile: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when updating user profile: " + err.Error())
		}

		// delete work experiences
		err = uc.WorkExperienceRepository.DeleteByUserProfileID(updatedProfile.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when deleting work experiences: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when deleting work experiences: " + err.Error())
		}

		// delete educations
		err = uc.EducationRepository.DeleteByUserProfileID(updatedProfile.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when deleting educations: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when deleting educations: " + err.Error())
		}

		// delete skills
		err = uc.SkillRepository.DeleteByUserProfileID(updatedProfile.ID)
		if err != nil {
			uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when deleting skills: %s", err.Error())
			return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when deleting skills: " + err.Error())
		}

		if len(req.WorkExperiences) > 0 {
			for _, we := range req.WorkExperiences {
				_, err := uc.WorkExperienceRepository.CreateWorkExperience(&entity.WorkExperience{
					UserProfileID:  updatedProfile.ID,
					CompanyName:    we.CompanyName,
					Name:           we.Name,
					YearExperience: we.YearExperience,
					JobDescription: we.JobDescription,
					Certificate:    we.CertificatePath,
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating work experience: %s", err.Error())
					return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating work experience: " + err.Error())
				}
			}
		}

		if len(req.Educations) > 0 {
			for _, ed := range req.Educations {
				_, err := uc.EducationRepository.CreateEducation(&entity.Education{
					UserProfileID: updatedProfile.ID,
					SchoolName:    ed.SchoolName,
					Major:         ed.Major,
					GraduateYear:  ed.GraduateYear,
					Certificate:   ed.CertificatePath,
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating education: %s", err.Error())
					return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating education: " + err.Error())
				}
			}
		}

		if len(req.Skills) > 0 {
			for _, s := range req.Skills {
				_, err := uc.SkillRepository.CreateSkill(&entity.Skill{
					UserProfileID: updatedProfile.ID,
					Name:          s.Name,
					Description:   s.Description,
					Certificate:   s.CertificatePath,
				})
				if err != nil {
					uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when creating skill: %s", err.Error())
					return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when creating skill: " + err.Error())
				}
			}
		}
	}
	// find user profile by user ID
	profile, err := uc.Repository.FindByUserID(userID)
	if err != nil {
		uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] error when finding user profile by user ID: %s", err.Error())
		return nil, errors.New("[UserProfileUseCase.FillUserProfile] error when finding user profile by user ID: " + err.Error())
	}

	if profile == nil {
		uc.Log.Errorf("[UserProfileUseCase.FillUserProfile] user profile not found")
		return nil, errors.New("[UserProfileUseCase.FillUserProfile] user profile not found")
	}

	return uc.DTO.ConvertEntityToResponse(profile), nil
}
