package auth

import (
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	
	"time"
	"fmt"
	"net/http"
	"strings"
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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:			"chirpy",
		IssuedAt:		jwt.NewNumericDate(time.Now()),
		ExpiresAt:		jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:		userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {return "", err}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {return uuid.Nil, err}
	
	user, err := token.Claims.GetSubject()
	if err != nil {return uuid.Nil, err}

	userID, err := uuid.Parse(user)
	if err != nil {return uuid.Nil, err}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))

	if token == "" {
		return "", fmt.Errorf("No token present")
	}

	return token, nil
}