// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"ServeBin"
	"ServeBin/data/response"
	"ServeBin/helper"
	"ServeBin/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

type APIController struct {
	apiService service.APIService
}

func NewAPIController(service service.APIService) *APIController {
	return &APIController{
		apiService: service,
	}
}

// About route handler
func (controller *APIController) About(ctx *gin.Context) {
	Response := response.AboutResponse{
		Version:    ServeBin.Version,
		ServerTime: time.Now().String(),
		Developer:  "Ayush Agnihotri <AyushAgnihotri2025>",
		Contact:    "contact@mrayush.me",
		SourceCode: "https://github.com/AyushAgnihotri2025/ServeBin",
	}
	ctx.Header("Content-Type", "application/json")
	ctx.PureJSON(http.StatusOK, Response)
}

// HeartBeat route handler
func (controller *APIController) HeartBeat(ctx *gin.Context) {
	stats := helper.GetHeartbeats()

	Response, err := json.Marshal(response.HeartbeatResponse{
		Stats:  stats,
		Status: "Up",
	})
	helper.ErrorPanic(err)

	ctx.Header("Content-Type", "application/json")
	ctx.Data(http.StatusOK, "application/json", Response)
}

// Redirection for the backup servers
func (controller *APIController) Redirect(ctx *gin.Context) {
	mainServerUrl := os.Getenv("MAIN_SERVER")
	RedirectUrl := mainServerUrl + ctx.Request.URL.Path

	ctx.Redirect(http.StatusTemporaryRedirect, RedirectUrl)
}

// Function to generate sitemap.xml
func (controller *APIController) GenerateSitemap(router *gin.Engine, ctx *gin.Context) {
	// Get the base URL
	var baseURL string
	ssl := os.Getenv("IS_SSL")
	if strings.ToLower(ssl) == "true" {
		baseURL = "https://" + ctx.Request.Host
	} else {
		baseURL = "http://" + ctx.Request.Host
	}

	// Initialize the XML string
	xmlStr := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:news="http://www.google.com/schemas/sitemap-news/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml" xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0" xmlns:image="http://www.google.com/schemas/sitemap-image/1.1" xmlns:video="http://www.google.com/schemas/sitemap-video/1.1">`

	// Loop through each registered route
	for _, route := range router.Routes() {
		// Construct the full URL
		fullURL := fmt.Sprintf("%s%s", baseURL, route.Path)

		// Skip any path with parameters
		skipPath := false
		cleanURL := strings.Split(fullURL, "/")
		for _, segment := range cleanURL {
			if strings.HasPrefix(segment, ":") || strings.HasPrefix(segment, "*") || strings.Contains(segment, "sitemap.xml") {
				skipPath = true
				continue
			}
		}
		if skipPath {
			continue
		} else {
			cleanFullURL := strings.Join(cleanURL, "/")
			currentTime := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")

			// Add the clean URL to the sitemap
			xmlStr += fmt.Sprintf("<url><loc>%s</loc><lastmod>%s</lastmod><changefreq>always</changefreq><priority>0.8</priority></url>", cleanFullURL, currentTime)
		}
	}

	// Close the XML string
	xmlStr += "</urlset>"

	ctx.Header("Content-Type", "text/xml; charset=utf-8")
	ctx.String(http.StatusOK, "%s", xmlStr)
}
