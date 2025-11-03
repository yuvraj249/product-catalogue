package utils

import (
	"fmt"
	"regexp"
	"strings"
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

func ProductValidate(
	name string,
	description string,
	cost int,
	categoryID *int,
	supplierID *int,
	discountType *string,
	discountValue *float64,
) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("product name required")
	}
	if len(name) > 50 {
		return fmt.Errorf("product name too long!! please enter valid name ")
	}
	if len(name) < 2 {
		return fmt.Errorf("product name too short!! please enter valid name ")
	}
	validName := regexp.MustCompile(`^[A-Za-z ] +$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("product name should only contain alphabets and spaces")
	}

	hasLetter := false

	for _, ch := range name {
		if unicode.IsLetter(ch) {
			hasLetter = true
			break
		}

	}

	if !hasLetter {
		return fmt.Errorf("product name must contain atleast one letter")
	}

	if cost <= 0 || cost > 9999999999 {
		return fmt.Errorf("please enter realistic cost for product ")
	}

	if strings.TrimSpace(description) != "" && len(description) > 1500 {
		return fmt.Errorf("description too long ")
	}

	if categoryID != nil && *categoryID <= 0 {
		return fmt.Errorf("invalid category id ")
	}

	if supplierID != nil && *supplierID <= 0 {
		return fmt.Errorf("invalid supplier id")
	}

	if discountType != nil {
		disc := strings.ToLower(strings.TrimSpace(*discountType))
		if disc != "flat" && disc != "percent" {
			return fmt.Errorf("discountType must be flat or percent")
		}
		if discountValue == nil {
			return fmt.Errorf("discountValue required if discountType selected")
		}
		if *discountValue < 0 {
			return fmt.Errorf("discountValue must be >= 0")
		}
		if disc == "percent" && *discountValue > 100 {
			return fmt.Errorf("discountValue must be >= 100")
		}
		if disc == "flat" && *discountValue > float64(cost) {
			return fmt.Errorf("discountValue cannot exceed product value")
		}
	}

	return nil

}
