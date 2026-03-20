package usecases

import (
	"context"
	"html"

	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/internal/dto"
)

// GetResultsUseCase агрегирует результаты прохождений по опросу.
type GetResultsUseCase struct {
	storage SurveyStorage
	log     AppLogger
}

func NewGetResultsUseCase(storage SurveyStorage, log AppLogger) *GetResultsUseCase {
	return &GetResultsUseCase{storage: storage, log: log}
}

// Execute собирает агрегированную статистику и текстовые ответы.
func (uc *GetResultsUseCase) Execute(ctx context.Context, surveyID string) (dto.SurveyResultOutput, error) {
	uc.log.InfoContext(ctx, "Сбор результатов опроса", "survey_id", surveyID)
	sessions, err := uc.storage.ListCompletedSessionsBySurvey(ctx, surveyID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка чтения завершенных сессий", "error", err)
		return dto.SurveyResultOutput{}, err
	}
	output := dto.SurveyResultOutput{SurveyID: surveyID, Total: len(sessions), Answers: map[string]map[string]int{}, Texts: map[string][]dto.EscapedTextAnswer{}}

	for _, session := range sessions {
		for questionID, values := range session.Answers {
			if len(values) == 0 {
				continue
			}
			if output.Answers[questionID] == nil {
				output.Answers[questionID] = map[string]int{}
			}
			for _, value := range values {
				if isPlainTextAnswer(value) {
					escaped := html.EscapeString(value)
					output.Texts[questionID] = append(output.Texts[questionID], dto.EscapedTextAnswer{SessionID: session.ID, Text: escaped})
					continue
				}
				output.Answers[questionID][value]++
			}
		}
	}

	uc.log.InfoContext(ctx, "Результаты опроса подготовлены", "survey_id", surveyID, "sessions", len(sessions))
	return output, nil
}

// isPlainTextAnswer определяет свободный текст для дополнительного экранирования.
func isPlainTextAnswer(value string) bool {
	for _, r := range value {
		if r == ' ' || r == '<' || r == '>' || r == '&' {
			return true
		}
	}
	return len(value) > 24
}

var _ = domain.ErrSurveyNotFound
