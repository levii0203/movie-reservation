package model

import "time"

type GENDER string

type Booking struct {
	ID        string	`json:"_id" bson:"_id"`
	UserID    string	`json:"user_id" bson:"user_id"`
	MovieID   string	`json:"movie_id" bson:"movie_id"`
	Gender    GENDER	`json:"gender" bson:"gender" validate:"required"`
	Age       int16		`json:"age" bson:"age" validate:"required,gte=1,lte=150"`
	Seat      string	`json:"seat" bson:"seat" validate:"required"`
	Type      string	`json:"type,omitempty" bson:"type"`
	CreatedAt time.Time	`json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type SeatLock struct {
	MovieID string
	UserID string
	Seat string
}