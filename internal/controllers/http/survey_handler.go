package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/internal/dto"
	"github.com/che1nov/quantitative-research-service/internal/usecases"
	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

// SurveyHandler объединяет обработчики endpoint сервиса опросов.
type SurveyHandler struct {
	createSurvey *usecases.CreateSurveyUseCase
	updateSurvey *usecases.UpdateSurveyUseCase
	deleteSurvey *usecases.DeleteSurveyUseCase
	listSurveys  *usecases.ListSurveysUseCase
	getPublic    *usecases.GetPublicSurveyUseCase
	startSession *usecases.StartSessionUseCase
	saveProgress *usecases.SaveProgressUseCase
	submit       *usecases.SubmitAnswersUseCase
	getResults   *usecases.GetResultsUseCase
	export       *usecases.ExportResultsUseCase
	log          logger.Logger
}

func NewSurveyHandler(
	createSurvey *usecases.CreateSurveyUseCase,
	updateSurvey *usecases.UpdateSurveyUseCase,
	deleteSurvey *usecases.DeleteSurveyUseCase,
	listSurveys *usecases.ListSurveysUseCase,
	getPublic *usecases.GetPublicSurveyUseCase,
	startSession *usecases.StartSessionUseCase,
	saveProgress *usecases.SaveProgressUseCase,
	submit *usecases.SubmitAnswersUseCase,
	getResults *usecases.GetResultsUseCase,
	export *usecases.ExportResultsUseCase,
	log logger.Logger,
) *SurveyHandler {
	return &SurveyHandler{
		createSurvey: createSurvey,
		updateSurvey: updateSurvey,
		deleteSurvey: deleteSurvey,
		listSurveys:  listSurveys,
		getPublic:    getPublic,
		startSession: startSession,
		saveProgress: saveProgress,
		submit:       submit,
		getResults:   getResults,
		export:       export,
		log:          log,
	}
}

// CreateSurvey создает новый опрос.
func (h *SurveyHandler) CreateSurvey(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос CreateSurvey")
	var input dto.CreateSurveyInput
	if err := decodeJSON(r, &input); err != nil {
		h.log.WarnContext(r.Context(), "Ошибка декодирования CreateSurvey", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.createSurvey.Execute(r.Context(), input)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.log.InfoContext(r.Context(), "HTTP CreateSurvey выполнен", "survey_id", output.ID)
	writeJSON(w, http.StatusCreated, output)
}

// ListSurveys возвращает список опросов.
func (h *SurveyHandler) ListSurveys(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос ListSurveys")
	output, err := h.listSurveys.Execute(r.Context())
	if err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.log.InfoContext(r.Context(), "HTTP ListSurveys выполнен", "count", len(output))
	writeJSON(w, http.StatusOK, output)
}

// UpdateSurvey обновляет опрос по идентификатору.
func (h *SurveyHandler) UpdateSurvey(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос UpdateSurvey")
	surveyID := strings.TrimPrefix(r.URL.Path, "/api/cabinet/surveys/")
	if surveyID == "" {
		h.log.WarnContext(r.Context(), "Отсутствует survey id в UpdateSurvey")
		http.Error(w, "survey id required", http.StatusBadRequest)
		return
	}

	var input dto.UpdateSurveyInput
	if err := decodeJSON(r, &input); err != nil {
		h.log.WarnContext(r.Context(), "Ошибка декодирования UpdateSurvey", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input.SurveyID = surveyID

	output, err := h.updateSurvey.Execute(r.Context(), input)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.log.InfoContext(r.Context(), "HTTP UpdateSurvey выполнен", "survey_id", output.ID)
	writeJSON(w, http.StatusOK, output)
}

// DeleteSurvey удаляет опрос.
func (h *SurveyHandler) DeleteSurvey(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос DeleteSurvey")
	surveyID := strings.TrimPrefix(r.URL.Path, "/api/cabinet/surveys/")
	if surveyID == "" {
		h.log.WarnContext(r.Context(), "Отсутствует survey id в DeleteSurvey")
		http.Error(w, "survey id required", http.StatusBadRequest)
		return
	}
	if err := h.deleteSurvey.Execute(r.Context(), surveyID); err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.log.InfoContext(r.Context(), "HTTP DeleteSurvey выполнен", "survey_id", surveyID)
	w.WriteHeader(http.StatusNoContent)
}

// GetPublicSurvey возвращает опрос по публичной ссылке.
func (h *SurveyHandler) GetPublicSurvey(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос GetPublicSurvey")
	token := strings.TrimPrefix(r.URL.Path, "/api/public/surveys/")
	output, err := h.getPublic.Execute(r.Context(), token)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.log.InfoContext(r.Context(), "HTTP GetPublicSurvey выполнен", "survey_id", output.ID)
	writeJSON(w, http.StatusOK, output)
}

// StartSession создает новую сессию прохождения.
func (h *SurveyHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос StartSession")
	var input dto.StartSessionInput
	if err := decodeJSON(r, &input); err != nil {
		h.log.WarnContext(r.Context(), "Ошибка декодирования StartSession", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.startSession.Execute(r.Context(), input)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.log.InfoContext(r.Context(), "HTTP StartSession выполнен", "session_id", output.SessionID)
	writeJSON(w, http.StatusCreated, output)
}

// SaveProgress сохраняет незавершенные ответы.
func (h *SurveyHandler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос SaveProgress")
	var input dto.SaveProgressInput
	if err := decodeJSON(r, &input); err != nil {
		h.log.WarnContext(r.Context(), "Ошибка декодирования SaveProgress", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.saveProgress.Execute(r.Context(), input); err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.log.InfoContext(r.Context(), "HTTP SaveProgress выполнен", "session_id", input.SessionID)
	w.WriteHeader(http.StatusNoContent)
}

// SubmitAnswers завершает прохождение опроса.
func (h *SurveyHandler) SubmitAnswers(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос SubmitAnswers")
	var input dto.SubmitAnswersInput
	if err := decodeJSON(r, &input); err != nil {
		h.log.WarnContext(r.Context(), "Ошибка декодирования SubmitAnswers", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.submit.Execute(r.Context(), input); err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.log.InfoContext(r.Context(), "HTTP SubmitAnswers выполнен", "session_id", input.SessionID)
	w.WriteHeader(http.StatusNoContent)
}

// GetResults возвращает агрегированные результаты опроса.
func (h *SurveyHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос GetResults")
	surveyID := strings.TrimPrefix(r.URL.Path, "/api/cabinet/surveys/")
	surveyID = strings.TrimSuffix(surveyID, "/results")
	output, err := h.getResults.Execute(r.Context(), surveyID)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.log.InfoContext(r.Context(), "HTTP GetResults выполнен", "survey_id", surveyID, "total", output.Total)
	writeJSON(w, http.StatusOK, output)
}

// GetExternalResults возвращает результаты для внешних интеграций.
func (h *SurveyHandler) GetExternalResults(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос GetExternalResults")
	surveyID := strings.TrimPrefix(r.URL.Path, "/api/external/surveys/")
	surveyID = strings.TrimSuffix(surveyID, "/results")
	output, err := h.getResults.Execute(r.Context(), surveyID)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.log.InfoContext(r.Context(), "HTTP GetExternalResults выполнен", "survey_id", surveyID, "total", output.Total)
	writeJSON(w, http.StatusOK, output)
}

// ExportResults выгружает результаты в xls.
func (h *SurveyHandler) ExportResults(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP запрос ExportResults")
	surveyID := strings.TrimPrefix(r.URL.Path, "/api/cabinet/surveys/")
	surveyID = strings.TrimSuffix(surveyID, "/export")
	content, err := h.export.Execute(r.Context(), surveyID)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.ms-excel")
	w.Header().Set("Content-Disposition", "attachment; filename=survey_results.xls")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
	h.log.InfoContext(r.Context(), "HTTP ExportResults выполнен", "survey_id", surveyID, "bytes", len(content))
}

func (h *SurveyHandler) writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrSurveyNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrSessionNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrInvalidQuestion), errors.Is(err, domain.ErrInvalidToken), errors.Is(err, domain.ErrInvalidSurveyTitle):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		h.log.Error("Внутренняя ошибка обработчика", "error", err)
		http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
	}
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
