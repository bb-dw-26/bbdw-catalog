package catalog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProducts(t *testing.T) {
	handler, err := NewHandler("http://localhost:3000")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		path   string
		status int
	}{
		{"/health", 200},
		{"/products", 200},
		{"/products/product-no-image", 200},
		{"/products/product-standard", 200},
		{"/products/missing", 404},
	} {
		t.Run(tc.path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tc.path, nil)
			request.Header.Set("X-Request-ID", "catalog-test")
			request.Header.Set("Origin", "http://localhost:3000")
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != tc.status {
				t.Fatalf("status = %d", response.Code)
			}
			if response.Header().Get("X-Request-ID") != "catalog-test" {
				t.Fatal("request ID perdido")
			}
			if response.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
				t.Fatal("CORS ausente")
			}
			var body map[string]any
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			switch tc.path {
			case "/products/product-no-image":
				if body["thumbnailUrl"] != "/demo-assets/caderno.png" {
					t.Fatal("caderno deve exibir sua fotografia")
				}
			case "/products/product-standard":
				if body["priceCents"] != float64(4000) || body["thumbnailUrl"] != "/demo-assets/caneca.png" {
					t.Fatal("produto incorreto")
				}
			case "/products/missing":
				if body["error"].(map[string]any)["code"] != "PRODUCT_NOT_FOUND" {
					t.Fatal("código incorreto")
				}
			}
		})
	}
}

func TestProductWithoutImages(t *testing.T) {
	// A cobertura do contrato sem imagem não depende dos produtos da vitrine.
	original := fixtures
	defer func() { fixtures = original }()
	for _, images := range []string{"[]", "null"} {
		fixtures = []byte(`[{"id":"without-photo","name":"Sem foto","priceCents":1000,"currency":"BRL","available":true,"imageUrls":` + images + `}]`)
		handler, err := NewHandler("http://localhost:3000")
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/products/without-photo", nil))
		var body Product
		if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if response.Code != 200 || body.ThumbnailURL != nil || body.ImageURLs == nil || len(body.ImageURLs) != 0 {
			t.Fatalf("contrato sem imagem inválido: %s", response.Body.String())
		}
	}
}
