package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
	"strconv"
)

// Функция конвертации данных для удаления комментария
func ConvertDeleteComment(r *http.Request) ([]byte, error) {
	commentId, err := strconv.Atoi(r.FormValue("commentId"))
	if err != nil {
		return nil, err
	}
	posttId, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		return nil, err
	}

	deleteComment := model.Comment{
		ID:   commentId,
		Post: posttId,
	}
	jsonData, err := json.Marshal(deleteComment)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
