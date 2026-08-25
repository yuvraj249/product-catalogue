package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"backend/internal/interface/storage/redis"
	"github.com/google/uuid"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}

func IdempotencyMiddleware(rdb *redis.RedisClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			idempotencyKey := r.Header.Get("X-Idempotency-Key")
			if idempotencyKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			tenantIDVal := ctx.Value(TenantContextKey)
			tenantIDStr := ""
			if tID, ok := tenantIDVal.(uuid.UUID); ok {
				tenantIDStr = tID.String()
			}

			redisKey := fmt.Sprintf("idempotency:%s:%s", tenantIDStr, idempotencyKey)

			// Atomic acquisition
			acquired, err := rdb.SetNX(ctx, redisKey, "IN_PROGRESS", 24*time.Hour)
			if err != nil {
				http.Error(w, `{"error":"internal distributed lock failure"}`, http.StatusInternalServerError)
				return
			}

			if !acquired {
				val, err := rdb.Get(ctx, redisKey)
				if err != nil || val == "IN_PROGRESS" {
					http.Error(w, `{"error":"transaction currently processing, retry shortly"}`, http.StatusConflict)
					return
				}
				// Return cached payload directly
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Cache-Hit", "true")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(val))
				return
			}

			rec := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           new(bytes.Buffer),
			}

			next.ServeHTTP(rec, r)

			if rec.statusCode >= 200 && rec.statusCode < 300 {
				_ = rdb.Set(ctx, redisKey, rec.body.String(), 24*time.Hour)
			} else {
				_ = rdb.Del(ctx, redisKey) // Release lock on processing error
			}
		})
	}
}
