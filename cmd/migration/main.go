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
		&entity.TemplateActivity{},
		&entity.TemplateActivityLine{},
		&entity.ProjectRecruitmentHeader{},
		&entity.ProjectRecruitmentLine{},
		&entity.JobPosting{},
		&entity.AnswerType{},
		&entity.TemplateQuestion{},
		&entity.Question{},
		&entity.QuestionOption{},
		&entity.QuestionResponse{},
	)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migration success")
	}

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
}
