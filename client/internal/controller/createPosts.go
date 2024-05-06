package controller

import (
	"bytes"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/internal/validation"
	"lzhuk/clients/pkg/errors"
	"lzhuk/clients/pkg/helpers"
	"net/http"
	"strings"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"):
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}
	// Создание шаблона для страницы создания поста
	t, err := template.ParseFiles("./ui/html/create_post.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка создания шаблона страницы для создания нового поста пользователем. Ошибка: %v", err)
		return
	}
	// Проверка метода запроса
	switch r.Method {
	case http.MethodGet:

		err = t.ExecuteTemplate(w, "create_post.html", nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при рендеринге страницы создания нового поста пользователем. Ошибка: %v", err)
			return
		}
	case http.MethodPost:
		// Проверка на валидность пользовательских данных
		validDatePost, _ := validation.ValidDatePost(r)
		if validDatePost == false {
			errorPage(w, errors.EmptyDatePost, http.StatusBadRequest)
			log.Printf("Произошла ошибка при рендеринге шаблона страницы создания нового поста пользователем при проверке на валидность данных. Ошибка: %v", err)
			return

		} else {
			// Конвертация данных при создании нового поста
			jsonData, err := convertor.ConvertCreatePost(r)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных в JSON для передачи на сервис forum-api при создании нового поста пользоваталем. Ошибка: %v", err)
				return
			}
			// Создание POST запроса на внесение информации о новом посте
			req, err := http.NewRequest("POST", createPosts, bytes.NewBuffer(jsonData))
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при формировании POST запроса на сервис forum-api при создании нового поста пользователем. Ошибка: %v", err)
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
				log.Printf("Произошла ошибка при передаче запроса на сервис forum-api при создании нового поста пользоваталем. Ошибка: %v", err)
				return
			}
			defer resp.Body.Close()
			switch resp.StatusCode {
			case http.StatusCreated:
				// Проверка на наличие загруженной картинки
				_, check, _ := r.FormFile("image")
				if check == nil {
					http.Redirect(w, r, "http://localhost:8082/userd3", 302)
					return
				}
				postId, err := convertor.ConvertPostId(*resp)
				// Проверка и конвертация данны о загруженной картинке к посту
				jsonDataImage, err := convertor.ConvertCreatePostImage(r, postId)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusBadRequest)
					log.Printf("Произошла ошибка при конвертации данных об картинке загруженной пользоваталем при создании нового поста: %v", err)
					return
				}
				// Если имеются данные об картинке
				if jsonDataImage != nil && err == nil {
					// Создание POST запроса на внесение информации о картинке к новому посту
					req, err := http.NewRequest("POST", UploadImagePosts, bytes.NewBuffer(jsonDataImage))
					if err != nil {
						errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
						log.Printf("Произошла ошибка при формировании POST запроса на сервис forum-api при внесении информации об картинке к новому посту. Ошибка: %v", err)
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
						log.Printf("Произошла ошибка при передаче запроса на сервис forum-api при внесении информации об картинке к новому посту Ошибка: %v", err)
						return
					}
					defer resp.Body.Close()
				}
				http.Redirect(w, r, "http://localhost:8082/userd3", 302)
				return
			case http.StatusInternalServerError:
				discriptionMsg, err := convertor.DecodeErrorResponse(resp)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при декодировании ответа ошибки и её описания от сервиса forum-api на запрос о создании нового поста пользователя. Ошибка: %v", err)
					return
				}
				switch {
				// Получена ошибка что почта уже используется
				case discriptionMsg.Discription == "Email already exist":
					errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
					log.Printf("Не используется для создания нового поста")
					return
					// Получена ошибка что введены неверные учетные данные
				case discriptionMsg.Discription == "Invalid Credentials":
					errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
					log.Printf("Не валидные данные при создании нового поста")
					return
				case discriptionMsg.Discription == "Not Found Any Data":
					errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
					log.Printf("Не используется при создании нового поста")
					return
				default:
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Получена не кастомная ошибка от сервиса forum-api при создании нового поста пользователем")
					return
				}
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получен статус-код не 201 или 500 от сервиса forum-api при создании нового поста пользователем")
				return
			}
		}
		// Метод запроса с браузера не POST и не GET
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-client на создание нового поста пользователем используется не верный метод")
		return
	}
}
