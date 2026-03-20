package domain

import "errors"

var (
	ErrSurveyNotFound     = errors.New("опрос не найден")
	ErrSessionNotFound    = errors.New("сессия не найдена")
	ErrInvalidQuestion    = errors.New("некорректный вопрос")
	ErrInvalidSurveyTitle = errors.New("название опроса не может быть пустым")
	ErrInvalidToken       = errors.New("некорректная ссылка опроса")
)
