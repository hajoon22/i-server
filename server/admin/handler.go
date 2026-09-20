package admin

import (
	"fmt"
	"server/client"

	"github.com/labstack/echo/v4"
)

func (h *handler) fetch(ctx echo.Context) error {
	key := ctx.QueryParam("key")
	if key == "" || key != h.cfg.AdminKey {
		return ctx.JSON(401, map[string]string{
			"status":  "error",
			"message": "who are you??",
		})
	}

	results := client.Fetch()
	return ctx.JSON(200, map[string]any{
		"status":  "succeed",
		"message": results,
	})
}

func (h *handler) newMessage(ctx echo.Context) error {
	key := ctx.QueryParam("key")
	if key == "" || key != h.cfg.AdminKey {
		return ctx.JSON(401, map[string]string{
			"status":  "error",
			"message": "who are you??",
		})
	}

	message := ctx.QueryParam("message")
	if message == "" {
		return ctx.JSON(400, map[string]string{
			"status":  "error",
			"message": "i am not message! where is the message??",
		})
	}

	cnt := client.NewMessage(h.cfg, message)
	return ctx.JSON(200, map[string]string{
		"status":  "succeed",
		"message": fmt.Sprintf("i sent %d clients", cnt),
	})
}
