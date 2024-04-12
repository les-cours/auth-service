package utils

func ValidateLoginInput(email, password, os string) bool {
	return len(email) == 0 || len(password) == 0
}
