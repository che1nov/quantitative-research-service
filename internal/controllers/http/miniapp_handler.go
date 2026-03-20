package http

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

//go:embed web/*
var miniAppFiles embed.FS

// MiniAppHandler раздает статику демо мини-приложения.
type MiniAppHandler struct {
	assets http.Handler
	log    logger.Logger
}

func NewMiniAppHandler(log logger.Logger) (*MiniAppHandler, error) {
	sub, err := fs.Sub(miniAppFiles, "web")
	if err != nil {
		return nil, err
	}

	return &MiniAppHandler{
		assets: http.StripPrefix("/miniapp/", http.FileServer(http.FS(sub))),
		log:    log,
	}, nil
}

// Index отдает главную страницу мини-приложения.
func (h *MiniAppHandler) Index(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "Выдача страницы miniapp")
	http.ServeFileFS(w, r, miniAppFiles, "web/index.html")
}

// Assets отдает css/js ресурсы мини-приложения.
func (h *MiniAppHandler) Assets(w http.ResponseWriter, r *http.Request) {
	h.log.DebugContext(r.Context(), "Выдача ассета miniapp", "path", r.URL.Path)
	h.assets.ServeHTTP(w, r)
}
