package api

import (
	cnt "AppPlaygroundService/constants"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func ResourceReview(c *gin.Context) {
	var (
		funcName = tkUtils.NameOfFunction().Name()
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(nil)

	var underReview = false
	hdrReview := strings.ToLower(c.GetHeader(cnt.HdrProjectExtraResourceReviewAPS))
	if hdrReview == "true" {
		underReview = true
	}
	c.Set(cnt.CtxResourceReview, underReview)
	c.Next()
}
