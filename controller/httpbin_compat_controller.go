// Copyright 2026 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"ServeBin/data/response"
	"ServeBin/helper"
	"crypto/md5"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const maxGeneratedBytes = 102400

// GetUUID			ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Returns a UUID4.
// @Description		Returns a randomly generated UUID4.
// @Success			200 {object} response.UUIDResponse
// @Failure      	500 {object} response.HTTPError
// @Router			/uuid [get]
func (controller *APIController) GetUUID(ctx *gin.Context) {
	uuid, err := helper.GenerateUUID()
	if err != nil {
		helper.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, response.UUIDResponse{
		UUID: uuid,
	})
}

// Anything			ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Returns request data for any method.
// @Description		Returns request data, including method used.
// @Success			200 {object} response.HTTPBinRequestResponse
// @Failure      	500 {object} response.HTTPError
// @Router			/anything [get]
// @Router			/anything [post]
// @Router			/anything [put]
// @Router			/anything [patch]
// @Router			/anything [delete]
// @Router			/anything [head]
// @Router			/anything [options]
// @Router			/anything/{anything} [get]
// @Router			/anything/{anything} [post]
// @Router			/anything/{anything} [put]
// @Router			/anything/{anything} [patch]
// @Router			/anything/{anything} [delete]
// @Router			/anything/{anything} [head]
// @Router			/anything/{anything} [options]
func (controller *APIController) Anything(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, controller.buildHTTPBinRequestResponse(ctx))
}

// GetBase64		ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Decodes a base64 string.
// @Description		Decodes base64url-encoded text.
// @Success			200
// @Failure      	400 {object} response.HTTPError
// @Router			/base64/{value} [get]
func (controller *APIController) GetBase64(ctx *gin.Context) {
	value := strings.TrimPrefix(ctx.Param("value"), "/")
	decoded, err := decodeBase64Value(value)
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Data(http.StatusOK, "text/plain; charset=utf-8", decoded)
}

// GetUTF8			ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Returns a UTF-8 sample page.
// @Description		Returns a page containing UTF-8 data.
// @Success			200
// @Failure      	500 {object} response.HTTPError
// @Router			/encoding/utf8 [get]
func (controller *APIController) GetUTF8(ctx *gin.Context) {
	utf8Data, err := os.ReadFile("templates/sample/utf8.html")
	if err != nil {
		helper.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", utf8Data)
}

// GetResponseHeaders	ServeBin
// @Tags				HTTPBin Compatibility
// @Summary				Returns arbitrary response headers.
// @Description			Returns given response headers.
// @Success				200
// @Failure      		400 {object} response.HTTPError
// @Router				/response-headers [get]
func (controller *APIController) GetResponseHeaders(ctx *gin.Context) {
	body := make(map[string]interface{})
	for key, values := range ctx.Request.URL.Query() {
		for _, value := range values {
			ctx.Writer.Header().Add(key, value)
		}

		if len(values) == 1 {
			body[key] = values[0]
		} else {
			body[key] = values
		}
	}

	body["Content-Type"] = "application/json"
	for attempts := 0; attempts < 3; attempts++ {
		payload, err := json.Marshal(body)
		if err != nil {
			helper.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		contentLength := strconv.Itoa(len(payload))
		if body["Content-Length"] == contentLength {
			break
		}
		body["Content-Length"] = contentLength
	}

	ctx.JSON(http.StatusOK, body)
}

// GetRedirect		ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Redirects multiple times.
// @Description		302 redirects n times.
// @Success			302
// @Failure      	400 {object} response.HTTPError
// @Router			/redirect/{n} [get]
func (controller *APIController) GetRedirect(ctx *gin.Context) {
	count, ok := parsePathInt(ctx, "n")
	if !ok {
		return
	}

	if count <= 1 {
		controller.writeRedirect(ctx, "/get", http.StatusFound)
		return
	}

	controller.writeRedirect(ctx, fmt.Sprintf("/relative-redirect/%d", count-1), http.StatusFound)
}

// GetRedirectTo	ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Redirects to the requested URL.
// @Description		Redirects to the url query parameter using the optional status_code.
// @Success			302
// @Failure      	400 {object} response.HTTPError
// @Router			/redirect-to [get]
func (controller *APIController) GetRedirectTo(ctx *gin.Context) {
	target := ctx.Query("url")
	if target == "" {
		helper.NewError(ctx, http.StatusBadRequest, fmt.Errorf("url query parameter is required"))
		return
	}

	statusCode, err := helper.ParsePositiveInt(ctx.Query("status_code"), http.StatusFound)
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if statusCode < 300 || statusCode > 399 {
		helper.NewError(ctx, http.StatusBadRequest, fmt.Errorf("status_code must be a redirect status"))
		return
	}

	controller.writeRedirect(ctx, target, statusCode)
}

// GetRelativeRedirect	ServeBin
// @Tags				HTTPBin Compatibility
// @Summary				Performs relative redirects.
// @Description			302 relative redirects n times.
// @Success				302
// @Failure      		400 {object} response.HTTPError
// @Router				/relative-redirect/{n} [get]
func (controller *APIController) GetRelativeRedirect(ctx *gin.Context) {
	count, ok := parsePathInt(ctx, "n")
	if !ok {
		return
	}

	target := "/get"
	if count > 1 {
		target = fmt.Sprintf("/relative-redirect/%d", count-1)
	}

	controller.writeRedirect(ctx, target, http.StatusFound)
}

// GetAbsoluteRedirect	ServeBin
// @Tags				HTTPBin Compatibility
// @Summary				Performs absolute redirects.
// @Description			302 absolute redirects n times.
// @Success				302
// @Failure      		400 {object} response.HTTPError
// @Router				/absolute-redirect/{n} [get]
func (controller *APIController) GetAbsoluteRedirect(ctx *gin.Context) {
	count, ok := parsePathInt(ctx, "n")
	if !ok {
		return
	}

	target := helper.GetBaseURL(ctx) + "/get"
	if count > 1 {
		target = helper.GetBaseURL(ctx) + fmt.Sprintf("/absolute-redirect/%d", count-1)
	}

	controller.writeRedirect(ctx, target, http.StatusFound)
}

// GetCookies		ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Returns cookie data.
// @Description		Returns cookie data from the request.
// @Success			200 {object} response.CookiesResponse
// @Router			/cookies [get]
func (controller *APIController) GetCookies(ctx *gin.Context) {
	cookies := make(map[string]string)
	for _, cookie := range ctx.Request.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	ctx.JSON(http.StatusOK, response.CookiesResponse{
		Cookies: cookies,
	})
}

// SetCookies		ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Sets one or more simple cookies.
// @Description		Sets cookies from query parameters and redirects to /cookies.
// @Success			302
// @Router			/cookies/set [get]
func (controller *APIController) SetCookies(ctx *gin.Context) {
	for key, values := range ctx.Request.URL.Query() {
		for _, value := range values {
			http.SetCookie(ctx.Writer, &http.Cookie{
				Name:  key,
				Value: value,
				Path:  "/",
			})
		}
	}

	controller.writeRedirect(ctx, "/cookies", http.StatusFound)
}

// DeleteCookies	ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Deletes one or more simple cookies.
// @Description		Deletes cookies named in query parameters and redirects to /cookies.
// @Success			302
// @Router			/cookies/delete [get]
func (controller *APIController) DeleteCookies(ctx *gin.Context) {
	for key := range ctx.Request.URL.Query() {
		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:    key,
			Value:   "",
			MaxAge:  -1,
			Expires: time.Unix(0, 0),
			Path:    "/",
		})
	}

	controller.writeRedirect(ctx, "/cookies", http.StatusFound)
}

// GetBasicAuth		ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Challenges HTTP basic auth.
// @Description		Returns authenticated user data when the credentials match.
// @Success			200 {object} response.AuthResponse
// @Failure      	401
// @Router			/basic-auth/{user}/{passwd} [get]
func (controller *APIController) GetBasicAuth(ctx *gin.Context) {
	user := ctx.Param("user")
	passwd := ctx.Param("passwd")

	username, password, ok := ctx.Request.BasicAuth()
	if ok &&
		subtle.ConstantTimeCompare([]byte(username), []byte(user)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(passwd)) == 1 {
		ctx.JSON(http.StatusOK, response.AuthResponse{
			Authenticated: true,
			User:          user,
		})
		return
	}

	ctx.Header("WWW-Authenticate", `Basic realm="Fake Realm"`)
	ctx.Status(http.StatusUnauthorized)
}

// GetHiddenBasicAuth	ServeBin
// @Tags				HTTPBin Compatibility
// @Summary				404'd basic auth.
// @Description			Returns authenticated user data when the credentials match, otherwise 404s.
// @Success				200 {object} response.AuthResponse
// @Failure      		404
// @Router				/hidden-basic-auth/{user}/{passwd} [get]
func (controller *APIController) GetHiddenBasicAuth(ctx *gin.Context) {
	user := ctx.Param("user")
	passwd := ctx.Param("passwd")

	username, password, ok := ctx.Request.BasicAuth()
	if ok &&
		subtle.ConstantTimeCompare([]byte(username), []byte(user)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(passwd)) == 1 {
		ctx.JSON(http.StatusOK, response.AuthResponse{
			Authenticated: true,
			User:          user,
		})
		return
	}

	ctx.Status(http.StatusNotFound)
}

// GetDigestAuth	ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Challenges HTTP digest auth.
// @Description		Supports qop=auth with MD5 or SHA-256.
// @Success			200 {object} response.AuthResponse
// @Failure      	400 {object} response.HTTPError
// @Failure      	401
// @Router			/digest-auth/{qop}/{user}/{passwd} [get]
// @Router			/digest-auth/{qop}/{user}/{passwd}/{algorithm} [get]
func (controller *APIController) GetDigestAuth(ctx *gin.Context) {
	qop := ctx.Param("qop")
	if !strings.EqualFold(qop, "auth") {
		helper.NewError(ctx, http.StatusBadRequest, fmt.Errorf("only qop=auth is supported"))
		return
	}

	algorithm := strings.ToUpper(ctx.Param("algorithm"))
	if algorithm == "" {
		algorithm = "MD5"
	}
	if algorithm != "MD5" && algorithm != "SHA-256" {
		helper.NewError(ctx, http.StatusBadRequest, fmt.Errorf("unsupported digest algorithm"))
		return
	}

	user := ctx.Param("user")
	passwd := ctx.Param("passwd")
	challenge := buildDigestChallenge(user, passwd, qop, algorithm)
	authorization := parseDigestAuthorization(ctx.GetHeader("Authorization"))

	if verifyDigestAuthorization(authorization, challenge, ctx.Request.Method, ctx.Request.URL.RequestURI(), user, passwd) {
		ctx.JSON(http.StatusOK, response.AuthResponse{
			Authenticated: true,
			User:          user,
		})
		return
	}

	ctx.Header("WWW-Authenticate", challenge.headerValue())
	ctx.Status(http.StatusUnauthorized)
}

// GetStream		ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Streams request data as JSON lines.
// @Description		Streams min(n, 100) JSON lines.
// @Success			200
// @Failure      	400 {object} response.HTTPError
// @Router			/stream/{n} [get]
func (controller *APIController) GetStream(ctx *gin.Context) {
	count, ok := parsePathInt(ctx, "n")
	if !ok {
		return
	}
	count = helper.ClampInt(count, 0, 100)

	webResponse := controller.buildHTTPBinRequestResponse(ctx)
	ctx.Header("Content-Type", "application/json")
	ctx.Status(http.StatusOK)

	for index := 0; index < count; index++ {
		line, err := json.Marshal(response.HTTPBinStreamResponse{
			HTTPBinRequestResponse: webResponse,
			Id:                     index,
		})
		if err != nil {
			helper.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		if _, err = ctx.Writer.Write(append(line, '\n')); err != nil {
			return
		}
		ctx.Writer.Flush()
	}
}

// GetDelay			ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Delays the response.
// @Description		Delays responding for min(n, 10) seconds.
// @Success			200 {object} response.HTTPBinRequestResponse
// @Failure      	400 {object} response.HTTPError
// @Router			/delay/{n} [get]
// @Router			/delay/{n} [post]
// @Router			/delay/{n} [put]
// @Router			/delay/{n} [patch]
// @Router			/delay/{n} [delete]
// @Router			/delay/{n} [head]
// @Router			/delay/{n} [options]
func (controller *APIController) GetDelay(ctx *gin.Context) {
	delay, ok := parsePathInt(ctx, "n")
	if !ok {
		return
	}
	delay = helper.ClampInt(delay, 0, 10)

	time.Sleep(time.Duration(delay) * time.Second)
	ctx.JSON(http.StatusOK, controller.buildHTTPBinRequestResponse(ctx))
}

// GetDrip			ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Drips bytes over a duration.
// @Description		Drips data over a duration after an optional initial delay.
// @Success			200
// @Failure      	400 {object} response.HTTPError
// @Router			/drip [get]
func (controller *APIController) GetDrip(ctx *gin.Context) {
	numBytes, err := helper.ParsePositiveInt(ctx.Query("numbytes"), 10)
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if numBytes > maxGeneratedBytes {
		helper.NewError(ctx, http.StatusBadRequest, fmt.Errorf("numbytes must be less than or equal to %d", maxGeneratedBytes))
		return
	}

	durationSeconds, err := helper.ParsePositiveFloat(ctx.Query("duration"), 2)
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	initialDelay, err := helper.ParsePositiveFloat(ctx.Query("delay"), 0)
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	statusCode, err := helper.ParsePositiveInt(ctx.Query("code"), http.StatusOK)
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if statusCode < http.StatusContinue || statusCode > http.StatusNetworkAuthenticationRequired {
		helper.NewError(ctx, http.StatusBadRequest, fmt.Errorf("code must be between %d and %d", http.StatusContinue, http.StatusNetworkAuthenticationRequired))
		return
	}

	if initialDelay > 0 {
		time.Sleep(time.Duration(initialDelay * float64(time.Second)))
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Length", strconv.Itoa(numBytes))
	ctx.Status(statusCode)

	if numBytes == 0 {
		return
	}

	intervals := helper.SleepIntervals(durationSeconds, numBytes)
	for index := 0; index < numBytes; index++ {
		if _, err = ctx.Writer.Write([]byte("*")); err != nil {
			return
		}
		ctx.Writer.Flush()

		if index < len(intervals) {
			time.Sleep(time.Duration(intervals[index] * float64(time.Second)))
		}
	}
}

// GetRange			ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Returns a range-enabled byte stream.
// @Description		Streams n bytes and honors Range requests.
// @Success			200
// @Failure      	400 {object} response.HTTPError
// @Failure      	416
// @Router			/range/{n} [get]
func (controller *APIController) GetRange(ctx *gin.Context) {
	totalLength, ok := parsePathInt(ctx, "n")
	if !ok {
		return
	}
	totalLength = helper.ClampInt(totalLength, 0, maxGeneratedBytes)

	chunkSize, err := helper.ParsePositiveInt(ctx.Query("chunk_size"), totalLength)
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if chunkSize == 0 {
		chunkSize = totalLength
	}

	durationSeconds, err := helper.ParsePositiveFloat(ctx.Query("duration"), 0)
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	body := helper.GenerateAlphabetBytes(totalLength)
	statusCode := http.StatusOK
	start := 0
	end := totalLength - 1

	rangeHeader := ctx.GetHeader("Range")
	if totalLength > 0 && rangeHeader != "" {
		start, end, err = helper.ParseRangeHeader(rangeHeader, totalLength)
		if err != nil {
			ctx.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		statusCode = http.StatusPartialContent
		body = body[start : end+1]
	}

	ctx.Header("Accept-Ranges", "bytes")
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("ETag", fmt.Sprintf("range%d", totalLength))
	ctx.Header("Content-Length", strconv.Itoa(len(body)))
	if statusCode == http.StatusPartialContent {
		ctx.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, totalLength))
	}
	ctx.Status(statusCode)

	if len(body) == 0 {
		return
	}

	writeBodyInChunks(ctx, body, chunkSize, durationSeconds)
}

// GetCache			ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Returns a cacheable response.
// @Description		Returns 304 when If-Modified-Since or If-None-Match is provided.
// @Success			200 {object} response.HTTPBinRequestResponse
// @Failure      	304
// @Router			/cache [get]
func (controller *APIController) GetCache(ctx *gin.Context) {
	etag := fmt.Sprintf("%x", md5.Sum([]byte(helper.GetRequestURL(ctx))))
	lastModified := time.Now().UTC().Format(http.TimeFormat)

	ctx.Header("ETag", etag)
	ctx.Header("Last-Modified", lastModified)

	if ctx.GetHeader("If-Modified-Since") != "" || ctx.GetHeader("If-None-Match") != "" {
		ctx.Status(http.StatusNotModified)
		return
	}

	ctx.JSON(http.StatusOK, controller.buildHTTPBinRequestResponse(ctx))
}

// GetETag			ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Returns an ETag-aware response.
// @Description		Responds to If-None-Match with 304 and If-Match with 412 when appropriate.
// @Success			200 {object} response.HTTPBinRequestResponse
// @Failure      	304
// @Failure      	412
// @Router			/etag/{etag} [get]
func (controller *APIController) GetETag(ctx *gin.Context) {
	etag := ctx.Param("etag")
	ctx.Header("ETag", etag)

	if ifMatch := ctx.GetHeader("If-Match"); ifMatch != "" && !etagMatches(ifMatch, etag) {
		ctx.Status(http.StatusPreconditionFailed)
		return
	}

	if ifNoneMatch := ctx.GetHeader("If-None-Match"); ifNoneMatch != "" && etagMatches(ifNoneMatch, etag) {
		ctx.Status(http.StatusNotModified)
		return
	}

	ctx.JSON(http.StatusOK, controller.buildHTTPBinRequestResponse(ctx))
}

// GetCacheFor		ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Sets a Cache-Control header.
// @Description		Sets Cache-Control for n seconds.
// @Success			200 {object} response.HTTPBinRequestResponse
// @Failure      	400 {object} response.HTTPError
// @Router			/cache/{n} [get]
func (controller *APIController) GetCacheFor(ctx *gin.Context) {
	maxAge, ok := parsePathInt(ctx, "n")
	if !ok {
		return
	}
	if maxAge < 0 {
		helper.NewError(ctx, http.StatusBadRequest, fmt.Errorf("n must be non-negative"))
		return
	}

	ctx.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	ctx.JSON(http.StatusOK, controller.buildHTTPBinRequestResponse(ctx))
}

// GetBytes			ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Generates random bytes.
// @Description		Generates n random bytes of binary data.
// @Success			200
// @Failure      	400 {object} response.HTTPError
// @Failure      	500 {object} response.HTTPError
// @Router			/bytes/{n} [get]
func (controller *APIController) GetBytes(ctx *gin.Context) {
	length, ok := parsePathInt(ctx, "n")
	if !ok {
		return
	}
	length = helper.ClampInt(length, 0, maxGeneratedBytes)

	seed, err := helper.ParseSeed(ctx.Query("seed"))
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	body, err := helper.GenerateRandomBytes(length, seed)
	if err != nil {
		helper.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Length", strconv.Itoa(len(body)))
	ctx.Data(http.StatusOK, "application/octet-stream", body)
}

// GetStreamBytes	ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Streams random bytes.
// @Description		Streams n random bytes of binary data in chunked encoding.
// @Success			200
// @Failure      	400 {object} response.HTTPError
// @Failure      	500 {object} response.HTTPError
// @Router			/stream-bytes/{n} [get]
func (controller *APIController) GetStreamBytes(ctx *gin.Context) {
	length, ok := parsePathInt(ctx, "n")
	if !ok {
		return
	}
	length = helper.ClampInt(length, 0, maxGeneratedBytes)

	chunkSize, err := helper.ParsePositiveInt(ctx.Query("chunk_size"), 10)
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if chunkSize == 0 {
		chunkSize = 10
	}

	seed, err := helper.ParseSeed(ctx.Query("seed"))
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	body, err := helper.GenerateRandomBytes(length, seed)
	if err != nil {
		helper.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Status(http.StatusOK)
	writeBodyInChunks(ctx, body, chunkSize, 0)
}

// GetLinks			ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Returns a page of HTML links.
// @Description		Returns a page containing n HTML links.
// @Success			200
// @Failure      	400 {object} response.HTTPError
// @Router			/links/{n} [get]
// @Router			/links/{n}/{offset} [get]
func (controller *APIController) GetLinks(ctx *gin.Context) {
	total, ok := parsePathInt(ctx, "n")
	if !ok {
		return
	}
	total = helper.ClampInt(total, 0, 256)

	offset := 0
	offsetValue := ctx.Param("offset")
	if offsetValue != "" {
		var err error
		offset, err = strconv.Atoi(offsetValue)
		if err != nil {
			helper.NewError(ctx, http.StatusBadRequest, err)
			return
		}
	}
	if total == 0 || offset < 0 || offset >= total {
		ctx.Status(http.StatusNotFound)
		return
	}

	var builder strings.Builder
	builder.WriteString("<html><head><title>Links</title></head><body>")
	for index := 0; index < total; index++ {
		if index > 0 {
			builder.WriteByte(' ')
		}
		if index == offset {
			builder.WriteString(strconv.Itoa(index))
			continue
		}
		builder.WriteString(fmt.Sprintf("<a href='/links/%d/%d'>%d</a>", total, index, index))
	}
	builder.WriteString(" </body></html>")

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(builder.String()))
}

// GetFormsPost		ServeBin
// @Tags			HTTPBin Compatibility
// @Summary			Returns an HTML form.
// @Description		Returns an HTML form that submits to /post.
// @Success			200
// @Failure      	500 {object} response.HTTPError
// @Router			/forms/post [get]
func (controller *APIController) GetFormsPost(ctx *gin.Context) {
	formData, err := os.ReadFile("templates/sample/forms_post.html")
	if err != nil {
		helper.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", formData)
}

func (controller *APIController) buildHTTPBinRequestResponse(ctx *gin.Context) response.HTTPBinRequestResponse {
	args, _ := controller.apiService.ReturnArguments(ctx)

	_ = ctx.Request.ParseMultipartForm(10 << 20)

	form, _ := controller.apiService.ReturnFormData(ctx)
	files, _ := controller.apiService.ReturnFormFile(ctx)
	rawData, _ := controller.apiService.ReturnJson_RawData(ctx)

	data := ""
	if rawValue, ok := rawData["rawData"].(string); ok {
		data = rawValue
	}

	var jsonValue interface{}
	if value, exists := rawData["json"]; exists {
		if stringValue, ok := value.(string); ok {
			if stringValue != "" {
				jsonValue = stringValue
			}
		} else {
			jsonValue = value
		}
	}

	return response.HTTPBinRequestResponse{
		Args:    normalizeInterfaceMap(args),
		Data:    data,
		Files:   normalizeInterfaceMap(files),
		Form:    normalizeInterfaceMap(form),
		Headers: helper.GetHeaders(ctx),
		Json:    jsonValue,
		Method:  ctx.Request.Method,
		Origin:  controller.apiService.FindIP(ctx),
		Url:     helper.GetRequestURL(ctx),
	}
}

func (controller *APIController) writeRedirect(ctx *gin.Context, target string, statusCode int) {
	ctx.Header("Location", target)
	ctx.Status(statusCode)
}

func normalizeInterfaceMap(value map[string]interface{}) map[string]interface{} {
	if value == nil {
		return map[string]interface{}{}
	}

	return value
}

func parsePathInt(ctx *gin.Context, name string) (int, bool) {
	value, err := strconv.Atoi(ctx.Param(name))
	if err != nil {
		helper.NewError(ctx, http.StatusBadRequest, err)
		return 0, false
	}

	return value, true
}

func decodeBase64Value(value string) ([]byte, error) {
	decoders := []func(string) ([]byte, error){
		base64.RawURLEncoding.DecodeString,
		base64.URLEncoding.DecodeString,
		base64.RawStdEncoding.DecodeString,
		base64.StdEncoding.DecodeString,
	}

	for _, decode := range decoders {
		decoded, err := decode(value)
		if err == nil {
			return decoded, nil
		}
	}

	return nil, fmt.Errorf("invalid base64 value")
}

func writeBodyInChunks(ctx *gin.Context, body []byte, chunkSize int, durationSeconds float64) {
	if chunkSize <= 0 || chunkSize >= len(body) {
		_, _ = ctx.Writer.Write(body)
		ctx.Writer.Flush()
		return
	}

	chunks := (len(body) + chunkSize - 1) / chunkSize
	intervals := helper.SleepIntervals(durationSeconds, chunks)
	chunkIndex := 0

	for start := 0; start < len(body); start += chunkSize {
		end := start + chunkSize
		if end > len(body) {
			end = len(body)
		}

		if _, err := ctx.Writer.Write(body[start:end]); err != nil {
			return
		}
		ctx.Writer.Flush()

		if chunkIndex < len(intervals) {
			time.Sleep(time.Duration(intervals[chunkIndex] * float64(time.Second)))
		}
		chunkIndex++
	}
}

type digestChallenge struct {
	Realm     string
	Nonce     string
	Opaque    string
	QOP       string
	Algorithm string
}

func buildDigestChallenge(user string, passwd string, qop string, algorithm string) digestChallenge {
	realm := "ServeBin"
	nonce := digestHash(algorithm, strings.Join([]string{"nonce", realm, user, passwd, qop}, ":"))
	opaque := digestHash(algorithm, strings.Join([]string{"opaque", realm, user, algorithm}, ":"))

	return digestChallenge{
		Realm:     realm,
		Nonce:     nonce,
		Opaque:    opaque,
		QOP:       qop,
		Algorithm: algorithm,
	}
}

func (challenge digestChallenge) headerValue() string {
	return fmt.Sprintf(
		`Digest realm="%s", nonce="%s", qop="%s", opaque="%s", algorithm=%s, stale=FALSE`,
		challenge.Realm,
		challenge.Nonce,
		challenge.QOP,
		challenge.Opaque,
		challenge.Algorithm,
	)
}

func parseDigestAuthorization(header string) map[string]string {
	if !strings.HasPrefix(header, "Digest ") {
		return map[string]string{}
	}

	content := strings.TrimPrefix(header, "Digest ")
	parts := make([]string, 0)
	var builder strings.Builder
	inQuotes := false

	for _, char := range content {
		switch char {
		case '"':
			inQuotes = !inQuotes
			builder.WriteRune(char)
		case ',':
			if inQuotes {
				builder.WriteRune(char)
				continue
			}
			parts = append(parts, builder.String())
			builder.Reset()
		default:
			builder.WriteRune(char)
		}
	}

	if builder.Len() > 0 {
		parts = append(parts, builder.String())
	}

	values := make(map[string]string, len(parts))
	for _, part := range parts {
		keyValue := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(keyValue) != 2 {
			continue
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.Trim(strings.TrimSpace(keyValue[1]), `"`)
		values[key] = value
	}

	return values
}

func verifyDigestAuthorization(values map[string]string, challenge digestChallenge, method string, requestURI string, user string, passwd string) bool {
	if len(values) == 0 {
		return false
	}

	if values["username"] != user ||
		values["realm"] != challenge.Realm ||
		values["nonce"] != challenge.Nonce ||
		values["uri"] != requestURI {
		return false
	}

	if opaque := values["opaque"]; opaque != "" && opaque != challenge.Opaque {
		return false
	}

	ha1 := digestHash(challenge.Algorithm, strings.Join([]string{user, challenge.Realm, passwd}, ":"))
	ha2 := digestHash(challenge.Algorithm, strings.Join([]string{method, requestURI}, ":"))

	expected := ""
	if values["qop"] != "" {
		if values["qop"] != challenge.QOP || values["cnonce"] == "" || values["nc"] == "" {
			return false
		}
		expected = digestHash(challenge.Algorithm, strings.Join([]string{
			ha1,
			challenge.Nonce,
			values["nc"],
			values["cnonce"],
			values["qop"],
			ha2,
		}, ":"))
	} else {
		expected = digestHash(challenge.Algorithm, strings.Join([]string{ha1, challenge.Nonce, ha2}, ":"))
	}

	return subtle.ConstantTimeCompare([]byte(expected), []byte(values["response"])) == 1
}

func digestHash(algorithm string, value string) string {
	var hasher hash.Hash
	switch algorithm {
	case "SHA-256":
		hasher = sha256.New()
	default:
		hasher = md5.New()
	}

	_, _ = hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}

func etagMatches(headerValue string, etag string) bool {
	if headerValue == "*" {
		return true
	}

	for _, value := range strings.Split(headerValue, ",") {
		trimmed := strings.TrimSpace(value)
		trimmed = strings.TrimPrefix(trimmed, "W/")
		trimmed = strings.Trim(trimmed, `"`)
		if trimmed == etag {
			return true
		}
	}

	return false
}
