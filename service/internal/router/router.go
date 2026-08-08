package router

import (
	"yanny-service/internal/handler"
	"yanny-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 注册路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.SanitizeInput())
	r.Use(middleware.IPBlacklistMiddleware())
	r.Use(middleware.RequestSizeLimit(100 << 20))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")

	// ========================================
	// 管理后台 /api/v1/admin/*
	// ========================================
	admin := v1.Group("/admin")
	{
		admin.POST("/login", middleware.LoginRateLimiter(), handler.AdminLogin)

		authAdmin := admin.Group("")
		authAdmin.Use(middleware.AdminAuthMiddleware(), middleware.AdminAPILimiter())
		{
			authAdmin.GET("/info", handler.AdminInfo)
			authAdmin.POST("/logout", handler.AdminLogout)
			authAdmin.PUT("/password", handler.ChangePassword)
		}

		entity := admin.Group("/entities")
		entity.Use(middleware.AdminAuthMiddleware())
		{
			entity.GET("", middleware.RequirePermission("entity:view"), handler.ListEntities)
			entity.GET("/:id", middleware.RequirePermission("entity:view"), handler.GetEntity)
			entity.POST("", middleware.RequirePermission("entity:create"), handler.CreateEntity)
			entity.PUT("/:id", middleware.RequirePermission("entity:edit"), handler.UpdateEntity)
			entity.DELETE("/:id", middleware.RequirePermission("entity:delete"), handler.DeleteEntity)
		}

		mp := admin.Group("/mp-accounts")
		mp.Use(middleware.AdminAuthMiddleware())
		{
			mp.GET("", middleware.RequirePermission("mp_account:view"), handler.ListMpAccounts)
			mp.POST("", middleware.RequirePermission("mp_account:create"), handler.CreateMpAccount)
			mp.PUT("/:id", middleware.RequirePermission("mp_account:edit"), handler.UpdateMpAccount)
		}

		binding := admin.Group("/bindings")
		binding.Use(middleware.AdminAuthMiddleware())
		{
			binding.GET("/entity/:id", middleware.RequirePermission("entity:view"), handler.ListEntityBindings)
			binding.POST("", middleware.RequirePermission("entity:edit"), handler.BindEntityMp)
			binding.DELETE("", middleware.RequirePermission("entity:edit"), handler.UnbindEntityMp)
		}

		cdn := admin.Group("/cdn")
		cdn.Use(middleware.AdminAuthMiddleware())
		{
			cdn.GET("", middleware.RequirePermission("cdn:view"), handler.ListCdnConfigs)
			cdn.POST("", middleware.RequirePermission("cdn:create"), handler.CreateCdnConfig)
			cdn.PUT("/:id", middleware.RequirePermission("cdn:edit"), handler.UpdateCdnConfig)
			cdn.DELETE("/:id", middleware.RequirePermission("cdn:delete"), handler.DeleteCdnConfig)
		}

		category := admin.Group("/categories")
		category.Use(middleware.AdminAuthMiddleware())
		{
			category.GET("", middleware.RequirePermission("video:view"), handler.ListCategories)
			category.POST("", middleware.RequirePermission("video:create"), handler.CreateCategory)
		}

		video := admin.Group("/videos")
		video.Use(middleware.AdminAuthMiddleware())
		{
			video.GET("", middleware.RequirePermission("video:view"), handler.ListVideos)
			video.POST("", middleware.RequirePermission("video:create"), handler.CreateVideo)
			video.PUT("/:id", middleware.RequirePermission("video:edit"), handler.UpdateVideo)
			video.DELETE("/:id", middleware.RequirePermission("video:delete"), handler.DeleteVideo)
		}

		user := admin.Group("/users")
		user.Use(middleware.AdminAuthMiddleware())
		{
			user.GET("", middleware.RequirePermission("user:view"), handler.ListUsers)
			user.PUT("/:id/status", middleware.RequirePermission("user:edit"), handler.UpdateUserStatus)
		}

		admins := admin.Group("/admins")
		admins.Use(middleware.AdminAuthMiddleware())
		{
			admins.GET("", middleware.RequirePermission("admin:view"), handler.ListAdmins)
			admins.POST("", middleware.RequirePermission("admin:create"), handler.CreateAdmin)
			admins.PUT("/:id", middleware.RequirePermission("admin:edit"), handler.UpdateAdmin)
			admins.DELETE("/:id", middleware.RequirePermission("admin:delete"), handler.DeleteAdmin)
		}

		// 文件上传（七牛云直传 Token）
		upload := admin.Group("/upload")
		upload.Use(middleware.AdminAuthMiddleware())
		{
			upload.POST("/token", middleware.RequirePermission("video:create"), handler.GetUploadToken)
			upload.POST("/check", middleware.RequirePermission("video:create"), handler.CheckMediaAsset)
			upload.POST("/confirm", middleware.RequirePermission("video:create"), handler.ConfirmMediaAsset)
		}

		// 数据统计
		stats := admin.Group("/stats")
		stats.Use(middleware.AdminAuthMiddleware())
		{
			stats.GET("/dashboard", middleware.RequirePermission("analytics:view"), handler.GetDashboardStats)
			stats.GET("/trend", middleware.RequirePermission("analytics:view"), handler.GetTrendData)
			stats.GET("/top-videos", middleware.RequirePermission("analytics:view"), handler.GetTopVideos)
			stats.GET("/regions", middleware.RequirePermission("analytics:view"), handler.GetRegionStats)
			stats.GET("/invites", middleware.RequirePermission("analytics:view"), handler.GetInviteStats)
		}

		roles := admin.Group("/roles")
		roles.Use(middleware.AdminAuthMiddleware())
		{
			roles.GET("", middleware.RequirePermission("role:view"), handler.ListRoles)
			roles.GET("/:id/permissions", middleware.RequirePermission("role:view"), handler.GetRolePermissions)
			roles.POST("", middleware.RequirePermission("role:create"), handler.CreateRole)
			roles.PUT("/:id", middleware.RequirePermission("role:edit"), handler.UpdateRole)
		}
	}

	// ========================================
	// 小程序端 /api/v1/mp/*
	// ========================================
	mpApi := v1.Group("/mp")
	{
		mpApi.POST("/login", middleware.MpAPILimiter(), handler.MpLogin)

		// 主体信息（需登录）
		mpApi.GET("/entity", middleware.MpAuthRequired(), handler.MpGetEntity)

		// 分类（需登录）
		mpApi.GET("/categories", middleware.MpAuthRequired(), handler.MpListCategories)

		// 视频（需登录，去重+防刷）
		mpVideos := mpApi.Group("/videos")
		mpVideos.Use(middleware.MpAuthRequired(), middleware.MpAPILimiter(), middleware.VideoViewAntiAbuse())
		{
			mpVideos.GET("", handler.MpGetVideos)
			mpVideos.GET("/:id", handler.MpGetVideoDetail)
		}

		// 评论（需登录）
		mpComments := mpApi.Group("/comments")
		mpComments.Use(middleware.MpAuthRequired())
		{
			mpComments.GET("", handler.MpListComments)
			mpComments.GET("/:id/replies", handler.MpListReplies)
			mpComments.POST("", handler.MpCreateComment)
		}

		// 互动（全部需登录）
		mpAuth := mpApi.Group("")
		mpAuth.Use(middleware.MpAuthRequired())
		{
			mpAuth.POST("/like", handler.MpToggleLike)
			mpAuth.POST("/favorite", handler.MpToggleFavorite)
			mpAuth.GET("/favorites", handler.MpMyFavorites)
			mpAuth.GET("/history", handler.MpMyHistory)
			mpAuth.POST("/share", handler.MpRecordShare)
			mpAuth.GET("/interaction-status", handler.MpInteractionStatus)
			mpAuth.POST("/follow", handler.MpToggleFollow)
			mpAuth.GET("/user/me", handler.MpGetUserInfo)
			mpAuth.POST("/user/phone", handler.MpUpdatePhone)
			mpAuth.PUT("/user/info", handler.MpUpdateUserInfo)
			mpAuth.POST("/upload/token", handler.MpGetUploadToken)
			mpAuth.POST("/upload/avatar", handler.MpUploadAvatar)
		}

		// 埋点上报（需登录 + 签名校验 + 限流）
		mpTrack := mpApi.Group("/track")
		mpTrack.Use(middleware.MpAuthRequired(), middleware.MpAPILimiter())
		{
			mpTrack.POST("/view", handler.MpReportView)
			mpTrack.POST("/action", handler.MpReportAction)
		}
	}

	return r
}
