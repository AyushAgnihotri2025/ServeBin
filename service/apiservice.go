// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"github.com/gin-gonic/gin"
)

type APIService interface {
	FindIP(ctx *gin.Context) string
	GeneratePNG() ([]byte, error)
	GenerateJPEG() ([]byte, error)
	GenerateSVG() ([]byte, error)
	GenerateGIF() ([]byte, error)
	GenerateWEBP() ([]byte, error)
	GenerateTIFF() ([]byte, error)
	GenerateBMP() ([]byte, error)
	GenerateAPNG() ([]byte, error)
	GenerateAVIF() ([]byte, error)
	GenerateICO() ([]byte, error)
}
