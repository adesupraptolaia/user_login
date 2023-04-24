package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/adesupraptolaia/user_login/config"
	"github.com/adesupraptolaia/user_login/internal/entity"
	"github.com/adesupraptolaia/user_login/pkg/jwt"
	"github.com/labstack/gommon/log"
)

type AuthRepo interface {
	CreateUser(entity.User) (*entity.User, error)
	DeleteUser(string) (*entity.User, error)
}

type authRepo struct{}

func NewAuthRepo() AuthRepo {
	return &authRepo{}
}

type AuthReponse struct {
	Status       string       `json:"status"`
	ErrorMessage string       `json:"error_message"`
	Data         *entity.User `json:"data"`
}

func (repo *authRepo) CreateUser(user entity.User) (*entity.User, error) {
	log.Info("create user to auth service")

	url := fmt.Sprintf("http://%s/user/create", getBaseURL())

	return doRequest(http.MethodPost, url, user)
}

func (repo *authRepo) DeleteUser(userKsuid string) (*entity.User, error) {
	log.Infof("delete user with ksuid %s to auth service", userKsuid)

	url := fmt.Sprintf("http://%s/user/%s", getBaseURL(), userKsuid)

	return doRequest(http.MethodDelete, url, nil)
}

func doRequest(httpMethod, url string, request interface{}) (*entity.User, error) {
	reqJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error when marshal create user request, err: %s", err.Error())
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("error when make http request, err: %s", err.Error())
	}

	accessToken, err := jwt.CreateAccessToken("2OokWa2yDw7yi7o9RpsAl58xuoW", entity.ADMIN)
	if err != nil {
		return nil, fmt.Errorf("error when create accessToken, err %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when calling to auth service, err: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error when read response body, err: %s", err.Error())
	}

	var response AuthReponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshal response body with err %s", err.Error())
	}

	if response.Status != "success" {
		return nil, fmt.Errorf(response.ErrorMessage)
	}

	return response.Data, nil
}

func getBaseURL() string {
	baseUrl := os.Getenv("AUTH_SERVICE_PRIVATE_URL")
	if baseUrl == "" {
		cfg := config.Config
		baseUrl = cfg.AuthServicePrivateUrl
	}

	return baseUrl
}
