package seller

import (
	"context"
	"net/http"
	"time"

	"github.com/Alfeus22/ecommerce-go/config"
	"github.com/Alfeus22/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson" // Konsisten pakai v2
)

func MakeProduct(c *gin.Context) {
	var newProduct models.Product

	// 1. Ambil ID user dari middleware (Pastikan key-nya sama dengan middleware)
	val, exists := c.Get("currentuser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak ditemukan"})
		return
	}

	userIDStr := val.(string)
	objID, err := bson.ObjectIDFromHex(userIDStr) // Pakai bson v2
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format ID User salah"})
		return
	}

	// 2. Bind JSON (Gunakan ShouldBindJSON)
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Gagal di binding JSON",
			"error":  err.Error()})
		return
	}

	// 3. Set SellerID dan ID Produk
	newProduct.SellerID = objID
	if newProduct.ID.IsZero() {
		newProduct.ID = bson.NewObjectID() // Pakai bson v2
	}

	// 4. Insert ke Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = config.ProductCollection.InsertOne(ctx, newProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan database: " + err.Error()})
		return
	}

	// 5. Kirim Respon Sukses (PENTING: Jangan biarkan kosong)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product berhasil ditambahkan",
		"data":    newProduct,
	})
}

func DeleteProduct(c *gin.Context) {
	productIDStr := c.Param("id")

	// Pastikan key ini "currentuser", sama dengan MakeProduct
	val, exists := c.Get("currentuser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	productObjID, errProd := bson.ObjectIDFromHex(productIDStr)
	sellerObjID, errSell := bson.ObjectIDFromHex(val.(string))

	if errProd != nil || errSell != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Filter agar seller hanya bisa hapus produk miliknya sendiri
	filter := bson.M{
		"_id":       productObjID,
		"seller_id": sellerObjID,
	}

	result, err := config.ProductCollection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan atau Anda bukan pemiliknya"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product berhasil dihapus"})
}
