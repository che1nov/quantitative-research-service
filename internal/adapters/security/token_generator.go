package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

// TokenGenerator генерирует криптостойкие токены.
type TokenGenerator struct {
	log logger.Logger
}

func NewTokenGenerator(log logger.Logger) *TokenGenerator {
	return &TokenGenerator{log: log}
}

// Generate создает безопасный токен указанного размера.
func (g *TokenGenerator) Generate(ctx context.Context, size int) (string, error) {
	g.log.DebugContext(ctx, "Генерация токена", "size", size)
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		g.log.ErrorContext(ctx, "Ошибка генерации токена", "error", err)
		return "", err
	}
	g.log.DebugContext(ctx, "Токен успешно сгенерирован")
	return base64.RawURLEncoding.EncodeToString(b), nil
}
