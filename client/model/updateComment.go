package model

type UpdateComment struct {
	ID          int    `json:"id"`
	Post        int    `json:"post_id"`
	Description string `json:"description"`
}
