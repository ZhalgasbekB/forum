package model

import "time"

type Post struct {
	PostId       int    `json:"post_id"`
	UserId       int    `json:"user_id"`
	CategoryName string `json:"category_name"`
	Title        string
	Description  string
	CreateDate   time.Time `json:"create_at"`

	Author  string `json:"name"`
	Like    int    `json:"likes"`
	Dislike int    `json:"dislikes"`
	PathImage string
}

type GetPost struct {
	Post     *Post
	Comments []*Comment `json:"comments"`
}

type GetPostDate struct {
	Post    *PostDate
	Comment []*CommentDate `json:"comments"`
}

type PostId struct {
	PostId int `json:"post_id"`
}


type PathImagePost struct {
	Path string `json:"path"`
}