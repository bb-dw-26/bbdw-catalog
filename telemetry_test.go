package catalog

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
)

func TestPanicTelemetry(t *testing.T) {
	var mu sync.Mutex
	var payloads []string
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		payloads = append(payloads, string(body))
		mu.Unlock()
		w.WriteHeader(200)
	}))
	defer collector.Close()
	t.Setenv("SENTRY_DSN", strings.Replace(collector.URL, "http://", "http://public@", 1)+"/1")
	t.Setenv("SENTRY_ENVIRONMENT", "test")
	t.Setenv("SENTRY_RELEASE", "test-release")
	t.Setenv("DEMO_SCENARIO", "SDK-SMOKE")
	t.Setenv("DEMO_RUN_ID", "local-test")
	if err := InitTelemetry(); err != nil {
		t.Fatal(err)
	}
	defer sentry.CurrentHub().BindClient(nil)
	handler := telemetryHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			problem(w, 404, "NOT_FOUND", "Ausente")
			return
		}
		if w.Header().Get("X-Request-ID") == "request-a" {
			panic(errors.New("Falha sintética de captura"))
		}
		panic("Falha sintética de captura")
	}))
	var wg sync.WaitGroup
	for _, id := range []string{"request-a", "request-b"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := httptest.NewRequest("GET", "/failure?secret=never-publish", nil)
			r.Header.Set("Authorization", "Bearer never-publish")
			w := httptest.NewRecorder()
			w.Header().Set("X-Request-ID", id)
			handler.ServeHTTP(w, r)
			if w.Code != 500 || !strings.Contains(w.Body.String(), id) {
				t.Errorf("Resposta incorreta: %s", w.Body.String())
			}
		}()
	}
	wg.Wait()
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/missing", nil))
	if w.Code != 404 {
		t.Fatal("404 alterado")
	}
	if !sentry.Flush(3 * time.Second) {
		t.Fatal("Eventos não enviados")
	}
	mu.Lock()
	defer mu.Unlock()
	ids := map[string]bool{}
	count := 0
	for _, payload := range payloads {
		if strings.Contains(payload, "never-publish") {
			t.Fatal("Dados HTTP vazaram")
		}
		for _, line := range strings.Split(payload, "\n") {
			var event map[string]any
			if json.Unmarshal([]byte(line), &event) != nil || (event["exception"] == nil && event["message"] == nil) {
				continue
			}
			count++
			tags := event["tags"].(map[string]any)
			ids[tags["request_id"].(string)] = true
			if tags["service"] != "bbdw-catalog" || event["release"] != "test-release" || event["request"] != nil || event["user"] != nil {
				t.Fatalf("Metadados incorretos")
			}
		}
	}
	if count != 2 || !ids["request-a"] || !ids["request-b"] {
		t.Fatalf("Eventos duplicados ou contexto misturado: %d, %v", count, ids)
	}
}
