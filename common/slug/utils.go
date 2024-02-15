package slug

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func GenerateSlugId() string {
	return uuid.NewString()[:7]
}

func GenerateSlug(title string, slugId string) string {
	slug := removeSpecialChars(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ToLower(slug)
	return slug + "-" + slugId
}

func removeSpecialChars(s string) string {
	regex := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	s = regex.ReplaceAllString(s, "")
	regex = regexp.MustCompile(`\s+`)
	s = regex.ReplaceAllString(s, " ")
	return s
}