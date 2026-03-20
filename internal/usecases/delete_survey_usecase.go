package usecases

import "context"

// DeleteSurveyUseCase удаляет опрос из личного кабинета.
type DeleteSurveyUseCase struct {
	storage SurveyStorage
	log     AppLogger
}

func NewDeleteSurveyUseCase(storage SurveyStorage, log AppLogger) *DeleteSurveyUseCase {
	return &DeleteSurveyUseCase{storage: storage, log: log}
}

// Execute удаляет опрос по идентификатору.
func (uc *DeleteSurveyUseCase) Execute(ctx context.Context, surveyID string) error {
	uc.log.InfoContext(ctx, "Удаление опроса", "survey_id", surveyID)
	if err := uc.storage.DeleteSurvey(ctx, surveyID); err != nil {
		uc.log.ErrorContext(ctx, "Ошибка удаления опроса", "error", err)
		return err
	}
	uc.log.InfoContext(ctx, "Опрос успешно удален", "survey_id", surveyID)
	return nil
}
