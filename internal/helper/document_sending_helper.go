package helper

import (
	"bytes"
	"html/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentSendingHelper interface {
	ReplacePlaceHoldersCoverLetter(htmlTemplate string, data DocumentDataCoverLetter) (*string, error)
	ReplacePlaceHoldersContract(htmlTemplate string, data DocumentDataContract) (*string, error)
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
	MaritalStatus  string
}

type DocumentDataContract struct {
	Name                 string  `json:"name"`
	Gender               string  `json:"gender"`
	BirthPlace           string  `json:"birth_place"`
	BirthDate            string  `json:"birth_date"`
	MaritalStatus        string  `json:"marital_status"`
	EducationLevel       string  `json:"education_level"`
	Major                string  `json:"major"`
	Religion             string  `json:"religion"`
	ApplicantAddress     string  `json:"applicant_address"`
	Position             string  `json:"position"`
	JobLevel             string  `json:"job_level"`
	Company              string  `json:"company"`
	Location             string  `json:"location"`
	BasicWage            float64 `json:"basic_wage"`
	PositionalAllowance  float64 `json:"positional_allowance"`
	OperationalAllowance float64 `json:"operational_allowance"`
	MealAllowance        float64 `json:"meal_allowance"`
	HomeTripTicket       string  `json:"hometrip_ticket"`
	JoinedDate           string  `json:"joined_date"`
	HiredStatus          string  `json:"hired_status"`
	ApprovalBy           string  `json:"approval_by"`
	DocumentDate         string  `json:"document_date"`
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

func (d *DocumentSendingHelper) ReplacePlaceHoldersContract(htmlTemplate string, data DocumentDataContract) (*string, error) {
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
