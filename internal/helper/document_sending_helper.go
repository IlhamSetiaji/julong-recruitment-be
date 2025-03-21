package helper

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/microcosm-cc/bluemonday"
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

func (d *DocumentSendingHelper) ReplacePlaceHoldersContract(htmlTemplate string, data DocumentDataContract) (*string, error) {
	// Use Bluemonday with an explicit policy to avoid stripping necessary placeholders
	p := bluemonday.UGCPolicy() // Allow common safe HTML elements
	htmlTemplate = p.Sanitize(htmlTemplate)
	// htmlTemplate = strings.ReplaceAll(htmlTemplate, "<", "&lt;")
	// htmlTemplate = strings.ReplaceAll(htmlTemplate, ">", "&gt;")

	// Parse the sanitized HTML template
	tmpl, err := template.New("document").Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %v", err)
	}

	// Execute the template with the provided data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %v", err)
	}

	// Wrap the executed template in a full HTML document
	fullHtml := fmt.Sprintf(`
	<html>
	<head>
			<title>Document</title>
	</head>
	<body>
			%s
	</body>
	</html>`, buf.String())

	result := fullHtml
	// result := buf.String()
	return &result, nil
}

func (d *DocumentSendingHelper) ReplacePlaceHoldersOfferLetter(htmlTemplate string, data DocumentDataOfferLetter) (*string, error) {
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

func escapeAmpersandsOutsideTemplates(htmlTemplate string) string {
	var result strings.Builder
	inTemplate := false

	for i := 0; i < len(htmlTemplate); i++ {
		char := htmlTemplate[i]

		// Check if we're entering a Go template placeholder
		if i+1 < len(htmlTemplate) && char == '{' && htmlTemplate[i+1] == '{' {
			inTemplate = true
			result.WriteByte(char)
			continue
		}

		// Check if we're exiting a Go template placeholder
		if i+1 < len(htmlTemplate) && char == '}' && htmlTemplate[i+1] == '}' {
			inTemplate = false
			result.WriteByte(char)
			continue
		}

		// Escape & characters only outside of Go template placeholders
		if char == '&' && !inTemplate {
			result.WriteString("&amp;")
		} else {
			result.WriteByte(char)
		}
	}

	return result.String()
}
