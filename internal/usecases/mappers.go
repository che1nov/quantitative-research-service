package usecases

import (
	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/internal/dto"
)

func mapQuestionsFromDTO(items []dto.CreateQuestionDTO) ([]domain.Question, error) {
	questions := make([]domain.Question, 0, len(items))
	for _, item := range items {
		q := domain.Question{ID: item.ID, Title: item.Title, Type: domain.QuestionType(item.Type)}
		for _, opt := range item.Options {
			q.Options = append(q.Options, domain.Option{ID: opt.ID, Text: opt.Text, Image: opt.Image})
		}
		if err := domain.ValidateQuestion(q); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func mapSurveyToDTO(baseURL string, survey domain.Survey) dto.SurveyOutput {
	questions := make([]dto.QuestionDTO, 0, len(survey.Questions))
	for _, q := range survey.Questions {
		opts := make([]dto.OptionDTO, 0, len(q.Options))
		for _, opt := range q.Options {
			opts = append(opts, dto.OptionDTO{ID: opt.ID, Text: opt.Text, Image: opt.Image})
		}
		questions = append(questions, dto.QuestionDTO{ID: q.ID, Title: q.Title, Type: string(q.Type), Options: opts})
	}

	return dto.SurveyOutput{
		ID:          survey.ID,
		Title:       survey.Title,
		Description: survey.Description,
		PublicLink:  baseURL + "/api/public/surveys/" + survey.PublicToken,
		Questions:   questions,
		CreatedAt:   survey.CreatedAt,
		UpdatedAt:   survey.UpdatedAt,
	}
}
