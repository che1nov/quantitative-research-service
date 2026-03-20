package http

import (
	"net/http"
	"strings"
)

// NewRouter собирает HTTP маршруты приложения.
func NewRouter(
	handler *SurveyHandler,
	miniApp *MiniAppHandler,
	auth *AuthMiddleware,
	csrf *CSRFMiddleware,
	external *ExternalAPIMiddleware,
	security *SecurityHeadersMiddleware,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /miniapp", miniApp.Index)
	mux.HandleFunc("GET /miniapp/", miniApp.Assets)

	mux.HandleFunc("GET /api/cabinet/surveys", handler.ListSurveys)
	mux.HandleFunc("POST /api/cabinet/surveys", handler.CreateSurvey)
	mux.HandleFunc("PUT /api/cabinet/surveys/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/results") || strings.HasSuffix(r.URL.Path, "/export") {
			http.NotFound(w, r)
			return
		}
		handler.UpdateSurvey(w, r)
	})
	mux.HandleFunc("DELETE /api/cabinet/surveys/", handler.DeleteSurvey)
	mux.HandleFunc("GET /api/cabinet/surveys/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/results") {
			handler.GetResults(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/export") {
			handler.ExportResults(w, r)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("GET /api/public/surveys/", handler.GetPublicSurvey)
	mux.HandleFunc("POST /api/public/sessions", handler.StartSession)
	mux.HandleFunc("PUT /api/public/sessions/progress", handler.SaveProgress)
	mux.HandleFunc("POST /api/public/sessions/submit", handler.SubmitAnswers)
	mux.HandleFunc("GET /api/external/surveys/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/results") {
			handler.GetExternalResults(w, r)
			return
		}
		http.NotFound(w, r)
	})

	wrapped := security.Wrap(external.Wrap(auth.WrapInternal(csrf.Wrap(mux))))
	return wrapped
}
