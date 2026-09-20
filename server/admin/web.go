package admin

import (
	"fmt"
	"server/config"

	"github.com/labstack/echo/v4"
)

type handler struct {
	cfg *config.Config
}

func ListenWeb(cfg *config.Config) error {
	e := echo.New()

	handler := &handler{
		cfg: cfg,
	}

	e.GET("/clients/fetch", handler.fetch)
	e.GET("/clients/message/new", handler.newMessage)

	err := e.Start(cfg.WebListen)
	return fmt.Errorf("start error: %w", err)
}
