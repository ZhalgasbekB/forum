package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
	"strconv"
)

// Функция для конвертации данных для редактирования комментариев
func ConvertUpdateComment(r *http.Request) (model.UpdateComment, error) {
	commentId, err := strconv.Atoi(r.FormValue("commentId"))
	if err != nil {
		return model.UpdateComment{}, err
	}
	posttId, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		return model.UpdateComment{}, err
	}

	updateComment := model.UpdateComment{
		ID:          commentId,
		Post:        posttId,
		Description: r.FormValue("description"),
	}

	return updateComment, nil
}

// Функция для конвертации данных с изменениями в комментариях
func ConvertUpdateCommentUser(r *http.Request) ([]byte, error) {
	commentId, err := strconv.Atoi(r.FormValue("commentId"))
	if err != nil {
		return nil, err
	}
	posttId, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		return nil, err
	}

	updateComment := model.UpdateComment{
		ID:          commentId,
		Post:        posttId,
		Description: r.FormValue("updatedComment"),
	}
	jsonData, err := json.Marshal(updateComment)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
