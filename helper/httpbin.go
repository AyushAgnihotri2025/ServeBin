// Copyright 2026 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package helper

import (
	"crypto/rand"
	"fmt"
	"math"
	mathrand "math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBaseURL(ctx *gin.Context) string {
	scheme := ctx.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		if ctx.Request.TLS != nil || strings.EqualFold(os.Getenv("IS_SSL"), "true") {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	return scheme + "://" + ctx.Request.Host
}

func GetRequestURL(ctx *gin.Context) string {
	requestURL := *ctx.Request.URL
	requestURL.Scheme = ""
	requestURL.Host = ""

	path := requestURL.String()
	if path == "" {
		path = ctx.Request.URL.Path
	}

	return GetBaseURL(ctx) + path
}

func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16],
	), nil
}

func GenerateRandomBytes(length int, seed *int64) ([]byte, error) {
	if length < 0 {
		length = 0
	}

	output := make([]byte, length)
	if seed != nil {
		random := mathrand.New(mathrand.NewSource(*seed))
		for index := range output {
			output[index] = byte(random.Intn(256))
		}
		return output, nil
	}

	_, err := rand.Read(output)
	return output, err
}

func GenerateAlphabetBytes(length int) []byte {
	if length < 0 {
		length = 0
	}

	alphabet := []byte("abcdefghijklmnopqrstuvwxyz")
	output := make([]byte, length)
	for index := range output {
		output[index] = alphabet[index%len(alphabet)]
	}

	return output
}

func ParseSeed(value string) (*int64, error) {
	if value == "" {
		return nil, nil
	}

	seed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}

	return &seed, nil
}

func ParsePositiveInt(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if parsed < 0 {
		return 0, fmt.Errorf("value must be non-negative")
	}

	return parsed, nil
}

func ParsePositiveFloat(value string, defaultValue float64) (float64, error) {
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	if parsed < 0 {
		return 0, fmt.Errorf("value must be non-negative")
	}

	return parsed, nil
}

func ClampInt(value int, minValue int, maxValue int) int {
	return max(minValue, min(value, maxValue))
}

func ParseRangeHeader(rangeHeader string, size int) (int, int, error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, fmt.Errorf("invalid range unit")
	}

	rangeValue := strings.TrimPrefix(rangeHeader, "bytes=")
	if strings.Contains(rangeValue, ",") {
		return 0, 0, fmt.Errorf("multiple ranges are not supported")
	}

	parts := strings.SplitN(rangeValue, "-", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range value")
	}

	if parts[0] == "" {
		length, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, err
		}
		if length <= 0 {
			return 0, 0, fmt.Errorf("invalid suffix range")
		}
		if length > size {
			length = size
		}
		return size - length, size - 1, nil
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	if start < 0 || start >= size {
		return 0, 0, fmt.Errorf("range start out of bounds")
	}

	if parts[1] == "" {
		return start, size - 1, nil
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	if end < start {
		return 0, 0, fmt.Errorf("range end before start")
	}
	if end >= size {
		end = size - 1
	}

	return start, end, nil
}

func BuildQueryRedirectURL(target string, query url.Values) string {
	if len(query) == 0 {
		return target
	}

	encoded := query.Encode()
	if encoded == "" {
		return target
	}

	return target + "?" + encoded
}

func SleepIntervals(durationSeconds float64, chunks int) []float64 {
	if durationSeconds <= 0 || chunks <= 1 {
		return nil
	}

	interval := durationSeconds / float64(chunks-1)
	intervals := make([]float64, 0, chunks-1)
	for index := 0; index < chunks-1; index++ {
		intervals = append(intervals, math.Max(interval, 0))
	}

	return intervals
}
