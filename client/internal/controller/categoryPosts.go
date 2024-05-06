package controller

import (
	"fmt"
	"log"
	"lzhuk/clients/internal/cahe"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/errors"
	"lzhuk/clients/pkg/helpers"
	"net/http"
	"strings"
)

func CategoryPosts(w http.ResponseWriter, r *http.Request) {
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
		// Передаем на сервис сервера запрос на получение постов по конкретной категории
		newCategoryGet := fmt.Sprintf(categoryGet+"%s", r.FormValue("category"))

		req, err := http.NewRequest("GET", newCategoryGet, nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании GET запроса на сервис forum-api об получении постов по конкретно выбранной категории. Ошибка: %v", err)
			return
		}
		// Добавление из браузера куки в запрос на сервер
		req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
		req.Header.Set("Content-Type", "application/json")
		// Создаем структуру нового клиента
		client := http.Client{}
		// Отправляем запрос на сервер
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса на сервис forum-api для получения постов по конкретно выбранной категории. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()
		// Конвертируем полученные данные об постах по конкретной категории
		result, err := convertor.ConvertCategoryPosts(resp)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации данных полученных из ответа от сервиса forum-api об постах по конкретно выбранной категории. Ошибка: %v", err)
			return
		}

		// Вносим в кеш результаты по выбранной категории
		cahe.CategoryPosts = result
		switch resp.StatusCode {
		case http.StatusOK:
			http.Redirect(w, r, "http://localhost:8082/userd3", 302)
			return
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и её описания от сервиса forum-api на запрос об постах по конкретно выбранной категории. Ошибка: %v", err)
				return
			}
			switch {
			// Получена ошибка что почта уже используется
			case discriptionMsg.Discription == "Email already exist":
				errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
				log.Printf("Не используется при получении постов по конкретно выбранной категории")
				return
				// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Не валидные данные при запросе на получение постов по конкретно выбранной категории")
				return
			case discriptionMsg.Discription == "Not Found Any Data":
				errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
				log.Printf("Запрашиваемые данные не найдены при запросе на получение постов по конкретно выбранной категории")
				return
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получена не кастомная ошибка от сервиса forum-api при запросе на получение постов по конекретно выбранной категории")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при запросе на получение постов по конкретно выбранной категории")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-client на получение постов по конкретно выбранной категории используется не верный метод")
		return
	}
}
