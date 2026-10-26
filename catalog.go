package catalog

import (
	"crypto/rand"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

//go:embed specs/fixtures.json
var fixtures []byte

type Product struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	PriceCents   int      `json:"priceCents"`
	Currency     string   `json:"currency"`
	Available    bool     `json:"available"`
	ImageURLs    []string `json:"imageUrls"`
	ThumbnailURL *string  `json:"thumbnailUrl"`
}

var validID = regexp.MustCompile(`^[a-zA-Z0-9._-]{1,128}$`)

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func problem(w http.ResponseWriter, status int, code, message string) {
	respond(w, status, map[string]any{
		"error":     map[string]string{"code": code, "message": message},
		"requestId": w.Header().Get("X-Request-ID"),
	})
}

func NewHandler(origin string) (http.Handler, error) {
	var products []Product
	if err := json.Unmarshal(fixtures, &products); err != nil {
		return nil, fmt.Errorf("fixtures inválidas: %w", err)
	}
	for i := range products {
		if products[i].ImageURLs == nil {
			products[i].ImageURLs = []string{}
		}
		if len(products[i].ImageURLs) > 0 {
			products[i].ThumbnailURL = &products[i].ImageURLs[0]
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, map[string]any{"items": products})
	})
	mux.HandleFunc("GET /products/{id}", func(w http.ResponseWriter, r *http.Request) {
		for _, product := range products {
			if product.ID == r.PathValue("id") {
				respond(w, http.StatusOK, product)
				return
			}
		}
		problem(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Produto não encontrado.")
	})
	instrumented := telemetryHandler(mux)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if !validID.MatchString(id) {
			id = rand.Text()
		}
		w.Header().Set("X-Request-ID", id)
		w.Header().Set("Vary", "Origin")
		if r.Header.Get("Origin") == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Padronizar inclusive rotas desconhecidas, sem expor respostas internas.
		_, pattern := mux.Handler(r)
		if pattern == "" {
			problem(w, http.StatusNotFound, "NOT_FOUND", "Rota não encontrada.")
			return
		}
		instrumented.ServeHTTP(w, r)
	}), nil
}
