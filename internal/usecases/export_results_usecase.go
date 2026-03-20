package usecases

import "context"

// ExportResultsUseCase выгружает результаты опроса в .xls файл.
type ExportResultsUseCase struct {
	storage  SurveyStorage
	exporter XLSExporter
	log      AppLogger
}

func NewExportResultsUseCase(storage SurveyStorage, exporter XLSExporter, log AppLogger) *ExportResultsUseCase {
	return &ExportResultsUseCase{storage: storage, exporter: exporter, log: log}
}

// Execute формирует бинарное содержимое xls-файла.
func (uc *ExportResultsUseCase) Execute(ctx context.Context, surveyID string) ([]byte, error) {
	uc.log.InfoContext(ctx, "Подготовка экспорта результатов", "survey_id", surveyID)
	survey, err := uc.storage.GetSurveyByID(ctx, surveyID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка чтения опроса для экспорта", "error", err)
		return nil, err
	}
	sessions, err := uc.storage.ListCompletedSessionsBySurvey(ctx, surveyID)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка чтения сессий для экспорта", "error", err)
		return nil, err
	}

	file, err := uc.exporter.BuildSurveyResultsXLS(ctx, survey, sessions)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка экспорта xls", "error", err)
		return nil, err
	}

	uc.log.InfoContext(ctx, "Файл xls подготовлен", "survey_id", surveyID, "rows", len(sessions))
	return file, nil
}
