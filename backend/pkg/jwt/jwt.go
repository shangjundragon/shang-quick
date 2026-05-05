package jwt

import (
	"backend/pkg/global_vars"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"user_id,string"`
	Username string `json:"username"`
	RoleCode string `json:"role_code"`
	jwt.RegisteredClaims
}

func getSecret() string {
	secret := global_vars.ConfigYml.GetString("Jwt.Secret")
	if secret == "" {
		return "shang-quick-admin-secret-key-2026"
	}
	return secret
}

func getExpireHours() int {
	hours := global_vars.ConfigYml.GetInt("Jwt.ExpireHours")
	if hours <= 0 {
		return 24
	}
	return hours
}

func GenerateToken(userID int64, username, roleCode string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleCode: roleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(getExpireHours()) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getSecret()))
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(getSecret()), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
