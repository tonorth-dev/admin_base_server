package middleware

import (
	"bytes"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	// 配置日志输出到文件
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	file, err := os.OpenFile("./log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(file),
		zap.InfoLevel,
	)

	logger = zap.New(core)
	zap.ReplaceGlobals(logger)
}

// ResponseWriterWrapper 用于捕获响应体
type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

// Write 实现了 ResponseWriter 接口
func (r *ResponseWriterWrapper) Write(b []byte) (int, error) {
	if r.Body == nil {
		r.Body = &bytes.Buffer{}
	}
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}

// LoggerMiddleware 记录请求和响应的日志
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成唯一标识符
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)

		// 记录请求开始时间
		start := time.Now()

		// 获取请求体
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		// 记录请求信息
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		query := c.Request.URL.RawQuery
		headers := c.Request.Header

		// 记录请求日志
		logger.Info("request_in_log",
			zap.String("client_ip", clientIP),
			zap.Int("status", 0), // 初始状态为 0
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ua", userAgent),
			zap.Any("headers", headers),
			zap.String("body", string(reqBody)),
			zap.Time("request_time", start),
			zap.String("trace_id", traceID),
		)

		// 捕获响应体
		rw := &ResponseWriterWrapper{ResponseWriter: c.Writer}
		c.Writer = rw

		// 处理请求
		c.Next()

		// 记录响应信息
		status := rw.Status()
		latency := time.Since(start)
		respBody := rw.Body.String()

		// 记录响应日志
		logger.Info("request_out_log",
			zap.String("client_ip", clientIP),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("body", respBody),
			zap.Time("request_time", start),
			zap.String("trace_id", traceID),
		)
	}
}
