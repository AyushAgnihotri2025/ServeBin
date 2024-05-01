// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package request

type StatusCodesRequest struct {
	StatusCode int `validate:"required,min=100,max=599" json:"statuscode"`
}
