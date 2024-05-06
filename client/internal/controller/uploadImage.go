package controller

import (
	"io/ioutil"
	"log"
	"lzhuk/clients/pkg/errors"
	"net/http"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
	}
	file, handler, err := r.FormFile("image")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при получении данных об картинке из формы запроса. Ошибка: %v", err)
		return
	}
	defer file.Close()

	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при чтении отправленной картинки пользователем. Ошибка: %v", err)
		return
	}

	err = ioutil.WriteFile("./uploads/"+handler.Filename, imageData, 0666)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при записи на сервер картинки. Ошибка: %v", err)
		return
	}
}
