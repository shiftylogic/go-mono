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
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"shiftylogic.dev/site-plat/internal/helpers"
	"shiftylogic.dev/site-plat/internal/services"
	"shiftylogic.dev/site-plat/internal/web"
)

const (
	kAuthCodeSize        = 32
	kAuthCodeCacheKey    = "auth_code"
	kAuthRequestIDSize   = 20
	kAuthRequestCacheKey = "auth_req"

	kLoginTemplate = "login.html"

	kLoginRoute   = "/login"
	kQRImageRoute = "/qrcode"

	kAccessDeniedError       = "access_denied"
	kServerError             = "server_error"
	kUnsupportedResponseType = "unsupported_response_type"
)

func WithOAuth2(config Config) web.RouterOptionFunc {
	templates := template.Must(template.ParseFS(os.DirFS(config.Templates), "*.html"))

	return func(root web.Router) {
		r := web.NewRouter()

		r.With(web.NoIFrame).Get("/authorize", Authorize(templates, config))
		r.Post("/login", Login(config))
		// r.Post("/token", )

		if config.QRScan.Enabled {
			r.Get(kQRImageRoute, QRGenerator(config.QRScan))
			r.Get("/do-a-thing", DoThing(config.QRScan.TTL))
		}

		root.Mount(config.Path, r)
	}
}

type loginRequestData struct {
	ClientID    string
	RedirectURI string
	Scope       string
	State       string
}

type loginViewData struct {
	RequestID string
	Scope     string
	QREnabled bool
}

func Login(config Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rid := r.FormValue("rid")
		user := r.FormValue("user")
		pwd := r.FormValue("pwd")

		svcs := services.ServicesFromContext(r.Context())
		if svcs == nil {
			log.Print("[Error] Failed to acquire active services")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rdataI, err := svcs.Cache().Get(kAuthRequestCacheKey, rid)
		if err != nil {
			log.Printf("[Error] Failed to fetch request data associated with request id (%s) - %v", rid, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// We don't need / want the request data in the cache anymore
		svcs.Cache().Remove(kAuthRequestCacheKey, rid)

		rdata, ok := rdataI.(loginRequestData)
		if !ok {
			log.Print("[Error] Request data is invalid.")
			redirectAuthError(w, r, kServerError, rdata)
			return
		}

		uid, err := svcs.Authorizer().ValidateUser(user, pwd, rdata.Scope)
		if err != nil {
			log.Printf("[Error] Authentication failed - %v", err)
			redirectAuthError(w, r, kAccessDeniedError, rdata)
			return
		}

		code, err := helpers.GenerateStringSecure(kAuthCodeSize, helpers.AlphaNumeric)
		if err != nil {
			log.Printf("[Error] Failed to generate secure random QR secret - %v", err)
			redirectAuthError(w, r, kServerError, rdata)
			return
		}

		state := rdata.State
		rdata.State = uid

		if err := svcs.Cache().Set(kAuthCodeCacheKey, code, rdata, config.RequestTTL); err != nil {
			log.Printf("[Error] Failed to write request data to cache = %v", err)
			redirectAuthError(w, r, kServerError, rdata)
			return
		}

		http.Redirect(
			w,
			r,
			fmt.Sprintf("%s?code=%s&state=%s", rdata.RedirectURI, code, state),
			http.StatusFound,
		)
	}
}

func Authorize(templates *template.Template, config Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rdata := loginRequestData{
			ClientID:    r.URL.Query().Get("client_id"),
			RedirectURI: r.URL.Query().Get("redirect_uri"),
			Scope:       r.URL.Query().Get("scope"),
			State:       r.URL.Query().Get("state"),
		}

		svcs := services.ServicesFromContext(r.Context())
		if svcs == nil {
			log.Print("[Error] Failed to acquire active services")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !svcs.Authorizer().ValidateClient(rdata.ClientID, rdata.RedirectURI) {
			log.Print("[Error] Invalid client and / or redirect URL in authorize call.")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if rtype := r.URL.Query().Get("response_type"); rtype != "code" {
			log.Print("[Error] Unsupported response_type in authorization request.")
			redirectAuthError(w, r, kUnsupportedResponseType, rdata)
			return
		}

		rid, err := helpers.GenerateStringSecure(kAuthRequestIDSize, helpers.AlphaNumeric)
		if err != nil {
			log.Printf("[Error] Failed to generate secure random QR secret - %v", err)
			redirectAuthError(w, r, kServerError, rdata)
			return
		}

		if err := svcs.Cache().Set(kAuthRequestCacheKey, rid, rdata, config.RequestTTL); err != nil {
			log.Printf("[Error] Failed to write request data to cache = %v", err)
			redirectAuthError(w, r, kServerError, rdata)
			return
		}

		data := loginViewData{
			RequestID: rid,
			Scope:     rdata.Scope,
			QREnabled: config.QRScan.Enabled,
		}

		if err := templates.ExecuteTemplate(w, kLoginTemplate, data); err != nil {
			log.Printf("[Error] Failed to execute 'login' template - %v", err)
			redirectAuthError(w, r, kServerError, rdata)
			return
		}
	}
}

func redirectAuthError(w http.ResponseWriter, r *http.Request, errS string, d loginRequestData) {
	http.Redirect(
		w,
		r,
		fmt.Sprintf("%s?error=%s&state=%s", d.RedirectURI, errS, d.State),
		http.StatusFound,
	)
}

/**
 * TODO: Remove
 * This is just for playing with the QR Code generator
 **/
func DoThing(ttl time.Duration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ts := r.URL.Query().Get(kQRTimestampName)
		token := r.URL.Query().Get(kQRTokenName)
		h := r.URL.Query().Get(kQRHashName)

		tsSecs, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			log.Printf("[Error] Timestamp parse failure - %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if time.Now().Sub(time.Unix(tsSecs, 0)) > ttl {
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

		key, err := svcs.Cache().Get(kQRCacheKey, token)
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
	}
}
