package apiutils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rakeranjan/image-service/api/models"
	"github.com/rakeranjan/image-service/internal/configuration"
)

const ()

func DecodeJWTToStruct(tokenString string) (*models.User, error) {
	conf, _ := configuration.NewConfiguration()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.SecretValue), nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &models.User{
			ID:          claims["id"].(string),
			FirstName:   claims["firstName"].(string),
			LastName:    claims["lastName"].(string),
			UserName:    claims["userName"].(string),
			PhoneNumber: claims["phoneNumber"].(string),
		}
		return user, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
