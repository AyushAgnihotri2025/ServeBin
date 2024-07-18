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
		router.GET("/sitemap.xml", apiController.Redirect)
	} else {
		// Root Path
		router.GET("", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "ServeBin", "Version": ServeBin.Version[1:]})
		})

		// Add Swagger
		router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Sitemap route
		router.GET("/sitemap.xml", func(context *gin.Context) {
			apiController.GenerateSitemap(router, context)
		})
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

	router.GET("/status", apiController.GetStatusCodes)
	router.Any("/status/:statuscode", apiController.GetStatusCodes)

	router.GET("/image", apiController.GetImages)
	router.GET("/image/:imagetype", apiController.GetImages)

	router.GET("/xml", apiController.GetXML)
	router.GET("/html", apiController.GetHTML)
	router.GET("/json", apiController.GetJson)
	router.GET("/deny", apiController.GetDenyPath)
	router.GET("/gzip", apiController.Getgzip)
	router.GET("/brotli", apiController.Getbrotli)
	router.GET("/deflate", apiController.Getdeflate)
	router.GET("/zstd", apiController.Getzstd)
	router.GET("/robots.txt", apiController.GetRobotsTxt)

	router.HEAD("/head", apiController.ResponseHeaderData)
	router.GET("/get", apiController.ResponseData)
	router.POST("/post", apiController.ResponseBodyData)
	router.PUT("/put", apiController.ResponseBodyData)
	router.DELETE("/delete", apiController.ResponseBodyData)
	router.PATCH("/patch", apiController.ResponseBodyData)
	router.OPTIONS("/options", apiController.ResponseHeaderData)

	return router
}
