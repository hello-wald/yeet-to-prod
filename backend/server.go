package main

import (
	"encoding/json"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

type response struct {
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
	Safe     bool   `json:"safe"`
	Reason   string `json:"reason"`
	Message  string `json:"message"`
}

// server holds the request-handling dependencies. Split out from main() so the
// HTTP handlers are testable with httptest (no real listener, no real clock).
type server struct {
	cfg            Config
	msgs           Messages
	lim            *limiter
	rng            *rand.Rand
	defaultCountry string
	nowOverride    string // RFC3339; demo only — fake the clock
}

func (s *server) handleDeploy(w http.ResponseWriter, r *http.Request) {
	if !s.lim.allow(clientIP(r), time.Now()) {
		writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "slow down"})
		return
	}

	country := r.URL.Query().Get("country")
	if country == "" {
		country = s.defaultCountry
	}
	cc, ok := s.cfg[country]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown country"})
		return
	}

	now := time.Now()
	if s.nowOverride != "" {
		if t, err := time.Parse(time.RFC3339, s.nowOverride); err == nil {
			now = t
		}
	}

	safe, code := decide(now, cc)
	v := s.msgs.lookup(code)
	writeJSON(w, http.StatusOK, response{
		Country:  country,
		Timezone: cc.Timezone,
		Safe:     safe,
		Reason:   v.Reason,
		Message:  v.pick(s.rng),
	})
}

// router wires routes and wraps them in CORS.
func (s *server) router(allowedOrigins []string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/should-i-deploy", s.handleDeploy)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return cors(allowedOrigins, mux)
}

// cors allows browser callers whose Origin is in the allowlist. A response can
// only name ONE origin, so we echo back the request's Origin when it matches
// (that's how multiple allowed origins work). "*" allows any. Browser-enforced
// only — does not stop curl/Postman. Hygiene, not real auth.
func cors(allowed []string, next http.Handler) http.Handler {
	allowSet := make(map[string]bool, len(allowed))
	wildcard := false
	for _, o := range allowed {
		if o == "*" {
			wildcard = true
		}
		allowSet[o] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		switch {
		case wildcard:
			w.Header().Set("Access-Control-Allow-Origin", "*")
		case origin != "" && allowSet[origin]:
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
