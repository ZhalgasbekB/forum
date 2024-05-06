package validation

import (
	"net/http"
	"strings"
)

type FormPost struct {
	Errors []string
}

func NewFormPost() *FormPost {
	return &FormPost{make([]string, 0, 7)}
}

func (f *FormPost) Valid() bool {
	return len(f.Errors) == 0
}

func (f *FormPost) EmptyFieldValid(val string) bool {
	return strings.TrimSpace(val) != ""
}

func (f *FormPost) CheckField(ok bool, err string) error {
	if !ok {
		f.Errors = append(f.Errors, err)
	}
	return nil
}

func ValidDatePost(r *http.Request) (bool, *FormPost) {
	var validDatePost bool
	form := NewFormPost()
	form.CheckField(form.EmptyFieldValid(r.FormValue("title")), "Пустой ввод темы!")
	form.CheckField(form.EmptyFieldValid(r.FormValue("description")), "Пустой ввод описания!")
	if form.Valid() {
		validDatePost = true
		return validDatePost, nil
	}
	return validDatePost, form
}
