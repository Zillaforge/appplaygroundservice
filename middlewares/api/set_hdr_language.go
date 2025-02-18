package api

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	mviper "pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func SetHdrLanguage(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().Name()
		requestID = utility.MustGetContextRequestID(c)
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(nil)

	if language := c.GetHeader(cnt.HdrLanguage); language == "" {
		defaultLanguage := mviper.GetString("app_playground_service.scopes.default_language")
		zap.L().With(
			zap.String(cnt.Middleware, "SetHdrLanguage"),
			zap.String(cnt.RequestID, requestID),
			zap.String("default_language", defaultLanguage),
		).Warn("use app_playground_service.scopes.default_language")
		c.Set(cnt.CtxLanguage, defaultLanguage)
	} else {
		c.Set(cnt.CtxLanguage, language)
	}
	c.Next()
}
