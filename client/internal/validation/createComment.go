package validation

import (
	"net/http"
	"strings"
)

type FormComment struct {
	Errors []string
}

func NewFormComment() *FormComment {
	return &FormComment{make([]string, 0, 7)}
}

func (f *FormComment) Valid() bool {
	return len(f.Errors) == 0
}

func (f *FormComment) EmptyFieldValid(val string) bool {
	return strings.TrimSpace(val) != ""
}

func (f *FormComment) CheckField(ok bool, err string) error {
	if !ok {
		f.Errors = append(f.Errors, err)
	}
	return nil
}

func ValidDateComment(r *http.Request) (bool, *FormComment) {
	var validDatePost bool
	form := NewFormComment()
	form.CheckField(form.EmptyFieldValid(r.FormValue("comment")), "Пустой ввод комментария!")
	if form.Valid() {
		validDatePost = true
		return validDatePost, nil
	}
	return validDatePost, form
}
