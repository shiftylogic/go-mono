// MIT License
//
// Copyright (c) 2023 Robert Anderson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"shiftylogic.dev/site-plat/internal/helpers"
	"shiftylogic.dev/site-plat/internal/services"
	"shiftylogic.dev/site-plat/internal/services/auth"
	"shiftylogic.dev/site-plat/internal/web"
)

const (
	kCorsOriginsEnvKey = "CORS_ORIGINS"
)

var (
	secret string
)

func init() {
	s, err := helpers.GenerateStringSecure(64, helpers.AlphaNumeric)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate secure string - %v", err))
	}

	secret = s
}

func main() {
	origins := strings.Split(helpers.ReadEnvWithDefault(kCorsOriginsEnvKey, "*"), ";")
	log.Printf("Launching with CORS origins: %v", origins)

	router := web.NewRouter(
		web.WithCors(web.CorsOptions{
			AllowedOrigins: origins,
			AllowedMethods: []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTION"},
			AllowedHeaders: []string{
				"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language",
				"Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer",
			},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
		web.WithLogging(),
		web.WithPanicRecovery(),
		web.WithNoCache(),
		web.WithNoIFrame(),

		// Must be after middleware
		web.WithProfiler(),

		// Authentication services
		auth.WithOAuth2(auth.Config{
			Endpoint:     "/auth",
			Secret:       secret,
			EnableQRScan: true,
			QRScanPrefix: "https://local.vroov.com:9443/do-a-thing",
		}),

		// Serve up all the "common" files
		web.WithStaticFiles("/", os.DirFS("./static")),
	)

	router.Get("/do-a-thing", func(w http.ResponseWriter, r *http.Request) {
		ts := r.URL.Query().Get("ts")
		token := r.URL.Query().Get("v")
		h := r.URL.Query().Get("h")

		tsSecs, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			log.Printf("[Error] Timestamp parse failure - %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		hm := hmac.New(sha256.New, []byte(secret))
		hm.Write([]byte(ts))
		hm.Write([]byte(token))
		computedH := hm.Sum(nil)

		expectedH, err := hex.DecodeString(h)
		if err != nil {
			log.Printf("[Error] Timestamp parse failure - %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "QR Code Scanned:\n")
		fmt.Fprintf(w, "  => Timestamp: %v (%s)\n", time.Unix(tsSecs, 0), ts)
		fmt.Fprintf(w, "  => Token: %s\n", token)
		fmt.Fprintf(w, "  => Hash: %s\n", h)
		fmt.Fprintf(w, "  => Verified? %v\n", hmac.Equal(computedH, expectedH))
	})

	services.Start(router)
}
