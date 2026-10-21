package catalog

import (
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
)

// InitTelemetry configura somente erros; a ausência de DSN desativa o envio.
func InitTelemetry() error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Environment:      os.Getenv("SENTRY_ENVIRONMENT"),
		Release:          os.Getenv("SENTRY_RELEASE"),
		SendDefaultPII:   false,
		AttachStacktrace: true,
		BeforeSend: func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
			event.Request = nil
			event.User = sentry.User{}
			event.Breadcrumbs = nil
			event.ServerName = "bbdw-catalog"
			return event
		},
	})
}

func telemetryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Um hub por requisição evita misturar os IDs de chamadas concorrentes.
		hub := sentry.CurrentHub().Clone()
		hub.Scope().SetTags(map[string]string{
			"service":    "bbdw-catalog",
			"request_id": w.Header().Get("X-Request-ID"),
			"scenario":   os.Getenv("DEMO_SCENARIO"),
			"demo_run":   os.Getenv("DEMO_RUN_ID"),
		})
		defer func() {
			if failure := recover(); failure != nil {
				hub.Recover(failure)
				problem(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Não foi possível concluir a operação.")
			}
		}()
		next.ServeHTTP(w, r.WithContext(sentry.SetHubOnContext(r.Context(), hub)))
	})
}
