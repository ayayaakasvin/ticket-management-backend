package helper

import (
	"fmt"
	"regexp"
)

type ValidationError string

func (v ValidationError) Error() string {
	return string(v)
}

var (
	uppercase = regexp.MustCompile(`[A-Z]`)
	lowercase = regexp.MustCompile(`[a-z]`)
	digit     = regexp.MustCompile(`[0-9]`)

	minLengthPassword = 8
	minLengthUsername = 3
)

func IsValidPassword(password string) error {
	if len(password) < minLengthPassword {
		return ValidationError(fmt.Sprintf("password must be at least %d characters", minLengthPassword))
	}

	if !uppercase.MatchString(password) {
		return ValidationError("password must contain at least one uppercase letter")
	}

	if !lowercase.MatchString(password) {
		return ValidationError("password must contain at least one lowercase letter")
	}

	if !digit.MatchString(password) {
		return ValidationError("password must contain at least one digit")
	}

	return nil
}

func IsValidUsername(username string) error {
	if len(username) < minLengthUsername {
		return ValidationError(fmt.Sprintf("username must be at least %d characters", minLengthUsername))
	}

	if !(lowercase.MatchString(username) || uppercase.MatchString(username)) {
		return ValidationError("username must contain at least one letter")
	}

	return nil
}
