package model

type CreateComment struct {
	Post        int    `json:"post_id"`
	Description string `json:"description"`
}
