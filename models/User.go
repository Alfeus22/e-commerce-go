package models

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRole string

const (
	RoleAdmin    UserRole = "ADMIN"
	RoleSeller   UserRole = "SELLER"
	RoleCustomer UserRole = "CUSTOMER"
)

type User struct {
	ID                   bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username             string        `bson:"username" json:"username" binding:"required"`
	Password             string        `bson:"password" json:"password" binding:"required,min=6"`
	Email                string        `bson:"email" json:"email" binding:"required,email"`
	Role                 UserRole      `bson:"role" json:"role"`
	SellerStatus         string        `bson:"seller_status" json:"seller_status"`
	jwt.RegisteredClaims `bson:"-" json:"-"`
}
