package usecases_test

import (
	"context"
	"testing"

	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/internal/dto"
	"github.com/che1nov/quantitative-research-service/internal/usecases"
	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

type storageStub struct {
	survey domain.Survey
}

func (s *storageStub) CreateSurvey(_ context.Context, survey domain.Survey) error {
	s.survey = survey
	return nil
}
func (s *storageStub) UpdateSurvey(context.Context, domain.Survey) error { return nil }
func (s *storageStub) DeleteSurvey(context.Context, string) error        { return nil }
func (s *storageStub) GetSurveyByID(context.Context, string) (domain.Survey, error) {
	return domain.Survey{}, nil
}
func (s *storageStub) GetSurveyByToken(context.Context, string) (domain.Survey, error) {
	return domain.Survey{}, nil
}
func (s *storageStub) ListSurveys(context.Context) ([]domain.Survey, error) { return nil, nil }
func (s *storageStub) CreateSession(context.Context, domain.AnswerSession) error {
	return nil
}
func (s *storageStub) SaveSessionAnswers(context.Context, string, map[string][]string, bool) error {
	return nil
}
func (s *storageStub) GetSessionByID(context.Context, string) (domain.AnswerSession, error) {
	return domain.AnswerSession{}, nil
}
func (s *storageStub) ListCompletedSessionsBySurvey(context.Context, string) ([]domain.AnswerSession, error) {
	return nil, nil
}

type tokenStub struct {
	calls int
}

func (t *tokenStub) Generate(_ context.Context, _ int) (string, error) {
	t.calls++
	if t.calls == 1 {
		return "survey-id", nil
	}
	return "public-token", nil
}

func TestCreateSurveyUseCase(t *testing.T) {
	storage := &storageStub{}
	tokens := &tokenStub{}
	uc := usecases.NewCreateSurveyUseCase(storage, tokens, logger.New("error"), "http://localhost:8080")

	out, err := uc.Execute(context.Background(), dto.CreateSurveyInput{
		Title:       "Новый опрос",
		Description: "Описание",
		Questions: []dto.CreateQuestionDTO{{
			ID:    "q1",
			Title: "Ваш ответ",
			Type:  string(domain.QuestionFreeText),
		}},
	})
	if err != nil {
		t.Fatalf("ошибка создания опроса: %v", err)
	}
	if out.ID == "" || storage.survey.PublicToken == "" {
		t.Fatalf("опрос создан некорректно: %+v", out)
	}
}
