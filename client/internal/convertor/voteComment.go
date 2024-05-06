package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
	"strconv"
)

// Функция для конвертации данных при постановке голоса на комментарий
func ConvertVoteComment(r *http.Request) ([]byte, error) {
	commentId, err := strconv.Atoi(r.FormValue("commentId"))
	if err != nil {
		return nil, err
	}
	var status bool
	switch r.FormValue("vote") {
	case "true":
		status = true
	case "false":
		status = false
	}
	voteComment := model.VoteComment{
		CommentId:  commentId,
		LikeStatus: status,
	}
	jsonData, err := json.Marshal(voteComment)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
