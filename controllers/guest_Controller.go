package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Alfeus22/ecommerce-go/config"
	"github.com/Alfeus22/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Product struct {
	ID primitive.ObjectID ``
}

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

func GetAllProduct(c *gin.Context) {
	var product []models.Product

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	// validasi agar ga minus

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 5
	}

	// hitung skip

	skip := (page - 1) * limit

	// siapkan option MongoDB
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{"id", -1}}) //urutkan dari yang terbaru

	// hitung total data (untuk info di response)
	total, err := config.ProductCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "gagal menghitung data"})
		return
	}

	// jalankan query find dengan options
	cursor, err := config.ProductCollection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil data dari produk"})
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &product); err != nil {
		c.JSON(500, gin.H{"error": "Gagal decode data"})
		return
	}

	// respon dengan metadata pagination
	c.JSON(200, gin.H{
		"status": "success",
		"pagination": gin.H{
			"total_item":   total,
			"current_page": page,
			"limit":        limit,
			"total_pages":  (int(total) + limit - 1) / limit,
		},
		"data": product,
	})

}

func SearchProduct(c *gin.Context) {
	searchQuery := c.Query("name")

	// filter
	filter := bson.M{}
	if searchQuery != "" {
		// $regex: mencari potongan kata, $options: "i" agar tidak sensitif huruf besar/kecil
		filter = bson.M{
			"name": bson.M{
				"$regex":   searchQuery,
				"$options": "i",
			},
		}
	}

	cursor, err := config.ProductCollection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	// 4. Decode hasil query ke dalam slice struct Product
	// Diinisialisasi dengan array kosong [] agar di JSON response tidak muncul nilai 'null'
	var products []models.Product = []models.Product{}
	if err := cursor.All(context.TODO(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data produk"})
		return
	}

	// kirim response ke guest
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mendapatkan Produk",
		"data":    products,
	})

}
