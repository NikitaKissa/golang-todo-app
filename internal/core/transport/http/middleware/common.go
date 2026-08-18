package core_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/NikitaKissa/golang-todo-app/internal/core/logger"
	core_http_response "github.com/NikitaKissa/golang-todo-app/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	RequestIdHeader = "X-Request-ID"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http://localhost:8080": {},
			}

			origin := r.Header.Get("Origin")
			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIdHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(RequestIdHeader, requestID)
			w.Header().Set(RequestIdHeader, requestID)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIdHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">> incoming http request",
				zap.Time("time", before.UTC()),
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
			)

			next.ServeHTTP(rw, r.WithContext(ctx))

			log.Debug(
				"<< done http request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Since(before)),
				zap.String("url", r.URL.String()),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHttpResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle http request got panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
