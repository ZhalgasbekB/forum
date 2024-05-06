package helpers

import "net/http"

func CheckCookie(c []*http.Cookie) bool {
	for _, v := range c {
		if v.Name == "CookieUUID" {
			return true
		}
	}
	return false
}

func CheckCookieIndex(c []*http.Cookie) int {
	for i, v := range c {
		if v.Name == "CookieUUID" {
			return i
		}
	}
	return len(c)
}
