package helper

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentSendingHelper interface {
	ReplacePlaceHoldersCoverLetter(htmlTemplate string, data DocumentDataCoverLetter) (*string, error)
	ReplacePlaceHoldersContract(htmlTemplate string, data DocumentDataContract) (*string, error)
	ReplacePlaceHoldersOfferLetter(htmlTemplate string, data DocumentDataOfferLetter) (*string, error)
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
	Name                 string `json:"name"`
	Gender               string `json:"gender"`
	BirthPlace           string `json:"birth_place"`
	BirthDate            string `json:"birth_date"`
	MaritalStatus        string `json:"marital_status"`
	EducationLevel       string `json:"education_level"`
	Major                string `json:"major"`
	Religion             string `json:"religion"`
	ApplicantAddress     string `json:"applicant_address"`
	Position             string `json:"position"`
	JobLevel             string `json:"job_level"`
	Company              string `json:"company"`
	Location             string `json:"location"`
	BasicWage            int    `json:"basic_wage"`
	PositionalAllowance  int    `json:"positional_allowance"`
	OperationalAllowance int    `json:"operational_allowance"`
	MealAllowance        int    `json:"meal_allowance"`
	HometripTicket       string `json:"hometrip_ticket"`
	JoinedDate           string `json:"joined_date"`
	HiredStatus          string `json:"hired_status"`
	ApprovalBy           string `json:"approval_by"`
	DocumentDate         string `json:"document_date"`
}

type DocumentDataOfferLetter struct {
	DocumentDate string `json:"document_date"`
	Name         string `json:"name"`
	Position     string `json:"position"`
	Company      string `json:"company"`
	ApprovalBy   string `json:"approval_by"`
	BasicWage    int    `json:"basic_wage"`
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

type DataContent struct {
	Content template.HTML
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
	fullHtml := fmt.Sprintf(`<!DOCTYPE html>
	<html>
	<head>
			<title>Document</title>
	</head>
	<body>
			%s
	</body>
	</html>`, template.HTML(htmlTemplate))

	tmpl, err := template.New("document").Parse(fullHtml)
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

func (d *DocumentSendingHelper) ReplacePlaceHoldersOfferLetter(htmlTemplate string, data DocumentDataOfferLetter) (*string, error) {
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
