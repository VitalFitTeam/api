package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource already exists")
)

type LogErrors struct {
	logger *zap.SugaredLogger
}

func NewLogErrors(logger *zap.SugaredLogger) LogErrors {
	return LogErrors{
		logger: logger,
	}
}

func (l *LogErrors) InternalServerError(c *gin.Context, err error) {

	l.logger.Errorw("internal server error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
	c.JSON(http.StatusInternalServerError, gin.H{"error": "the server encountered a problem"})
}

func (l *LogErrors) BadRequestResponse(c *gin.Context, err error) {
	l.logger.Warnw("bad request error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func (l *LogErrors) NotFoundResponse(c *gin.Context) {
	l.logger.Warnw("not found error", "method", c.Request.Method, "path", c.Request.URL.Path)
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func (l *LogErrors) ConflictResponse(c *gin.Context, err error) {
	l.logger.Errorw("conflict error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
	c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
}

func (l *LogErrors) UnauthorizedBasicErrorResponse(c *gin.Context) {
	l.logger.Errorw("unauthorized basic error", "method", c.Request.Method, "path", c.Request.URL.Path)
	c.Header("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}

func (l *LogErrors) ForbiddenResponse(c *gin.Context) {
	l.logger.Warnw("forbidden", "method", c.Request.Method, "path", c.Request.URL.Path)
	c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
}
