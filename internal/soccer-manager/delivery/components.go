package delivery

import (
	evbus "github.com/asaskevich/EventBus"
	"github.com/bwmarrin/snowflake"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/refresh"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/hexley21/soccer-manager/pkg/config"
	"github.com/hexley21/soccer-manager/pkg/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type JWTManagers struct {
	Access  jwt.ManagerWithTTL[access.Data]
	Refresh jwt.ManagerWithTTL[refresh.Data]
}

type Services struct {
	GlobeService service.GlobeService

	AuthService service.AuthService
	UserService service.UserService

	TeamService service.TeamService

	PlayerPosService service.PlayerPositionService
	PlayerService    service.PlayerService

	TransferService service.TransferService
	TransferRecordService service.TransferRecordService
}

type Components struct {
	Cfg           *config.Config
	Logger        echo.Logger
	Validator     validator.Validator
	SnowflakeNode *snowflake.Node
	EventBus      evbus.Bus

	DbPool *pgxpool.Pool
	// redisCluster *redis.ClusterClient

	JWTManagers *JWTManagers
	Services    *Services
}

type Middlewares struct {
	JWTMiddleware  echo.MiddlewareFunc
	IsAdmin        echo.MiddlewareFunc
	AcceptLanguage echo.MiddlewareFunc
}
