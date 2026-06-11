package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// healthHandler reports liveness. It pings the database when a pool is
// available but never fails the endpoint, so the server reports healthy even
// before Postgres is reachable.
func healthHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := "down"
		if pool != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			defer cancel()
			if err := pool.Ping(ctx); err == nil {
				db = "up"
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "db": db})
	}
}
