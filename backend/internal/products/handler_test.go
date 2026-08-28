package products

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewProductHandler(db)

	router.GET("/products", handler.List)
	router.GET("/products/:id", handler.Get)
	router.POST("/products", handler.Create)
	router.PUT("/products/:id", handler.Update)
	router.DELETE("/products/:id", handler.Delete)

	return router
}

func TestProductHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("rejects invalid product data", func(t *testing.T) {
		router := setupTestRouter(nil)

		product := Product{
			Name: "",
		}
		body, _ := json.Marshal(product)

		req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("rejects invalid id format", func(t *testing.T) {
		router := setupTestRouter(nil)

		req, _ := http.NewRequest("GET", "/products/invalid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
