// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package request

type Request struct {
	Name string `validate:"required,min=1,max=20" json:"name"`
}
