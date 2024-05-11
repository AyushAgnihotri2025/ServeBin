// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"ServeBin/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetStatusCodes 	ServeBin
// @Tags			Status Codes
// @Summary			Return status code or random status code if more than one are given.
// @Description		Returns the requester's IP Address.
// @Param        	statuscode   path  int  true  "Status Code"
// @Default			200
// @Success			200
// @Failure      	400
// @Failure      	404
// @Failure      	500  	{object}  	response.HTTPError
// @Router			/status/{statuscode} [get]
// @Router			/status/{statuscode} [post]
// @Router			/status/{statuscode} [delete]
// @Router			/status/{statuscode} [put]
// @Router			/status/{statuscode} [patch]
func (controller *APIController) GetStatusCodes(ctx *gin.Context) {
	statusCode := ctx.Param("statuscode")
	if statusCode == "" {
		statusCode = "200"
	}
	status, err := strconv.ParseInt(statusCode, 10, 64)
	if err != nil {
		helper.NewError(ctx, http.StatusNotFound, err)
		return
	}
	helper.ErrorPanic(err)
	ctx.Status(int(status))
}
