package usecases

import (
	"context"

	"github.com/che1nov/quantitative-research-service/internal/dto"
)

// SubmitAnswersUseCase завершает сессию и фиксирует итоговые ответы.
type SubmitAnswersUseCase struct {
	storage SurveyStorage
	log     AppLogger
}

func NewSubmitAnswersUseCase(storage SurveyStorage, log AppLogger) *SubmitAnswersUseCase {
	return &SubmitAnswersUseCase{storage: storage, log: log}
}

// Execute завершает прохождение опроса.
func (uc *SubmitAnswersUseCase) Execute(ctx context.Context, input dto.SubmitAnswersInput) error {
	uc.log.InfoContext(ctx, "Завершение опроса", "session_id", input.SessionID)
	if err := uc.storage.SaveSessionAnswers(ctx, input.SessionID, input.Answers, true); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка завершения опроса", "error", err)
		return err
	}
	uc.log.InfoContext(ctx, "Опрос успешно завершен", "session_id", input.SessionID)
	return nil
}
