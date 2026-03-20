package export

import (
	"bytes"
	"context"
	"fmt"
	"html"

	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

// XLSExporter формирует табличный файл в формате, совместимом с .xls.
type XLSExporter struct {
	log logger.Logger
}

func NewXLSExporter(log logger.Logger) *XLSExporter {
	return &XLSExporter{log: log}
}

// BuildSurveyResultsXLS формирует содержимое файла для выгрузки результатов.
func (e *XLSExporter) BuildSurveyResultsXLS(ctx context.Context, survey domain.Survey, sessions []domain.AnswerSession) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("survey_id\tsurvey_title\tsession_id\tquestion_id\tanswer\n")
	for _, session := range sessions {
		for qid, values := range session.Answers {
			for _, value := range values {
				line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", survey.ID, sanitizeTSV(survey.Title), session.ID, qid, sanitizeTSV(value))
				b.WriteString(line)
			}
		}
	}
	e.log.InfoContext(ctx, "Экспорт xls сформирован", "survey_id", survey.ID)
	return b.Bytes(), nil
}

func sanitizeTSV(value string) string {
	return html.EscapeString(value)
}
