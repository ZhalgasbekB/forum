package model

type VotePost struct {
	PostId     int  `json:"post_id"`
	LikeStatus bool `json:"status"`
}
