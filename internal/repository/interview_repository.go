package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IInterviewRepository interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.Interview, int64, error)
	CreateInterview(interview *entity.Interview) (*entity.Interview, error)
	FindByID(id uuid.UUID) (*entity.Interview, error)
	FindByIDForMyself(id uuid.UUID, userProfile uuid.UUID) (*entity.Interview, error)
	FindByIDForMyselfAndAssessorFix(id uuid.UUID, userProfile uuid.UUID, interviewAssessorID uuid.UUID) (*entity.Interview, error)
	UpdateInterview(interview *entity.Interview) (*entity.Interview, error)
	DeleteInterview(id uuid.UUID) error
	GetHighestDocumentNumberByDate(date string) (int, error)
	FindByKeys(keys map[string]interface{}) (*entity.Interview, error)
	FindAllByKeys(keys map[string]interface{}) (*[]entity.Interview, error)
	FindByIDsForMyselfAssessor(ids []uuid.UUID, interviewAssessorID uuid.UUID) (*[]entity.Interview, error)
	FindByIDForAnswer(id, jobPostingID uuid.UUID) (*entity.Interview, error)
}

type InterviewRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewInterviewRepository(log *logrus.Logger, db *gorm.DB) *InterviewRepository {
	return &InterviewRepository{
		Log: log,
		DB:  db,
	}
}

func InterviewRepositoryFactory(log *logrus.Logger) IInterviewRepository {
	db := config.NewDatabase()
	return NewInterviewRepository(log, db)
}

func (r *InterviewRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.Interview, int64, error) {
	var interviews []entity.Interview
	var total int64

	query := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine")

	if search != "" {
		query = query.Where("document_number ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&interviews).Error; err != nil {
		r.Log.Error("[InterviewRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[InterviewRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &interviews, total, nil
}

func (r *InterviewRepository) CreateInterview(interview *entity.Interview) (*entity.Interview, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[InterviewRepository.CreateInterview] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(interview).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[InterviewRepository.CreateInterview] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Error("[InterviewRepository.CreateInterview] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine").First(interview).Error; err != nil {
		r.Log.Error("[InterviewRepository.CreateInterview] " + err.Error())
		return nil, err
	}

	return interview, nil
}

func (r *InterviewRepository) FindByID(id uuid.UUID) (*entity.Interview, error) {
	var interview entity.Interview

	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine").
		Preload("InterviewApplicants.UserProfile").Preload("InterviewApplicants.InterviewResults.InterviewAssessor").Preload("InterviewAssessors").Where("id = ?", id).First(&interview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[InterviewRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &interview, nil
}

func (r *InterviewRepository) FindByIDForMyself(id uuid.UUID, userProfile uuid.UUID) (*entity.Interview, error) {
	var interview entity.Interview

	if err := r.DB.Preload("JobPosting").
		Preload("ProjectPic").
		Preload("ProjectRecruitmentHeader").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.AnswerType").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionOptions").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionResponses", "user_profile_id = ?", userProfile).
		Preload("InterviewApplicants", "user_profile_id = ?", userProfile).
		Preload("InterviewApplicants.UserProfile").
		Preload("InterviewAssessors").
		Where("id = ?", id).First(&interview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[InterviewRepository.FindByIDForMyself] " + err.Error())
			return nil, err
		}
	}

	return &interview, nil
}

func (r *InterviewRepository) FindByIDForMyselfAndAssessorFix(id uuid.UUID, userProfile uuid.UUID, interviewAssessorID uuid.UUID) (*entity.Interview, error) {
	var interview entity.Interview
	if err := r.DB.Preload("JobPosting").
		Preload("ProjectPic").
		Preload("ProjectRecruitmentHeader").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.AnswerType").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionOptions").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionResponses", "user_profile_id = ? AND interview_assessor_id = ?", userProfile, interviewAssessorID).
		Preload("InterviewApplicants", "user_profile_id = ?", userProfile).
		Preload("InterviewApplicants.UserProfile").
		Preload("InterviewAssessors").
		Where("id = ?", id).First(&interview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[InterviewRepository.FindByIDForMyself] " + err.Error())
			return nil, err
		}
	}

	return &interview, nil
}

func (r *InterviewRepository) FindByIDsForMyselfAssessor(ids []uuid.UUID, interviewAssessorID uuid.UUID) (*[]entity.Interview, error) {
	var interviews []entity.Interview

	if err := r.DB.Preload("JobPosting").
		Preload("ProjectPic").
		Preload("ProjectRecruitmentHeader").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.AnswerType").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionOptions").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionResponses", "interview_assessor_id = ?", interviewAssessorID).
		Preload("InterviewApplicants.UserProfile").
		Preload("InterviewAssessors", "id = ?", interviewAssessorID).
		Where("id IN ?", ids).Find(&interviews).Error; err != nil {
		r.Log.Error("[InterviewRepository.FindByIDsForMyselfAssessor] " + err.Error())
		return nil, err
	}

	return &interviews, nil
}

func (r *InterviewRepository) UpdateInterview(interview *entity.Interview) (*entity.Interview, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[InterviewRepository.UpdateInterview] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.Interview{}).Where("id = ?", interview.ID).Updates(interview).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[InterviewRepository.UpdateInterview] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Error("[InterviewRepository.UpdateInterview] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine").First(interview).Error; err != nil {
		r.Log.Error("[InterviewRepository.UpdateInterview] " + err.Error())
		return nil, err
	}

	return interview, nil
}

func (r *InterviewRepository) DeleteInterview(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[InterviewRepository.DeleteInterview] " + tx.Error.Error())
		return tx.Error
	}

	var interview entity.Interview
	if err := tx.Where("id = ?", id).First(&interview).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&interview).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[InterviewRepository.DeleteInterview] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Error("[InterviewRepository.DeleteInterview] " + err.Error())
		return err
	}

	return nil
}

func (r *InterviewRepository) GetHighestDocumentNumberByDate(date string) (int, error) {
	var maxNumber int
	err := r.DB.Raw(`
			SELECT COALESCE(MAX(CAST(SUBSTRING(document_number FROM '[0-9]+$') AS INTEGER)), 0)
			FROM interviews
			WHERE DATE(created_at) = ?
	`, date).Scan(&maxNumber).Error
	if err != nil {
		r.Log.Errorf("[TestScheduleHeaderRepository.GetHighestDocumentNumberByDate] error when querying max document number: %v", err)
		return 0, err
	}
	return maxNumber, nil
}

func (r *InterviewRepository) FindByKeys(keys map[string]interface{}) (*entity.Interview, error) {
	var interview entity.Interview
	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine").
		Preload("InterviewApplicants.UserProfile").Preload("InterviewApplicants.InterviewResults.InterviewAssessor").Preload("InterviewAssessors").Where(keys).First(&interview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[InterviewRepository.FindByKeys] " + err.Error())
			return nil, err
		}
	}
	return &interview, nil
}

func (r *InterviewRepository) FindAllByKeys(keys map[string]interface{}) (*[]entity.Interview, error) {
	var interviews []entity.Interview
	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine").
		Preload("InterviewApplicants.UserProfile").Preload("InterviewApplicants.InterviewResults.InterviewAssessor").Preload("InterviewAssessors").Where(keys).Find(&interviews).Error; err != nil {
		r.Log.Error("[InterviewRepository.FindAllByKeys] " + err.Error())
		return nil, err
	}
	return &interviews, nil
}

func (r *InterviewRepository) FindByIDForAnswer(id, jobPostingID uuid.UUID) (*entity.Interview, error) {
	var interview entity.Interview
	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine").
		Preload("InterviewApplicants.UserProfile").Preload("InterviewApplicants.InterviewResults.InterviewAssessor").Preload("InterviewAssessors").
		Where("id = ? AND job_posting_id = ?", id, jobPostingID).First(&interview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[InterviewRepository.FindByIDForAnswer] " + err.Error())
			return nil, err
		}
	}
	return &interview, nil
}
