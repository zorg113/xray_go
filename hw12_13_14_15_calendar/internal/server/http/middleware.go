package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)

		resp := fmt.Sprintf("%s %s %s %s %d %v %s", r.RemoteAddr, r.Method, r.URL.Path, r.Proto, http.StatusOK, time.Since(startTime), r.UserAgent())
		s.logger.Info(resp)
	})
}
