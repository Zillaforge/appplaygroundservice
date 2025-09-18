package api

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/mviper"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

type VerifyNamespaceInput struct {
	Namespace string
	_         struct{}
}

func VerifyNamespace(c *gin.Context) {

	var (
		funcName   = tkUtils.NameOfFunction().Name()
		requestID  = utility.MustGetContextRequestID(c)
		statusCode = http.StatusOK
		err        error
		input      = &VerifyNamespaceInput{
			Namespace: c.GetHeader(cnt.HdrNamespace),
		}
	)

	f := tracer.StartWithGinContext(c, funcName)

	defer f(tracer.Attributes{
		"input":      input,
		"err":        &err,
		"statusCode": &statusCode,
	})

	allowNS := mviper.GetStringSlice("app_playground_service.scopes.allow_namespaces")
	if len(allowNS) > 0 {
		if !slices.Contains(allowNS, input.Namespace) {
			err = tkErr.New(cnt.MidNamespaceNotAllowErr)
			statusCode = http.StatusForbidden
			zap.L().With(
				zap.String(cnt.Middleware, "!slices.Contains"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("allowNS", allowNS),
				zap.String("namespace", input.Namespace),
			).Error(err.Error())
			utility.ResponseWithType(c, statusCode, err)
			return
		}
	}
	c.Set(cnt.CtxNamespace, input.Namespace)
	c.Next()
}
