package types

type ArticleResponse struct {
	ID                 int
	Title              string
	Description        string
	Slug               string
	TagList            []string
	PublishedAt        string
	ReadingTimeMinutes int
	BodyMarkdown       string
}

type GetArticlesResponse struct {
	ID                 int      `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Slug               string   `json:"slug"`
	TagList            []string `json:"tag_list"`
	PublishedAt        string   `json:"published_at"`
	ReadingTimeMinutes int      `json:"reading_time_minutes"`
	BodyMarkdown       string   `json:"body_markdown"`
}

type GetArticleByPathResponse struct {
	ID                 int      `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Slug               string   `json:"slug"`
	TagList            []string `json:"tags"`
	PublishedAt        string   `json:"published_at"`
	ReadingTimeMinutes int      `json:"reading_time_minutes"`
	BodyMarkdown       string   `json:"body_markdown"`
}
