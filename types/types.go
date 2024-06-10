package types

type CredentialsUser struct {
	Name string
	ID   float64
}

type PostData struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
