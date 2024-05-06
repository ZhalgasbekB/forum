package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"lzhuk/clients/internal/cahe"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/model"
	"lzhuk/clients/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

type googleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Sub   string `json:"sub"`
}


func Google(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=profile email", googleAuthEndPoint, "16807684949-bp5bhvp85ar5sfj2iuasmfsf4l6bj4up.apps.googleusercontent.com", "http://localhost:8082/google/callback")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorPage(w, "Попробуйте войти или зарегестрироваться через другой сервис", http.StatusInternalServerError)
		log.Printf("При входе через Google пользователь не получил код для обмена на токен")
		return
	}
	s := strings.NewReader(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, "16807684949-bp5bhvp85ar5sfj2iuasmfsf4l6bj4up.apps.googleusercontent.com", "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", "http://localhost:8082/google/callback"))
	response, err := http.Post(googleAuthEndPointAccessToken, "application/x-www-form-urlencoded", s)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при отправке POST запроса на получения токена от сервера Google. Ошибка: %v", err)
		return
	}
	defer response.Body.Close()

	resp, err := io.ReadAll(response.Body)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при чтении ответа тела от сервера Google со значением токена. Ошибка: %v", err)
		return
	}

	token := ExtractValueFromBody(resp, "access_token")

	if token == "" {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Значение токена от сервера Google пустое. Ошибка: %v", err)
		return
	}
	info, err := getUserInfo(token, googleUserInfoURL)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Ошибка при получении данных о пользователе с помощью полученного от сервера Google токена. Ошибка: %v", err)
		return
	}

	var googleUserInfo googleUserInfo
	if err = json.Unmarshal(info, &googleUserInfo); err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Ошибка при конвертации полученных данных из JSON о конкретном пользователе от сервера Google. Ошибка: %v", err)
		return
	}

	us := model.UserReq{Name: googleUserInfo.Name, Email: googleUserInfo.Email, Password: googleUserInfo.Sub}
	us1, err := json.Marshal(us)
	fmt.Println(string(us1))
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Ошибка при конвертации полученных в JSON о конкретном пользователе для передачи на сервис forum-api. Ошибка: %v", err)
		return
	}
	req, err := http.NewRequest("POST", registry, bytes.NewBuffer(us1))
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при формировании POST запроса на регистрацию пользователя по данным от сервера Google. Ошибка: %v", err)
		return
	}
	// Записываем тип контента в заголовок запроса
	req.Header.Set("Content-Type", "application/json")
	// Создаем структуру клиента для передачи запроса
	client := http.Client{}
	// Отправляем запрос на сервис сервера
	respApi, err := client.Do(req)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api при регистрации нового пользователя  по данным от сервера Google. Ошибка: %v", err)
		return
	}
	defer respApi.Body.Close()
	// Проверка кода статуса ответа сервера
	switch respApi.StatusCode {
	// Получен статус код 201 об успешной регистрации пользователя в системе
	case http.StatusCreated:
		// Формирования POST запроса на вход пользователя на сервере
		req, err := http.NewRequest("POST", login, bytes.NewBuffer(us1))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании POST запроса на вход пользователя  по данным от сервера Google. Ошибка: %v", err)
			return
		}
		// Записываем тип контента в заголовок запроса
		req.Header.Set("Content-Type", "application/json")
		// Создаем структуру клиента для передачи запроса
		client := http.Client{}
		// Отправляем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api при входе пользователя  по данным от сервера Google. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			var clientName string
			// Получение сгенерированных сервером куки
			cookie, err := convertor.ConvertFirstCookie(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации куки из ответа сервиса forum-api на вход пользователя  по данным от сервера Google. Ошибка: %v", err)
				return
			}
			// Получение в глобальную переменную имени вошедшего пользователя
			clientName, err = convertor.DecodeClientName(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при получении имени пользователя из ответа сервиса forum-api на вход пользователя  по данным от сервера Google. Ошибка: %v", err)
				return
			}
			// Записываем клиента в хеш-таблицу
			cahe.Username[cookie.Value] = clientName
			// Записываем в ответ браузеру полученный экземпляр куки от сервера
			http.SetCookie(w, cookie)
			// Переход на домашнюю страницу пользователя
			http.Redirect(w, r, "http://localhost:8082/userd3", 302)
			return
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос о входе пользователя  по данным от сервера Google")
				return
			}
			switch {
			// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Получены не валидные данные при входе по данным от сервера Google")
				return
			default:
				errorPage(w, "Сервис Google не доступен, попробуйте другой сервис. Убедитесь что Вы не зарегестрированы другим сервисом!", http.StatusInternalServerError)
				log.Printf("Получена не кастомная ошибка от сервиса forum-api при входе пользователя по данным от сервера Google")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при входе пользователя по данным от сервера Google")
			return
		}
	case http.StatusInternalServerError:
		discriptionMsg, err := convertor.DecodeErrorResponse(respApi)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос об регистрации нового пользователя по данным от сервера Google. Ошибка: %v", err)
			return
		}
		switch {
		// Получена ошибка что почта уже используется
		case discriptionMsg.Discription == "Email already exist":
			log.Printf("Почта уже зарегестрирована сервисом Google и осуществляется автоматический вход")
			// Формирования POST запроса на вход пользователя на сервере
			req, err := http.NewRequest("POST", login, bytes.NewBuffer(us1))
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при формировании POST запроса на вход пользователя по данным от сервера Google. Ошибка: %v", err)
				return
			}
			// Записываем тип контента в заголовок запроса
			req.Header.Set("Content-Type", "application/json")
			// Создаем структуру клиента для передачи запроса
			client := http.Client{}
			// Отправляем запрос на сервис сервера
			resp, err := client.Do(req)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api при входе пользователя по данным от сервера Google. Ошибка: %v", err)
				return
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK:
				var clientName string
				// Получение сгенерированных сервером куки
				cookie, err := convertor.ConvertFirstCookie(resp)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при конвертации куки из ответа сервиса forum-api на вход пользователя по данным от сервера Google. Ошибка: %v", err)
					return
				}
				// Получение в глобальную переменную имени вошедшего пользователя
				clientName, err = convertor.DecodeClientName(resp)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при получении имени пользователя из ответа сервиса forum-api на вход пользователя по данным от сервера Google. Ошибка: %v", err)
					return
				}
				// Записываем клиента в хеш-таблицу
				cahe.Username[cookie.Value] = clientName
				// Записываем в ответ браузеру полученный экземпляр куки от сервера
				http.SetCookie(w, cookie)
				// Переход на домашнюю страницу пользователя
				http.Redirect(w, r, "http://localhost:8082/userd3", 302)
				return
			case http.StatusInternalServerError:
				discriptionMsg, err := convertor.DecodeErrorResponse(resp)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос о входе пользователя по данным от сервера Google")
					return
				}
				switch {
				// Получена ошибка что введены неверные учетные данные
				case discriptionMsg.Discription == "Invalid Credentials":
					errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
					log.Printf("Получены не валидные данные при входе пользователя по данным от сервера Google")
					return
				default:
					errorPage(w, "Сервис Google не доступен, попробуйте другой сервис. Убедитесь что Вы не зарегестрированы другим сервисом!", http.StatusInternalServerError)
					log.Printf("Получена не кастомная ошибка от сервиса forum-api при входе пользователя по данным от сервера Google")
					return
				}
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при входе пользователя по данным от сервера Google")
				return
			}
			// Получена ошибка что введены неверные учетные данные
		case discriptionMsg.Discription == "Invalid Credentials":
			errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
			log.Printf("Получены не валидные данные при регистрации по данным от сервера Google")
			return
		case discriptionMsg.Discription == "Not Found Any Data":
			errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
			log.Printf("Нет запрашиваемых данных по данным от сервера Google")
			return
		default:
			errorPage(w, "Сервис Google не доступен, попробуйте другой сервис. Убедитесь что Вы не зарегестрированы другим сервисом!", http.StatusInternalServerError)
			log.Printf("Получена не кастомная ошибка от сервиса forum-api при регистрации пользователя по данным от сервера Google")
			return
		}
		// Получен статус код 405 об неверном методе запроса с сервера
	case http.StatusMethodNotAllowed:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-api на регистрацию нового пользователя используется не верный метод по данным от сервера Google")
		return
		// Получен статус код не 201, 405 или 500
	default:
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Получен статус-код не 201, 405 и 500 от сервиса forum-api при регистрации нового пользователя по данным от сервера Google")
		return
	}
}

func ExtractValueFromBody(body []byte, key string) string {
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return ""
	}

	value, ok := response[key].(string)
	if !ok {
		return ""
	}
	return value
}

func ExtractAccessTokenFromResponse(response string) (string, error) {
	params, err := url.ParseQuery(response)
	if err != nil {
		return "", err
	}
	accessToken := params.Get("access_token")
	return accessToken, nil
}

func getUserInfo(accessToken string, userInfoURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
