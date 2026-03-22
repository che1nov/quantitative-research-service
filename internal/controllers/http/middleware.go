package http

import (
	"net/http"
	"strings"

	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

// AuthMiddleware проверяет доступ внутренних пользователей в личный кабинет.
type AuthMiddleware struct {
	log logger.Logger
}

func NewAuthMiddleware(log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{log: log}
}

// WrapInternal ограничивает доступ к внутренним endpoint по служебному заголовку.
func (m *AuthMiddleware) WrapInternal(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/cabinet") {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			m.log.DebugContext(r.Context(), "Доступ в кабинет по JWT", "path", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("X-Internal-User") == "" {
			m.log.WarnContext(r.Context(), "Доступ в кабинет без внутреннего пользователя")
			http.Error(w, "неавторизован", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CSRFMiddleware валидирует CSRF токен в методах изменения состояния.
type CSRFMiddleware struct {
	token string
	log   logger.Logger
}

func NewCSRFMiddleware(token string, log logger.Logger) *CSRFMiddleware {
	return &CSRFMiddleware{token: token, log: log}
}

// Wrap проверяет CSRF токен только для внутренних маршрутов и небезопасных методов.
func (m *CSRFMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/cabinet") {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			m.log.DebugContext(r.Context(), "CSRF проверка пропущена для JWT запроса", "path", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("X-CSRF-Token") != m.token {
			m.log.WarnContext(r.Context(), "Невалидный CSRF токен")
			http.Error(w, "csrf token required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ExternalAPIMiddleware проверяет токен интеграций внешних систем.
type ExternalAPIMiddleware struct {
	token string
	log   logger.Logger
}

func NewExternalAPIMiddleware(token string, log logger.Logger) *ExternalAPIMiddleware {
	return &ExternalAPIMiddleware{token: token, log: log}
}

// Wrap валидирует заголовок внешнего API.
func (m *ExternalAPIMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/external") {
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("X-External-Token") != m.token {
			m.log.WarnContext(r.Context(), "Невалидный токен внешнего API")
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware добавляет CORS заголовки для браузерных запросов.
type CORSMiddleware struct {
	allowOrigin string
	log         logger.Logger
}

func NewCORSMiddleware(allowOrigin string, log logger.Logger) *CORSMiddleware {
	return &CORSMiddleware{allowOrigin: allowOrigin, log: log}
}

// Wrap выставляет CORS заголовки и обрабатывает preflight.
func (m *CORSMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := m.allowOrigin
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-CSRF-Token,X-Internal-User,X-External-Token")
		w.Header().Set("Access-Control-Max-Age", "600")

		if r.Method == http.MethodOptions {
			m.log.DebugContext(r.Context(), "CORS preflight обработан", "path", r.URL.Path)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware выставляет базовые заголовки защиты.
type SecurityHeadersMiddleware struct {
	log logger.Logger
}

func NewSecurityHeadersMiddleware(log logger.Logger) *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{log: log}
}

// Wrap добавляет безопасные заголовки для всех HTTP ответов.
func (m *SecurityHeadersMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'self' https://vk.com https://*.vk.com https://vk.ru https://*.vk.ru https://ok.ru https://*.ok.ru",
		)
		m.log.DebugContext(r.Context(), "Проставлены security headers", "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
