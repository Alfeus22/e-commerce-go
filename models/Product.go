package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Product struct {
	ID    bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name  string        `bson:"name" json:"name" binding:"required"`
	Stock int           `bson:"stock" json:"stock" binding:"required,min=1"`
	Image string        `bson:"image_url" json:"image_url" binding:"required"`
	Price int           `bson:"price" json:"price" binding:"required"`
	Desc  string        `bson:"desc" json:"desc" binding:"required"`

	SellerID bson.ObjectID `bson:"seller_id" json:"seller_id"`
}
