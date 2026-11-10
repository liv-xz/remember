//go:build wireinject
// +build wireinject

package wire

import (
	"time"

	"server/config"
	"server/configmodel"
	"server/database/connection"
	"server/middleware"
	"server/routes"
	"server/service/upload"
	"server/utils/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	// repositories (all use package name "repository")
	deceasedRepo "server/repository/deceased"
	graveRepo "server/repository/grave"
	memberRepo "server/repository/grave_member"
	gslRepo "server/repository/grave_share_link"
	kdrRepo "server/repository/kinship_deceased_relation"
	krRepo "server/repository/kinship_relation"
	loginLogRepo "server/repository/login_log"
	eventRepo "server/repository/memorial_event"
	recordRepo "server/repository/memorial_record"
	mrRepo "server/repository/memorial_relation"
	oplogRepo "server/repository/operation_log"
	tokenRepo "server/repository/token"
	userRepo "server/repository/user"

	// handlers
	adminHandler "server/handler/admin"
	authHandler "server/handler/auth"
	deceasedHandler "server/handler/deceased"
	graveHandler "server/handler/grave"
	memberHandler "server/handler/grave_member"
	gslHandler "server/handler/grave_share_link"
	kdrHandler "server/handler/kinship_deceased_relation"
	krHandler "server/handler/kinship_relation"
	eventHandler "server/handler/memorial_event"
	recordHandler "server/handler/memorial_record"
	mrHandler "server/handler/memorial_relation"
	graphHandler "server/handler/relation_graph"
	uploadHandler "server/handler/upload"
	userHandler "server/handler/user"

	// services (all use package name "service")
	deceasedSvc "server/service/deceased"
	graveSvc "server/service/grave"
	memberSvc "server/service/grave_member"
	gslSvc "server/service/grave_share_link"
	kdrSvc "server/service/kinship_deceased_relation"
	krSvc "server/service/kinship_relation"
	eventSvc "server/service/memorial_event"
	recordSvc "server/service/memorial_record"
	mrSvc "server/service/memorial_relation"
	rgSvc "server/service/relation_graph"
	userSvc "server/service/user"

	authService "server/service/auth"
)

// ── Config Providers ──

func ProvideDB() *gorm.DB {
	return connection.DB
}

func ProvideJWTConfig() configmodel.JwtConfig {
	return config.Global.JWT
}

func ProvideLoginProtector() *middleware.LoginProtector {
	return middleware.NewLoginProtector(5, 15*time.Minute)
}

func ProvideUploader() upload.Uploader {
	cfg := config.Global.Upload
	return upload.NewLocalUploader(cfg.Dir, cfg.MaxSize, cfg.AllowedTypes)
}

func ProvideAdminHandler(db *gorm.DB) *adminHandler.AdminHandler {
	return adminHandler.NewAdminHandler(db, time.Now())
}

// ── ProviderSets ──

var RepoSet = wire.NewSet(
	userRepo.NewUserRepository,
	tokenRepo.NewRefreshTokenRepository,
	loginLogRepo.NewLoginLogRepository,
	oplogRepo.NewOperationLogRepository,
	graveRepo.NewGraveRepository,
	memberRepo.NewGraveMemberRepository,
	deceasedRepo.NewDeceasedRepository,
	kdrRepo.NewDeceasedRelationRepository,
	krRepo.NewKinshipRelationRepository,
	mrRepo.NewMemorialRelationRepository,
	recordRepo.NewMemorialRecordRepository,
	eventRepo.NewMemorialEventRepository,
	gslRepo.NewGraveShareLinkRepository,
)

var ServiceSet = wire.NewSet(
	jwt.NewJWTService,
	authService.NewAuthService,
	graveSvc.NewGraveService,
	deceasedSvc.NewDeceasedService,
	recordSvc.NewMemorialRecordService,
	eventSvc.NewMemorialEventService,
	kdrSvc.NewDeceasedRelationService,
	memberSvc.NewGraveMemberService,
	userSvc.NewUserService,
	gslSvc.NewGraveShareLinkService,
	krSvc.NewKinshipRelationService,
	mrSvc.NewMemorialRelationService,
	rgSvc.NewRelationGraphService,
)

var HandlerSet = wire.NewSet(
	authHandler.NewAuthHandler,
	graveHandler.NewGraveHandler,
	deceasedHandler.NewDeceasedHandler,
	recordHandler.NewMemorialRecordHandler,
	eventHandler.NewMemorialEventHandler,
	kdrHandler.NewDeceasedRelationHandler,
	memberHandler.NewGraveMemberHandler,
	userHandler.NewUserHandler,
	gslHandler.NewGraveShareLinkHandler,
	krHandler.NewKinshipRelationHandler,
	mrHandler.NewMemorialRelationHandler,
	graphHandler.NewRelationGraphHandler,
	uploadHandler.NewUploadHandler,
)

// ── Injector ──

func InitApp() (*gin.Engine, error) {
	wire.Build(
		ProvideDB,
		ProvideJWTConfig,
		ProvideLoginProtector,
		ProvideUploader,
		ProvideAdminHandler,
		RepoSet,
		ServiceSet,
		HandlerSet,
		routes.SetupRouter,
	)
	return nil, nil
}
