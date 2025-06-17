package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	evbus "github.com/asaskevich/EventBus"
	"github.com/bwmarrin/snowflake"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/event"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/refresh"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/server/middleware"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/hexley21/soccer-manager/pkg/config"
	"github.com/hexley21/soccer-manager/pkg/hasher"
	"github.com/hexley21/soccer-manager/pkg/json/jsoniter_json"
	"github.com/hexley21/soccer-manager/pkg/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	router        *echo.Echo
	metricsRouter *echo.Echo
	mux           *http.Server
	metricsMux    *http.Server

	*delivery.Components
}

func NewServer(
	cfg *config.Config,
	logger echo.Logger,
	validator validator.Validator,
	snowflakeNode *snowflake.Node,
	hasher hasher.Hasher,
	dbPool *pgxpool.Pool,
) *Server {
	jsonProcessor := jsoniter_json.New()

	router := echo.New()
	router.Debug = !cfg.IsProd
	mux := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      router,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	metricsRouter := echo.New()
	metricsRouter.Debug = !cfg.IsProd
	metricsMux := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Metrics.Port),
		Handler:      metricsRouter,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	router.Logger = logger
	router.JSONSerializer = jsoniter_json.NewEcho(jsonProcessor)
	router.Validator = validator

	metricsRouter.Logger = logger
	metricsRouter.JSONSerializer = jsoniter_json.NewEcho(jsonProcessor)

	userRepo := repository.NewUserRepo(dbPool, snowflakeNode)
	globeRepo := repository.NewGlobeRepo(dbPool, cfg.Globe.TTL)

	teamRepo := repository.NewTeamRepository(dbPool, snowflakeNode)
	teamTranslationRepo := repository.NewTeamTranslationsRepository(dbPool)

	playerPosRepo := repository.NewPlayerPositionRepo(dbPool)
	playerRepo := repository.NewPlayerRepository(dbPool, snowflakeNode)

	services := delivery.Services{
		GlobeService: service.NewGlobeService(globeRepo),

		AuthService: service.NewAuthService(userRepo, hasher),
		UserService: service.NewUserService(userRepo, hasher),

		TeamService: service.NewTeamService(teamRepo, teamTranslationRepo),

		PlayerPosService: service.NewPlayerPositionService(playerPosRepo),
		PlayerService:    service.NewPlayerService(playerRepo),
	}

	jwtManagers := delivery.JWTManagers{
		Access:  access.NewManager(cfg.JWT.Access),
		Refresh: refresh.NewManager(cfg.JWT.Refresh),
	}

	return &Server{
		Components: &delivery.Components{
			Cfg:           cfg,
			Logger:        logger,
			Validator:     validator,
			SnowflakeNode: snowflakeNode,
			DbPool:        dbPool,

			JWTManagers: &jwtManagers,
			Services:    &services,
			EventBus:    evbus.New(),
		},

		mux:           &mux,
		router:        router,
		metricsMux:    &metricsMux,
		metricsRouter: metricsRouter,
	}
}

func (s *Server) Run() error {
	s.router.Use(echo_middleware.LoggerWithConfig(echo_middleware.LoggerConfig{
		Skipper: echo_middleware.DefaultSkipper,
		Output:  os.Stdout,
	}))
	s.router.Use(echo_middleware.CORSWithConfig(echo_middleware.CORSConfig{
		AllowOrigins: []string{s.Cfg.HTTP.CorsOrigins}, // default config - { * }
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"Accept-Language",
		},
	}))
	s.router.Use(echo_middleware.Recover())

	middlewares := delivery.Middlewares{
		JWTMiddleware:  middleware.JWTAuth(s.JWTManagers.Access),
		IsAdmin:        middleware.IsAdmin(),
		AcceptLanguage: middleware.AcceptLanguage(),
	}

	apiGroup := s.router.Group("/api")

	v1Group := apiGroup.Group("/v1")

	// register event handlers
	event.RegisterEventHandlers(s.EventBus, s.Components)

	// register api handlers
	v1.RegisterRoutes(v1Group, s.Components, &middlewares)

	// register metric handling
	s.metricsRouter.GET("/metrics", echoprometheus.NewHandler())

	var wg sync.WaitGroup
	var httpErrs error
	var mu sync.Mutex

	wg.Add(2)

	go startServer(&wg, &mu, s.mux, "main", &httpErrs)
	go startServer(&wg, &mu, s.metricsMux, "metrics", &httpErrs)

	wg.Wait()
	return httpErrs
}

func (s *Server) Close() error {
	var closeErrs error
	var mu sync.Mutex

	ctx, cancel := context.WithTimeout(context.Background(), s.Cfg.Server.ShutdownTimeout)
	var wg sync.WaitGroup

	wg.Add(3)

	go shutdownServer(ctx, &wg, &mu, s.mux, "main", &closeErrs)
	go shutdownServer(ctx, &wg, &mu, s.metricsMux, "metrics", &closeErrs)

	go func() {
		s.DbPool.Close()
		wg.Done()
	}()

	go func() {
		wg.Wait()
		cancel()
	}()

	<-ctx.Done()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		closeErrs = errors.Join(
			closeErrs,
			fmt.Errorf("shutdown timed out after %v", s.Cfg.Server.ShutdownTimeout),
		)
	}

	return closeErrs
}

func startServer(
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	server *http.Server,
	serverName string,
	errs *error,
) {
	defer wg.Done()
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		mu.Lock()
		*errs = errors.Join(*errs, fmt.Errorf("%s server: %w", serverName, err))
		mu.Unlock()
	}
}

func shutdownServer(
	ctx context.Context,
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	server *http.Server,
	serverName string,
	errs *error,
) {
	defer wg.Done()
	if err := server.Shutdown(ctx); err != nil {
		mu.Lock()
		*errs = errors.Join(*errs, fmt.Errorf("%s shutdown: %w", serverName, err))
		mu.Unlock()
	}
}
