package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(pswrd string) (string, error) {
	hashed_pswrd, err := argon2id.CreateHash(pswrd, argon2id.DefaultParams)
	if err != nil {return "", err}

	return hashed_pswrd, nil
}

func CheckPasswordHash(pswrd, hash string) (bool, error) {
	ok, err := argon2id.ComparePasswordAndHash(pswrd, hash)
	if err != nil {return false, err}

	return ok, nil
}