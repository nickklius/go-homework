package main

import (
	"go-homework/internal/checker"
	"go-homework/internal/handler"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	c := checker.New()
	h := handler.NewHandler(logger, c)

	router := handler.NewRouter(h)
	router.Run(":8080")
}
