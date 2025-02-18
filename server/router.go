package server

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/eventpublish"
	mid "AppPlaygroundService/middlewares/api"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	tkMid "pegasus-cloud.com/aes/toolkits/middleware"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

func router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	if mviper.GetBool("app_playground_service.developer") {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()
	router.UseRawPath = true
	router.UnescapePathValues = false
	router.Use(mid.GinLogger)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     mviper.GetStringSlice("app_playground_service.http.access_control.allow_origins"),
		AllowMethods:     mviper.GetStringSlice("app_playground_service.http.access_control.allow_methods"),
		AllowHeaders:     mviper.GetStringSlice("app_playground_service.http.access_control.allow_headers"),
		ExposeHeaders:    mviper.GetStringSlice("app_playground_service.http.access_control.expose_headers"),
		AllowCredentials: mviper.GetBool("app_playground_service.http.access_control.allow_credentials"),
	}))

	router.Use(mid.APIOperationMiddleware, tkMid.RequestIDGenerator, mid.SetExtraHeaders, mid.AccessLoggerMiddleware)
	// Base Path: /aps/api/v1
	rootRG := router.Group(path.Join(cnt.APIPrefix, cnt.APIVersion))
	enableUserAppPlaygroundServiceRouter(rootRG.Group(""))
	enableAdminAppPlaygroundServiceRouter(rootRG.Group("admin"))
	eventpublish.EnableHTTPRouters(rootRG)

	return router
}
