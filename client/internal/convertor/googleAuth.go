package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
)

func ConvertRegisterGoogle(userInfo map[string]interface{}) ([]byte, error) {
	register := model.Register{
		Name:     userInfo["given_name"].(string),
		Email:    userInfo["email"].(string),
		Password: userInfo["sub"].(string),
	}

	jsonData, err := json.Marshal(register)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
