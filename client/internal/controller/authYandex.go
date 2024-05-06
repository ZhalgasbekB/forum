package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"lzhuk/clients/internal/cahe"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

const (
	redirectURI = "http://localhost:8082/yandex/callback"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func YandexAuth(w http.ResponseWriter, r *http.Request) {
	url := "https://oauth.yandex.ru/authorize?response_type=code&client_id=" + "884ca44c4a5c47f6b462f355d2cc0505" + "&redirect_uri=" + redirectURI
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func YandexCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := exchangeCodeForToken(code)
	if err != nil {
		errorPage(w, "Попробуйте войти или зарегестрироваться через другой сервис", http.StatusInternalServerError)
		log.Printf("Произошла ошибка при обмене кода на токен у сервера Яндекса при регистрации и входе нового пользователя. Ошибка: %v", err)
		return
	}
	// По полученному токену запрашиваем с сервера Яндекс информация о пользователем
	userInfoYandex, err := getUserInfoYandex(token)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при получении данных о пользователе от сервера Яндекса. Ошибка: %v", err)
		return
	}

	jsonData, err := convertor.ConvertRegisterYandex(userInfoYandex)

	// Формирования POST запроса на регистрацию нового пользователя на сервере
	req, err := http.NewRequest("POST", registry, bytes.NewBuffer(jsonData))
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при формировании POST запроса на регистрацию пользователя по данным от сервера Яндекса. Ошибка: %v", err)
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
		log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api при регистрации нового пользователя  по данным от сервера Яндекса. Ошибка: %v", err)
		return
	}
	defer resp.Body.Close()
	// Проверка кода статуса ответа сервера
	switch resp.StatusCode {
	// Получен статус код 201 об успешной регистрации пользователя в системе
	case http.StatusCreated:
		// Формирования POST запроса на вход пользователя на сервере
		req, err := http.NewRequest("POST", login, bytes.NewBuffer(jsonData))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании POST запроса на вход пользователя  по данным от сервера Яндекса. Ошибка: %v", err)
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
			log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api при входе пользователя  по данным от сервера Яндекса. Ошибка: %v", err)
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
				log.Printf("Произошла ошибка при конвертации куки из ответа сервиса forum-api на вход пользователя  по данным от сервера Яндекса. Ошибка: %v", err)
				return
			}
			// Получение в глобальную переменную имени вошедшего пользователя
			clientName, err = convertor.DecodeClientName(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при получении имени пользователя из ответа сервиса forum-api на вход пользователя  по данным от сервера Яндекса. Ошибка: %v", err)
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
				log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос о входе пользователя  по данным от сервера Яндекса")
				return
			}
			switch {
			// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Получены не валидные данные при входе по данным от сервера Яндекса")
				return
			default:
				errorPage(w, "Сервис Яндекс не доступен, попробуйте другой сервис. Убедитесь что Вы не зарегестрированы другим сервисом!", http.StatusInternalServerError)
				log.Printf("Получена не кастомная ошибка от сервиса forum-api при входе пользователя по данным от сервера Яндекса")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при входе пользователя по данным от сервера Яндекса")
			return
		}
	case http.StatusInternalServerError:
		discriptionMsg, err := convertor.DecodeErrorResponse(resp)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос об регистрации нового пользователя по данным от сервера Яндекса. Ошибка: %v", err)
			return
		}
		switch {
		// Получена ошибка что почта уже используется
		case discriptionMsg.Discription == "Email already exist":
			log.Printf("Почта уже зарегестрирована сервисом Яндекса и осуществляется автоматический вход")
			// Формирования POST запроса на вход пользователя на сервере
			req, err := http.NewRequest("POST", login, bytes.NewBuffer(jsonData))
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при формировании POST запроса на вход пользователя по данным от сервера Яндекса. Ошибка: %v", err)
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
				log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api при входе пользователя по данным от сервера Яндекса. Ошибка: %v", err)
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
					log.Printf("Произошла ошибка при конвертации куки из ответа сервиса forum-api на вход пользователя по данным от сервера Яндекса. Ошибка: %v", err)
					return
				}
				// Получение в глобальную переменную имени вошедшего пользователя
				clientName, err = convertor.DecodeClientName(resp)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при получении имени пользователя из ответа сервиса forum-api на вход пользователя по данным от сервера Яндекса. Ошибка: %v", err)
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
					log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос о входе пользователя по данным от сервера Яндекса")
					return
				}
				switch {
				// Получена ошибка что введены неверные учетные данные
				case discriptionMsg.Discription == "Invalid Credentials":
					errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
					log.Printf("Получены не валидные данные при входе пользователя по данным от сервера Яндекса")
					return
				default:
					errorPage(w, "Сервис Яндекс не доступен, попробуйте другой сервис. Убедитесь что Вы не зарегестрированы другим сервисом!", http.StatusInternalServerError)
					log.Printf("Получена не кастомная ошибка от сервиса forum-api при входе пользователя по данным от сервера Яндекса")
					return
				}
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при входе пользователя по данным от сервера Яндекса")
				return
			}
			// Получена ошибка что введены неверные учетные данные
		case discriptionMsg.Discription == "Invalid Credentials":
			errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
			log.Printf("Получены не валидные данные при регистрации по данным от сервера Яндекса")
			return
		case discriptionMsg.Discription == "Not Found Any Data":
			errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
			log.Printf("Нет запрашиваемых данных по данным от сервера Яндекса")
			return
		default:
			errorPage(w, "Сервис Яндекс не доступен, попробуйте другой сервис. Убедитесь что Вы не зарегестрированы другим сервисом!", http.StatusInternalServerError)
			log.Printf("Получена не кастомная ошибка от сервиса forum-api при регистрации пользователя по данным от сервера Яндекса")
			return
		}
		// Получен статус код 405 об неверном методе запроса с сервера
	case http.StatusMethodNotAllowed:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-api на регистрацию нового пользователя используется не верный метод по данным от сервера Яндекса")
		return
		// Получен статус код не 201, 405 или 500
	default:
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Получен статус-код не 201, 405 и 500 от сервиса forum-api при регистрации нового пользователя по данным от сервера Яндекса")
		return
	}
}

func getUserInfoYandex(accessToken string) (map[string]interface{}, error) {
	// Формируем GET запрос к Google API для получения информации о пользователе
	req, err := http.NewRequest("GET", "https://login.yandex.ru/info?", nil)
	if err != nil {
		return nil, err
	}

	// Устанавливаем заголовок Authorization с токеном доступа
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("format", "json")
	req.Header.Set("jwt_secret", "a56f2664095a4d289c3520a34cd37538")

	// Отправляем запрос
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем и декодируем ответ
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func exchangeCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", "884ca44c4a5c47f6b462f355d2cc0505")
	data.Set("client_secret", "a56f2664095a4d289c3520a34cd37538")
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.Post("https://oauth.yandex.ru/token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Читаем ответ
	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	// Возвращаем токен
	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("Access token not found in response")
	}
	return accessToken, nil
}
