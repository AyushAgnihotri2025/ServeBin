// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package helper

import (
	"ServeBin/data/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// It create raise server error while it counters any error
func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// It create send error as response while it counters any failure in parsing of json
func WebErrorPanic(err error, ctx *gin.Context) {
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to marshal JSON")
		return
	}
}

// It create send error as response while it counters any error
func NewError(ctx *gin.Context, status int, err error) {
	er := response.HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}
