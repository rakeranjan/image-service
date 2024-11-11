package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	SECRET_KEY = "SECRET_KEY"
)

type User struct {
	ID          string `json:"id" dynamodbav:"id"`
	FirstName   string `json:"firstName" dynamodbav:"firstName" form:"firstName" binding:"required,alpha"`
	LastName    string `json:"lastName" dynamodbav:"lastName" form:"lastName" binding:"required,alpha"`
	UserName    string `json:"userName" dynamodbav:"userName" form:"userName" binding:"required,userNameFormat"`
	Password    string `json:"-" dynamodbav:"password" form:"password" binding:"required,userNameFormat`
	PhoneNumber string `json:"phoneNumber" dynamodbav:"phoneNumber" form:"phoneNumber" binding:"required,phoneNumberFormat"`
}

func (u *User) Encrypt(secret string) (string, error) {
	claims := jwt.MapClaims{
		"id":          u.ID,
		"firstName":   u.FirstName,
		"lastName":    u.LastName,
		"userName":    u.UserName,
		"phoneNumber": u.PhoneNumber,
		"exp":         time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	}

	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// func (u *User) DecodeJWTToStruct(tokenString string) error {
// 	conf, _ := configuration.NewConfiguration()
// 	secretKey := os.Getenv(conf.SecretValue)
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(secretKey), nil
// 	})

// 	if err != nil {
// 		return fmt.Errorf("failed to parse token: %w", err)
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		u.ID = claims["id"].(string)
// 		u.FirstName = claims["firstName"].(string)
// 		u.LastName:    claims["lastName"].(string),
// 		// UserName:    claims["userName"].(string),
// 		// PhoneNumber: claims["phoneNumber"].(string),

// 		// ID:          claims["id"].(string),
// 		// FirstName:   claims["firstName"].(string),
// 		// LastName:    claims["lastName"].(string),
// 		// UserName:    claims["userName"].(string),
// 		// PhoneNumber: claims["phoneNumber"].(string),
// 		// }
// 		return nil
// 	} else {
// 		return fmt.Errorf("invalid token")
// 	}
// }
