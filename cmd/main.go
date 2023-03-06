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
	"log"
	"net/http"
	"os"
	"path"

	"shiftylogic.dev/site-plat/internal/web"
)

const (
	kServerAddr   = ":9443"
	kCertRootPath = "./.cert"
	kCertEnvKey   = "TLS_CERT_DOMAIN"
)

func main() {
	tlsCert, tlsKey := getTLSFiles()

	log.Printf("Launching server on port %s\n", kServerAddr)

	web.StartWithOptions(
		web.WithAddress(kServerAddr),
		web.WithHandler(buildRoutes()),
		web.WithTLS(tlsCert, tlsKey),
	)
}

func getTLSFiles() (string, string) {
	domain, ok := os.LookupEnv(kCertEnvKey)
	if !ok {
		log.Fatalf("TLS certificate domain environment ('%s') not set.\n", kCertEnvKey)
	}

	certDir := path.Join(kCertRootPath, "config", "live", domain)
	tlsCert := path.Join(certDir, "fullchain.pem")
	tlsKey := path.Join(certDir, "privkey.pem")

	return tlsCert, tlsKey
}

func buildRoutes() web.Router {
	r := web.NewRouter(
		web.WithLogging(),
		web.WithProfiler(),
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello\n"))
	})

	return r
}
