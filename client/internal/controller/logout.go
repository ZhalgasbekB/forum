package controller

import (
	"bytes"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/helpers"
	"lzhuk/clients/pkg/errors"
	"net/http"
	"strings"
)

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"):
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Конвертация данных при выходе пользователя
		jsonData, err := convertor.ConvertLogout(r)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации данных при выходе пользователя. Ошибка: %v", err)
			return
		}
		// Формируем новый запрос на сервис сервера
		req, err := http.NewRequest("POST", logoutUsers, bytes.NewBuffer(jsonData))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании запроса при выходе пользователя. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос на сервер
		req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
		req.Header.Set("Content-Type", "application/json")
		// Создаем структуру нового клиента
		client := http.Client{}
		// Отправляем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при отправке запроса при выходе пользователя. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusSeeOther:
			expiredCookie := &http.Cookie{
				Name:   "CookieUUID",
				Value:  "",
				MaxAge: -1,
				Path:   "/",
			}
			http.SetCookie(w, expiredCookie)
			http.Redirect(w, r, "http://localhost:8082/userd3", 302)

		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и описания от сервера на запрос об выходе пользователя")
				return
			}
			switch {
			// Получена ошибка что почта уже используется
			case discriptionMsg.Discription == "Email already exist":
				errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
				log.Printf("Пользователь пытается зарегестировать почту которая используется под другим аккаунтом")
				return
				// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Не валидные данные")
				return
			case discriptionMsg.Discription == "Not Found Any Data":
				errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
				log.Printf("Не найдено")
				return
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на выход пользователя")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на выход пользователя")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("Не верный метод запроса при выходе пользователя")
		return
	}
}
