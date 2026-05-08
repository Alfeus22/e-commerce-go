package controllers

import (
	"context"
	"net/http"

	"github.com/Alfeus22/ecommerce-go/config"
	"github.com/Alfeus22/ecommerce-go/models"
	"github.com/Alfeus22/ecommerce-go/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(400, gin.H{"error ": err.Error()})
		return
	}

	// hashing password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(bytes)

	// create

	user.Role = "customer"
	user.SellerStatus = "none"
	_, err := config.UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"message ": "user berhasil dibuat",
	})

}

func Login(ctx *gin.Context) {
	var input models.User

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.Username == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username dna paswword wajib diisi"})
		return
	}

	var foundUser models.User
	filter := bson.M{"username": input.Username}

	// cari by username
	err := config.UserCollection.FindOne(context.TODO(), filter).Decode(&foundUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username tidak ditemukan"})
		return
	}

	// bandingkn password
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), ([]byte(input.Password)))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Password Salah"})
		return
	}
	token, err := utils.GenerateToken(foundUser.ID, foundUser.Role)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "gada token"})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "Berhasil Login",
		"token :": token,
		"role":    foundUser.Role,
	})
}
