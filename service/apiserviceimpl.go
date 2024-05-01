// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type APIServiceImpl struct {
	Validate *validator.Validate
}

func NewAPIServiceImpl(validate *validator.Validate) APIService {
	return &APIServiceImpl{
		Validate: validate,
	}
}

// FindIP implements APIService
func (t *APIServiceImpl) FindIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	return ip
}
