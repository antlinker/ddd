package log

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	isatty "github.com/mattn/go-isatty"
)

const (
	traceLogKey = "__ginlogkey__"
)

var (
	green        = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white        = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow       = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	red          = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue         = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta      = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan         = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset        = string([]byte{27, 91, 48, 109})
	disableColor = false
)

// DisableConsoleColor disables color output in the console.
func DisableConsoleColor() {
	disableColor = true
}

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
var (
	DefaultWriter = os.Stdout
)

// DebugLogger 日志
func DebugLogger() gin.HandlerFunc {
	return LoggerWithWriter(DefaultWriter)
}

// LoggerWithWriter instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer, notlogged ...string) gin.HandlerFunc {

	isTerm := true

	if w, ok := out.(*os.File); !ok ||
		(os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()))) ||
		disableColor {
		isTerm = false
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		traceid := c.GetHeader("X-Request-Id")
		l := newTraceLog(traceid)

		l.traceID = traceid
		c.Set(traceLogKey, l)
		c.Next()
		releaseTraceLog(l)

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			var statusColor, methodColor, resetColor string
			if isTerm {
				statusColor = colorForStatus(statusCode)
				methodColor = colorForMethod(method)
				resetColor = reset
			}
			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			l.Infof("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s %s",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, resetColor,
				latency,
				clientIP,
				methodColor, method, resetColor,
				path,
				comment,
			)
		}
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

// ProdLogger 生产日志
func ProdLogger() gin.HandlerFunc {
	return ProdLoggerWithWriter(DefaultWriter)
}

// ProdLoggerWithWriter  不打印请求
func ProdLoggerWithWriter(out io.Writer, notlogged ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceid := c.GetHeader("X-Request-Id")
		l := newTraceLog(traceid)
		l.traceID = traceid
		c.Set(traceLogKey, l)
		c.Next()
		releaseTraceLog(l)
	}
}

// WrapUIDForContext 将uid放入到日志
func WrapUIDForContext(c *gin.Context, uid string) {
	if uid != "" {
		logger := FromContext(c)
		if l, ok := logger.(*traceLog); ok {
			l.uid = uid
		}
	}
}

// FromContext 通过context获取ILogger
func FromContext(c *gin.Context) ILogger {
	if l, ok := c.Get(traceLogKey); ok {
		if log, ok := l.(ILogger); ok {
			return log
		}
	}
	return logger
}
