package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
)

func DecodeErrorResponse(resp *http.Response) (*model.Error, error) {
	errResp := &model.Error{}
	if err := json.NewDecoder(resp.Body).Decode(errResp); err != nil {
		return nil, err
	}
	return &model.Error{
		Status: errResp.Status,
		Discription: errResp.Discription,
	}, nil
}
