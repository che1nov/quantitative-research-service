package dto

import "time"

// CreateSurveyInput описывает запрос на создание опроса.
type CreateSurveyInput struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Questions   []CreateQuestionDTO `json:"questions"`
}

// UpdateSurveyInput описывает запрос на обновление опроса.
type UpdateSurveyInput struct {
	SurveyID    string              `json:"-"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Questions   []CreateQuestionDTO `json:"questions"`
}

// CreateQuestionDTO описывает вопрос во входных данных.
type CreateQuestionDTO struct {
	ID      string            `json:"id"`
	Title   string            `json:"title"`
	Type    string            `json:"type"`
	Options []CreateOptionDTO `json:"options"`
}

// CreateOptionDTO описывает вариант ответа во входных данных.
type CreateOptionDTO struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Image string `json:"image"`
}

// SurveyOutput описывает ответ с опросом.
type SurveyOutput struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	PublicLink  string        `json:"public_link"`
	Questions   []QuestionDTO `json:"questions"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// QuestionDTO описывает вопрос в ответе.
type QuestionDTO struct {
	ID      string      `json:"id"`
	Title   string      `json:"title"`
	Type    string      `json:"type"`
	Options []OptionDTO `json:"options"`
}

// OptionDTO описывает вариант ответа в ответе.
type OptionDTO struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Image string `json:"image"`
}

// StartSessionInput описывает запрос на начало прохождения опроса.
type StartSessionInput struct {
	PublicToken string `json:"public_token"`
}

// StartSessionOutput описывает созданную сессию прохождения.
type StartSessionOutput struct {
	SessionID string `json:"session_id"`
	SurveyID  string `json:"survey_id"`
}

// SaveProgressInput описывает частичное сохранение ответов.
type SaveProgressInput struct {
	SessionID string              `json:"session_id"`
	Answers   map[string][]string `json:"answers"`
}

// SubmitAnswersInput описывает завершение опроса.
type SubmitAnswersInput struct {
	SessionID string              `json:"session_id"`
	Answers   map[string][]string `json:"answers"`
}

// SurveyResultOutput описывает агрегированный результат по опросу.
type SurveyResultOutput struct {
	SurveyID string                         `json:"survey_id"`
	Total    int                            `json:"total"`
	Answers  map[string]map[string]int      `json:"answers"`
	Texts    map[string][]EscapedTextAnswer `json:"texts"`
}

// EscapedTextAnswer хранит экранированный текстовый ответ.
type EscapedTextAnswer struct {
	SessionID string `json:"session_id"`
	Text      string `json:"text"`
}
