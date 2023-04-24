package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/adesupraptolaia/user_login/internal/entity"
	"github.com/golang-jwt/jwt"
)

// create JWT Access Token, valid until 1 hour
func CreateAccessToken(userKsuid, role string) (string, error) {
	return createToken(
		userKsuid,
		role,
		ACCESS_TOKEN_SECRET,
		time.Now().Add(1*time.Hour).Unix(),
	)
}

// create JWT Refressh Token, valid until 24 hour
func CreateRefreshToken(userKsuid, role string) (string, error) {
	return createToken(
		userKsuid,
		role,
		REFRESH_TOKEN_SECRET,
		time.Now().Add(24*time.Hour).Unix(),
	)
}

// validate Admin
func IsAdmin(tokenString string) error {
	claims, err := validateToken(tokenString, ACCESS_TOKEN_SECRET)
	if err != nil {
		return err
	}

	if claims.Role != entity.ADMIN {
		return fmt.Errorf("unauthorize")
	}

	return nil
}

// validate JWT Access Token
func ValidateAccessToken(tokenString, userKsuid string) error {
	claims, err := validateToken(tokenString, ACCESS_TOKEN_SECRET)
	if err != nil {
		return err
	}

	if !(claims.Role == entity.ADMIN || claims.UserKsuid == userKsuid) {
		return fmt.Errorf("unauthorize")
	}

	return nil
}

// validate JWT Refresh Token
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := validateToken(tokenString, REFRESH_TOKEN_SECRET)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func createToken(userKsuid, role string, secretKey []byte, expiresAt int64) (string, error) {
	// Create the claims for the JWT token
	claims := &Claims{
		UserKsuid: userKsuid,
		Role:      role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}
	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the JWT token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func validateToken(tokenString string, secretKey []byte) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Extract the custom claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if the token is expired
	if time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
		return nil, errors.New("token has expired")
	}

	// Return the custom claims
	return claims, nil
}
