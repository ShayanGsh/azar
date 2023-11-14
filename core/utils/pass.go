package utils

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type PasswordStrength string

const (
	VeryWeak   PasswordStrength = "very weak"
	Weak       PasswordStrength = "weak"
	Normal      PasswordStrength = "normal"
	Strong      PasswordStrength = "strong"
	Unbeatable  PasswordStrength = "unbeatable"
)

func GetPasswordStrength(password string) PasswordStrength {
	length := len(password)

	if length < 8 {
		return VeryWeak
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()]`).MatchString(password)

	score := 0
	if hasLower {
		score++
	}
	if hasUpper {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}
	if length >= 12 {
		score++
	}

	switch score {
	case 1, 2:
		return Weak
	case 3:
		return Normal
	case 4:
		return Strong
	default:
		return Unbeatable
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}