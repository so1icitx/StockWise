package httpapi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/so1icitx/StockWise/internal/domain"
)

const actorIDKey = "actorID"

// ActorMiddleware reads the optional X-User-ID header into the Gin context.
func ActorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		raw := ctx.GetHeader("X-User-ID")
		if raw == "" {
			ctx.Next()
			return
		}

		value, err := strconv.ParseUint(raw, 10, 64)
		if err == nil && value > 0 {
			ctx.Set(actorIDKey, domain.ID(value))
		}

		ctx.Next()
	}
}

func actorID(ctx *gin.Context) (domain.ID, error) {
	value, ok := ctx.Get(actorIDKey)
	if !ok {
		return 0, validationHTTPError("X-User-ID header is required")
	}

	id, ok := value.(domain.ID)
	if !ok || id.IsZero() {
		return 0, validationHTTPError("X-User-ID header is invalid")
	}

	return id, nil
}
