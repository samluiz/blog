package article

import (
	"regexp"

	"github.com/samluiz/blog/api/types"
)

func getReadTime(article types.ArticleResponse) int {
	re := regexp.MustCompile(`/([^A-Za-z0-9])+/g`)
	formattedContent := re.ReplaceAllString(article.BodyMarkdown, "")

	return len(formattedContent) / 200
}
