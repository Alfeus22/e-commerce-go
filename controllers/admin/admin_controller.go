package controllers

import (
	"context"
	"net/http"

	"github.com/Alfeus22/ecommerce-go/config"
	"github.com/Alfeus22/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetPendingSeller(ctx *gin.Context) {
	var user []models.User
	cursor, err := config.UserCollection.Find(context.TODO(), bson.M{"seller_status": "pending"})

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	cursor.All(context.TODO(), &user)
	ctx.JSON(200, user)

}

// setujui pengajuan user menjadi seller

func UpdateUser(c *gin.Context) {
	targetUserID := c.Param("id")
	objId, _ := bson.ObjectIDFromHex(targetUserID)

	filter := bson.M{"_id": objId}
	update := bson.M{
		"$set": bson.M{
			"role":          "seller",
			"seller_status": "approved",
		},
	}

	_, err := config.UserCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "gagal menyetujui seller",
		})
		return
	}
	c.JSON(200, gin.H{"message": "User berhasil menjadi seller"})

}
func GetSeller(c *gin.Context) {
	var seller []models.User
	cursor, err := config.UserCollection.Find(context.TODO(), bson.M{"role": "seller"})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	cursor.All(context.TODO(), &seller)
	c.JSON(200, seller)

}

func DeleteSeller(c *gin.Context) {
	idSeller := c.Param("id")

	objId, err := bson.ObjectIDFromHex(idSeller)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format id tidak valid"})
		return
	}

	filter := bson.M{"_id": objId}

	delete, err := config.UserCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gagal menghapus seller"})
		return
	}

	if delete.DeletedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seller tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "seller berhasil dihapus",
		"user_id": idSeller,
	})

}
