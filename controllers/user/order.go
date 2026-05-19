package user

import (
	"context"
	"os"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"

	"github.com/Alfeus22/ecommerce-go/config"
	"github.com/Alfeus22/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateOrder(c *gin.Context) {
	// 1. Ambil User ID dari context (setelah login)
	val, _ := c.Get("currentuser")
	userId, _ := bson.ObjectIDFromHex(val.(string))

	var input struct {
		Items []models.OrderItem `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var totalAmount int
	var orderItems []models.OrderItem
	var midtransItems []midtrans.ItemDetails // Untuk detail ke Midtrans

	// 2. Validasi produk, hitung total, dan siapkan item detail
	for _, item := range input.Items {
		var product models.Product
		err := config.ProductCollection.FindOne(context.TODO(), bson.M{"_id": item.ProductID}).Decode(&product)
		if err != nil {
			c.JSON(404, gin.H{"error": "Produk " + item.ProductID.Hex() + " tidak ditemukan"})
			return
		}

		if product.Stock < item.Quantity {
			c.JSON(400, gin.H{"error": "Stok tidak cukup untuk: " + product.Name})
			return
		}

		totalItem := product.Price * item.Quantity
		totalAmount += totalItem

		// Masukkan ke slice orderItems (untuk simpan di DB)
		item.Price = product.Price
		item.SellerID = product.SellerID
		orderItems = append(orderItems, item)

		// Masukkan ke slice midtransItems (untuk kirim ke Midtrans)
		midtransItems = append(midtransItems, midtrans.ItemDetails{
			ID:    item.ProductID.Hex(),
			Name:  product.Name,
			Price: int64(product.Price),
			Qty:   int32(item.Quantity),
		})

		// Potong stok produk
		newStok := product.Stock - item.Quantity
		config.ProductCollection.UpdateOne(context.TODO(),
			bson.M{"_id": product.ID},
			bson.M{"$set": bson.M{"stock": newStok}})
	}

	// 3. Ambil data User untuk Customer Detail Midtrans
	var user models.User
	config.UserCollection.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&user)

	// 4. Simpan orderan ke database
	newOrder := models.Order{
		ID:          bson.NewObjectID(),
		UserID:      userId,
		Items:       orderItems,
		TotalAmount: totalAmount,
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
	}

	_, err := config.OrderCollection.InsertOne(context.TODO(), newOrder)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal menyimpan orderan ke database"})
		return
	}

	// 5. Inisialisasi Midtrans Client
	var s = snap.Client{}
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")

	s.New(serverKey, midtrans.Sandbox) // Ganti 'secretbos' dengan Server Key kamu

	// 6. Buat request Snap Midtrans
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  newOrder.ID.Hex(),
			GrossAmt: int64(newOrder.TotalAmount),
		},
		Items: &midtransItems, // Detail item sudah terisi dari loop di atas
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Username, // Mengambil username dari DB
			Email: user.Email,    // Mengambil email dari DB
		},
	}

	// 7. Minta token ke Midtrans
	snapResp, snapErr := s.CreateTransaction(req)
	if snapErr != nil {
		c.JSON(500, gin.H{
			"error":   "Gagal membuat transaksi ke Midtrans",
			"details": snapErr.Message,
		})
		return
	}

	// 8. Kirim redirect url ke frontend (Hanya satu response sukses di akhir)
	c.JSON(201, gin.H{
		"message":      "Order dibuat, silakan dibayar",
		"order_id":     newOrder.ID.Hex(),
		"redirect_url": snapResp.RedirectURL,
		"token":        snapResp.Token,
	})
}

func HandleMidtransNotification(c *gin.Context) {
	var notificationPayLoad map[string]interface{}
	if err := c.ShouldBindJSON(&notificationPayLoad); err != nil {
		c.JSON(400, gin.H{"error": "Payload tidak valid"})
		return
	}

	orderId := notificationPayLoad["order_id"].(string)
	transactionStatus := notificationPayLoad["transaction_status"].(string)
	fraudStatus := notificationPayLoad["fraud_status"].(string)

	objId, _ := bson.ObjectIDFromHex(orderId)
	statusUpdate := ""

	// logika status midtrans
	if transactionStatus == "capture" || transactionStatus == "settlement" {
		if fraudStatus == "accept" {
			statusUpdate = "PAID"
		}
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		statusUpdate = "CANCELLED"
	}
	if statusUpdate != "" {
		// update status di mongo
		filter := bson.M{"_id": objId}
		update := bson.M{"$set": bson.M{"status": statusUpdate}}
		config.OrderCollection.UpdateOne(context.TODO(), filter, update)

	}
}
