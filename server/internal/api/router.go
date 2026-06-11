package api

import (
	"context"

	"github.com/HIP-2026/MediCabinet/server/internal/inventory"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InventoryStore is the interface the api layer needs from the inventory subsystem.
type InventoryStore interface {
	AddItem(context.Context, inventory.AddItemParams) (inventory.Item, error)
	ListItems(context.Context) ([]inventory.Item, error)
}

// NewRouter builds the gin engine with middleware and routes. The pool and store
// may be nil; handlers report errors gracefully in that case.
func NewRouter(pool *pgxpool.Pool, store InventoryStore) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	engine.Use(cors.New(corsConfig))

	engine.GET("/health", healthHandler(pool))

	inv := engine.Group("/inventory")
	inv.POST("", createInventoryItemHandler(store))
	inv.GET("", listInventoryItemsHandler(store))

	return engine
}
