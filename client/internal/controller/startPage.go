package controller

import (
	"html/template"
	"log"
	"lzhuk/clients/pkg/errors"
	"net/http"
)

func StartPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("Неверный метод в запросе получения стартовой страницы")
		return
	}

	t, err := template.ParseFiles("./ui/html/start.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при создании шаблона стартовой страницы")
		return
	}

	err = t.ExecuteTemplate(w, "start.html", nil)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при рендеринге стартовой страницы")
		return
	}
}
