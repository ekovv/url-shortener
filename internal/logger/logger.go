package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Sugar zap.SugaredLogger

func HTTPLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		Sugar.Info(
			"url", param.Path,
			"method", param.Method,
			"duration", param.Latency,
			"statusCode", param.StatusCode,
			"bodySize", param.BodySize,
		)
		return ""
	})

}
