package helper

import (
	"bytes"
	"html/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentSendingHelper interface {
	ReplacePlaceHoldersCoverLetter(htmlTemplate string, data DocumentDataCoverLetter) (*string, error)
}

type DocumentSendingHelper struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

type DocumentDataCoverLetter struct {
	Company        string
	DocumentDate   string
	Name           string
	Gender         string
	BirthPlace     string
	BirthDate      string
	EducationLevel string
	Major          string
	Position       string
	JobLevel       string
	JoinedDate     string
	HiredStatus    string
}

func NewDocumentSendingHelper(
	log *logrus.Logger,
	viper *viper.Viper,
) IDocumentSendingHelper {
	return &DocumentSendingHelper{
		Log:   log,
		Viper: viper,
	}
}

func DocumentSendingHelperFactory(log *logrus.Logger, viper *viper.Viper) IDocumentSendingHelper {
	return NewDocumentSendingHelper(log, viper)
}

func (d *DocumentSendingHelper) ReplacePlaceHoldersCoverLetter(htmlTemplate string, data DocumentDataCoverLetter) (*string, error) {
	tmpl, err := template.New("document").Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	result := buf.String()
	return &result, nil
}
