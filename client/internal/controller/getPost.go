package controller

import (
	"fmt"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/errors"
	"lzhuk/clients/pkg/helpers"
	"net/http"
	"strings"
)

func GetPost(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"):
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}
	// Создаем шаблон страницы с конкретным постом
	t, err := template.ParseFiles("./ui/html/post.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при создании шаблона страницы с конкретным постом. Ошибка: %v", err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Разбиваем путь URL на срез по признаку слеша
		parts := strings.Split(r.URL.Path, "/")
		// Формируем URL запроса на сервис сервера с конкретным id поста
		getUserPostId := fmt.Sprintf(getUserPost+"%s", parts[len(parts)-1])
		// Формируем запрос
		req, err := http.NewRequest("GET", getUserPostId, nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании GET запроса на получение данных конкретного поста. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос к серверу
		// req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
		// Создаем структуру нового клиента
		client := http.Client{}
		// Передаем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api на получение данных о конкретном посте. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()
		switch resp.StatusCode {
		case http.StatusOK:
			result, err := convertor.ConvertGetPosts(r, resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных из ответа сервиса forum-api на получение данных о конкретном посте. Ошибка: %v", err)
				return
			}
			err = t.ExecuteTemplate(w, "post.html", result)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при рендеренге страницы с конкретным постом. Ошибка: %v", err)
				return
			}
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и её описания от сервиса forum-api на запрос об получении данных конкретного поста")
				return
			}
			switch {
			// Получена ошибка что почта уже используется
			case discriptionMsg.Discription == "Email already exist":
				errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
				log.Printf("Не используется для получения конкретного поста")
				return
				// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Не используется для получения конкретного поста")
				return
			case discriptionMsg.Discription == "Not Found Any Data":
				errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
				log.Printf("Нет запрашиваемых данных о конкретном посте")
				return
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получена не кастомная ошибка от сервиса forum-api при получении данных о конкретном посте 1")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при получении данных о конкретно посте")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-client на получение данных о конкретном посте используется не верный метод")
		return
	}
}
