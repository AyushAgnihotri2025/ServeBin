// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package helper

import "github.com/gin-gonic/gin"

func GetHeaders(ctx *gin.Context) map[string]string {
	hdr := make(map[string]string, len(ctx.Request.Header))
	for key, value := range ctx.Request.Header {
		hdr[key] = value[0]
	}
	return hdr
}
