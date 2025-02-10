package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IFgdScheduleRepository interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.FgdSchedule, int64, error)
	CreateFgdSchedule(fgdSchedule *entity.FgdSchedule) (*entity.FgdSchedule, error)
	FindByID(id uuid.UUID) (*entity.FgdSchedule, error)
	FindByIDForMyself(id uuid.UUID, userProfile uuid.UUID) (*entity.FgdSchedule, error)
	FindByIDForMyselfAndAssessor(id uuid.UUID, userProfile uuid.UUID, fgdScheduleAssessorID uuid.UUID) (*entity.FgdSchedule, error)
	UpdateFgdSchedule(fgdSchedule *entity.FgdSchedule) (*entity.FgdSchedule, error)
	DeleteFgdSchedule(id uuid.UUID) error
	GetHighestDocumentNumberByDate(date string) (int, error)
	FindByKeys(keys map[string]interface{}) (*entity.FgdSchedule, error)
	FindAllByKeys(keys map[string]interface{}) (*[]entity.FgdSchedule, error)
	FindByIDForMyselfAssessor(id uuid.UUID, fgdScheduleAssessorID uuid.UUID) (*entity.FgdSchedule, error)
	FindByIDsForMyselfAssessor(ids []uuid.UUID, fgdScheduleAssessorID uuid.UUID) (*[]entity.FgdSchedule, error)
	FindByIDForAnswer(id, jobPostingID uuid.UUID) (*entity.FgdSchedule, error)
}

type FgdScheduleRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewFgdScheduleRepository(log *logrus.Logger, db *gorm.DB) *FgdScheduleRepository {
	return &FgdScheduleRepository{
		Log: log,
		DB:  db,
	}
}

func FgdScheduleRepositoryFactory(log *logrus.Logger) IFgdScheduleRepository {
	db := config.NewDatabase()
	return NewFgdScheduleRepository(log, db)
}

func (r *FgdScheduleRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.FgdSchedule, int64, error) {
	var fgdSchedules []entity.FgdSchedule
	var total int64

	query := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine")

	if search != "" {
		query = query.Where("document_number ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&fgdSchedules).Error; err != nil {
		r.Log.Error("[FgdScheduleRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[FgdScheduleRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &fgdSchedules, total, nil
}

func (r *FgdScheduleRepository) CreateFgdSchedule(fgdSchedule *entity.FgdSchedule) (*entity.FgdSchedule, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[FgdScheduleRepository.CreateFgdSchedule] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(fgdSchedule).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[FgdScheduleRepository.CreateFgdSchedule] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Error("[FgdScheduleRepository.CreateFgdSchedule] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine").First(fgdSchedule).Error; err != nil {
		r.Log.Error("[FgdScheduleRepository.CreateFgdSchedule] " + err.Error())
		return nil, err
	}

	return fgdSchedule, nil
}

func (r *FgdScheduleRepository) FindByID(id uuid.UUID) (*entity.FgdSchedule, error) {
	var fgdSchedule entity.FgdSchedule

	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine").
		Preload("InterviewApplicants.UserProfile").Preload("InterviewApplicants.InterviewResults.InterviewAssessor").Preload("InterviewAssessors").Where("id = ?", id).First(&fgdSchedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[FgdScheduleRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &fgdSchedule, nil
}

func (r *FgdScheduleRepository) FindByIDForMyself(id uuid.UUID, userProfile uuid.UUID) (*entity.FgdSchedule, error) {
	var fgdSchedule entity.FgdSchedule

	if err := r.DB.Preload("JobPosting").
		Preload("ProjectPic").
		Preload("ProjectRecruitmentHeader").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.AnswerType").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionOptions").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionResponses", "user_profile_id = ?", userProfile).
		Preload("FgdApplicants", "user_profile_id = ?", userProfile).
		Preload("FgdApplicants.UserProfile").
		Preload("FgdAssessors").
		Where("id = ?", id).First(&fgdSchedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[FgdScheduleRepository.FindByIDForMyself] " + err.Error())
			return nil, err
		}
	}

	return &fgdSchedule, nil
}

func (r *FgdScheduleRepository) FindByIDForMyselfAndAssessor(id uuid.UUID, userProfile uuid.UUID, fgdScheduleAssessorID uuid.UUID) (*entity.FgdSchedule, error) {
	var fgdSchedule entity.FgdSchedule

	if err := r.DB.Preload("JobPosting").
		Preload("ProjectPic").
		Preload("ProjectRecruitmentHeader").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.AnswerType").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionOptions").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionResponses", "user_profile_id = ? AND fgd_assessor_id = ?", userProfile, fgdScheduleAssessorID).
		Preload("FgdApplicants", "user_profile_id = ?", userProfile).
		Preload("FgdApplicants.UserProfile").
		Preload("FgdAssessors").
		Where("id = ?", id).First(&fgdSchedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[FgdScheduleRepository.FindByIDForMyself] " + err.Error())
			return nil, err
		}
	}

	return &fgdSchedule, nil
}

func (r *FgdScheduleRepository) FindByIDForMyselfAssessor(id uuid.UUID, fgdScheduleAssessorID uuid.UUID) (*entity.FgdSchedule, error) {
	var fgdSchedule entity.FgdSchedule

	if err := r.DB.Preload("JobPosting").
		Preload("ProjectPic").
		Preload("ProjectRecruitmentHeader").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.AnswerType").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionOptions").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionResponses", "fgd_assessor_id = ?", fgdScheduleAssessorID).
		Preload("FgdApplicants.UserProfile").
		Preload("FgdAssessors", "id = ?", fgdScheduleAssessorID).
		Where("id = ?", id).First(&fgdSchedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[FgdScheduleRepository.FindByIDForMyself] " + err.Error())
			return nil, err
		}
	}

	return &fgdSchedule, nil
}

func (r *FgdScheduleRepository) FindByIDsForMyselfAssessor(ids []uuid.UUID, fgdScheduleAssessorID uuid.UUID) (*[]entity.FgdSchedule, error) {
	var fgdSchedules []entity.FgdSchedule

	if err := r.DB.Preload("JobPosting").
		Preload("ProjectPic").
		Preload("ProjectRecruitmentHeader").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.AnswerType").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionOptions").
		Preload("ProjectRecruitmentLine.TemplateActivityLine.TemplateQuestion.Questions.QuestionResponses", "interview_assessor_id = ?", fgdScheduleAssessorID).
		Preload("FgdApplicants.UserProfile").
		Preload("FgdAssessors", "id = ?", fgdScheduleAssessorID).
		Where("id IN ?", ids).Find(&fgdSchedules).Error; err != nil {
		r.Log.Error("[FgdScheduleRepository.FindByIDsForMyselfAssessor] " + err.Error())
		return nil, err
	}

	return &fgdSchedules, nil
}

func (r *FgdScheduleRepository) UpdateFgdSchedule(fgdSchedule *entity.FgdSchedule) (*entity.FgdSchedule, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[FgdScheduleRepository.UpdateFgdSchedule] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.Interview{}).Where("id = ?", fgdSchedule.ID).Updates(fgdSchedule).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[FgdScheduleRepository.UpdateFgdSchedule] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Error("[FgdScheduleRepository.UpdateFgdSchedule] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine").First(fgdSchedule).Error; err != nil {
		r.Log.Error("[FgdScheduleRepository.UpdateFgdSchedule] " + err.Error())
		return nil, err
	}

	return fgdSchedule, nil
}

func (r *FgdScheduleRepository) DeleteFgdSchedule(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[FgdScheduleRepository.DeleteFgdSchedule] " + tx.Error.Error())
		return tx.Error
	}

	var fgdSchedule entity.FgdSchedule
	if err := tx.Where("id = ?", id).First(&fgdSchedule).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&fgdSchedule).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[FgdScheduleRepository.DeleteFgdSchedule] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Error("[FgdScheduleRepository.DeleteFgdSchedule] " + err.Error())
		return err
	}

	return nil
}

func (r *FgdScheduleRepository) GetHighestDocumentNumberByDate(date string) (int, error) {
	var maxNumber int
	err := r.DB.Raw(`
			SELECT COALESCE(MAX(CAST(SUBSTRING(document_number FROM '[0-9]+$') AS INTEGER)), 0)
			FROM fgd_schedules
			WHERE DATE(created_at) = ?
	`, date).Scan(&maxNumber).Error
	if err != nil {
		r.Log.Errorf("[FgdScheduleRepository.GetHighestDocumentNumberByDate] error when querying max document number: %v", err)
		return 0, err
	}
	return maxNumber, nil
}

func (r *FgdScheduleRepository) FindByKeys(keys map[string]interface{}) (*entity.FgdSchedule, error) {
	var fgdSchedule entity.FgdSchedule
	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine").
		Preload("FgdApplicants.UserProfile").Preload("FgdApplicants.FgdResults.FgdAssessor").Preload("FgdAssessors").Where(keys).First(&fgdSchedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[FgdScheduleRepository.FindByKeys] " + err.Error())
			return nil, err
		}
	}
	return &fgdSchedule, nil
}

func (r *FgdScheduleRepository) FindAllByKeys(keys map[string]interface{}) (*[]entity.FgdSchedule, error) {
	var fgdSchedules []entity.FgdSchedule
	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine").
		Preload("FgdApplicants.UserProfile").Preload("FgdApplicants.FgdResults.FgdAssessor").Preload("FgdAssessors").Where(keys).Find(&fgdSchedules).Error; err != nil {
		r.Log.Error("[FgdScheduleRepository.FindAllByKeys] " + err.Error())
		return nil, err
	}
	return &fgdSchedules, nil
}

func (r *FgdScheduleRepository) FindByIDForAnswer(id, jobPostingID uuid.UUID) (*entity.FgdSchedule, error) {
	var fgdSchedule entity.FgdSchedule
	if err := r.DB.Preload("JobPosting").Preload("ProjectPic").Preload("ProjectRecruitmentHeader").Preload("ProjectRecruitmentLine.TemplateActivityLine").
		Preload("Fgdpplicants.UserProfile").Preload("Fgdpplicants.InterviewResults.Fgdssessor").Preload("Fgdssessors").
		Where("id = ? AND job_posting_id = ?", id, jobPostingID).First(&fgdSchedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[FgdScheduleRepository.FindByIDForAnswer] " + err.Error())
			return nil, err
		}
	}
	return &fgdSchedule, nil
}
