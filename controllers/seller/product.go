package seller

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
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

func GetProduct(c *gin.Context) {
	var product []models.Product
	val, exists := c.Get("currentuser")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak ditemukan"})
		return
	}
	id := val.(string)
	sellerID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id tidak ditemukan"})
		return
	}

	filter := bson.M{"seller_id": sellerID}

	cursor, err := config.ProductCollection.Find(context.TODO(), filter)

	if err != nil {
		c.JSON(500, gin.H{"error": "tidak ditemukan"})
		return
	}
	cursor.All(context.TODO(), &product)
	c.JSON(200, gin.H{
		"message": "berhasil mengirim data",
		"data":    product,
	})

}

func EditProduct(c *gin.Context) {
	// ambil id parameter
	productId := c.Param("id")
	objProductID, _ := bson.ObjectIDFromHex(productId)

	// ambil id seller
	val, _ := c.Get("currentuser")
	userID := val.(string)
	objUserID, _ := bson.ObjectIDFromHex(userID)

	// tangkap data baru
	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// validasi edit product sendiri
	filter := bson.M{
		"_id":       objProductID,
		"seller_id": objUserID,
	}

	update := bson.M{
		"$set": bson.M{
			"name":        updateData.Name,
			"price":       updateData.Price,
			"description": updateData.Desc,
			"stock":       updateData.Stock,
		},
	}

	result, err := config.ProductCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(403, gin.H{"error": "anda tidak punya akses atau produk tidak ditemukan"})
		return
	}
	c.JSON(200, gin.H{"message": "berhasil update data"})

}

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File gambar wajib diisi"})
		return
	}
	// validasi ijinkan hanya file berupa gambar
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file harus beupa jpg, png, jpeg, atau webp"})
		return
	}
	// modifikasi nama file agar unik (menggunakan timeStamp dan tidak tertimpa dengan file lama

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// tentukan jalur penyimpanan lokal
	dst := filepath.Join("uploads", fileName)

	// simpan file ke folder './uploads
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Menyimpan gambar ke server"})
		return
	}

	imageUrl := fmt.Sprintf("http://localhost:8080/uploads/%s", fileName)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Gambar berhasil diunggah",
		"image_url": imageUrl,
	})
}
