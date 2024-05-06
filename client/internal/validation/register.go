package validation

import (
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

var (
	nameRegex  = regexp.MustCompile("^[a-zA-Zа-яА-Я0-9_]{3,16}$")
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Form struct {
	Errors []string
}

func New() *Form {
	return &Form{make([]string, 0, 7)}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Проверка на валидность почты
func (f *Form) EmailValid(email string) bool {
	return emailRegex.Match([]byte(email))
}

// Проверка на валидность имени пользователя
func (f *Form) NameValid(name string) bool {
	return nameRegex.Match([]byte(name))
}

func (f *Form) PasswordValid(password string) bool {
	// Проверка длины пароля не менее 8 символов
	if len(password) < 8 {
		return false
	}

	var (
		hasUpperCase bool
		hasLowerCase bool
		hasDigit     bool
		hasSpecial   bool
		hasAscii     bool
	)

	for _, char := range password {
		// Проверка на ASCII
		if char > 32 && char < 127 {
			hasAscii = true
		}
		switch {
		// Поиск на заглавную букву
		case unicode.IsUpper(char):
			hasUpperCase = true
			// Поиск на прописную букву
		case unicode.IsLower(char):
			hasLowerCase = true
			// Поиск на число
		case unicode.IsDigit(char):
			hasDigit = true
			// Поиск на знак пунктуации или символ
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true

		}
	}
	// Должен содержать как минимум 8 символов, включая по крайней мере одну заглавную букву,
	// одну строчную букву, одну цифру и один специальный символ. Пароль должен содержать
	// только английские символы
	return hasUpperCase && hasLowerCase && hasDigit && hasSpecial && hasAscii
}

func (f *Form) EmptyFieldValid(val string) bool {
	return strings.TrimSpace(val) != ""
}

func (f *Form) CheckField(ok bool, err string) error {
	if !ok {
		f.Errors = append(f.Errors, err)
	}
	return nil
}

func ValidDate(r *http.Request) (bool, *Form) {
	var validDate bool
	form := New()
	form.CheckField(form.EmptyFieldValid(r.FormValue("email")), "Пустой ввод почты!")
	form.CheckField(form.EmptyFieldValid(r.FormValue("name")), "Пустой ввод имени!")
	form.CheckField(form.EmptyFieldValid(r.FormValue("password")), "Пустой ввод пароля!")
	form.CheckField(form.EmailValid(r.FormValue("email")), "Проверьте корректность ввода почты!")
	form.CheckField(form.NameValid(r.FormValue("name")), "Имя должно состоять из русских или английских букв и может содержать символ подчеркивания!")
	form.CheckField(form.PasswordValid(r.FormValue("password")), "Пароль должен содержать только английские буквы, минимум 8 символов, 1 заглавную и 1 строчную буквы, 1 цифру и 1 спец.символ!")
	if form.Valid() {
		validDate = true
		return validDate, nil
	}
	return validDate, form
}
