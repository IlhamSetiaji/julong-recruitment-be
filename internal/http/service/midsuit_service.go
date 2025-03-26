package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type IMidsuitService interface {
	AuthOneStep() (*AuthOneStepResponse, error)
	SyncEmployeeMidsuit(payload request.SyncEmployeeMidsuitRequest, jwtToken string) (*string, error)
	SyncEmployeeJobMidsuit(payload request.SyncEmployeeJobMidsuitRequest, jwtToken string) (*string, error)
	SyncEmployeeWorkExperienceMidsuit(payload request.SyncEmployeeWorkExperienceMidsuitRequest, jwtToken string) (*string, error)
	SyncEmployeeEducationMidsuit(payload request.SyncEmployeeEducationMidsuitRequest, jwtToken string) (*string, error)
	SyncEmployeeAllowanceMidsuit(payload request.SyncEmployeeAllowanceMidsuitRequest, jwtToken string) (*string, error)
}

type MidsuitService struct {
	Viper *viper.Viper
	Log   *logrus.Logger
	DB    *gorm.DB
}

func NewMidsuitService(
	viper *viper.Viper,
	log *logrus.Logger,
	db *gorm.DB,
) IMidsuitService {
	return &MidsuitService{
		Viper: viper,
		Log:   log,
		DB:    db,
	}
}

func MidsuitServiceFactory(
	viper *viper.Viper,
	log *logrus.Logger,
) IMidsuitService {
	db := config.NewDatabase()
	return NewMidsuitService(viper, log, db)
}

type AuthOneStepResponse struct {
	UserID       int    `json:"userId"`
	Language     string `json:"language"`
	MenuTreeID   int    `json:"menuTreeId"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type SyncEmployeeMidsuitResponse struct {
	ID string `json:"id"`
}

type SyncEmployeeJobMidsuitResponse struct {
	ID string `json:"id"`
}

func (s *MidsuitService) AuthOneStep() (*AuthOneStepResponse, error) {
	payload := map[string]interface{}{
		"userName": s.Viper.GetString("midsuit.username"),
		// "password": s.Viper.GetString("midsuit.username") + "321!",
		"password": "JgiMidsuit123!",
		"parameters": map[string]interface{}{
			"clientId":       s.Viper.GetString("midsuit.client_id"),
			"roleId":         s.Viper.GetString("midsuit.role_id"),
			"organizationId": 0,
		},
	}

	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + s.Viper.GetString("midsuit.auth_endpoint")
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.AuthOneStep] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.AuthOneStep] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.AuthOneStep] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.AuthOneStep] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var authResponse AuthOneStepResponse
	if err := json.Unmarshal(bodyBytes, &authResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.AuthOneStep] Error when unmarshalling response: " + err.Error())
	}

	return &authResponse, nil
}

func (s *MidsuitService) SyncEmployeeMidsuit(payload request.SyncEmployeeMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/hc_employee"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeMidsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeMidsuit] Error when unmarshalling response: " + err.Error())
	}

	return &syncResponse.ID, nil
}

func (s *MidsuitService) SyncEmployeeJobMidsuit(payload request.SyncEmployeeJobMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/HC_EmployeeJob"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeJobMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeJobMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeJobMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeJobMidsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeJobMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeJobMidsuit] Error when unmarshalling response: " + err.Error())
	}

	return &syncResponse.ID, nil
}

func (s *MidsuitService) SyncEmployeeWorkExperienceMidsuit(payload request.SyncEmployeeWorkExperienceMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/HC_WorkHistory"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeWorkExperienceMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeWorkExperienceMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeWorkExperienceMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeWorkExperienceMidsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeWorkExperienceMidsuit] Error when unmarshalling response: " + err.Error())
	}

	return &syncResponse.ID, nil
}

func (s *MidsuitService) SyncEmployeeEducationMidsuit(payload request.SyncEmployeeEducationMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/HC_EmployeeEducation"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeEducationMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeEducationMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeEducationMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeEducationMidsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeEducationMidsuit] Error when unmarshalling response: " + err.Error())
	}

	return &syncResponse.ID, nil
}

func (s *MidsuitService) SyncEmployeeAllowanceMidsuit(payload request.SyncEmployeeAllowanceMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/HC_AllowanceProvision"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeAllowanceMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeAllowanceMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeAllowanceMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeAllowanceMidsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeAllowanceMidsuit] Error when unmarshalling response: " + err.Error())
	}

	return &syncResponse.ID, nil
}
