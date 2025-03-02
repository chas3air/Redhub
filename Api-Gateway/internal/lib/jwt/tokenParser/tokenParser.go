package tokenparser

import (
	"apigateway/internal/domain/models"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenString string, secretKey []byte) (*models.Claims, error) {
	fmt.Println("Начинаем проверку токена...")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("Ошибка верификации токена:", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Токен действителен. Извлекаем данные...")
		customClaims := &models.Claims{
			Uid:  claims["uid"].(string),
			Role: claims["role"].(string),
		}

		// Отладочный вывод типа exp
		expClaim := claims["exp"]
		fmt.Printf("Тип exp: %T, значение: %v\n", expClaim, expClaim)

		fmt.Println("Токен действителен. Данные:", customClaims)
		return customClaims, nil
	}

	fmt.Println("Токен недействителен.")
	return nil, fmt.Errorf("token is invalid")
}
