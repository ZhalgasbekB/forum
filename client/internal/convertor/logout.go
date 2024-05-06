package convertor

import (
	"encoding/json"
	"fmt"
	"lzhuk/clients/model"
	"net/http"
)

// Функция для конвертации данных при выходе пользователя
func ConvertLogout(r *http.Request) ([]byte, error) {
	cookie := r.Cookies()
	uuid := cookie[0].String()
	fmt.Println(uuid)
	logoutUser := model.LogoutUser{
		UUID: uuid,
	}
	jsonData, err := json.Marshal(logoutUser)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
