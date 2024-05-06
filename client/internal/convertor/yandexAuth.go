package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
)

func ConvertRegisterYandex(userInfo map[string]interface{}) ([]byte, error) {
	register := model.Register{
		Name:     userInfo["real_name"].(string),
		Email:    userInfo["default_email"].(string),
		Password: userInfo["client_id"].(string),
	}

	jsonData, err := json.Marshal(register)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
