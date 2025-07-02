package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
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

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresDB(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, p *Property) (string, error) {
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	availabilityJSON, _ := json.Marshal(p.Availability)
	imageURLsJSON, _ := json.Marshal(p.ImageURLs)

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO properties (
			id, owner_id, title, description, city, address_line, lat, lng,
			property_type, rooms, area, floor, total_floors, price_per_month,
			currency, main_image_url, image_urls, has_wifi, has_parking,
			has_elevator, is_verified, rating, reviews_count, availability,
			status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19,
			$20, $21, $22, $23, $24,
			$25, $26, $27
		)`,
		p.ID, p.OwnerID, p.Title, p.Description, p.City, p.AddressLine, p.Lat, p.Lng,
		p.PropertyType, p.Rooms, p.Area, p.Floor, p.TotalFloors, p.PricePerMonth,
		p.Currency, p.MainImageURL, imageURLsJSON, p.HasWiFi, p.HasParking,
		p.HasElevator, p.IsVerified, p.Rating, p.ReviewsCount, availabilityJSON,
		p.Status, p.CreatedAt, p.UpdatedAt,
	)
	return p.ID, err
}

func (r *PostgresRepo) Get(ctx context.Context, id string) (*Property, error) {
	row := r.db.QueryRowContext(ctx, `SELECT 
		id, owner_id, title, description, city, address_line, lat, lng,
		property_type, rooms, area, floor, total_floors, price_per_month,
		currency, main_image_url, image_urls, has_wifi, has_parking,
		has_elevator, is_verified, rating, reviews_count, availability,
		status, created_at, updated_at
	FROM properties WHERE id = $1`, id)

	var p Property
	var imageURLsJSON, availabilityJSON []byte

	err := row.Scan(
		&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.City, &p.AddressLine, &p.Lat, &p.Lng,
		&p.PropertyType, &p.Rooms, &p.Area, &p.Floor, &p.TotalFloors, &p.PricePerMonth,
		&p.Currency, &p.MainImageURL, &imageURLsJSON, &p.HasWiFi, &p.HasParking,
		&p.HasElevator, &p.IsVerified, &p.Rating, &p.ReviewsCount, &availabilityJSON,
		&p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(imageURLsJSON, &p.ImageURLs)
	_ = json.Unmarshal(availabilityJSON, &p.Availability)

	return &p, nil
}

func (r *PostgresRepo) Update(ctx context.Context, id string, p *Property) error {
	availabilityJSON, _ := json.Marshal(p.Availability)
	imageURLsJSON, _ := json.Marshal(p.ImageURLs)

	_, err := r.db.ExecContext(ctx, `
		UPDATE properties SET
		owner_id = $1, title = $2, description = $3, city = $4, address_line = $5, lat = $6, lng = $7,
		property_type = $8, rooms = $9, area = $10, floor = $11, total_floors = $12, price_per_month = $13,
		currency = $14, main_image_url = $15, image_urls = $16, has_wifi = $17, has_parking = $18,
		has_elevator = $19, is_verified = $20, rating = $21, reviews_count = $22, availability = $23,
		status = $24, updated_at = $25
		WHERE id = $26`,
		p.OwnerID, p.Title, p.Description, p.City, p.AddressLine, p.Lat, p.Lng,
		p.PropertyType, p.Rooms, p.Area, p.Floor, p.TotalFloors, p.PricePerMonth,
		p.Currency, p.MainImageURL, imageURLsJSON, p.HasWiFi, p.HasParking,
		p.HasElevator, p.IsVerified, p.Rating, p.ReviewsCount, availabilityJSON,
		p.Status, p.UpdatedAt, id,
	)
	return err
}

func (r *PostgresRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM properties WHERE id = $1`, id)
	return err
}