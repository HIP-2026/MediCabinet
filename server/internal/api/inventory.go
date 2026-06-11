package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/HIP-2026/MediCabinet/server/internal/inventory"
)

type createItemRequest struct {
	Name       string  `json:"name"`
	Quantity   *int32  `json:"quantity"`
	ExpiryDate *string `json:"expiry_date"`
	Location   *string `json:"location"`
}

func createInventoryItemHandler(store InventoryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		if req.Quantity == nil || *req.Quantity < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be a non-negative integer"})
			return
		}
		if req.ExpiryDate != nil {
			if _, err := time.Parse("2006-01-02", *req.ExpiryDate); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "expiry_date must be in YYYY-MM-DD format"})
				return
			}
		}

		item, err := store.AddItem(c.Request.Context(), inventory.AddItemParams{
			Name:       req.Name,
			Quantity:   *req.Quantity,
			ExpiryDate: req.ExpiryDate,
			Location:   req.Location,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add item"})
			return
		}
		c.JSON(http.StatusCreated, item)
	}
}

func listInventoryItemsHandler(store InventoryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := store.ListItems(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list items"})
			return
		}
		if items == nil {
			items = []inventory.Item{}
		}
		c.JSON(http.StatusOK, items)
	}
}
