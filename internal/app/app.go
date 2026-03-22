package app

import (
	"github.com/che1nov/quantitative-research-service/config"
	exportadapter "github.com/che1nov/quantitative-research-service/internal/adapters/export"
	"github.com/che1nov/quantitative-research-service/internal/adapters/memory"
	"github.com/che1nov/quantitative-research-service/internal/adapters/security"
	httpcontroller "github.com/che1nov/quantitative-research-service/internal/controllers/http"
	"github.com/che1nov/quantitative-research-service/internal/usecases"
	httpserver "github.com/che1nov/quantitative-research-service/pkg/http"
	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

// App инкапсулирует собранное приложение.
type App struct {
	server *httpserver.Server
	log    logger.Logger
}

// New собирает зависимости и настраивает маршруты.
func New(cfg config.Config) (*App, error) {
	log := logger.New(cfg.LogLevel)
	log.Info("Инициализация приложения")

	storage := memory.NewSurveyAdapter(log.With("layer", "adapter", "component", "memory"))
	tokens := security.NewTokenGenerator(log.With("layer", "adapter", "component", "token_generator"))
	jwtSigner := security.NewJWTSigner(cfg.JWTSecret, log.With("layer", "adapter", "component", "jwt_signer"))
	xlsExporter := exportadapter.NewXLSExporter(log.With("layer", "adapter", "component", "xls"))

	publicURL := cfg.PublicBaseURL
	authenticateVKUC := usecases.NewAuthenticateVKUseCase(storage, tokens, jwtSigner, log.With("usecase", "authenticate_vk"))
	createSurveyUC := usecases.NewCreateSurveyUseCase(storage, tokens, log.With("usecase", "create_survey"), publicURL)
	updateSurveyUC := usecases.NewUpdateSurveyUseCase(storage, log.With("usecase", "update_survey"), publicURL)
	deleteSurveyUC := usecases.NewDeleteSurveyUseCase(storage, log.With("usecase", "delete_survey"))
	listSurveysUC := usecases.NewListSurveysUseCase(storage, log.With("usecase", "list_surveys"), publicURL)
	getPublicSurveyUC := usecases.NewGetPublicSurveyUseCase(storage, log.With("usecase", "get_public_survey"), publicURL)
	startSessionUC := usecases.NewStartSessionUseCase(storage, tokens, log.With("usecase", "start_session"))
	saveProgressUC := usecases.NewSaveProgressUseCase(storage, log.With("usecase", "save_progress"))
	submitAnswersUC := usecases.NewSubmitAnswersUseCase(storage, log.With("usecase", "submit_answers"))
	getResultsUC := usecases.NewGetResultsUseCase(storage, log.With("usecase", "get_results"))
	exportResultsUC := usecases.NewExportResultsUseCase(storage, xlsExporter, log.With("usecase", "export_results"))

	handler := httpcontroller.NewSurveyHandler(
		authenticateVKUC,
		createSurveyUC,
		updateSurveyUC,
		deleteSurveyUC,
		listSurveysUC,
		getPublicSurveyUC,
		startSessionUC,
		saveProgressUC,
		submitAnswersUC,
		getResultsUC,
		exportResultsUC,
		log.With("layer", "controller"),
	)
	miniAppHandler, err := httpcontroller.NewMiniAppHandler(log.With("layer", "controller", "component", "miniapp"))
	if err != nil {
		log.Error("Ошибка инициализации mini app", "error", err)
		return nil, err
	}
	authMiddleware := httpcontroller.NewAuthMiddleware(log.With("layer", "middleware", "type", "auth"))
	csrfMiddleware := httpcontroller.NewCSRFMiddleware(cfg.CSRFToken, log.With("layer", "middleware", "type", "csrf"))
	externalMiddleware := httpcontroller.NewExternalAPIMiddleware(cfg.ExternalAPIToken, log.With("layer", "middleware", "type", "external_api"))
	corsMiddleware := httpcontroller.NewCORSMiddleware(cfg.CORSAllowOrigin, log.With("layer", "middleware", "type", "cors"))
	securityHeadersMiddleware := httpcontroller.NewSecurityHeadersMiddleware(log.With("layer", "middleware", "type", "security_headers"))
	router := httpcontroller.NewRouter(handler, miniAppHandler, authMiddleware, csrfMiddleware, externalMiddleware, corsMiddleware, securityHeadersMiddleware)

	server := httpserver.New(":"+cfg.Port, router, log.With("layer", "http_server"))
	log.Info("Приложение инициализировано")
	return &App{server: server, log: log}, nil
}

// Start запускает HTTP сервер.
func (a *App) Start() error {
	a.log.Info("Запуск приложения")
	return a.server.Start()
}
