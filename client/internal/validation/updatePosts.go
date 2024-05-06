package validation

import (
	"net/http"
	"strings"
)

type FormPostUpdate struct {
	Errors []string
}

func NewFormPostUpdate() *FormPostUpdate {
	return &FormPostUpdate{make([]string, 0, 7)}
}

func (f *FormPostUpdate) Valid() bool {
	return len(f.Errors) == 0
}

func (f *FormPostUpdate) EmptyFieldValid(val string) bool {
	return strings.TrimSpace(val) != ""
}

func (f *FormPostUpdate) CheckField(ok bool, err string) error {
	if !ok {
		f.Errors = append(f.Errors, err)
	}
	return nil
}

func ValidDatePostUpdate(r *http.Request) (bool, *FormPostUpdate) {
	var validDatePost bool
	form := NewFormPostUpdate()
	form.CheckField(form.EmptyFieldValid(r.FormValue("title")), "Пустой ввод темы!")
	form.CheckField(form.EmptyFieldValid(r.FormValue("description")), "Пустой ввод описания!")
	if form.Valid() {
		validDatePost = true
		return validDatePost, nil
	}
	return validDatePost, form
}
