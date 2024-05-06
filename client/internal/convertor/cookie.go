package convertor

import (
	"errors"
	"fmt"
	"net/http"
)

func ConvertFirstCookie(resp *http.Response) (*http.Cookie, error) {
	cookie := resp.Cookies()[0]
	if cookie == nil {
		return nil, errors.New("В ответе нет кук")
	}
	fmt.Println(cookie)
	return cookie, nil
}
