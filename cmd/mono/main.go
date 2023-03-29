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
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"shiftylogic.dev/site-plat/internal/helpers"
	"shiftylogic.dev/site-plat/internal/services"
	"shiftylogic.dev/site-plat/internal/services/auth"
	"shiftylogic.dev/site-plat/internal/web"
)

const (
	kConfigFileEnvKey = "SL_MONO_CONFIG"
)

type ServicesConfig struct {
	Auth auth.Config `json:"auth" yaml:"Auth"`
}

type MonoConfig struct {
	Base     services.Config `json:"root" yaml:"Root"`
	Services ServicesConfig  `json:"services" yaml:"Services"`
}

func loadConfig() MonoConfig {
	config := MonoConfig{
		services.DefaultConfig(),
		ServicesConfig{
			Auth: auth.DefaultConfig(),
		},
	}

	if configFile := helpers.ReadEnvWithDefault(kConfigFileEnvKey, ""); configFile != "" {
		services.LoadConfig(configFile, &config)
	}

	return config
}

func loadServices(ctx context.Context) services.Services {
	return &services.ServicesImpl{
		KVStore: services.NewMemoryStore(ctx),
	}
}

func selectMiddleware(config services.Config) []web.RouterOptionFunc {
	options := []web.RouterOptionFunc{
		web.WithLogging(),
		web.WithPanicRecovery(),
		web.WithNoIFrame(),
		web.WithNoCache(), // TODO: Remove this at some point later
	}

	if config.CORS.Enabled() {
		options = append(options, web.WithCors(config.CORS.Options()))
	}

	// This needs to be the last thing added (as middleware) before we start
	// adding other handlers
	if config.Profiler {
		options = append(options, web.WithProfiler())
	}

	return options
}

func loadRoutes(config ServicesConfig) []web.RouterOptionFunc {
	return []web.RouterOptionFunc{
		auth.WithOAuth2(config.Auth),
	}
}

func run() {
	ctx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	svcs := loadServices(ctx)
	config := loadConfig()

	options := selectMiddleware(config.Base)
	options = append(options, services.WithServices(svcs))
	options = append(options, loadRoutes(config.Services)...)
	options = append(
		options,
		web.WithStaticFiles("/", os.DirFS("./static")),
	)

	router := web.NewRouter(options...)

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

		if time.Now().Sub(time.Unix(tsSecs, 0)) > config.Services.Auth.TokenTTL {
			log.Printf("[Error] Token expired")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		svcs := services.ServicesFromContext(r.Context())
		if svcs == nil {
			log.Print("[Error] Failed to acquire active services")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		key, err := svcs.Cache().Get("qrscan", token)
		if err != nil {
			log.Printf("[Error] Failed to fetch key associated with token (%s) - %v", token, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		hm := hmac.New(sha256.New, []byte(key.(string)))
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

	services.Start(config.Base, router)
}

func main() {
	run()
	time.Sleep(250 * time.Millisecond)
	log.Print("Bye for realz!")
}
