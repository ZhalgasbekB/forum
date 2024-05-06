package controller

import (
	"bytes"
	"fmt"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/helpers"
	"lzhuk/clients/internal/validation"
	"lzhuk/clients/pkg/errors"
	"net/http"
	"strings"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"):
		fmt.Println(strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"))
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}
	switch r.Method {
	case http.MethodPost:
		// Проверка на валидность данных комментария
		validDateComment, _ := validation.ValidDateComment(r)
		if validDateComment == false {
			errorPage(w, errors.EmptyComments, http.StatusBadRequest)
			log.Printf("Пользователь ввел пустой комментарий")
		} else {
			// Конвертируем данные при создании нового комментария
			jsonData, err := convertor.ConvertCreateComment(r)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка конвертации данных в JSON при создании нового комментария пользователем. Ошибка: %v", err)
				return
			}
			// Формирование нового запроса
			req, err := http.NewRequest("POST", createComments, bytes.NewBuffer(jsonData))
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при формирование POST запроса на сервис forum-api для создания нового комментария пользователя. Ошибка: %v", err)
				return
			}
			// Запись куки из браузера в запрос для сервера
			req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
			// Записываем тип передаваемого контента
			req.Header.Set("Content-Type", "application/json")
			// Создаем структуру нового клиента
			client := http.Client{}
			// Передаем запрос на сервер
			resp, err := client.Do(req)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при отправке запроса на сервис forum-api для создания нового комментария пользователя. Ошибка: %v", err)
				return
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusCreated:
				newReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8082/userd3/post/%s", r.FormValue("postId")), nil)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при формировании GET запроса на переход к странице с постом где пользователь оставил комментарий. Ошибка: %v", err)
					return
				}
				http.Redirect(w, r, newReq.URL.String(), 302)
			case http.StatusInternalServerError:
				discriptionMsg, err := convertor.DecodeErrorResponse(resp)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при декодировании ответа ошибки и её описания от сервиса forum-api на запрос об создании нового комментария пользователя. Ошибка: %v", err)
					return
				}
				switch {
				// Получена ошибка что почта уже используется
				case discriptionMsg.Discription == "Email already exist":
					errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
					log.Printf("Не используется при создании нового комментария")
					return
					// Получена ошибка что введены неверные учетные данные
				case discriptionMsg.Discription == "Invalid Credentials":
					errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
					log.Printf("Не валидные данные при создании нового комментария")
					return
				case discriptionMsg.Discription == "Not Found Any Data":
					errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
					log.Printf("Не используется при создании нового комментария")
					return
				default:
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Получена не кастомная ошибка от сервиса forum-api при создании нового комментария")
					return
				}
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получен статус-код не 201 или 500 от сервиса forum-api при создании нового комментария")
				return
			}
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-client на создание нового комментария используется не верный метод")
		return
	}
}
