package middleware

import "net/http"

// const allowedDomain = "http://localhost:8010"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// w.Header().Set("Access-Control-Allow-Origin", allowedDomain)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PUTCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PUTCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			return
		}

		next.ServeHTTP(w, r)
	})
}
