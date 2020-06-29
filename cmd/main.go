package main

import (
	"fmt"

	"github.com/duyanghao/coredns-dynapi-adapter/pkg/config"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/log"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/route"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/util"
	"github.com/gin-gonic/gin"
)

func main() {
	util.SetupSigusr1Trap()
	r := gin.Default()
	m := config.GetString(config.FLAG_KEY_GIN_MODE)
	gin.SetMode(m)

	route.InstallRoutes(r)
	serverBindAddr := fmt.Sprintf("%s:%d", config.GetString(config.FLAG_KEY_SERVER_HOST), config.GetInt(config.FLAG_KEY_SERVER_PORT))
	log.Infof("Run server at %s", serverBindAddr)
	r.Run(serverBindAddr) // listen and serve
}
