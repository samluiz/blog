package integrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/samluiz/blog/api/types"
)

const GITHUB_BASE_URL = "https://github.com"
const GITHUB_API_BASE_URL = "https://api.github.com"

var (
	GITHUB_CLIENT_ID    = os.Getenv("GITHUB_CLIENT_ID")
	GITHUB_SECRET_KEY   = os.Getenv("GITHUB_SECRET_KEY")
	GITHUB_REDIRECT_URI = os.Getenv("GITHUB_REDIRECT_URI")
)

func GetGithubAuthURL() string {
	return fmt.Sprintf("%s/login/oauth/authorize?client_id=%s&redirect_uri=%s", GITHUB_BASE_URL, GITHUB_CLIENT_ID, GITHUB_REDIRECT_URI)
}

func ExchangeGithubToken(code string) (*types.GithubOAuthResponse, error) {
	var githubResponse types.GithubOAuthResponse

	log.Default().Println("exchanging github code for token...")

	queryString := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s", GITHUB_CLIENT_ID, GITHUB_SECRET_KEY, code)

	request := fiber.Get(GITHUB_BASE_URL + "/login/oauth/access_token")
	request.Request().Header.Set("Accept", "application/json")
	request.QueryString(queryString)

	status, response, err := request.Bytes()

	if (status != 200) || (err != nil) {
		return nil, errors.New("error exchanging code for token: " + string(response))
	}

	jsonErr := json.Unmarshal(response, &githubResponse)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return &githubResponse, nil
}

func GetGithubUserInfo(accessToken string) (*types.GithubUserResponse, error) {
	var githubUserResponse types.GithubUserResponse

	log.Default().Println("getting user info from github")

	request := fiber.Get(GITHUB_API_BASE_URL + "/user")
	request.Request().Header.Set("Accept", "application/json")
	request.Request().Header.Set("Authorization", "Bearer "+accessToken)

	status, response, err := request.Bytes()

	log.Default().Printf("Status: %v", status)

	if (status != 200) || (err != nil) {
		return nil, errors.New("error getting user info from github: " + string(response))
	}

	jsonErr := json.Unmarshal(response, &githubUserResponse)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return &githubUserResponse, nil
}
