package domain

import (
	"database/sql"
	"time"
)

// Ingredient represents a food item in the refrigerator
type Ingredient struct {
	ID           int64      `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Quantity     string     `json:"quantity" db:"quantity"`
	PurchaseDate *time.Time `json:"purchase_date,omitempty" db:"purchase_date"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// NullableTime is a helper type for handling nullable time fields in database
type NullableTime struct {
	sql.NullTime
}

// MarshalJSON implements json.Marshaler interface for NullableTime
func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return nt.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler interface for NullableTime
func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}

	var t time.Time
	if err := t.UnmarshalJSON(data); err != nil {
		return err
	}

	nt.Time = t
	nt.Valid = true
	return nil
}
