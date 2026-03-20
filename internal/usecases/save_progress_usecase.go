package usecases

import (
	"context"

	"github.com/che1nov/quantitative-research-service/internal/dto"
)

// SaveProgressUseCase сохраняет промежуточный прогресс прохождения.
type SaveProgressUseCase struct {
	storage SurveyStorage
	log     AppLogger
}

func NewSaveProgressUseCase(storage SurveyStorage, log AppLogger) *SaveProgressUseCase {
	return &SaveProgressUseCase{storage: storage, log: log}
}

// Execute сохраняет незавершенные ответы по сессии.
func (uc *SaveProgressUseCase) Execute(ctx context.Context, input dto.SaveProgressInput) error {
	uc.log.InfoContext(ctx, "Сохранение прогресса", "session_id", input.SessionID)
	if err := uc.storage.SaveSessionAnswers(ctx, input.SessionID, input.Answers, false); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка сохранения прогресса", "error", err)
		return err
	}
	uc.log.InfoContext(ctx, "Прогресс успешно сохранен", "session_id", input.SessionID)
	return nil
}
