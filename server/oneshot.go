package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/run4w4y/consul-cleanup/common"
)

func OneshotCleanupAll(c echo.Context) error {
	cc := c.(*ApplicationContext)

	err := common.OneshotCleanup(cc.Request().Context(), common.OneshotCleanupConfig{
		CleanupConfig: cc.CleanupConfig,
	})

	if err != nil {
		return cc.JSON(http.StatusInternalServerError, map[string]any{
			"ok":    false,
			"error": fmt.Sprintf("%v", err),
		})
	}

	return cc.JSON(http.StatusOK, map[string]bool{"ok": true})
}

func OneshotCleanupService(c echo.Context) error {
	cc := c.(*ApplicationContext)

	err := common.OneshotCleanup(cc.Request().Context(), common.OneshotCleanupConfig{
		CleanupConfig: cc.CleanupConfig,
		ServiceName:   cc.Param("service"),
	})

	if err != nil {
		return cc.JSON(http.StatusInternalServerError, map[string]any{
			"ok":    false,
			"error": fmt.Sprintf("%v", err),
		})
	}

	return cc.JSON(http.StatusOK, map[string]bool{"ok": true})
}
