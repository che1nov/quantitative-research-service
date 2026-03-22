package security

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

// JWTSigner подписывает JWT токен для пользовательской сессии.
type JWTSigner struct {
	secret []byte
	log    logger.Logger
}

func NewJWTSigner(secret string, log logger.Logger) *JWTSigner {
	return &JWTSigner{secret: []byte(secret), log: log}
}

// SignUserToken формирует и подписывает JWT HS256.
func (s *JWTSigner) SignUserToken(ctx context.Context, user domain.User) (string, error) {
	header, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		s.log.ErrorContext(ctx, "Ошибка сериализации JWT header", "error", err)
		return "", err
	}

	now := time.Now().UTC()
	payload, err := json.Marshal(map[string]any{
		"sub":   user.ID,
		"vk_id": user.VKID,
		"iat":   now.Unix(),
		"exp":   now.Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		s.log.ErrorContext(ctx, "Ошибка сериализации JWT payload", "error", err)
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(header)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	unsigned := fmt.Sprintf("%s.%s", encodedHeader, encodedPayload)

	h := hmac.New(sha256.New, s.secret)
	if _, err := h.Write([]byte(unsigned)); err != nil {
		s.log.ErrorContext(ctx, "Ошибка подписи JWT", "error", err)
		return "", err
	}
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s.%s", unsigned, signature), nil
}
