package memory

import (
	"context"
	"sync"
	"time"

	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

// SurveyAdapter хранит данные опросов в оперативной памяти.
type SurveyAdapter struct {
	mu       sync.RWMutex
	surveys  map[string]domain.Survey
	byToken  map[string]string
	sessions map[string]domain.AnswerSession
	log      logger.Logger
}

func NewSurveyAdapter(log logger.Logger) *SurveyAdapter {
	return &SurveyAdapter{
		surveys:  map[string]domain.Survey{},
		byToken:  map[string]string{},
		sessions: map[string]domain.AnswerSession{},
		log:      log,
	}
}

// CreateSurvey сохраняет новый опрос.
func (a *SurveyAdapter) CreateSurvey(ctx context.Context, survey domain.Survey) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.surveys[survey.ID] = survey
	a.byToken[survey.PublicToken] = survey.ID
	a.log.DebugContext(ctx, "Опрос сохранен в памяти", "survey_id", survey.ID)
	return nil
}

// UpdateSurvey обновляет существующий опрос.
func (a *SurveyAdapter) UpdateSurvey(ctx context.Context, survey domain.Survey) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.surveys[survey.ID]; !ok {
		a.log.WarnContext(ctx, "Опрос для обновления не найден", "survey_id", survey.ID)
		return domain.ErrSurveyNotFound
	}
	a.surveys[survey.ID] = survey
	a.log.DebugContext(ctx, "Опрос обновлен в памяти", "survey_id", survey.ID)
	return nil
}

// DeleteSurvey удаляет опрос.
func (a *SurveyAdapter) DeleteSurvey(ctx context.Context, surveyID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	survey, ok := a.surveys[surveyID]
	if !ok {
		a.log.WarnContext(ctx, "Опрос для удаления не найден", "survey_id", surveyID)
		return domain.ErrSurveyNotFound
	}
	delete(a.byToken, survey.PublicToken)
	delete(a.surveys, surveyID)
	a.log.DebugContext(ctx, "Опрос удален из памяти", "survey_id", surveyID)
	return nil
}

// GetSurveyByID возвращает опрос по идентификатору.
func (a *SurveyAdapter) GetSurveyByID(ctx context.Context, surveyID string) (domain.Survey, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	survey, ok := a.surveys[surveyID]
	if !ok {
		a.log.WarnContext(ctx, "Опрос по id не найден", "survey_id", surveyID)
		return domain.Survey{}, domain.ErrSurveyNotFound
	}
	a.log.DebugContext(ctx, "Опрос получен по id", "survey_id", surveyID)
	return survey, nil
}

// GetSurveyByToken возвращает опрос по публичному токену.
func (a *SurveyAdapter) GetSurveyByToken(ctx context.Context, token string) (domain.Survey, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	surveyID, ok := a.byToken[token]
	if !ok {
		a.log.WarnContext(ctx, "Опрос по токену не найден")
		return domain.Survey{}, domain.ErrSurveyNotFound
	}
	survey, ok := a.surveys[surveyID]
	if !ok {
		a.log.WarnContext(ctx, "Опрос по токену найден в индексе, но отсутствует в основном хранилище", "survey_id", surveyID)
		return domain.Survey{}, domain.ErrSurveyNotFound
	}
	a.log.DebugContext(ctx, "Опрос получен по токену", "survey_id", surveyID)
	return survey, nil
}

// ListSurveys возвращает список всех опросов.
func (a *SurveyAdapter) ListSurveys(ctx context.Context) ([]domain.Survey, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]domain.Survey, 0, len(a.surveys))
	for _, survey := range a.surveys {
		result = append(result, survey)
	}
	a.log.DebugContext(ctx, "Список опросов получен", "count", len(result))
	return result, nil
}

// CreateSession сохраняет новую сессию прохождения.
func (a *SurveyAdapter) CreateSession(ctx context.Context, session domain.AnswerSession) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.sessions[session.ID] = session
	a.log.DebugContext(ctx, "Сессия сохранена в памяти", "session_id", session.ID, "survey_id", session.SurveyID)
	return nil
}

// SaveSessionAnswers обновляет ответы и состояние сессии.
func (a *SurveyAdapter) SaveSessionAnswers(ctx context.Context, sessionID string, answers map[string][]string, completed bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	session, ok := a.sessions[sessionID]
	if !ok {
		a.log.WarnContext(ctx, "Сессия для сохранения ответов не найдена", "session_id", sessionID)
		return domain.ErrSessionNotFound
	}

	session.Answers = answers
	session.UpdatedAt = time.Now().UTC()
	session.Completed = completed
	if completed {
		finished := time.Now().UTC()
		session.FinishedAt = &finished
	}
	a.sessions[sessionID] = session
	a.log.DebugContext(ctx, "Ответы по сессии сохранены", "session_id", sessionID, "completed", completed)
	return nil
}

// GetSessionByID получает сессию по идентификатору.
func (a *SurveyAdapter) GetSessionByID(ctx context.Context, sessionID string) (domain.AnswerSession, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	session, ok := a.sessions[sessionID]
	if !ok {
		a.log.WarnContext(ctx, "Сессия не найдена", "session_id", sessionID)
		return domain.AnswerSession{}, domain.ErrSessionNotFound
	}
	a.log.DebugContext(ctx, "Сессия получена", "session_id", sessionID)
	return session, nil
}

// ListCompletedSessionsBySurvey получает завершенные прохождения.
func (a *SurveyAdapter) ListCompletedSessionsBySurvey(ctx context.Context, surveyID string) ([]domain.AnswerSession, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]domain.AnswerSession, 0)
	for _, session := range a.sessions {
		if session.SurveyID == surveyID && session.Completed {
			result = append(result, session)
		}
	}
	a.log.DebugContext(ctx, "Список завершенных сессий получен", "survey_id", surveyID, "count", len(result))
	return result, nil
}
