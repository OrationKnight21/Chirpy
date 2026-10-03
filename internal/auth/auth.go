package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    "chirpy-access",
		Subject:   userID.String(),
	})
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}
func ValidateJWT(tokenstring, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.UUID{}, errors.New("invalid type of claims")
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header included in request")
	}
	authSplit := strings.Split(authHeader, " ")
	if len(authSplit) < 2 || authSplit[0] != "Bearer" {
		return "", errors.New("malformed authorization Header")
	}
	return authSplit[1], nil
}
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header included")
	}
	authSplit := strings.Split(authHeader, " ")
	if len(authSplit) < 2 || authSplit[0] != "ApiKey" {
		return "", errors.New("malformed authorization Header")
	}
	return authSplit[1], nil
}
