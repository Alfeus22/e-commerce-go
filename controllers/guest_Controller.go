package controllers

import (
	"context"

	"github.com/Alfeus22/ecommerce-go/config"
	"github.com/Alfeus22/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func DetailProduct(c *gin.Context) {
	idParam := c.Param("id")

	objId, err := bson.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "format id tidak valid"})
		return
	}
	var Product models.Product
	filter := bson.M{"_id": objId}

	err = config.ProductCollection.FindOne(context.TODO(), filter).Decode(&Product)
	if err != nil {
		c.JSON(404, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Berhasil kirim data",
		"data":    Product,
	})

}
