package model

import (
	"time"
)


type Movie struct {
	ID          		string `json:"id" bson:"_id"`
	Title	   			string `json:"title" bson:"title" validate:"required"`
	Genre	  	  	  []string `json:"genre,omitempty" bson:"genre"`
	Rating	   	       float64 `json:"rating" bson:"rating,omitempty" validate:"gte=1.0,lte=10"`
	Available     	      bool `json:"available" bson:"available" validate:"default=true"`
	ReleaseDate 		string `json:"release_date" bson:"release_date" validate:"required"`
	FilledSeats       []string `json:"filled_seats,omitempty" bson:"filled_seats,omitempty"`
	AvailableSeats    []string `json:"available_seats,omitempty" bson:"available_seats,omitempty"`
	CreatedAt	  	 time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt	  	 time.Time `json:"updated_at" bson:"updated_at"`
	DurationInMinutes 	   int `json:"duration_in_minutes" bson:"duration_in_minutes" validate:"required,gte=1,lte=300"`
	Hall				string `json:"hall" bson:"hall" validate:"required"`
	PosterURL			string `json:"poster_url,omitempty" bson:"poster_url,omitempty" validate:"url"`
	Start            time.Time `json:"start" bson:"start" validate:"required"`
	End              time.Time `json:"end" bson:"end" validate:"required"`
	Cast			  []string `json:"cast,omitempty" bson:"cast,omitempty"`
	Cinema				string `json:"cinema" bson:"cinema" validate:"required"`
	City				string `json:"city" bson:"city" validate:"required,min=1,max=255"`
	Over				string `json:"over,omitempty" bson:"over,omitempty"`
	ViewCount			   int `json:"view_count,omitempty" bson:"view_count,omitempty" validate:"gte=0"`
}