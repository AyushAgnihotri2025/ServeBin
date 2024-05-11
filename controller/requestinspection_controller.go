// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"ServeBin/data/response"
	"ServeBin/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetIP 			ServeBin
// @Tags			Request inspection
// @Summary			Get Request IP.
// @Description		Returns the requester's IP Address.
// @Success			200 {object} response.IPResponse{}
// @Router			/ip [get]
func (controller *APIController) GetIP(ctx *gin.Context) {
	ipResponse := controller.apiService.FindIP(ctx)
	webResponse := response.IPResponse{
		IP: []interface{}{ipResponse},
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// GetHeaders 		ServeBin
// @Tags			Request inspection
// @Summary			Return the incoming request's HTTP headers.
// @Description		It returns the incoming request's HTTP headers.
// @Success			200 {object} response.HeaderResponse{}
// @Router			/headers [get]
func (controller *APIController) GetHeaders(ctx *gin.Context) {
	header := helper.GetHeaders(ctx)
	webResponse := response.HeaderResponse{
		Header: header,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// GetUserAgent 	ServeBin
// @Tags			Request inspection
// @Summary			Return the incoming request's User-Agent header.
// @Description		User-Agent header.
// @Success			200 {object} response.UserAgentResponse{}
// @Router			/user-agent [get]
func (controller *APIController) GetUserAgent(ctx *gin.Context) {
	webResponse := response.UserAgentResponse{
		UserAgent: ctx.Request.UserAgent(),
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
