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

	"shiftylogic.dev/site-plat/internal/services"
	"shiftylogic.dev/site-plat/internal/services/auth"
	"shiftylogic.dev/site-plat/internal/web"
)

func loadServices(ctx context.Context) services.Services {
	return &services.ServicesImpl{
		KVStore: services.NewMemoryStore(ctx),
		Authy:   &fixedAuthorizer{},
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

func getRoutes(config ServicesConfig) []web.RouterOptionFunc {
	return []web.RouterOptionFunc{
		auth.WithOAuth2(config.Auth),
	}
}
