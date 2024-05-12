package controllers

import (
	"module/models"
	"regexp"
	"strings"
	"unicode"
)

// IsValidEmail checks if the given email is in a valid format.
func IsValidEmail(email string) string {
	// Regular expression pattern for email validation.
	// This pattern allows for a wide range of valid email formats.
	// Modify it according to your specific requirements.
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex pattern.
	regex := regexp.MustCompile(pattern)

	// Check if the email matches the pattern.
	if regex.MatchString(email) {
		return ""
	}
	return "Invalid email format"
}

func RequiredField(field string) string {

	field = strings.ReplaceAll(field, " ", "")
	if len(field) > 0 {
		return ""
	}
	return "This field is required"
}

func VarifyPassword(password string) string {
	password = strings.ReplaceAll(password, " ", "")
	if len(password) < 1 {
		return "This field is required"
	}

	// Pssword must be at least 8 characters long, have 1 uppercase letter, 1 lowercase letter, 1 number and 1 special character
	letters := len(password)
	if letters < 8 {
		return "Password must have at least 8 characters"
	}

	number := false
	upper := false
	special := false

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		default:
		}
	}

	if !number {
		return "Password must have at least 1 number"
	}
	if !upper {
		return "Password must have at least 1 uppercase letter"
	}
	if !special {
		return "Password must have at least 1 special character"
	}

	return ""
}

func UserAlredyExists(email string) bool {
	usr, _ := models.GetUserByEmail(email)
	return usr != nil
}
