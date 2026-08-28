package products

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Quantity    int            `gorm:"not null;default:0" json:"quantity"`
	SKU         string         `gorm:"uniqueIndex;not null" json:"sku"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Price < 0 {
		return errors.New("price must be non-negative")
	}
	if p.Quantity < 0 {
		return errors.New("quantity must be non-negative")
	}
	if p.SKU == "" {
		return errors.New("sku is required")
	}
	return nil
}
