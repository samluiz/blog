package types

type SessionUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`
	Provider string `json:"provider"`
}
