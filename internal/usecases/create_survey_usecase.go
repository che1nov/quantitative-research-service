package usecases

import (
	"context"

	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/internal/dto"
)

// CreateSurveyUseCase создает новый опрос для личного кабинета.
type CreateSurveyUseCase struct {
	storage   SurveyStorage
	tokens    TokenGenerator
	log       AppLogger
	publicURL string
}

func NewCreateSurveyUseCase(storage SurveyStorage, tokens TokenGenerator, log AppLogger, publicURL string) *CreateSurveyUseCase {
	return &CreateSurveyUseCase{storage: storage, tokens: tokens, log: log, publicURL: publicURL}
}

// Execute создает новый опрос с уникальной публичной ссылкой.
func (uc *CreateSurveyUseCase) Execute(ctx context.Context, input dto.CreateSurveyInput) (dto.SurveyOutput, error) {
	uc.log.InfoContext(ctx, "Создание опроса", "title", input.Title)

	questions, err := mapQuestionsFromDTO(input.Questions)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка валидации вопросов", "error", err)
		return dto.SurveyOutput{}, err
	}

	surveyID, err := uc.tokens.Generate(ctx, 12)
	if err != nil {
		return dto.SurveyOutput{}, err
	}
	publicToken, err := uc.tokens.Generate(ctx, 32)
	if err != nil {
		return dto.SurveyOutput{}, err
	}

	survey, err := domain.NewSurvey(surveyID, input.Title, input.Description, publicToken, questions)
	if err != nil {
		return dto.SurveyOutput{}, err
	}
	if err := uc.storage.CreateSurvey(ctx, survey); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка сохранения опроса", "error", err)
		return dto.SurveyOutput{}, err
	}

	uc.log.InfoContext(ctx, "Опрос успешно создан", "survey_id", survey.ID)
	return mapSurveyToDTO(uc.publicURL, survey), nil
}
