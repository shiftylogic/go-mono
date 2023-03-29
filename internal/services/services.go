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

package services

import (
	"context"
	"time"

	"shiftylogic.dev/site-plat/internal/web"
)

const (
	ServicesContextKey = "sl.services"
)

type KVStore interface {
	Get(ns, key string) (any, error)

	CheckAndSet(ns, key string, value any, ttl time.Duration) error
	Refresh(ns, key string, ttl time.Duration) error
	Set(ns, key string, value any, ttl time.Duration) error

	Remove(ns, key string)
}

type Services interface {
	Cache() KVStore
}

func ServicesFromContext(ctx context.Context) Services {
	return ctx.Value(ServicesContextKey).(Services)
}

func WithServices(svcs Services) web.RouterOptionFunc {
	return func(r web.Router) {
		r.Use(web.InjectContext(ServicesContextKey, svcs))
	}
}

/**
 *
 * A simple implementation of the Services interface that stores the various
 * service instances directly and just provides references as needed.
 *
 **/

type ServicesImpl struct {
	KVStore KVStore
}

func (svcs ServicesImpl) Cache() KVStore {
	return svcs.KVStore
}
