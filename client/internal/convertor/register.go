package convertor

import (
	"encoding/json"
	"net/http"

	"lzhuk/clients/model"
)

// Функция конвертации пользовательских данных об регистрации в формат JSON
func ConvertRegister(r *http.Request) ([]byte, error) {
	register := model.Register{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	jsonData, err := json.Marshal(register)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
