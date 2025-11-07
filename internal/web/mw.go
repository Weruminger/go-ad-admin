package web

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const ctxReqID ctxKey = "reqid"

func withReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), ctxReqID, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func reqIDFrom(r *http.Request) string {
	if v := r.Context().Value(ctxReqID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
