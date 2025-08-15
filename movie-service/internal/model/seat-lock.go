package model

type SeatLock struct {
	MovieID string 	`json:"movie_id" validate:"required`
	UserID string	`json:"user_id" validate:"required`
	Seat string		`json:"seat" validate:"required`
}