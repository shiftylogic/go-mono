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
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"shiftylogic.dev/site-plat/internal/helpers"
)

const (
	kTokenSize = 16
)

func makeQRGenerator(prefix, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		token, err := helpers.GenerateStringSecure(kTokenSize, helpers.AlphaNumeric)
		if err != nil {
			log.Printf("[Error] Failed to generate secure random token - %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		hm := hmac.New(sha256.New, []byte(secret))
		hm.Write([]byte(ts))
		hm.Write([]byte(token))
		hash := hm.Sum(nil)

		code, err := qrcode.New(fmt.Sprintf("%s?ts=%s&v=%s&h=%x", prefix, ts, token, hash), qrcode.Low)
		if err != nil {
			log.Printf("[Error] Failed to generate QR Code - %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		png, err := code.PNG(512)
		if err != nil {
			log.Printf("[Error] Failed to generate QR Code PNG - %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(png)
	}
}
