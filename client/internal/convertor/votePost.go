package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
	"strconv"
)

// Функция конвертации данных при постановке голосов на пост
func ConvertVotePost(r *http.Request) ([]byte, error) {
	postId, err := strconv.Atoi(r.FormValue("postId"))
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
	votePost := model.VotePost{
		PostId:     postId,
		LikeStatus: status,
	}
	jsonData, err := json.Marshal(votePost)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
