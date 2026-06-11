package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewRouter builds the gin engine with middleware and routes. The pool may be
// nil; the health check reports the database as down in that case.
func NewRouter(pool *pgxpool.Pool) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	engine.Use(cors.New(corsConfig))

	engine.GET("/health", healthHandler(pool))

	return engine
}
