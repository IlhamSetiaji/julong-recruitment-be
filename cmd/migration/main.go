package main

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)
	db := config.NewDatabase()

	// migrate the schema
	err := db.AutoMigrate(
		&entity.MPRequest{},
		&entity.ProjectRecruitmentHeader{},
		&entity.ProjectRecruitmentLine{},
		&entity.JobPosting{},
		&entity.ProjectPic{},
		&entity.AnswerType{},
		&entity.TemplateQuestion{},
		&entity.Question{},
		&entity.QuestionOption{},
		&entity.QuestionResponse{},
		&entity.UserProfile{},
		&entity.Applicant{},
		&entity.DocumentType{},
		&entity.MailTemplate{},
		&entity.DocumentSetup{},
		&entity.DocumentVerification{},
		&entity.TemplateActivity{},
		&entity.TemplateActivityLine{},
		&entity.University{},
		&entity.WorkExperience{},
		&entity.Education{},
		&entity.Skill{},
		&entity.TestType{},
		&entity.TestScheduleHeader{},
		&entity.TestApplicant{},
		&entity.SavedJob{},
		&entity.AdministrativeSelection{},
		&entity.AdministrativeResult{},
		&entity.Interview{},
		&entity.InterviewApplicant{},
		&entity.InterviewAssessor{},
		&entity.InterviewResult{},
		&entity.FgdSchedule{},
		&entity.FgdApplicant{},
		&entity.FgdAssessor{},
		&entity.FgdResult{},
		&entity.DocumentSending{},
		&entity.DocumentAgreement{},
		&entity.DocumentVerificationHeader{},
		&entity.DocumentVerificationLine{},
	)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migration success")
	}

	// seed answer type
	answerTypes := []entity.AnswerType{
		{Name: "Text"},
		{Name: "Dropdown"},
		{Name: "Multiple Choice"},
		{Name: "Checkbox"},
		{Name: "Attachment"},
	}

	for _, answerType := range answerTypes {
		err := db.Create(&answerType).Error
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Info("Seed AnswerType success")

	// seed document type
	documentTypes := []entity.DocumentType{
		{Name: "OFFERING_LETTER"},
		{Name: "PKWT"},
		{Name: "PKWTT"},
		{Name: "SURAT_PENGANTAR_MASUK"},
		{Name: "SURAT_IZIN_ORANG_TUA"},
		{Name: "DOCUMENT_CHECKING"},
		{Name: "KARYAWAN_TETAP"},
	}

	for _, documentType := range documentTypes {
		err := db.Create(&documentType).Error
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Info("Seed DocumentType success")
}
