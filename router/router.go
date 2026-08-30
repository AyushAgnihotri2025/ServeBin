// Copyright 2026 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package router

import (
	"ServeBin"
	"ServeBin/controller"
	"ServeBin/middleware"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.GET("/uuid", apiController.GetUUID)
	router.GET("/headers", apiController.GetHeaders)
	router.GET("/user-agent", apiController.GetUserAgent)

	router.GET("/status", apiController.GetStatusCodes)
	router.Any("/status/:statuscode", apiController.GetStatusCodes)

	router.GET("/image", apiController.GetImages)
	router.GET("/image/:imagetype", apiController.GetImages)

	router.Any("/anything", apiController.Anything)
	router.Any("/anything/*anything", apiController.Anything)
	router.GET("/base64/*value", apiController.GetBase64)
	router.GET("/encoding/utf8", apiController.GetUTF8)
	router.GET("/xml", apiController.GetXML)
	router.GET("/html", apiController.GetHTML)
	router.GET("/json", apiController.GetJson)
	router.GET("/deny", apiController.GetDenyPath)
	router.GET("/gzip", apiController.Getgzip)
	router.GET("/brotli", apiController.Getbrotli)
	router.GET("/deflate", apiController.Getdeflate)
	router.GET("/zstd", apiController.Getzstd)
	router.GET("/robots.txt", apiController.GetRobotsTxt)
	router.GET("/response-headers", apiController.GetResponseHeaders)
	router.GET("/redirect/:n", apiController.GetRedirect)
	router.GET("/redirect-to", apiController.GetRedirectTo)
	router.GET("/relative-redirect/:n", apiController.GetRelativeRedirect)
	router.GET("/absolute-redirect/:n", apiController.GetAbsoluteRedirect)
	router.GET("/cookies", apiController.GetCookies)
	router.GET("/cookies/set", apiController.SetCookies)
	router.GET("/cookies/delete", apiController.DeleteCookies)
	router.GET("/basic-auth/:user/:passwd", apiController.GetBasicAuth)
	router.GET("/hidden-basic-auth/:user/:passwd", apiController.GetHiddenBasicAuth)
	router.GET("/digest-auth/:qop/:user/:passwd", apiController.GetDigestAuth)
	router.GET("/digest-auth/:qop/:user/:passwd/:algorithm", apiController.GetDigestAuth)
	router.GET("/stream/:n", apiController.GetStream)
	router.Any("/delay/:n", apiController.GetDelay)
	router.GET("/drip", apiController.GetDrip)
	router.GET("/range/:n", apiController.GetRange)
	router.GET("/cache", apiController.GetCache)
	router.GET("/etag/:etag", apiController.GetETag)
	router.GET("/cache/:n", apiController.GetCacheFor)
	router.GET("/bytes/:n", apiController.GetBytes)
	router.GET("/stream-bytes/:n", apiController.GetStreamBytes)
	router.GET("/links/:n", apiController.GetLinks)
	router.GET("/links/:n/:offset", apiController.GetLinks)
	router.GET("/forms/post", apiController.GetFormsPost)

	router.HEAD("/head", apiController.ResponseHeaderData)
	router.GET("/get", apiController.ResponseData)
	router.POST("/post", apiController.ResponseBodyData)
	router.PUT("/put", apiController.ResponseBodyData)
	router.DELETE("/delete", apiController.ResponseBodyData)
	router.PATCH("/patch", apiController.ResponseBodyData)
	router.OPTIONS("/options", apiController.ResponseHeaderData)

	return router
}
