package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
)

// Функция для конвертации данных при входе пользователя
func ConvertLogin(r *http.Request) ([]byte, error) {
	login := model.Login{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	jsonData, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// Функция для конвертации из ответа имени клиента
func DecodeClientName(resp *http.Response) (string, error) {
	client := &model.UserResponseDTO{} 
	if err := json.NewDecoder(resp.Body).Decode(client); err != nil {
		return "", err
	}
	return client.Name, nil
}