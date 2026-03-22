package usecases

import (
	"context"

	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

// SurveyStorage описывает атомарные операции доступа к хранилищу опросов.
type SurveyStorage interface {
	CreateSurvey(ctx context.Context, survey domain.Survey) error
	UpdateSurvey(ctx context.Context, survey domain.Survey) error
	DeleteSurvey(ctx context.Context, surveyID string) error
	GetSurveyByID(ctx context.Context, surveyID string) (domain.Survey, error)
	GetSurveyByToken(ctx context.Context, token string) (domain.Survey, error)
	ListSurveys(ctx context.Context) ([]domain.Survey, error)
	CreateSession(ctx context.Context, session domain.AnswerSession) error
	SaveSessionAnswers(ctx context.Context, sessionID string, answers map[string][]string, completed bool) error
	GetSessionByID(ctx context.Context, sessionID string) (domain.AnswerSession, error)
	ListCompletedSessionsBySurvey(ctx context.Context, surveyID string) ([]domain.AnswerSession, error)
}

// TokenGenerator генерирует криптостойкие токены ссылок и идентификаторы.
type TokenGenerator interface {
	Generate(ctx context.Context, size int) (string, error)
}

// XLSExporter экспортирует результаты в формат .xls.
type XLSExporter interface {
	BuildSurveyResultsXLS(ctx context.Context, survey domain.Survey, sessions []domain.AnswerSession) ([]byte, error)
}

// UserStorage описывает атомарные операции с пользователями.
type UserStorage interface {
	GetUserByVKID(ctx context.Context, vkID int64) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User) error
}

// JWTSigner создает JWT токен для пользователя.
type JWTSigner interface {
	SignUserToken(ctx context.Context, user domain.User) (string, error)
}

// AppLogger описывает логгер use case слоя.
type AppLogger = logger.Logger
