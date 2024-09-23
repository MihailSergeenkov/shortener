package routes

import (
	"net"
	"net/http"

	"go.uber.org/zap"
)

func checkSubnetMiddleware(l *zap.Logger, trustedSubnet *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			ipStr := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(ipStr)
			if ip == nil {
				w.WriteHeader(http.StatusInternalServerError)
				l.Error("failed to parse ip address in header")
				return
			}

			ok := trustedSubnet.Contains(ip)
			if !ok {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
