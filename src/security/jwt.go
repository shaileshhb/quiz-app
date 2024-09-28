package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shaileshhb/quiz/src/db/models"
)

var jwtKey = "This should be in .env"

// GenerateJWT will generate a JWT token for user login
func GenerateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
	})
	return token.SignedString([]byte(jwtKey))
}

// ValidateJWT will validate JWT token and return user details if valid
func ValidateJWT(t string) (*models.User, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	if exp.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID: userID,
	}, nil
}
