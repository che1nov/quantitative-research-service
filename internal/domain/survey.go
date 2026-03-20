package domain

import (
	"strings"
	"time"
)

// QuestionType определяет тип вопроса.
type QuestionType string

const (
	QuestionSingleChoice QuestionType = "single_choice"
	QuestionMultiChoice  QuestionType = "multi_choice"
	QuestionFreeText     QuestionType = "free_text"
	QuestionImageChoice  QuestionType = "image_choice"
)

// Survey описывает опрос в системе.
type Survey struct {
	ID          string
	Title       string
	Description string
	Questions   []Question
	PublicToken string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Question описывает один вопрос опроса.
type Question struct {
	ID      string
	Title   string
	Type    QuestionType
	Options []Option
}

// Option описывает вариант ответа в вопросе.
type Option struct {
	ID    string
	Text  string
	Image string
}

// AnswerSession хранит незавершенные и завершенные ответы пользователя.
type AnswerSession struct {
	ID         string
	SurveyID   string
	Answers    map[string][]string
	Completed  bool
	StartedAt  time.Time
	UpdatedAt  time.Time
	FinishedAt *time.Time
}

// NewSurvey создает опрос и валидирует обязательные поля.
func NewSurvey(id, title, description, publicToken string, questions []Question) (Survey, error) {
	if strings.TrimSpace(title) == "" {
		return Survey{}, ErrInvalidSurveyTitle
	}
	if strings.TrimSpace(publicToken) == "" {
		return Survey{}, ErrInvalidToken
	}
	if len(questions) == 0 {
		return Survey{}, ErrInvalidQuestion
	}

	now := time.Now().UTC()
	return Survey{
		ID:          id,
		Title:       title,
		Description: description,
		Questions:   questions,
		PublicToken: publicToken,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ValidateQuestion проверяет корректность вопроса.
func ValidateQuestion(q Question) error {
	if strings.TrimSpace(q.ID) == "" || strings.TrimSpace(q.Title) == "" {
		return ErrInvalidQuestion
	}

	switch q.Type {
	case QuestionSingleChoice, QuestionMultiChoice, QuestionFreeText, QuestionImageChoice:
	default:
		return ErrInvalidQuestion
	}

	if q.Type != QuestionFreeText && len(q.Options) == 0 {
		return ErrInvalidQuestion
	}

	return nil
}
