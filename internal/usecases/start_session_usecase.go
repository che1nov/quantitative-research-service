package usecases

import (
	"context"
	"time"

	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/internal/dto"
)

// StartSessionUseCase запускает сессию прохождения опроса.
type StartSessionUseCase struct {
	storage SurveyStorage
	tokens  TokenGenerator
	log     AppLogger
}

func NewStartSessionUseCase(storage SurveyStorage, tokens TokenGenerator, log AppLogger) *StartSessionUseCase {
	return &StartSessionUseCase{storage: storage, tokens: tokens, log: log}
}

// Execute создает новую сессию для публичного опроса.
func (uc *StartSessionUseCase) Execute(ctx context.Context, input dto.StartSessionInput) (dto.StartSessionOutput, error) {
	uc.log.InfoContext(ctx, "Создание сессии прохождения")
	survey, err := uc.storage.GetSurveyByToken(ctx, input.PublicToken)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка поиска опроса по токену для создания сессии", "error", err)
		return dto.StartSessionOutput{}, err
	}

	sessionID, err := uc.tokens.Generate(ctx, 18)
	if err != nil {
		return dto.StartSessionOutput{}, err
	}
	now := time.Now().UTC()
	session := domain.AnswerSession{
		ID:        sessionID,
		SurveyID:  survey.ID,
		Answers:   map[string][]string{},
		StartedAt: now,
		UpdatedAt: now,
	}

	if err := uc.storage.CreateSession(ctx, session); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка создания сессии", "error", err)
		return dto.StartSessionOutput{}, err
	}

	uc.log.InfoContext(ctx, "Сессия успешно создана", "session_id", sessionID, "survey_id", survey.ID)
	return dto.StartSessionOutput{SessionID: sessionID, SurveyID: survey.ID}, nil
}
