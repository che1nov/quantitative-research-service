package domain_test

import (
	"testing"

	"github.com/che1nov/quantitative-research-service/internal/domain"
)

func TestNewSurveySuccess(t *testing.T) {
	questions := []domain.Question{{
		ID:      "q1",
		Title:   "Выберите один вариант",
		Type:    domain.QuestionSingleChoice,
		Options: []domain.Option{{ID: "o1", Text: "Да"}},
	}}

	survey, err := domain.NewSurvey("survey1", "Тестовый опрос", "Описание", "public-token", questions)
	if err != nil {
		t.Fatalf("ожидалось успешное создание, получена ошибка: %v", err)
	}
	if survey.Title != "Тестовый опрос" {
		t.Fatalf("некорректный title: %s", survey.Title)
	}
}

func TestValidateQuestionInvalidType(t *testing.T) {
	err := domain.ValidateQuestion(domain.Question{ID: "q1", Title: "Текст", Type: "invalid"})
	if err == nil {
		t.Fatal("ожидалась ошибка валидации")
	}
}
