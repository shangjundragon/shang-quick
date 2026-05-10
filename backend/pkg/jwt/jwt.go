package jwt

import (
	"backend/pkg/global_vars"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	Issuer   = "shang-quick"
	Audience = "api"
)

type Claims struct {
	UserID   int64  `json:"user_id,string"`
	Username string `json:"username"`
	RoleCode string `json:"role_code"`
	jwt.RegisteredClaims
}

func getSecret() (string, error) {
	secret := global_vars.ConfigYml.GetString("Jwt.Secret")
	if secret == "" {
		return "", errors.New("JWT Secret 未配置")
	}
	return secret, nil
}

func getExpireHours() int {
	hours := global_vars.ConfigYml.GetInt("Jwt.ExpireHours")
	if hours <= 0 {
		return 24
	}
	return hours
}

func GenerateToken(userID int64, username, roleCode string) (string, error) {
	secret, err := getSecret()
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleCode: roleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(getExpireHours()) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    Issuer,
			Audience:  []string{Audience},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr string) (*Claims, error) {
	secret, err := getSecret()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithIssuer(Issuer), jwt.WithAudience(Audience))
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}
	return nil, errors.New("无效的token")
}
