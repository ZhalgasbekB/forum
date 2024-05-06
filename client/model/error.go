package model

type Error struct {
	Status      int    `json:"status"`
	Discription string `json:"message"`
}
