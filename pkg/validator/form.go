package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func CreatePost(title, content, image string) map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	if strings.TrimSpace(image) == "" {
		errors["image"] = "This field cannot be blank"
	}

	return errors
}

func Signup(name, email, password string) map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(name) == "" {
		errors["name"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(name) < 3 {
		errors["name"] = "This field is too short (minimum is 3 characters)"
	}

	if strings.TrimSpace(email) == "" {
		errors["email"] = "This field cannot be blank"
	} else if !EmailRX.MatchString(email) {
		errors["email"] = "This field is invalid"
	}

	if strings.TrimSpace(password) == "" {
		errors["password"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(password) < 6 {
		errors["password"] = "This field is too short (minimum is 6 characters)"
	}

	return errors
}

func Login() map[string]string {
	errors := make(map[string]string)
	errors["generic"] = ""
	return errors
}

func EditPost(title, content, image string) map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	return errors
}

func ContactForm(name, email, subject, message string) map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(name) == "" {
		errors["name"] = "Это поле не может быть пустым"
	} else if utf8.RuneCountInString(name) < 3 {
		errors["name"] = "Это поле слишком короткое (минимум 3 символа)"
	}

	if strings.TrimSpace(email) == "" {
		errors["email"] = "Это поле не может быть пустым"
	} else if !EmailRX.MatchString(email) {
		errors["email"] = "Это поле недопустимо"
	}

	if strings.TrimSpace(subject) == "" {
		errors["password"] = "Это поле не может быть пустым"
	}

	if strings.TrimSpace(message) == "" {
		errors["password"] = "Это поле не может быть пустым"
	} 

	return errors
}