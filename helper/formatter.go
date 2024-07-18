// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

func GetHeaders(ctx *gin.Context) map[string]string {
	hdr := make(map[string]string, len(ctx.Request.Header))
	for key, value := range ctx.Request.Header {
		hdr[key] = value[0]
	}
	return hdr
}

func GetFileHeaders(fileHeader *multipart.FileHeader) map[string]string {
	hdr := make(map[string]string, len(fileHeader.Header))
	for key, value := range fileHeader.Header {
		hdr[key] = value[0]
	}
	return hdr
}

func HumanBytes(size int64) string {
	if size == 0 {
		return ""
	}
	const power = 1024
	n := 0
	dicPowerN := map[int]string{0: "", 1: "Ki", 2: "Mi", 3: "Gi", 4: "Ti"}
	fsize := float64(size)
	for fsize > power {
		fsize /= power
		n++
	}
	return fmt.Sprintf("%.2f %sB", fsize, dicPowerN[n])
}
