package usecases

import (
	"context"
	"time"

	"github.com/che1nov/quantitative-research-service/internal/dto"
)

// UpdateSurveyUseCase обновляет существующий опрос.
type UpdateSurveyUseCase struct {
	storage   SurveyStorage
	log       AppLogger
	publicURL string
}

func NewUpdateSurveyUseCase(storage SurveyStorage, log AppLogger, publicURL string) *UpdateSurveyUseCase {
	return &UpdateSurveyUseCase{storage: storage, log: log, publicURL: publicURL}
}

// Execute обновляет заголовок, описание и вопросы опроса.
func (uc *UpdateSurveyUseCase) Execute(ctx context.Context, input dto.UpdateSurveyInput) (dto.SurveyOutput, error) {
	uc.log.InfoContext(ctx, "Обновление опроса", "survey_id", input.SurveyID)

	survey, err := uc.storage.GetSurveyByID(ctx, input.SurveyID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка чтения опроса перед обновлением", "error", err)
		return dto.SurveyOutput{}, err
	}
	questions, err := mapQuestionsFromDTO(input.Questions)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка валидации вопросов при обновлении", "error", err)
		return dto.SurveyOutput{}, err
	}

	survey.Title = input.Title
	survey.Description = input.Description
	survey.Questions = questions
	survey.UpdatedAt = time.Now().UTC()

	if err := uc.storage.UpdateSurvey(ctx, survey); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка обновления опроса", "error", err)
		return dto.SurveyOutput{}, err
	}

	uc.log.InfoContext(ctx, "Опрос успешно обновлен", "survey_id", survey.ID)
	return mapSurveyToDTO(uc.publicURL, survey), nil
}
