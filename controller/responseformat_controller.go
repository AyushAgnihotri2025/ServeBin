// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"ServeBin/data/response"
	"ServeBin/helper"
	"compress/flate"
	"compress/gzip"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"github.com/klauspost/compress/zstd"
	"net/http"
	"os"
)

// GetXML	 	 	ServeBin
// @Tags			Response formats
// @Summary			Returns a simple XML document.
// @Description		Returns a simple XML document.
// @produce			application/xml
// @Success			200
// @Failure      	406
// @Failure      	404
// @Failure      	500  		{object}  	response.HTTPError
// @Router			/xml 		[get]
func (controller *APIController) GetXML(ctx *gin.Context) {
	xmlData, err := os.ReadFile("templates/sample/sample.xml")
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to read XML file")
		return
	}
	ctx.Header("Content-Type", "application/xml")
	ctx.String(http.StatusOK, "%s", xmlData)
}

// GetHTML	 		ServeBin
// @Tags			Response formats
// @Summary			Returns a simple HTML document.
// @Description		Returns a simple HTML document.
// @produce			text/plain
// @Success			200
// @Failure      	406
// @Failure      	404
// @Failure      	500  		{object}  	response.HTTPError
// @Router			/html 		[get]
func (controller *APIController) GetHTML(ctx *gin.Context) {
	htmlData, err := os.ReadFile("templates/sample/sample.html")
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to read HTML file")
		return
	}
	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, "%s", htmlData)
}

// GetJson	 		ServeBin
// @Tags			Response formats
// @Summary			Returns a simple JSON document.
// @Description		Returns a simple JSON document.
// @produce			text/plain
// @Success			200
// @Failure      	406
// @Failure      	404
// @Failure      	500  		{object}  	response.HTTPError
// @Router			/json 		[get]
func (controller *APIController) GetJson(ctx *gin.Context) {
	jsonData, err := os.ReadFile("templates/sample/sample.json")
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to read JSON file")
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, "%s", jsonData)
}

// GetDenyPath	 	ServeBin
// @Tags			Response formats
// @Summary			Returns page denied by robots.txt rules.
// @Description		Returns page denied by robots.txt rules.
// @produce			text/plain
// @Success			200
// @Failure      	406
// @Failure      	404
// @Failure      	500  		{object}  	response.HTTPError
// @Router			/deny 		[get]
func (controller *APIController) GetDenyPath(ctx *gin.Context) {
	denyData, err := os.ReadFile("templates/sample/deny.txt")
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to read Deny file")
		return
	}
	ctx.Header("Content-Type", "text/plain")
	ctx.String(http.StatusOK, "%s", denyData)
}

// GetRobotsTxt	 	ServeBin
// @Tags			Response formats
// @Summary			Returns some robots.txt rules.
// @Description		Returns some robots.txt rules.
// @produce			text/plain
// @Success			200
// @Failure      	406
// @Failure      	404
// @Failure      	500  				{object}  	response.HTTPError
// @Router			/robots.txt 		[get]
func (controller *APIController) GetRobotsTxt(ctx *gin.Context) {
	robotsData, err := os.ReadFile("templates/sample/robots.txt")
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to read Robots.txt file")
		return
	}
	ctx.Header("Content-Type", "text/plain")
	ctx.String(http.StatusOK, "%s", robotsData)
}

// Getgzip	 		ServeBin
// @Tags			Response formats
// @Summary			Returns GZip-encoded data.
// @Description		Returns GZip-encoded data.
// @produce			text/plain
// @Success			200			{object} 	response.GzipResponse
// @Failure      	406
// @Failure      	404
// @Failure      	500  		{object}  	response.HTTPError
// @Router			/gzip 		[get]
func (controller *APIController) Getgzip(ctx *gin.Context) {
	header := helper.GetHeaders(ctx)
	ipResponse := controller.apiService.FindIP(ctx)

	webResponse := response.GzipResponse{
		HeaderResponse: response.HeaderResponse{Header: header},
		IPResponse:     response.IPResponse{IP: []interface{}{ipResponse}},
		Gzipped:        true,
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Encoding", "gzip")
	ww := gzip.NewWriter(ctx.Writer)
	defer ww.Close() // flush
	if err := helper.WriteJSON(ww, webResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write JSON"})
		return
	}
}

// Getbrotli	 	ServeBin
// @Tags			Response formats
// @Summary			Returns Brotli-encoded data.
// @Description		Returns Brotli-encoded data.
// @produce			text/plain
// @Success			200			{object} 	response.BrotliResponse
// @Failure      	406
// @Failure      	404
// @Failure      	500  		{object}  	response.HTTPError
// @Router			/brotli 	[get]
func (controller *APIController) Getbrotli(ctx *gin.Context) {
	header := helper.GetHeaders(ctx)
	ipResponse := controller.apiService.FindIP(ctx)

	webResponse := response.BrotliResponse{
		HeaderResponse: response.HeaderResponse{Header: header},
		IPResponse:     response.IPResponse{IP: []interface{}{ipResponse}},
		Compressed:     true,
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Encoding", "br")
	ww := brotli.NewWriter(ctx.Writer)
	defer ww.Close() // flush
	if err := helper.WriteJSON(ww, webResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write JSON"})
		return
	}
}

// Getdeflate	 	ServeBin
// @Tags			Response formats
// @Summary			Returns Deflate-encoded data.
// @Description		Returns Deflate-encoded data.
// @produce			text/plain
// @Success			200			{object} 	response.DeflateResponse
// @Failure      	406
// @Failure      	404
// @Failure      	500  		{object}  	response.HTTPError
// @Router			/deflate 	[get]
func (controller *APIController) Getdeflate(ctx *gin.Context) {
	header := helper.GetHeaders(ctx)
	ipResponse := controller.apiService.FindIP(ctx)

	webResponse := response.DeflateResponse{
		HeaderResponse: response.HeaderResponse{Header: header},
		IPResponse:     response.IPResponse{IP: []interface{}{ipResponse}},
		Deflated:       true,
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Encoding", "deflate")
	ww, _ := flate.NewWriter(ctx.Writer, flate.BestCompression)
	defer ww.Close() // flush
	if err := helper.WriteJSON(ww, webResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write JSON"})
		return
	}
}

// Getzstd	 		ServeBin
// @Tags			Response formats
// @Summary			Returns Zstd-compressed data.
// @Description		Returns Zstd-compressed data.
// @produce			text/plain
// @Success			200			{object}	response.ZstdResponse
// @Failure      	406
// @Failure      	404
// @Failure      	500  		{object}  	response.HTTPError
// @Router			/zstd 		[get]
func (controller *APIController) Getzstd(ctx *gin.Context) {
	header := helper.GetHeaders(ctx)
	ipResponse := controller.apiService.FindIP(ctx)

	webResponse := response.ZstdResponse{
		HeaderResponse: response.HeaderResponse{Header: header},
		IPResponse:     response.IPResponse{IP: []interface{}{ipResponse}},
		Compressed:     true,
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Content-Encoding", "zstd")

	bestLevel := zstd.WithEncoderLevel(zstd.SpeedBestCompression)
	ww, _ := zstd.NewWriter(ctx.Writer, bestLevel)

	defer ww.Close() // flush
	if err := helper.WriteJSON(ww, webResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write JSON"})
		return
	}
}
