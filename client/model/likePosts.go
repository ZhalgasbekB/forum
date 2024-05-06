package model

type PostDate struct {
	PostId       int
	UserId       int
	CategoryName string
	Title        string
	Description  string
	CreateDate   string

	Author  string
	Like    int
	Dislike int
	PathImage    string
}
