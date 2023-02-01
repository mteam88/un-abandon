package database

type User struct {
	Username string
	Token	string
	GithubID int64
}

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"html_url"`
	ID          int64  `json:"id"`
	Token		string `json:"token"`
}