package models

import (
	"net/mail"
	"regexp"
	"unicode"

	customErrors "github.com/lsmltesting/MicroBlog/internal/errors"
)

var usernameRegex = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`)

type User struct {
	Username string
	Email    string
	Password string
	ID       int
}

func NewUser(username string, email string, password string) (*User, error) {
	user := &User{}

	// Set username for user after validating
	if err := user.SetUsername(username); err != nil {
		return nil, err
	}

	// Set email for user after validating
	if err := user.SetUserEmail(email); err != nil {
		return nil, err
	}

	// Set password for user after validating
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

func (user *User) SetUsername(username string) error {
	if usernameRegex.MatchString(username) {
		user.Username = username
		return nil
	}
	return customErrors.ErrWrongUserName
}

func (user *User) SetUserEmail(email string) error {
	addr, err := mail.ParseAddress(email)

	if err != nil {
		return customErrors.ErrWrongEmailAddress
	}

	user.Email = addr.Address
	return nil
}

func (user *User) SetPassword(password string) error {
	if len(password) < 8 {
		return customErrors.ErrInvalidPassword
	}

	var (
		hasLower = false
		hasDigit = false
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasLower || !hasDigit {
		return customErrors.ErrInvalidPassword
	}

	user.Password = password
	return nil
}
