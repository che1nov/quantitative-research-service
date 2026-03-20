package usecases

import (
	"context"

	"github.com/che1nov/quantitative-research-service/internal/dto"
)

// GetPublicSurveyUseCase получает публичный опрос по токену ссылки.
type GetPublicSurveyUseCase struct {
	storage   SurveyStorage
	log       AppLogger
	publicURL string
}

func NewGetPublicSurveyUseCase(storage SurveyStorage, log AppLogger, publicURL string) *GetPublicSurveyUseCase {
	return &GetPublicSurveyUseCase{storage: storage, log: log, publicURL: publicURL}
}

// Execute получает публичный опрос по токену.
func (uc *GetPublicSurveyUseCase) Execute(ctx context.Context, token string) (dto.SurveyOutput, error) {
	uc.log.InfoContext(ctx, "Запрос публичного опроса")
	survey, err := uc.storage.GetSurveyByToken(ctx, token)
	if err != nil {
		uc.log.WarnContext(ctx, "Публичный опрос не найден", "token", token)
		return dto.SurveyOutput{}, err
	}
	uc.log.InfoContext(ctx, "Публичный опрос успешно выдан", "survey_id", survey.ID)
	return mapSurveyToDTO(uc.publicURL, survey), nil
}
