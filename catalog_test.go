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
				if body["thumbnailUrl"] != nil {
					t.Fatal("miniatura deve ser null")
				}
				if images, ok := body["imageUrls"].([]any); !ok || len(images) != 0 {
					t.Fatal("imagens devem ser []")
				}
			case "/products/product-standard":
				if body["priceCents"] != float64(4000) || body["thumbnailUrl"] != "/demo-assets/caneca.svg" {
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
