// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package router

import (
	"ServeBin"
	"ServeBin/controller"
	"ServeBin/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
)

func NewRouter(apiController *controller.APIController) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.LoadHTMLGlob("templates/**/*")

	if os.Getenv("IS_BACKUP_SERVER") == "true" {
		// Redirect to the Main Server
		router.GET("", apiController.Redirect)
		router.GET("/docs/*any", apiController.Redirect)
	} else {
		// Root Path
		router.GET("", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "ServeBin", "Version": ServeBin.Version[1:]})
		})

		// Add Swagger
		router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.GET("/docs", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})

	router.StaticFile("/favicon.ico", "./static/logo/favicon.ico")
	router.GET("/about", apiController.About)
	router.GET("/heartbeat", apiController.HeartBeat)

	router.GET("/ip", apiController.GetIP)
	router.GET("/headers", apiController.GetHeaders)
	router.GET("/user-agent", apiController.GetUserAgent)

	router.Any("/status/:statuscode", apiController.GetStatusCodes)

	return router
}
