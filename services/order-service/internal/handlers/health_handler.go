package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Liveness: is the process itself alive? No dependency checks — if this
// fails, K8s restarts the pod. Keep it trivially cheap on purpose.
func Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

// Readiness: can this pod actually serve traffic right now? Checks the DB —
// if Postgres is unreachable, K8s stops routing traffic to this pod (but
// does NOT restart it) until the check passes again. This is the same
// liveness/readiness split as the Notification service's Actuator probes.
func Readiness(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "reason": "db_unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
}
