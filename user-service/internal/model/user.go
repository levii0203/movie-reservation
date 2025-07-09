package model 

import (

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct {

	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Username   string               `bson:"username" json:"username" validate:"required,min=3,max=255,alphanum"`
	Password   string               `bson:"password" json:"-" validate:"required,min=8,alphanum"`
	Email      string               `bson:"email" json:"email" validate:"required,email"`
	Age        uint8      			`bson:"age" json:"age" validate:"gte=0,lte=130"`
	Phone      string               `bson:"phone,omitempty" json:"phone,omitempty" validate:"omitempty,phone"`
	Gender     string     			`bson:"gender,omitempty" json:"gender" validate:"oneof=male female prefer_not_to"`
	Bookings   []primitive.ObjectID `bson:"bookings,omitempty" json:"bookings,omitempty" validate:"dive,objectid"`
	WatchList  []primitive.ObjectID `bson:"watchlist,omitempty" json:"watchlist,omitempty" validate:"dive,objectid"`
	CreatedAt  time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time            `bson:"updated_at" json:"updated_at"`

}
