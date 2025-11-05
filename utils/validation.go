package utils

import (
	"context"
	"database/sql"
	"fmt"
	"product-catalogue/config"
	"regexp"
	"strings"
	"time"
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
	validName := regexp.MustCompile(`^[A-Za-z ]+$`)
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

func CategoryExists(categoryID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var exists bool
	err := config.DB.QueryRowContext(ctx, "select exists(select 1 from categories where category_id=?)", categoryID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err

	}
	return exists, nil

}

func SupplierValidate(
	name string,
	contact_info string,
	email string,
	company string,

) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("supplier name required")
	}
	if len(name) > 50 {
		return fmt.Errorf("supplier name too long!! please enter valid name ")
	}
	if len(name) < 2 {
		return fmt.Errorf("supplier name too short!! please enter valid name ")
	}
	validName := regexp.MustCompile(`^[A-Za-z ]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("supplier name should only contain alphabets and spaces")
	}

	hasLetter := false

	for _, ch := range name {
		if unicode.IsLetter(ch) {
			hasLetter = true
			break
		}

	}

	if !hasLetter {
		return fmt.Errorf("supplier name must contain atleast one letter")
	}

	if contact_info == "" {
		return fmt.Errorf("contact_info (phone) required")
	}

	validPhone := regexp.MustCompile(`^[0-9+\-() ]{7,15}$`)
	if !validPhone.MatchString(contact_info) {
		return fmt.Errorf("contact_info should only contain (+, - or spaces)")
	}

	if email == "" {
		return fmt.Errorf("please enter email")
	}

	if !IsValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	company = strings.TrimSpace(company)

	if company == "" {
		return fmt.Errorf("supplier company name required")
	}
	if len(company) > 50 {
		return fmt.Errorf("supplier company name too long. please enter valid name ")
	}
	if len(company) < 2 {
		return fmt.Errorf("supplier company name too short. please enter valid name ")
	}
	validCompany := regexp.MustCompile(`^[A-Za-z ]+$`)
	if !validCompany.MatchString(company) {
		return fmt.Errorf("supplier company name should only contain alphabets and spaces")
	}

	hasLetter2 := false

	for _, ch := range company {
		if unicode.IsLetter(ch) {
			hasLetter2 = true
			break
		}

	}

	if !hasLetter2 {
		return fmt.Errorf("supplier company name must contain atleast one letter")
	}

	return nil

}

func CategoryValidate(
	categoryName string,
	description string,
) error {
	categoryName = strings.TrimSpace(categoryName)
	if categoryName == "" {
		return fmt.Errorf("category name required")
	}
	if len(categoryName) > 50 {
		return fmt.Errorf("category name too long. please enter valid name ")
	}
	if len(categoryName) < 2 {
		return fmt.Errorf("category name too short. please enter valid name ")
	}
	validName := regexp.MustCompile(`^[A-Za-z ]+$`)
	if !validName.MatchString(categoryName) {
		return fmt.Errorf("category name should only contain alphabets and spaces")
	}

	hasLetter3 := false

	for _, ch := range categoryName {
		if unicode.IsLetter(ch) {
			hasLetter3 = true
			break
		}

	}

	if !hasLetter3 {
		return fmt.Errorf("category name must contain atleast one letter")
	}

	if description != "" {
		if len(description) > 1000 {
			return fmt.Errorf("description too long.. ")
		}
		if len(description) < 2 {
			return fmt.Errorf("description too short.. ")
		}
		validName := regexp.MustCompile(`^[A-Za-z0-9 ]+$`)
		if !validName.MatchString(description) {
			return fmt.Errorf("category name should only contain alphabets and spaces and numbers")
		}
		hasLetter4 := false

		for _, ch := range description {
			if unicode.IsLetter(ch) {
				hasLetter4 = true
				break
			}

		}

		if !hasLetter4 {
			return fmt.Errorf("description must contain atleast one letter")
		}

	}

	return nil

}
