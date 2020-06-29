package route

import (
	"os"

	_ "github.com/duyanghao/coredns-dynapi-adapter/docs"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/controller"
	"github.com/duyanghao/coredns-dynapi-adapter/pkg/log"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Swagger coredns-dynapi-adapter
// @version 0.1.0
// @description This is a coredns-dynapi-adapter.
// @contact.name duyanghao
// @contact.url https://duyanghao.github.io
// @contact.email 1294057873@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
func InstallRoutes(r *gin.Engine) {
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// a ping api test
	r.GET("/ping", controller.Ping)

	// get coredns-dynapi-adapter version
	r.GET("/version", controller.Version)

	// config reload
	r.Any("/-/reload", func(c *gin.Context) {
		log.Info("===== Server Stop! Cause: Config Reload. =====")
		os.Exit(1)
	})

	rootGroup := r.Group("/api/v1")

	{
		// a ping api to test basic auth
		rootGroup.GET("/ping", controller.Ping)
	}

	{
		corednsController := controller.NewCorednsController()
		rootGroup.POST("/node", corednsController.RegisterNode)
		rootGroup.GET("/node", corednsController.GetNode)
		rootGroup.POST("/domain", corednsController.AddDomain)
	}
}
