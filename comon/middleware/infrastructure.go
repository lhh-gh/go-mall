package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github/lhh-gh/go-mall/comon/logger"
	"github/lhh-gh/go-mall/comon/util"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// infrastructure 中存放项目运行需要的基础中间件

// StartTrace 启动追踪中间件
// 用于处理请求的追踪信息，包括 traceId、spanId 和 pSpanId
func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("traceid")
		pSpanId := c.Request.Header.Get("spanid")
		spanId := util.GenerateSpanID(c.Request.RemoteAddr)
		if traceId == "" { // 如果traceId为空，说明是链路的起始端，将其设置为当前的spanId，起始端的spanId作为root spanId
			traceId = spanId // trace用于标识整个请求链路，span则用于标识链路中的不同服务
		}
		c.Set("traceid", traceId)
		c.Set("spanid", spanId)
		c.Set("pspanid", pSpanId)
		c.Next()
	}
}

// bodyLogWriter 包装 gin.ResponseWriter
// 用于拦截响应写入，以便记录响应内容
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写 Write 方法
// 实现响应内容的拦截，同时写入到原始响应和缓冲区
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LogAccess 访问日志中间件
// 记录请求的详细信息，包括请求体、响应体、处理时间等
func LogAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 保存请求体
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
		start := time.Now()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		accessLog(c, "access_start", time.Since(start), reqBody, nil)
		defer func() {
			accessLog(c, "access_end", time.Since(start), reqBody, blw.body.String())
		}()
		c.Next()

		return
	}
}

// accessLog 记录访问日志
// 参数说明：
// - c: gin上下文
// - accessType: 访问类型（开始/结束）
// - dur: 处理时长
// - body: 请求体
// - dataOut: 响应数据
func accessLog(c *gin.Context, accessType string, dur time.Duration, body []byte, dataOut interface{}) {
	req := c.Request
	bodyStr := string(body)
	query := req.URL.RawQuery
	path := req.URL.Path
	// TODO: 实现Token认证后再把访问日志里也加上token记录
	// token := c.Request.Header.Get("token")
	logger.New(c).Info("AccessLog",
		"type", accessType,
		"ip", c.ClientIP(),
		//"token", token,
		"method", req.Method,
		"path", path,
		"query", query,
		"body", bodyStr,
		"output", dataOut,
		"time(ms)", int64(dur/time.Millisecond))
}

// GinPanicRecovery 自定义的 gin panic 恢复中间件
// 用于捕获并处理 panic，记录错误信息和调用栈
func GinPanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否是连接断开导致的错误
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.New(c).Error("http request broken pipe", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest))
					// 如果连接已断开，无法写入状态码
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				logger.New(c).Error("http_request_panic", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest), "stack", string(debug.Stack()))

				c.AbortWithError(http.StatusInternalServerError, err.(error))
			}
		}()
		c.Next()
	}
}

// 对于文件上传下载，还有httptool中的文件上传下载，日志都不能打，特别是在一些流式场景
//
// 关于文件上传下载中间件：
func LogAccessFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// multipart form data 不打印body
		contentType := c.GetHeader("Content-Type")
		var reqBody []byte
		if !strings.Contains(contentType, "multipart/form-data") {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(reqBody))
		}
		start := time.Now()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		accessLog(c, "access_start", time.Since(start), reqBody, nil)
		defer func() {
			accessLog(c, "access_end", time.Since(start), reqBody, blw.body.String())
		}()
		c.Next()

		return
	}
}

// LogAccessFile2 文件访问日志中间件
// 用于处理文件上传下载等特殊场景的访问日志记录
// 特点：
// 1. 对于 multipart/form-data 类型的请求不记录请求体
// 2. 对于大于 10KB 的响应不记录响应体
// 3. 记录请求开始和结束的完整信息
func LogAccessFile2() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 初始化请求体变量
		var reqBody []byte

		// 检查请求类型，multipart/form-data 类型不记录请求体
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "multipart/form-data") {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		// 记录请求开始时间
		start := time.Now()

		// 创建响应体写入器，用于记录响应内容
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// 记录请求开始日志
		accessLog(c, "access_start", time.Since(start), reqBody, nil)

		// 使用 defer 确保在函数结束时记录请求结束日志
		defer func() {
			var responseLogging string
			// 响应体大于 10KB 时不记录具体内容
			if c.Writer.Size() > 10*1024 {
				responseLogging = "响应数据过大，已超过10KB，不记录具体内容"
			} else {
				responseLogging = blw.body.String()
			}
			accessLog(c, "access_end", time.Since(start), reqBody, responseLogging)
		}()

		c.Next()
	}
}
