package usecases

import (
	"context"

	"github.com/che1nov/quantitative-research-service/internal/dto"
)

// ListSurveysUseCase возвращает список опросов для личного кабинета.
type ListSurveysUseCase struct {
	storage   SurveyStorage
	log       AppLogger
	publicURL string
}

func NewListSurveysUseCase(storage SurveyStorage, log AppLogger, publicURL string) *ListSurveysUseCase {
	return &ListSurveysUseCase{storage: storage, log: log, publicURL: publicURL}
}

// Execute возвращает все опросы пользователя.
func (uc *ListSurveysUseCase) Execute(ctx context.Context) ([]dto.SurveyOutput, error) {
	uc.log.InfoContext(ctx, "Запрос списка опросов")
	surveys, err := uc.storage.ListSurveys(ctx)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка чтения опросов", "error", err)
		return nil, err
	}

	result := make([]dto.SurveyOutput, 0, len(surveys))
	for _, survey := range surveys {
		result = append(result, mapSurveyToDTO(uc.publicURL, survey))
	}
	uc.log.InfoContext(ctx, "Список опросов успешно подготовлен", "count", len(result))
	return result, nil
}
