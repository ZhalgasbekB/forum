package controller

import (
	"bytes"
	"fmt"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/helpers"
	"lzhuk/clients/pkg/errors"
	"net/http"
	"strings"
)

func VoteComment(w http.ResponseWriter, r *http.Request) {
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
		// Конвертируем данные полуения голоса комментария
		jsonData, err := convertor.ConvertVoteComment(r)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации данных в JSON о голосе поставленного на комменатрий. Ошибка: %v", err)
			return
		}
		// Формируем запрос на сервис сервера
		req, err := http.NewRequest("POST", voteComments, bytes.NewBuffer(jsonData))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании POST запроса на сервис forum-api для передачи данных о голосе поставленного на комменатрий. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос на сервис сервера
		req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
		req.Header.Set("Content-Type", "application/json")
		// Создаем структуру нового клиента
		client := http.Client{}
		// Отправляем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса на сервис forun-api с данными о голосе поставленного на комменатрий. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			newReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8082/userd3/post/%s", r.FormValue("postId")), nil)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при формировании GET запроса на получение поста с комментарием где поставлен был голос. Ошибка: %v", err)
				return
			}
			http.Redirect(w, r, newReq.URL.String(), 302)
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос об рроставлении голоса на комментарий. Ошибка: %v", err)
				return
			}
			switch {
			// Получена ошибка что почта уже используется
			case discriptionMsg.Discription == "Email already exist":
				errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
				log.Printf("Не используется при постановке голоса на комментарий")
				return
				// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Не валидные данные при постановке голоса на комментарий")
				return
			case discriptionMsg.Discription == "Not Found Any Data":
				errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
				log.Printf("Не используется при постановке голоса на комментари")
				return
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Получена не кастомная ошибка от сервиса forum-api при постановке голоса на комментарий")
					return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при постановке голоса на комментарий")
				return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-client на постановку голоса используется не верный метод")
		return
	}
}
