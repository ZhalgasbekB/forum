package validation

import (
	"net/http"
	"strings"
)

type FormCommentUpdate struct {
	Errors []string
}

func NewFormCommentUpdate() *FormCommentUpdate {
	return &FormCommentUpdate{make([]string, 0, 7)}
}

func (f *FormCommentUpdate) Valid() bool {
	return len(f.Errors) == 0
}

func (f *FormCommentUpdate) EmptyFieldValid(val string) bool {
	return strings.TrimSpace(val) != ""
}

func (f *FormCommentUpdate) CheckField(ok bool, err string) error {
	if !ok {
		f.Errors = append(f.Errors, err)
	}
	return nil
}

func ValidDateCommentUpdate(r *http.Request) (bool, *FormCommentUpdate) {
	var validDatePost bool
	form := NewFormCommentUpdate()
	form.CheckField(form.EmptyFieldValid(r.FormValue("updatedComment")), "Пустой ввод комментария!")
	if form.Valid() {
		validDatePost = true
		return validDatePost, nil
	}
	return validDatePost, form
}
