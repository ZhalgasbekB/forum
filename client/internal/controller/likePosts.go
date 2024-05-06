package controller

import (
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/helpers"
	"lzhuk/clients/pkg/errors"
	"net/http"
	"strings"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"):
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}
	// Создаем шаблон страницы понравившимися постами пользователю
	t, err := template.ParseFiles("./ui/html/user_likeposts.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при создании шаблона страницы с понравившимися постами пользователя. Ошибка: %v", err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		// Формируем запрос на сервис сервера
		req, err := http.NewRequest("GET", likePosts, nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании GET запроса для получения понравившихся постов пользователя. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос на сервис сервера
		req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
		// Формируем структуру нового клиента
		client := http.Client{}
		// Передаем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса для получения понравившихся постов пользователя. Ошибка: %v", err)
			return
		}
		switch resp.StatusCode {
		case http.StatusOK:
			// Конвертация данных для получения всех понравившихся постов пользователя
			userLikePosts, err := convertor.ConvertUserLikePosts(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных для получения понравившихся постов пользователя. Ошибка: %v", err)
				return
			}
			// Рендеринг страницы с понравившимися постами пользователя
			err = t.ExecuteTemplate(w, "user_likeposts.html", userLikePosts)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при рендеринге страницы понравившихся постов пользователя. Ошибка: %v", err)
				return
			}
			defer resp.Body.Close()
			return
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и ёё описания от сервиса forum-api на запрос об данных постов понравившихся пользователю. Ошибка: %v", err)
				return
			}
			switch {
			// Получена ошибка что почта уже используется
			case discriptionMsg.Discription == "Email already exist":
				errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
				log.Printf("Не используются при получении данных об понравившихся постов")
				return
				// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Не используются при получении данных об понравившихся постов")
				return
			case discriptionMsg.Discription == "Not Found Any Data":
				errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
				log.Printf("Нет запрашиваемых данных об понравившихся постах")
				return
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получена не кастомная ошибка от сервиса forum-api на запрос об данных понравившихся пользователю постов")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при запросе об данных понравившихся пользователю постов")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-client на получение данных о всех понравившихся постах пользователю используется не верный метод")
		return
	}
}
