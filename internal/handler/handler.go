package handler

import (
	"encoding/json"
	"errors"
	"go-homework/internal/checker"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
)

var (
	ErrWrongRequestBodyFormat = errors.New("wrong format")
	ErrFixSpellsInBatchMode   = errors.New("can't fix spells in batch mode")
)

type Handler struct {
	logger  zap.Logger
	checker checker.YandexSpellChecker
}

func NewHandler(l *zap.Logger, c *checker.YandexSpellChecker) *Handler {
	return &Handler{
		logger:  *l,
		checker: *c,
	}
}

func NewRouter(h *Handler) *gin.Engine {
	router := gin.New()

	router.Use(h.zapLoggerMiddleware())
	router.Use(gin.Recovery())

	router.POST("/", h.checkSpellsBatchHandler)

	return router
}

func (h Handler) zapLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				h.logger.Error(e)
			}
		} else {
			fields := []zapcore.Field{
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
			}

			h.logger.Info(path, fields...)
		}
	}
}

func (h Handler) checkSpellsBatchHandler(c *gin.Context) {
	var data map[string][]string

	dec := json.NewDecoder(c.Request.Body)
	err := dec.Decode(&data)
	if err != nil {
		c.Error(ErrWrongRequestBodyFormat)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrWrongRequestBodyFormat.Error()})
		return
	}

	_, ok := data["texts"]

	if ok {
		corrected, err := h.checker.FixSpellsInBatchMode(data["texts"])
		if err != nil {
			c.Error(ErrFixSpellsInBatchMode)
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrFixSpellsInBatchMode.Error()})
			return
		}

		c.JSON(http.StatusOK, corrected)
	} else {
		c.Error(ErrWrongRequestBodyFormat)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrWrongRequestBodyFormat.Error()})
	}
}
