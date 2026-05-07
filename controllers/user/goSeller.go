package user

import (
	"context"

	"github.com/Alfeus22/ecommerce-go/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func UpdateToSeller(c *gin.Context) {
	userId, _ := c.Get("currentuser")
	objId, _ := bson.ObjectIDFromHex(userId.(string))

	filter := bson.M{"_id": objId}
	update := bson.M{
		"$set": bson.M{"seller_status": "pending"},
	}

	_, err := config.UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "gagal mengajikan menjadi seller",
		})
		return
	}
	c.JSON(200, gin.H{"message": "pengajuan diproses admin"})
}
