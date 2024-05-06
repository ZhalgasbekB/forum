package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/helpers"
	"lzhuk/clients/internal/validation"
	"lzhuk/clients/pkg/errors"
	"net/http"
	"strings"
)

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"):
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}
	// Создаем шаблон страницы с постом для его изменения
	t, err := template.ParseFiles("./ui/html/update_post.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при создании шаблона страницы с постом для его изменения. Ошибка: %v", err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		// Разбиваем путь URL на срез по признаку слеша
		parts := strings.Split(r.URL.Path, "/")
		// Формируем URL запроса на сервис сервера с конкретным id поста
		getUserPostId := fmt.Sprintf(getUserPost+"%s", parts[len(parts)-1])
		// Формируем запрос на сервис сервера
		req, err := http.NewRequest("GET", getUserPostId, nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при создании GET запроса на сервис forum-api для изменения поста. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос к серверу
		req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
		// Создаем структуру нового клиента
		client := http.Client{}
		// Передаем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса на сервис forum-api для изменения поста. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			// Конвертируем данные о посте из БД который пользователь будет редактировать
			result, err := convertor.ConvertGetPosts(r, resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных из ответа сервиса forum-api который пользователь будет редактировать. Ошибка: %v", err)
				return
			}

			err = t.ExecuteTemplate(w, "update_post.html", result)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при рендеринге страницы с данными о посте который пользователь будет редактировать. Ошибка: %v", err)
				return
			}
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и её описания от сервиса forum-api на запрос об данных о посте который пользователь будет редактировать")
				return
			}
			switch {
			// Получена ошибка что почта уже используется
			case discriptionMsg.Discription == "Email already exist":
				errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
				log.Printf("Не используется при изменении поста")
				return
				// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Не валидные данные при изменении поста")
				return
			case discriptionMsg.Discription == "Not Found Any Data":
				errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
				log.Printf("Нет запрашиваемых данных при изменении поста")
				return
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Получена не кастомная ошибка от сервиса forum-api при изменении поста")
					return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получен статус-код не 200 и 500 от сервиса forum-api при изменении поста")
				return
		}
	case http.MethodPost:
		// Проверка на валидность пользовательских данных
		validDatePost, _ := validation.ValidDatePostUpdate(r)
		if validDatePost == false {
			errorPage(w, errors.EmptyDatePost, http.StatusBadRequest)
			return
		} else {
			// Конвертируем данные с изменениями поста от пользователя для передачи на сервер
			jsonData, err := convertor.ConvertUpdatePost(r)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных c изменениями внесенными пользователем о посте для передачи на сервис forum-api. Ошибка: %v", err)
				return
			}
			// Формирование URL  c конкретным id поста для редактирования
			updatePostsId := fmt.Sprintf(updatePosts+"%s", r.FormValue("postId"))
			// Формируем запрос на измение данных в посте
			req, err := http.NewRequest("PUT", updatePostsId, bytes.NewBuffer(jsonData))
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при формировании PUT запроса об изменениях данных в посте. Ошибка: %v", err)
				return
			}
			// Записываем куки из браузера в запрос на сервер
			req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
			// Указываем формат передачи данных
			req.Header.Set("Content-Type", "application/json")
			// Создаем структуру нового клиента
			client := http.Client{}
			// Передаем запрос на сервис сервера
			resp, err := client.Do(req)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при передаче запроса на сервис forum-api об изменениях данных в посте. Ошибка: %v", err)
				return
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusAccepted:
				// Формируем URL запроса на сервис сервера с конкретным id поста
				newReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8082/userd3/post/%s", r.FormValue("postId")), nil)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при формировании GET запроса на сервис forum-api об изменениях данных в посте. Ошибка: %v", err)
					return
				}
				http.Redirect(w, r, newReq.URL.String(), 302)
			case http.StatusInternalServerError:
				discriptionMsg, err := convertor.DecodeErrorResponse(resp)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при декодировании ответа ошибки и её описания от сервиса forum-api на запрос об данных о посте который пользователь будет редактировать")
					return
				}
				switch {
				// Получена ошибка что почта уже используется
				case discriptionMsg.Discription == "Email already exist":
					errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
					log.Printf("Не используется при изменении поста")
					return
					// Получена ошибка что введены неверные учетные данные
				case discriptionMsg.Discription == "Invalid Credentials":
					errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
					log.Printf("Не валидные данные при изменении поста")
					return
				case discriptionMsg.Discription == "Not Found Any Data":
					errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
					log.Printf("Нет запрашиваемых данных при изменении поста")
					return
				default:
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
						log.Printf("Получена не кастомная ошибка от сервиса forum-api при внесении изменений поста")
						return
				}
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получен статус-код не 202 и 500 от сервиса forum-api при внесении изменений поста")
				return
			}
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-client на изменение поста используется не верный метод")
		return
	}
}
