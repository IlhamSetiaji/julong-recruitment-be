package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
	SyncUpdateEmployeeNationalDataMidsuit(midsuitId int, payload request.SyncUpdateEmployeeNationalDataMidsuitRequest, jwtToken string) (*string, error)
	SyncUpdateEmployeeNationalData1Midsuit(midsuitId int, payload request.SyncUpdateEmployeeNationalData1MidsuitRequest, jwtToken string) (*string, error)
	SyncUpdateEmployeeNationalData3Midsuit(midsuitId int, payload request.SyncUpdateEmployeeNationalData3MidsuitRequest, jwtToken string) (*string, error)
	SyncUpdateEmployeeNationalData4Midsuit(midsuitId int, payload request.SyncUpdateEmployeeNationalData4MidsuitRequest, jwtToken string) (*string, error)
	SyncUpdateEmployeeNationalData5Midsuit(midsuitId int, payload request.SyncUpdateEmployeeNationalData5MidsuitRequest, jwtToken string) (*string, error)
	SyncEmployeeImageMidsuit(payload request.SyncEmployeeImageMidsuitRequest, jwtToken string) (*string, error)
	SyncUpdateEmployeeImageMidsuit(midsuitId int, payload request.SyncUpdateEmployeeImageMidsuitRequest, jwtToken string) (*string, error)
	RecruitmentTypeMidsuitAPIWithoutFilter(jwtToken string) (*RecruitmentTypeMidsuitAPIResponse, error)
	RecruitmentTypeMidsuitAPI(filter string, jwtToken string) (*RecruitmentTypeMidsuitAPIResponse, error)
	SyncGenerateUserMidsuit(empMidsuitID int, jwtToken string) (*string, error)
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
	ID int `json:"id"`
}

type SyncEmployeeJobMidsuitResponse struct {
	ID int `json:"id"`
}

type RecruitmentTypeMidsuitResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type RecruitmentTypeMidsuitAPIResponse struct {
	PageCount   int                              `json:"page-count"`
	RecordsSize int                              `json:"records-size"`
	SkipRecords int                              `json:"skip-records"`
	RowCount    int                              `json:"row-count"`
	ArrayCount  int                              `json:"array-count"`
	Records     []RecruitmentTypeMidsuitResponse `json:"records"`
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

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
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

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
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

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
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

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
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

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
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

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
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

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
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

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
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

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
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

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}

func (s *MidsuitService) SyncUpdateEmployeeNationalDataMidsuit(midsuitId int, payload request.SyncUpdateEmployeeNationalDataMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/hc_employee/" + strconv.Itoa(midsuitId)
	method := "PUT"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalDataMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalDataMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalDataMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalDataMidsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalDataMidsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}

func (s *MidsuitService) SyncUpdateEmployeeNationalData1Midsuit(midsuitId int, payload request.SyncUpdateEmployeeNationalData1MidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/hc_employee/" + strconv.Itoa(midsuitId)
	method := "PUT"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData1Midsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData1Midsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData1Midsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData1Midsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData1Midsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}

func (s *MidsuitService) SyncUpdateEmployeeNationalData3Midsuit(midsuitId int, payload request.SyncUpdateEmployeeNationalData3MidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/hc_employee/" + strconv.Itoa(midsuitId)
	method := "PUT"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData3Midsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData3Midsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData3Midsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData3Midsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData3Midsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}

func (s *MidsuitService) SyncUpdateEmployeeNationalData4Midsuit(midsuitId int, payload request.SyncUpdateEmployeeNationalData4MidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/hc_employee/" + strconv.Itoa(midsuitId)
	method := "PUT"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData4Midsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData4Midsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData4Midsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData4Midsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData4Midsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}

func (s *MidsuitService) SyncUpdateEmployeeNationalData5Midsuit(midsuitId int, payload request.SyncUpdateEmployeeNationalData5MidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/hc_employee/" + strconv.Itoa(midsuitId)
	method := "PUT"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData5Midsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData5Midsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData5Midsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData5Midsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeNationalData5Midsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}

func (s *MidsuitService) SyncEmployeeImageMidsuit(payload request.SyncEmployeeImageMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/AD_Image"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeImageMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeImageMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeImageMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeImageMidsuit] Error when fetching response cok: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncEmployeeImageMidsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}

func (s *MidsuitService) SyncUpdateEmployeeImageMidsuit(midsuitId int, payload request.SyncUpdateEmployeeImageMidsuitRequest, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/models/hc_employee/" + strconv.Itoa(midsuitId)
	method := "PUT"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeImageMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeImageMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeImageMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeImageMidsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncEmployeeMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncUpdateEmployeeImageMidsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.ID)
	return &idStr, nil
}

func (s *MidsuitService) RecruitmentTypeMidsuitAPIWithoutFilter(jwtToken string) (*RecruitmentTypeMidsuitAPIResponse, error) {
	baseURL := strings.TrimRight(s.Viper.GetString("midsuit.url"), "/")
	endpoint := strings.TrimLeft(s.Viper.GetString("midsuit.api_endpoint"), "/")

	urlStr := fmt.Sprintf("%s/%s/models/HC_RecruitmentType",
		baseURL,
		endpoint)

	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		s.Log.Error(err)
		return nil, fmt.Errorf("[MidsuitService.RecruitmentTypeMidsuitAPIWithoutFilter] Error when creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, fmt.Errorf("[MidsuitService.RecruitmentTypeMidsuitAPIWithoutFilter] Error when fetching response: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		s.Log.Error(err)
		return nil, fmt.Errorf("[MidsuitService.RecruitmentTypeMidsuitAPIWithoutFilter] Error reading response body: %w", err)
	}

	s.Log.Info("[MidsuitService.RecruitmentTypeMidsuitAPIWithoutFilter] Response Status Code: ", res.StatusCode)
	s.Log.Info("[MidsuitService.RecruitmentTypeMidsuitAPIWithoutFilter] Response Header: ", res.Header)
	s.Log.Info("[MidsuitService.RecruitmentTypeMidsuitAPIWithoutFilter] Response Body: ", string(bodyBytes))

	if res.StatusCode != http.StatusOK {
		s.Log.Error(string(bodyBytes))
		return nil, fmt.Errorf("[MidsuitService.RecruitmentTypeMidsuitAPIWithoutFilter] Error response from API: %s", string(bodyBytes))
	}

	var recruitmentTypeResponse RecruitmentTypeMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &recruitmentTypeResponse); err != nil {
		s.Log.Error(err)
		return nil, fmt.Errorf("[MidsuitService.RecruitmentTypeMidsuitAPIWithoutFilter] Error when unmarshalling response: %w", err)
	}

	return &recruitmentTypeResponse, nil
}

func (s *MidsuitService) RecruitmentTypeMidsuitAPI(filter string, jwtToken string) (*RecruitmentTypeMidsuitAPIResponse, error) {
	// Properly encode the filter value to handle special characters
	encodedFilter := url.QueryEscape(filter)

	baseURL := strings.TrimRight(s.Viper.GetString("midsuit.url"), "/")
	endpoint := strings.TrimLeft(s.Viper.GetString("midsuit.api_endpoint"), "/")

	urlStr := fmt.Sprintf("%s/%s/models/HC_RecruitmentType?$filter=Value eq '%s'",
		baseURL,
		endpoint,
		encodedFilter)

	method := "GET"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		s.Log.Error(err)
		return nil, fmt.Errorf("[MidsuitService.RecruitmentTypeMidsuitAPI] Error when creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, fmt.Errorf("[MidsuitService.RecruitmentTypeMidsuitAPI] Error when fetching response: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		s.Log.Error(err)
		return nil, fmt.Errorf("[MidsuitService.RecruitmentTypeMidsuitAPI] Error reading response body: %w", err)
	}

	s.Log.Info("[MidsuitService.RecruitmentTypeMidsuitAPI] Response Status Code: ", res.StatusCode)
	s.Log.Info("[MidsuitService.RecruitmentTypeMidsuitAPI] Response Header: ", res.Header)
	s.Log.Info("[MidsuitService.RecruitmentTypeMidsuitAPI] Response Body: ", string(bodyBytes))

	if res.StatusCode != http.StatusOK {
		s.Log.Error(string(bodyBytes))
		return nil, fmt.Errorf("[MidsuitService.RecruitmentTypeMidsuitAPI] Error response from API: %s", string(bodyBytes))
	}

	var recruitmentTypeResponse RecruitmentTypeMidsuitAPIResponse
	if err := json.Unmarshal(bodyBytes, &recruitmentTypeResponse); err != nil {
		s.Log.Error(err)
		return nil, fmt.Errorf("[MidsuitService.RecruitmentTypeMidsuitAPI] Error when unmarshalling response: %w", err)
	}

	return &recruitmentTypeResponse, nil
}

type SummaryResponse struct {
	HcEmployeeID int `json:"HC_Employee_ID"`
}

type SyncGenerateUserMidsuitResponse struct {
	AdPinstanceID int             `json:"AD_PInstance_ID"`
	Summary       SummaryResponse `json:"summary"`
}

func (s *MidsuitService) SyncGenerateUserMidsuit(empMidsuitID int, jwtToken string) (*string, error) {
	url := s.Viper.GetString("midsuit.url") + s.Viper.GetString("midsuit.api_endpoint") + "/processes/hcm_generateuser"
	method := "POST"

	payload := map[string]interface{}{
		"HC_Employee_ID": empMidsuitID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncGenerateUserMidsuit] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncGenerateUserMidsuit] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncGenerateUserMidsuit] Error when fetching response: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncGenerateUserMidsuit] Error when fetching response: " + string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	var syncResponse SyncGenerateUserMidsuitResponse
	if err := json.Unmarshal(bodyBytes, &syncResponse); err != nil {
		s.Log.Error(err)
		return nil, errors.New("[MidsuitService.SyncGenerateUserMidsuit] Error when unmarshalling response: " + err.Error())
	}

	idStr := strconv.Itoa(syncResponse.Summary.HcEmployeeID)
	return &idStr, nil
}
