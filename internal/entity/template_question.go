package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateQuestionStatus string

const (
	TEMPLATE_QUESTION_STATUS_ACTIVE   TemplateQuestionStatus = "ACTIVE"
	TEMPLATE_QUESTION_STATUS_INACTIVE TemplateQuestionStatus = "INACTIVE"
)

type TemplateQuestionFormType string

const (
	TQ_FORM_TYPE_ADMINISTRATIVE_SELECTION TemplateQuestionFormType = "ADMINISTRATIVE_SELECTION"
	TQ_FORM_TYPE_TEST                     TemplateQuestionFormType = "TEST"
	TQ_FORM_TYPE_INTERVIEW                TemplateQuestionFormType = "INTERVIEW"
	TQ_FORM_TYPE_FGD                      TemplateQuestionFormType = "FGD"
	TQ_FORM_TYPE_FINAL_INTERVIEW          TemplateQuestionFormType = "FINAL_INTERVIEW"
	TQ_FORM_TYPE_OFFERING_LETTER          TemplateQuestionFormType = "OFFERING_LETTER"
	TQ_FORM_TYPE_CONTRACT_DOCUMENT        TemplateQuestionFormType = "CONTRACT_DOCUMENT"
	TQ_FORM_TYPE_DOCUMENT_CHECKING        TemplateQuestionFormType = "DOCUMENT_CHECKING"
	TQ_FORM_TYPE_SURAT_IZIN_ORTU          TemplateQuestionFormType = "SURAT_IZIN_ORANG_TUA"
	TQ_FORM_TYPE_FINAL_RESULT          	  TemplateQuestionFormType = "FINAL_RESULT"
	TQ_FORM_TYPE_KARYAWAN_TETAP           TemplateQuestionFormType = "KARYAWAN_TETAP"
	TQ_FORM_TYPE_PKWT                     TemplateQuestionFormType = "PKWT"
	TQ_FORM_TYPE_PKWTT                    TemplateQuestionFormType = "PKWTT"
)

type TemplateQuestion struct {
	gorm.Model      `json:"-"`
	ID              uuid.UUID              `json:"id" gorm:"type:char(36);primaryKey;"`
	DocumentSetupID *uuid.UUID             `json:"document_setup_id" gorm:"type:char(36);default:null"`
	Name            string                 `json:"name" gorm:"type:varchar(255);not null"`
	FormType        string                 `json:"form_type" gorm:"type:varchar(255);default:null"`
	Description     string                 `json:"description" gorm:"type:text;default:null"`
	Duration        *int                   `json:"duration" gorm:"default:0"`
	Status          TemplateQuestionStatus `json:"status" gorm:"default:'ACTIVE'"`

	Questions             []Question             `json:"questions" gorm:"foreignKey:TemplateQuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DocumentSetup         *DocumentSetup         `json:"document_setup" gorm:"foreignKey:DocumentSetupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DocumentVerifications []DocumentVerification `json:"document_verifications" gorm:"foreignKey:TemplateQuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Applicants            []Applicant            `json:"applicants" gorm:"foreignKey:TemplateQuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TemplateActivityLines []TemplateActivityLine `json:"template_activity_lines" gorm:"foreignKey:QuestionTemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (tq *TemplateQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	tq.ID = uuid.New()
	tq.CreatedAt = time.Now()
	tq.UpdatedAt = time.Now()
	return
}

func (tq *TemplateQuestion) BeforeUpdate(tx *gorm.DB) (err error) {
	tq.UpdatedAt = time.Now()
	return
}

func (TemplateQuestion) TableName() string {
	return "template_questions"
}

func GetAllFormTypes() []TemplateQuestionFormType {
	return []TemplateQuestionFormType{
		TQ_FORM_TYPE_ADMINISTRATIVE_SELECTION,
		TQ_FORM_TYPE_TEST,
		TQ_FORM_TYPE_INTERVIEW,
		TQ_FORM_TYPE_FGD,
		TQ_FORM_TYPE_FINAL_INTERVIEW,
		TQ_FORM_TYPE_OFFERING_LETTER,
		TQ_FORM_TYPE_CONTRACT_DOCUMENT,
		TQ_FORM_TYPE_DOCUMENT_CHECKING,
		TQ_FORM_TYPE_FINAL_RESULT,
	}
}
