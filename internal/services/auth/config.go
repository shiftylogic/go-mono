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

package auth

import (
	"time"

	"shiftylogic.dev/site-plat/internal/helpers"
)

const (
	kDefaultTokenTTL = 15 * time.Minute
	kMinimumTokenTTL = 5 * time.Minute
)

type Config struct {
	Endpoint string

	// OAuth2.0 Configuration
	Secret   string
	TokenTTL time.Duration

	// Features
	EnableQRScan bool
	QRScanPrefix string
}

func (cfg *Config) validate() error {
	// If a secret wasn't set, generate one
	if cfg.Secret == "" {
		secret, err := helpers.GenerateStringSecure(32, helpers.AlphaNumeric)
		if err != nil {
			return err
		}
		cfg.Secret = secret
	}

	// Make sure we have a TTL that makes sense
	if cfg.TokenTTL < kMinimumTokenTTL {
		cfg.TokenTTL = kMinimumTokenTTL
	}

	return nil
}
