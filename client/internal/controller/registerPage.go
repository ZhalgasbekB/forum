package controller

import (
	"bytes"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/internal/validation"
	"lzhuk/clients/pkg/errors"
	"net/http"
)

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	// Создание шаблона для страницы регистрации
	t, err := template.ParseFiles("./ui/html/register.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка создания шаблона страницы регистрации пользователя. Ошибка: %v", err)
		return
	}
	// Проверка метода запроса
	switch r.Method {
	case http.MethodGet:
		err = t.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при рендеринге шаблона страницы регистрации пользователя. Ошибка: %v", err)
			return
		}
	case http.MethodPost:
		// Проверка на валидность пользовательских данных
		validDate, errValid := validation.ValidDate(r)
		// Рендеринг страницы при невалидных данных регистрации пользователем
		if validDate == false {
			w.WriteHeader(http.StatusBadRequest)
			err = t.ExecuteTemplate(w, "register.html", &errValid)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при рендеринге шаблона страницы регистрации пользователя при проверке данных на валидность. Ошибка: %v", err)
				return
			}
		} else {
			// Конвертор полученных пользовательских данных из формы html в формат JSON
			jsonData, err := convertor.ConvertRegister(r)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных о регистрации пользователя в JSON. Ошибка: %v", err)
				return
			}
			// Формирования POST запроса на регистрацию нового пользователя на сервере
			req, err := http.NewRequest("POST", registry, bytes.NewBuffer(jsonData))
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при формировании POST запроса на регистрацию пользователя. Ошибка: %v", err)
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
				log.Printf("Произошла ошибка при передаче запроса от клиента к сервису forum-api при регистрации нового пользователя. Ошибка: %v", err)
				return
			}
			defer resp.Body.Close()
			// Проверка кода статуса ответа сервера
			switch resp.StatusCode {
			// Получен статус код 201 об успешной регистрации пользователя в системе
			case http.StatusCreated:
				http.Redirect(w, r, "http://localhost:8082/login", 303)
				// Получен статус код 500 об ошибке на сервере
			case http.StatusInternalServerError:
				discriptionMsg, err := convertor.DecodeErrorResponse(resp)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при декодировании ответа ошибки и ее описания от сервиса forum-api на запрос об регистрации нового пользователя. Ошибка: %v", err)
					return
				}
				switch {
				// Получена ошибка что почта уже используется
				case discriptionMsg.Discription == "Email already exist":
					errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
					log.Printf("Пользователь пытается зарегестировать почту которая уже используется")
					return
					// Получена ошибка что введены неверные учетные данные
				case discriptionMsg.Discription == "Invalid Credentials":
					errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
					log.Printf("Пользователь ввел не валидные данные при регистрации")
					return
				case discriptionMsg.Discription == "Not Found Any Data":
					errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
					log.Printf("Нет запрашиваемых данных")
					return
				case discriptionMsg.Discription == "Name already exist":
					errorPage(w, errors.NotName, http.StatusBadRequest)
					log.Printf("Пользователь пытается зарегестрировать имя которое уже имеется")
					return
				default:
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Получена не кастомная ошибка от сервиса forum-api при регистрации пользователя")
					return
				}
				// Получен статус код 405 об неверном методе запроса с сервера
			case http.StatusMethodNotAllowed:
				errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
				log.Printf("При передаче запроса сервису forum-api на регистрацию нового пользователя используется не верный метод")
				return
				// Получен статус код не 201, 405 или 500
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получен статус-код не 201, 405 и 500 от сервиса forum-api при регистрации нового пользователя")
				return
			}
		}
		// Метод запроса с браузера не POST и не GET
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса сервису forum-client на регистрацию нового пользователя используется не верный метод")
		return
	}
}
