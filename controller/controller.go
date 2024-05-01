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
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
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
		Developer:  "Coder",
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, Response)
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
