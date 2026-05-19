package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING_PAYMENT"
	StatusPaid      OrderStatus = "PAID"
	StatusCancelled OrderStatus = "CANCELLED"
)

type OrderItem struct {
	ProductID bson.ObjectID `bson:"product_id" json:"product_id"`
	Quantity  int           `bson:"quantity" json:"quantity"`
	Price     int           `bson:"price" json:"price"`
	SellerID  bson.ObjectID `bson:"seller_id" json:"seller_id"`
}

type Order struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      bson.ObjectID `bson:"user_id" json:"user_id"`
	Items       []OrderItem   `bson:"items" json:"items"`
	TotalAmount int           `bson:"total_amount" json:"total_amount"`
	Status      OrderStatus   `bson:"status" json:"status"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
}
