// Package password 提供 bcrypt 密码哈希、验证和随机密码生成
package password

import (
	"crypto/rand"
	"errors"
	"math/big"
	mrand "math/rand"

	"golang.org/x/crypto/bcrypt"
)

// passwordChars 随机密码字符集：小写字母(26)+大写字母(26)+数字(10)+特殊字符(10)
const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度不能少于8位")
	}
	return nil
}

func ValidatePasswordStrong(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度不能少于8位")
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	if !hasUpper {
		return errors.New("密码必须包含大写字母")
	}
	if !hasLower {
		return errors.New("密码必须包含小写字母")
	}
	if !hasDigit {
		return errors.New("密码必须包含数字")
	}
	if !hasSpecial {
		return errors.New("密码必须包含特殊字符")
	}
	return nil
}

// GenerateRandomPassword 生成指定长度的随机密码，确保至少包含大写字母、小写字母、数字和特殊字符各一个
func GenerateRandomPassword(length int) string {
	if length < 8 {
		length = 12
	}

	result := []byte{
		passwordChars[randInt(26, 52)],
		passwordChars[randInt(0, 26)],
		passwordChars[randInt(52, 62)],
		passwordChars[randInt(62, len(passwordChars))],
	}

	for i := 4; i < length; i++ {
		result = append(result, passwordChars[randInt(0, len(passwordChars))])
	}

	mrand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return string(result)
}

func randInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(n.Int64()) + min
}
