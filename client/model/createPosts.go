package model

type CreatePost struct {
	CategoryName string `json:"category_name"`
	Title        string
	Description  string
}
