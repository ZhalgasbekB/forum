package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
	"strconv"
)

// Функция конвертации данных при создании комментария
func ConvertCreateComment(r *http.Request) ([]byte, error) {
	postId, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		return nil, err
	}
	createComment := model.CreateComment{
		Post:        postId,
		Description: r.FormValue("comment"),
	}
	jsonData, err := json.Marshal(createComment)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
