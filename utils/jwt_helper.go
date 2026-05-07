package utils

import (
	"time"

	"github.com/Alfeus22/ecommerce-go/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var secretKey = []byte("secretbos")

func GenerateToken(userID bson.ObjectID, role models.UserRole) (string, error) {
	claims := jwt.MapClaims{
		"seller_id": userID.Hex(),
		"role":      role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
