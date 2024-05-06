package model

import "time"

type AllPosts []struct {
	PostID       int       `json:"post_id"`
	UserID       int       `json:"user_id"`
	CategoryName string    `json:"category_name"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"create_at"`
	Author       string    `json:"name"`
	Like         int       `json:"likes"`
	Dislike      int       `json:"dislikes"`
	PathImage    string
}

type AllPostsConvertDate []struct {
	PostID       int
	UserID       int
	CategoryName string
	Title        string
	Description  string
	CreatedAt    string
	Author       string
	Like         int
	Dislike      int
	PathImage    string
}
