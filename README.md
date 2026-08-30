<h1 align="center">ServeBin</h1>
<p align="center">
    <a href="https://github.com/AyushAgnihotri2025/ServeBin">
        <img src="https://cdn.mrayush.me/ServeBin/favicon.png" height="100" width="100" alt="ServeBin Logo" title="ServeBin" />
    </a>
</p>
<p align="center">
    Elevating HTTP Testing to <b>Effortless Debugging</b>.
    <br />
</p>
<hr>

ServeBin is a Go-based HTTP testing and debugging server with two layers of functionality:

- native ServeBin endpoints for request inspection, status-code testing, sample responses, and dynamic image generation
- a growing set of `httpbin`-compatible utilities for redirects, cookies, auth challenges, caching, ranges, delays, streaming, and more

The project also ships a landing page, Swagger UI, sample templates, and static assets so you can run the service locally and inspect behavior quickly.

## Highlights

- Inspect request IPs, headers, query parameters, form fields, JSON payloads, raw bodies, and uploaded files
- Return specific HTTP status codes for any method
- Exercise common response formats: HTML, JSON, XML, plain text, gzip, brotli, deflate, and zstd
- Generate images on the fly in `png`, `jpeg`, `svg`, `gif`, `webp`, `tiff`, `bmp`, `apng`, `avif`, and `ico`
- Use `httpbin`-style helpers for redirects, cookies, auth, caching, streaming, delayed responses, byte ranges, and random byte generation
- Browse the API through Swagger UI at [`/docs`](#api-docs)

## Requirements

- Go 1.22 or newer
- [`cwebp`](https://developers.google.com/speed/webp/download) in `PATH` if you want `GET /image/webp`
- [`avifenc`](https://github.com/AOMediaCodec/libavif?tab=readme-ov-file#installation) in `PATH` if you want `GET /image/avif`
- write access to a local `.cache/` directory for AVIF generation

ServeBin reads templates and static files through relative paths, so run it from the repository root.

## Quick Start

1. Copy the example environment file:

    ```bash
    cp .env.example .env
    ```

2. Adjust values if needed.

3. Start the server from the repository root:

    ```bash
    go run ./cmd/ServeBin
    ```

4. Open the app:
    - landing page: `http://127.0.0.1:8888/`
    - Swagger UI: `http://127.0.0.1:8888/docs`

If `.env` does not exist, the binary falls back to already-exported environment variables and then to built-in defaults.

## Configuration

ServeBin uses environment variables for runtime configuration.

| Variable           | Purpose                | Effective behavior                                                                                   |
| ------------------ | ---------------------- | ---------------------------------------------------------------------------------------------------- |
| `HOST`             | Listen host            | Falls back to `127.0.0.1` when unset. `.env.example` uses `0.0.0.0`.                                 |
| `PORT`             | Listen port            | Falls back to `8888` when unset.                                                                     |
| `ENV`              | General runtime mode   | If `ENV != production`, ServeBin attempts to open the app URL in your default browser after startup. |
| `GIN_MODE`         | Gin mode               | Pass `debug` or `release` through the usual Gin environment variable.                                |
| `IS_SSL`           | URL/scheme hint        | Used when building sitemap URLs and absolute redirects. It does not enable TLS by itself.            |
| `IS_BACKUP_SERVER` | Backup-node behavior   | When `true`, `/`, `/docs/*`, and `/sitemap.xml` redirect to `MAIN_SERVER`.                           |
| `MAIN_SERVER`      | Backup redirect target | Required if `IS_BACKUP_SERVER=true`.                                                                 |

### Important runtime notes

- ServeBin always starts with `http.Server.ListenAndServe()`. Setting `IS_SSL=true` changes generated URLs, not socket-level TLS. If you need HTTPS locally or in production, terminate TLS in front of the app.
- In headless environments, set `ENV=production` to avoid the browser auto-open path.
- CORS is permissive by default: `Access-Control-Allow-Origin: *` plus common request methods and headers.

## API Docs

- Live Swagger UI: `/docs`
- Repo specs: [`docs/swagger.yaml`](docs/swagger.yaml) and [`docs/swagger.json`](docs/swagger.json)

## Endpoint Overview

### App and meta routes

| Route              | What it does                                                                            |
| ------------------ | --------------------------------------------------------------------------------------- |
| `GET /`            | Landing page rendered from `templates/landing/index.html` unless backup mode is enabled |
| `GET /about`       | Returns version, server time, developer/contact details, and source URL                 |
| `GET /heartbeat`   | Returns host CPU, RAM, disk, and network-latency stats                                  |
| `GET /sitemap.xml` | Builds a sitemap from registered non-parameterized routes                               |
| `GET /favicon.ico` | Serves the favicon from `static/logo/favicon.ico`                                       |

`/heartbeat` does more than a simple health check: it inspects local machine stats and performs outbound latency checks against `https://ping.atishir.co`.

### Request inspection and method testing

| Route                     | What it does                                                                                  |
| ------------------------- | --------------------------------------------------------------------------------------------- |
| `GET /ip`                 | Returns the client IP                                                                         |
| `GET /uuid`               | Returns a generated UUIDv4                                                                    |
| `GET /headers`            | Returns request headers                                                                       |
| `GET /user-agent`         | Returns the request `User-Agent`                                                              |
| `GET /status`             | Returns `200 OK`                                                                              |
| `ANY /status/:statuscode` | Returns the requested status code                                                             |
| `GET /get`                | Echoes query args, headers, origin, URL, and method                                           |
| `POST /post`              | Echoes args, multipart fields, files, raw body, parsed JSON, headers, origin, URL, and method |
| `PUT /put`                | Same body echo behavior as `/post`                                                            |
| `PATCH /patch`            | Same body echo behavior as `/post`                                                            |
| `DELETE /delete`          | Same body echo behavior as `/post`                                                            |
| `HEAD /head`              | Copies incoming request headers into the response headers and adds origin/URL/method headers  |

Multipart parsing is initialized with a `10 MB` memory limit.

### Response-format routes

| Route             | What it does                                  |
| ----------------- | --------------------------------------------- |
| `GET /json`       | Returns a sample JSON document                |
| `GET /xml`        | Returns a sample XML document                 |
| `GET /html`       | Returns a sample HTML document                |
| `GET /deny`       | Returns the sample `deny.txt` file            |
| `GET /robots.txt` | Returns sample crawler rules                  |
| `GET /gzip`       | Returns JSON with `Content-Encoding: gzip`    |
| `GET /brotli`     | Returns JSON with `Content-Encoding: br`      |
| `GET /deflate`    | Returns JSON with `Content-Encoding: deflate` |
| `GET /zstd`       | Returns JSON with `Content-Encoding: zstd`    |

### Image routes

| Route                   | What it does                                                                                     |
| ----------------------- | ------------------------------------------------------------------------------------------------ |
| `GET /image`            | Chooses an image format from the request `Accept` header and falls back to PNG for unknown types |
| `GET /image/:imagetype` | Forces a specific image type                                                                     |

Supported `imagetype` values:

- `png`
- `jpeg`
- `svg`
- `gif`
- `webp`
- `tiff`
- `bmp`
- `apng`
- `avif`
- `ico`

Notes:

- `webp` generation uses the external `cwebp` encoder through `go-webpbin`
- `avif` generation writes temporary files under `.cache/` and shells out to `avifenc`
- unknown explicit image types return `406 Not Acceptable`

### HTTPBin-compatible routes

| Route                                            | What it does                                                            |
| ------------------------------------------------ | ----------------------------------------------------------------------- |
| `ANY /anything` and `ANY /anything/*anything`    | Echoes request data in `httpbin`-style format                           |
| `GET /base64/*value`                             | Base64-decodes a path value                                             |
| `GET /encoding/utf8`                             | Returns a UTF-8 sample page                                             |
| `GET /response-headers`                          | Copies query parameters into response headers and echoes them as JSON   |
| `GET /redirect/:n`                               | Performs chained redirects                                              |
| `GET /redirect-to`                               | Redirects to `?url=...` with optional `?status_code=`                   |
| `GET /relative-redirect/:n`                      | Performs chained relative redirects                                     |
| `GET /absolute-redirect/:n`                      | Performs chained absolute redirects                                     |
| `GET /cookies`                                   | Returns request cookies                                                 |
| `GET /cookies/set`                               | Sets cookies from query parameters, then redirects to `/cookies`        |
| `GET /cookies/delete`                            | Deletes cookies named in the query string, then redirects to `/cookies` |
| `GET /basic-auth/:user/:passwd`                  | Basic-auth challenge                                                    |
| `GET /hidden-basic-auth/:user/:passwd`           | Returns `404` when auth fails instead of `401`                          |
| `GET /digest-auth/:qop/:user/:passwd`            | Digest-auth challenge                                                   |
| `GET /digest-auth/:qop/:user/:passwd/:algorithm` | Digest-auth challenge with `MD5` or `SHA-256`                           |
| `GET /stream/:n`                                 | Streams JSON lines                                                      |
| `ANY /delay/:n`                                  | Waits before responding                                                 |
| `GET /drip`                                      | Drips bytes over time                                                   |
| `GET /range/:n`                                  | Returns a range-enabled byte stream                                     |
| `GET /cache`                                     | Returns a cacheable response with `ETag` and `Last-Modified`            |
| `GET /etag/:etag`                                | Implements simple `ETag` semantics                                      |
| `GET /cache/:n`                                  | Sets `Cache-Control: public, max-age=n`                                 |
| `GET /bytes/:n`                                  | Returns random bytes                                                    |
| `GET /stream-bytes/:n`                           | Streams random bytes in chunks                                          |
| `GET /links/:n` and `GET /links/:n/:offset`      | Returns an HTML page of links                                           |
| `GET /forms/post`                                | Returns a sample HTML form that submits to `/post`                      |

Behavior limits implemented in code:

- `/stream/:n` is capped at `100` lines
- `/delay/:n` is capped at `10` seconds
- `/bytes/:n`, `/stream-bytes/:n`, and `/range/:n` are capped at `102400` bytes
- `/links/:n` is capped at `256` links
- `/range/:n` supports a single `Range` header value, not multiple ranges
- `/drip` accepts `numbytes`, `duration`, `delay`, and `code` query parameters

### OPTIONS behavior

The router registers `OPTIONS /options`, but the global CORS middleware short-circuits all `OPTIONS` requests with `204 No Content`. In practice, `OPTIONS` is handled as a CORS preflight response rather than by the echo controller.

## Example Requests

```bash
# Inspect request metadata
curl http://127.0.0.1:8888/get?demo=servebin

# Send JSON and have it echoed back
curl -X POST http://127.0.0.1:8888/post \
  -H 'Content-Type: application/json' \
  -d '{"hello":"world"}'

# Force a specific status code
curl -i http://127.0.0.1:8888/status/418

# Ask /image to negotiate the response type from Accept
curl -H 'Accept: image/webp' http://127.0.0.1:8888/image --output image.webp

# Set and inspect cookies
curl -i 'http://127.0.0.1:8888/cookies/set?theme=dark'
curl http://127.0.0.1:8888/cookies

# Test redirects
curl -i http://127.0.0.1:8888/redirect/3

# Exercise byte-range support
curl -i http://127.0.0.1:8888/range/64 -H 'Range: bytes=0-15'
```

## Project Structure

| Path                   | Purpose                                                                                 |
| ---------------------- | --------------------------------------------------------------------------------------- |
| `cmd/ServeBin/main.go` | Application entrypoint                                                                  |
| `router/`              | Route registration                                                                      |
| `controller/`          | HTTP handlers                                                                           |
| `service/`             | Request parsing and image-generation logic                                              |
| `helper/`              | Shared utilities for server startup, formatting, heartbeat stats, and `httpbin` helpers |
| `data/request/`        | Request models                                                                          |
| `data/response/`       | Response models                                                                         |
| `templates/`           | Landing page and sample payload templates                                               |
| `static/`              | Favicon and sample image assets                                                         |
| `docs/`                | Swagger/OpenAPI artifacts                                                               |

## License

ServeBin is released under the BSD 3-Clause license. See [`LICENSE`](LICENSE).
