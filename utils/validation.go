package utils

import (
	"regexp"
	"unicode"
)

func IsValidEmail(email string) bool {
	email_check := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return email_check.MatchString(email)
}

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasSpecial := false
	hasNumber := false

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		case unicode.IsNumber(ch):
			hasNumber = true
		}

	}

	return hasUpper && hasLower && hasSpecial && hasNumber

}
