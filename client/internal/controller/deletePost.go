package controller

import (
	"fmt"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/helpers"
	"lzhuk/clients/pkg/errors"
	"net/http"
	"strings"
)

func DeletePost(w http.ResponseWriter, r *http.Request) {
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
	case http.MethodPost:
		// Формируем запрос на удаление конкретного поста по id
		req, err := http.NewRequest("DELETE", fmt.Sprintf(deletePosts+"%s", r.FormValue("postId")), nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании DELETE запроса на удаление конкретного поста пользователя. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос на сервис сервера
		req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
		// Создаем структуру нового клиента
		client := http.Client{}
		// Отправляем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при отправке запроса на сервис forum-api на удаление конкретного поста пользователя. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusAccepted:
			http.Redirect(w, r, "http://localhost:8082/userd3/myposts", 302)
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и её описания от сервиса forum-api на запрос об удалении конкретного поста пользователя")
				return
			}
			switch {
			// Получена ошибка что почта уже используется
			case discriptionMsg.Discription == "Email already exist":
				errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
				log.Printf("Не используется при удалении конкретного поста пользователя")
				return
				// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Не валидные данные при удалении конкретного поста пользователя")
				return
			case discriptionMsg.Discription == "Not Found Any Data":
				errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
				log.Printf("Не найдены данные об удаляемом посте пользователя")
				return
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получена не кастомная ошибка от сервиса forum-api при удалении конкретного поста пользователя")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получен статус-код не 202 или 500 от сервиса forum-api при удалении конкретного поста пользователя")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-client на удаление конкретного поста пользователя используется не верный метод")
		return
	}
}
