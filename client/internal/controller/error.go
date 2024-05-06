package controller

import (
	"html/template"
	"log"
	"net/http"
)

// Функция обработчика ошибок
func errorPage(w http.ResponseWriter, statusText string, statusCode int) {
	w.WriteHeader(statusCode)

	data := struct {
		StatusText string
		StatusCode int
	}{
		statusText,
		statusCode,
	}

	t, err := template.ParseFiles("./ui/html/error.html")
	if err != nil {
		http.Error(w, "Произошла ошибка на сервере, зайдите на форум позже", http.StatusInternalServerError)
		log.Printf("Произошла ошибка создании шаблона страницы ошибок. Ошибка: %v", err)
		return
	}
	err = t.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		log.Printf("Произошла ошибка при рендеринге страницы ошибок. Ошибка: %v", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusInternalServerError)
		return
	}
}
