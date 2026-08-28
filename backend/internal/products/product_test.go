package products

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct(t *testing.T) {
	t.Run("creates product with valid fields", func(t *testing.T) {
		product := Product{
			Name:        "Widget",
			Description: "A useful widget",
			Price:       29.99,
			Quantity:    100,
			SKU:         "WDG-001",
		}

		assert.Equal(t, "Widget", product.Name)
		assert.Equal(t, "A useful widget", product.Description)
		assert.Equal(t, 29.99, product.Price)
		assert.Equal(t, 100, product.Quantity)
		assert.Equal(t, "WDG-001", product.SKU)
	})

	t.Run("validates required fields", func(t *testing.T) {
		product := Product{}
		err := product.Validate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("validates name is required", func(t *testing.T) {
		product := Product{Price: 10.0, SKU: "TEST-001"}
		err := product.Validate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("validates price is non-negative", func(t *testing.T) {
		product := Product{Name: "Test", Price: -5.0, SKU: "TEST-001"}
		err := product.Validate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "price")
	})

	t.Run("validates quantity is non-negative", func(t *testing.T) {
		product := Product{Name: "Test", Price: 10.0, Quantity: -1, SKU: "TEST-001"}
		err := product.Validate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "quantity")
	})

	t.Run("validates SKU is required", func(t *testing.T) {
		product := Product{Name: "Test", Price: 10.0}
		err := product.Validate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sku")
	})

	t.Run("passes validation with valid data", func(t *testing.T) {
		product := Product{
			Name:     "Valid Product",
			Price:    19.99,
			Quantity: 50,
			SKU:      "VP-001",
		}
		err := product.Validate()

		assert.NoError(t, err)
	})
}
