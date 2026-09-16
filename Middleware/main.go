package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const authLBURL = "http://authlb:8080"
const logsLBURL = "http://logslb:8080"
const usersLBURL = "http://userslb:8080"

func routeTarget(path string) string {
	switch {
	case strings.HasPrefix(path, "/auth"):
		after, _ := strings.CutPrefix(path, "/auth")
		return (authLBURL + after)
	case strings.HasPrefix(path, "/logs"):
		after, _ := strings.CutPrefix(path, "/logs")
		return (logsLBURL + after)
	case strings.HasPrefix(path, "/users"):
		after, _ := strings.CutPrefix(path, "/users")
		return (usersLBURL + after)
	default:
		return ""
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	target := routeTarget(r.URL.Path)
	if target == "" {
		http.Error(w, "No target service found for this path", http.StatusNotFound)
		return
	}

	forwardReq, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to build upstream request: %v", err), http.StatusBadGateway)
		return
	}

	resp, err := http.DefaultClient.Do(forwardReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to reach backend service: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("failed to write response body: %v", err)
	}
}

// definir CORS para permitir solicitudes desde el frontend
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirectHandler)

	log.Println("Middleware running on :80")
	if err := http.ListenAndServe(":80", enableCORS(mux)); err != nil {
		log.Fatal(err)
	}
}
