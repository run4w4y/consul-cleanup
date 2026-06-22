package server

import (
	"github.com/labstack/echo/v4"
	"github.com/run4w4y/consul-cleanup/common"
)

type ServerCleanupConfig struct {
	common.CleanupConfig
	Port        uint
	AccessToken string
}

type ApplicationContext struct {
	echo.Context
	ServerCleanupConfig
}
