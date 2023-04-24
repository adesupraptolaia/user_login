package jwt

import (
	"github.com/adesupraptolaia/user_login/config"
	"github.com/golang-jwt/jwt"
)

var (
	ACCESS_TOKEN_SECRET  []byte
	REFRESH_TOKEN_SECRET []byte
)

type Claims struct {
	UserKsuid string `json:"user_ksuid"`
	Role      string `json:"role"`
	jwt.StandardClaims
}

func init() {
	cfg := config.Config
	ACCESS_TOKEN_SECRET = []byte(cfg.Secret.AccessToken)
	REFRESH_TOKEN_SECRET = []byte(cfg.Secret.RefreshToken)
}
