package routes

import (
	"time"

	"server/config"
	adminhandler "server/handler/admin"
	authhandler "server/handler/auth"
	deceasedhandler "server/handler/deceased"
	gravehandler "server/handler/grave"
	memberhandler "server/handler/grave_member"
	gslhandler "server/handler/grave_share_link"
	kdrhandler "server/handler/kinship_deceased_relation"
	krhandler "server/handler/kinship_relation"
	eventhandler "server/handler/memorial_event"
	recordhandler "server/handler/memorial_record"
	mrhandler "server/handler/memorial_relation"
	graphhandler "server/handler/relation_graph"
	uploadhandler "server/handler/upload"
	userhandler "server/handler/user"
	"server/middleware"
	oplogrepo "server/repository/operation_log"
	myjwt "server/utils/jwt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHdl *authhandler.AuthHandler,
	graveHdl *gravehandler.GraveHandler,
	deceasedHdl *deceasedhandler.DeceasedHandler,
	recordHdl *recordhandler.MemorialRecordHandler,
	eventHdl *eventhandler.MemorialEventHandler,
	kdrHdl *kdrhandler.DeceasedRelationHandler,
	memberHdl *memberhandler.GraveMemberHandler,
	userHdl *userhandler.UserHandler,
	opLogRepo *oplogrepo.OperationLogRepository,
	jwtSvc *myjwt.JWTService,
	gslHdl *gslhandler.GraveShareLinkHandler,
	krHdl *krhandler.KinshipRelationHandler,
	mrHdl *mrhandler.MemorialRelationHandler,
	graphHdl *graphhandler.RelationGraphHandler,
	uploadHdl *uploadhandler.UploadHandler,
	adminHdl *adminhandler.AdminHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.Global.App.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.RequestLog())

	middleware.RegisterPprof(r)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.Use(middleware.AuthRateLimit())
		auth.POST("/register", authHdl.Register)
		auth.POST("/login", authHdl.Login)
		auth.POST("/wx-login", authHdl.WxLogin)
		auth.POST("/refresh", authHdl.Refresh)
		auth.POST("/logout", authHdl.Logout)

		public := v1.Group("")
		public.Use(middleware.GeneralRateLimit())
		public.GET("/share-links/:code", gslHdl.Validate)

		shareAuth := v1.Group("")
		shareAuth.Use(middleware.GeneralRateLimit())
		shareAuth.Use(middleware.JWTAuth(jwtSvc))
		shareAuth.POST("/share-links/:code/join", gslHdl.Join)

		protected := v1.Group("")
		protected.Use(middleware.GeneralRateLimit())
		protected.Use(middleware.JWTAuth(jwtSvc))
		protected.Use(middleware.OperationLog(opLogRepo))
		protected.PUT("/auth/password", authHdl.ChangePassword)

		user := protected.Group("/user")
		user.GET("/profile", userHdl.GetProfile)
		user.PUT("/profile", userHdl.UpdateProfile)

		graves := protected.Group("/graves")
		graves.POST("/with-deceased", graveHdl.CreateWithDeceased)
		graves.POST("", graveHdl.Create)
		graves.GET("", graveHdl.ListMy)
		graves.GET("/:gid", graveHdl.GetByID)
		graves.PUT("/:gid", graveHdl.Update)
		graves.DELETE("/:gid", graveHdl.Delete)

		deceased := protected.Group("/graves/:gid/deceased")
		deceased.POST("", deceasedHdl.Create)
		deceased.GET("", deceasedHdl.ListByGrave)

		deceasedByID := protected.Group("/deceased")
		deceasedByID.PUT("/:id", deceasedHdl.Update)
		deceasedByID.DELETE("/:id", deceasedHdl.Delete)

		records := protected.Group("/graves/:gid/records")
		records.POST("", recordHdl.Create)
		records.GET("", recordHdl.ListByGrave)

		deceasedRecords := protected.Group("/deceased/:did/records")
		deceasedRecords.GET("", recordHdl.ListByDeceased)

		recordByID := protected.Group("/records")
		recordByID.DELETE("/:id", recordHdl.Delete)

		events := protected.Group("/graves/:gid/events")
		events.POST("", eventHdl.Create)
		events.GET("", eventHdl.ListByGrave)

		deceasedEvents := protected.Group("/deceased/:did/events")
		deceasedEvents.GET("", eventHdl.ListByDeceased)

		eventByID := protected.Group("/events")
		eventByID.PUT("/:id", eventHdl.Update)
		eventByID.DELETE("/:id", eventHdl.Delete)

		kdr := protected.Group("/graves/:gid/kinship-deceased-relations")
		kdr.POST("", kdrHdl.Create)
		kdr.GET("", kdrHdl.ListByGrave)

		deceasedKDR := protected.Group("/deceased/:did/kinship-deceased-relations")
		deceasedKDR.GET("", kdrHdl.ListByDeceased)

		kdrByID := protected.Group("/kinship-deceased-relations")
		kdrByID.DELETE("/:id", kdrHdl.Delete)

		members := protected.Group("/graves/:gid/members")
		members.POST("", memberHdl.Add)
		members.GET("", memberHdl.ListByGrave)

		memberOperations := protected.Group("/members")
		memberOperations.PUT("/:id/role", memberHdl.UpdateRole)
		memberOperations.PUT("/:id/relation-desc", memberHdl.UpdateRelationDesc)
		memberOperations.DELETE("/:id", memberHdl.Remove)
		memberOperations.POST("/transfer", memberHdl.TransferOwnership)

		shareLinks := protected.Group("/graves/:gid/share-links")
		shareLinks.POST("", gslHdl.Create)
		shareLinks.GET("", gslHdl.ListByGrave)
		shareLinks.PUT("/:id/deactivate", gslHdl.Deactivate)

		kinshipRelations := protected.Group("/graves/:gid/kinship-relations")
		kinshipRelations.POST("", krHdl.CreatePublic)
		kinshipRelations.GET("", krHdl.ListByGrave)

		kinshipOps := protected.Group("/kinship-relations")
		kinshipOps.PUT("/:id/confirm", krHdl.Confirm)
		kinshipOps.PUT("/:id/reject", krHdl.Reject)

		memorialRelations := protected.Group("/graves/:gid/memorial-relations")
		memorialRelations.POST("", mrHdl.CreatePublic)
		memorialRelations.GET("", mrHdl.ListByGrave)

		memorialOps := protected.Group("/memorial-relations")
		memorialOps.PUT("/:id/confirm", mrHdl.Confirm)
		memorialOps.PUT("/:id/reject", mrHdl.Reject)

		protected.GET("/graves/:gid/relation-graph", graphHdl.GetGraph)
		protected.POST("/upload", uploadHdl.Upload)

		admin := protected.Group("/admin")
		admin.GET("/stats", adminHdl.StatsHandler)
		admin.GET("/stats/overview", adminHdl.OverviewHandler)
		admin.GET("/stats/trend", adminHdl.TrendHandler)
		admin.GET("/stats/hot", adminHdl.HotHandler)
		admin.GET("/ws/stats", adminHdl.StatsWS)
		admin.GET("/users", adminHdl.ListUsersHandler)
		admin.GET("/graves", adminHdl.ListGravesHandler)
		admin.PUT("/users/:id/status", adminHdl.UpdateUserStatusHandler)
		admin.DELETE("/users/:id", adminHdl.DeleteUserHandler)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Static("/uploads", config.Global.Upload.Dir)

	return r
}
