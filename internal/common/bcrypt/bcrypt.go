package bcrypt

import "golang.org/x/crypto/bcrypt"

func Hash(text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	return string(bytes), err
}

func Compare(hash string, text string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text)); err != nil {
		return false
	}
	return true
}
