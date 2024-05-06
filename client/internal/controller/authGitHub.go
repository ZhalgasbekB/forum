package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"lzhuk/clients/internal/cahe"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/model"
	"lzhuk/clients/pkg/errors"
	"net/http"
	"strings"
)

type githubUserInfo struct {
	Name     string `json:"name"`
	Email    string `json:"login"`
	Password string `json:"node_id"`
}

func GitHub(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=user:email", githubAuthEndPoint, "21c2671efe47648ceedd", "http://localhost:8082/github/callback")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorPage(w, "Попробуйте войти или зарегестрироваться через другой сервис", http.StatusInternalServerError)
		log.Printf("При входе через GitHub пользователь не получил код для обмена на токен")
		return
	}
	s := strings.NewReader(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, "21c2671efe47648ceedd", "acb9aa72c6f829fa1262760243287dfd71566859", "http://localhost:8082/github/callback"))
	response, err := http.Post(githubAuthEndAccessToken, "application/x-www-form-urlencoded", s)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при отправке POST запроса на получения токена от сервера GitHub. Ошибка: %v", err)
		return
	}
	defer response.Body.Close()

	resp, err := io.ReadAll(response.Body)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при чтении ответа тела от сервера GitHub со значением токена. Ошибка: %v", err)
		return
	}

	token, err := ExtractAccessTokenFromResponse(string(resp))
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при чтении полученного токена от сервера GitHub. Ошибка: %v", err)
		return
	}

	if token == "" {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Значение токена от сервера GitHub пустое. Ошибка: %v", err)
		return
	}

	info, err := getUserInfo(token, githubUserInfoURL)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Ошибка при получении данных о пользователе с помощью полученного от сервера GitHub токена. Ошибка: %v", err)
		return
	}

	var github githubUserInfo
	if err := json.Unmarshal(info, &github); err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Ошибка при конвертации полученных данных из JSON о конкретном пользователе от сервера GitHub. Ошибка: %v", err)
		return
	}
	if github.Name == "" {
		github.Name = github.Email
	}
	user := model.UserReq{Name: github.Name, Email: github.Email, Password: github.Password}
	jsonData, err := json.Marshal(user)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Ошибка при конвертации полученных данных из JSON о конкретном пользователе от сервера GitHub. Ошибка: %v", err)
		return
	}

	// Формирования POST запроса на регистрацию нового пользователя на сервере
	req, err := http.NewRequest("POST", registry, bytes.NewBuffer(jsonData))
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при формировании POST запроса на регистрацию пользователя по данным от сервера GitHub. Ошибка: %v", err)
		return
	}
	// Записываем тип контента в заголовок запроса
	req.Header.Set("Content-Type", "application/json")
	// Создаем структуру клиента для передачи запроса
	clientApi := http.Client{}
	// Отправляем запрос на сервис сервера
	respApi, err := clientApi.Do(req)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api при регистрации нового пользователя  по данным от сервера GitHub. Ошибка: %v", err)
		return
	}
	defer respApi.Body.Close()
	// Проверка кода статуса ответа сервера
	switch respApi.StatusCode {
	// Получен статус код 201 об успешной регистрации пользователя в системе
	case http.StatusCreated:
		// Формирования POST запроса на вход пользователя на сервере
		req, err := http.NewRequest("POST", login, bytes.NewBuffer(jsonData))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании POST запроса на вход пользователя  по данным от сервера GitHub. Ошибка: %v", err)
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
			log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api при входе пользователя  по данным от сервера GitHub. Ошибка: %v", err)
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
				log.Printf("Произошла ошибка при конвертации куки из ответа сервиса forum-api на вход пользователя  по данным от сервера GitHub. Ошибка: %v", err)
				return
			}
			// Получение в глобальную переменную имени вошедшего пользователя
			clientName, err = convertor.DecodeClientName(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при получении имени пользователя из ответа сервиса forum-api на вход пользователя  по данным от сервера GitHub. Ошибка: %v", err)
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
				log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос о входе пользователя  по данным от сервера GitHub")
				return
			}
			switch {
			// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Получены не валидные данные при входе по данным от сервера GitHub")
				return
			default:
				errorPage(w, "Сервис GitHub не доступен, попробуйте другой сервис. Убедитесь что Вы не зарегестрированы другим сервисом!", http.StatusInternalServerError)
				log.Printf("Получена не кастомная ошибка от сервиса forum-api при входе пользователя по данным от сервера GitHub")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при входе пользователя по данным от сервера GitHub")
			return
		}
	case http.StatusInternalServerError:
		discriptionMsg, err := convertor.DecodeErrorResponse(respApi)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос об регистрации нового пользователя по данным от сервера GitHub. Ошибка: %v", err)
			return
		}
		switch {
		// Получена ошибка что почта уже используется
		case discriptionMsg.Discription == "Email already exist":
			log.Printf("Почта уже зарегестрирована сервисом GitHub и осуществляется автоматический вход")
			// Формирования POST запроса на вход пользователя на сервере
			req, err := http.NewRequest("POST", login, bytes.NewBuffer(jsonData))
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при формировании POST запроса на вход пользователя по данным от сервера GitHub. Ошибка: %v", err)
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
				log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api при входе пользователя по данным от сервера GitHub. Ошибка: %v", err)
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
					log.Printf("Произошла ошибка при конвертации куки из ответа сервиса forum-api на вход пользователя по данным от сервера GitHub. Ошибка: %v", err)
					return
				}
				// Получение в глобальную переменную имени вошедшего пользователя
				clientName, err = convertor.DecodeClientName(resp)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при получении имени пользователя из ответа сервиса forum-api на вход пользователя по данным от сервера GitHub. Ошибка: %v", err)
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
					log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос о входе пользователя по данным от сервера GitHub")
					return
				}
				switch {
				// Получена ошибка что введены неверные учетные данные
				case discriptionMsg.Discription == "Invalid Credentials":
					errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
					log.Printf("Получены не валидные данные при входе пользователя по данным от сервера GitHub")
					return
				default:
					errorPage(w, "Сервис GitHub не доступен, попробуйте другой сервис. Убедитесь что Вы не зарегестрированы другим сервисом!", http.StatusInternalServerError)
					log.Printf("Получена не кастомная ошибка от сервиса forum-api при входе пользователя по данным от сервера GitHub")
					return
				}
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при входе пользователя по данным от сервера GitHub")
				return
			}
			// Получена ошибка что введены неверные учетные данные
		case discriptionMsg.Discription == "Invalid Credentials":
			errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
			log.Printf("Получены не валидные данные при регистрации по данным от сервера GitHub")
			return
		case discriptionMsg.Discription == "Not Found Any Data":
			errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
			log.Printf("Нет запрашиваемых данных по данным от сервера GitHub")
			return
		default:
			errorPage(w, "Сервис GitHub не доступен, попробуйте другой сервис. Убедитесь что Вы не зарегестрированы другим сервисом!", http.StatusInternalServerError)
			log.Printf("Получена не кастомная ошибка от сервиса forum-api при регистрации пользователя по данным от сервера GitHub")
			return
		}
		// Получен статус код 405 об неверном методе запроса с сервера
	case http.StatusMethodNotAllowed:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-api на регистрацию нового пользователя используется не верный метод по данным от сервера GitHub")
		return
		// Получен статус код не 201, 405 или 500
	default:
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Получен статус-код не 201, 405 и 500 от сервиса forum-api при регистрации нового пользователя по данным от сервера GitHub")
		return
	}
}
