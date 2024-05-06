package model

type VoteComment struct {
	CommentId  int  `json:"comment_id"`
	LikeStatus bool `json:"status"`
}
