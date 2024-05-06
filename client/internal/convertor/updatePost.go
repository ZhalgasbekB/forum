package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
)

// Функция преобразования данных с изменениями в посте от пользователя для передачи на сервер
func ConvertUpdatePost(r *http.Request) ([]byte, error) {
	updatePost := model.Post{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	jsonData, err := json.Marshal(updatePost)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
