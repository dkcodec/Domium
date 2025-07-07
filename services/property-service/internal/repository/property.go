package repository

import (
	"context"
	"time"

	_ "github.com/lib/pq"
)

type Property struct {
	ID             string
	OwnerID        string
	Title          string
	Description    string
	City           string
	AddressLine    string
	Lat            float64
	Lng            float64
	PropertyType   string
	Rooms          int32
	Area           float64
	Floor          int32
	TotalFloors    int32
	PricePerMonth  int32
	Currency       string
	MainImageURL   string
	ImageURLs      []string
	HasWiFi        bool
	HasParking     bool
	HasElevator    bool
	IsVerified     bool
	Rating         float64
	ReviewsCount   int32
	Availability   []AvailabilityPeriod
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type AvailabilityPeriod struct {
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}

type Repository interface {
	Create(ctx context.Context, p *Property) (string, error)
	Get(ctx context.Context, id string) (*Property, error)
	Update(ctx context.Context, id string, p *Property) error
	Delete(ctx context.Context, id string) error
}

