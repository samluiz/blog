package integrations

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/samluiz/blog/api/types"
	"github.com/samluiz/blog/common/date"
)

const DEV_TO_API_BASE_URL = "https://dev.to/api"
const DEV_TO_USERNAME = "samluiz"

func GetArticleBySlugDevTo(slug string) (*types.ArticleResponse, error) {
	var getArticleResponse types.GetArticleByPathResponse
	var articleResponse types.ArticleResponse

	log.Default().Println("Getting article from dev.to")

	request := fiber.Get(DEV_TO_API_BASE_URL + "/articles/" + DEV_TO_USERNAME + "/" + slug)

	status, response, err := request.Bytes()

	log.Default().Printf("Status: %v", status)

	if (status != 200) || (err != nil) {
		log.Default().Printf("Error: %v", err)
		return nil, errors.New("error getting article from dev.to: " + string(response))
	}

	jsonErr := json.Unmarshal(response, &getArticleResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	getArticleResponse.PublishedAt = date.FormatDate(getArticleResponse.PublishedAt)
	articleResponse = types.ArticleResponse(getArticleResponse)

	return &articleResponse, nil
}

func GetArticlesFromDevTo() ([]types.ArticleResponse, error) {
	var articles []types.GetArticleByPathResponse
	articlesResponse := make([]types.ArticleResponse, len(articles))

	log.Default().Println("Getting articles from dev.to")

	request := fiber.Get(DEV_TO_API_BASE_URL + "/articles/me/published")
	request.Set("api-key", os.Getenv("DEV_TO_API_KEY"))
	request.Request().URI().SetQueryString("page=1&per_page=5")

	status, response, err := request.Bytes()

	log.Default().Printf("Status: %v", status)

	if (status != 200) || (err != nil) {
		log.Default().Printf("Error: %v", err)
		return nil, errors.New("error getting articles from dev.to: " + string(response))
	}

	jsonErr := json.Unmarshal(response, &articles)
	if jsonErr != nil {
		return nil, jsonErr
	}

	for _, a := range articles {
		a.PublishedAt = date.FormatDate(a.PublishedAt)
		articlesResponse = append(articlesResponse, types.ArticleResponse(a))
	}

	return articlesResponse, nil
}
