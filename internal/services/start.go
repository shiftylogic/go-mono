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
	"log"
	"path"

	"shiftylogic.dev/site-plat/internal/helpers"
	"shiftylogic.dev/site-plat/internal/web"
)

const (
	kDefaultAddress   = ":9443"
	kAddrEnvKey       = "SERVER_ADDRESS"
	kCertRootEnvKey   = "TLS_CERT_ROOT"
	kCertDomainEnvKey = "TLS_CERT_DOMAIN"
)

func Start(router web.Router) {
	serverAddr := helpers.ReadEnvWithDefault(kAddrEnvKey, kDefaultAddress)
	tlsCert, tlsKey := loadTLS()

	log.Printf("Launching server on port %s\n", serverAddr)

	web.StartWithOptions(
		web.WithAddress(serverAddr),
		web.WithTLS(tlsCert, tlsKey),
		web.WithHandler(router),
	)

	log.Print("Server stopped.")
}

func loadTLS() (string, string) {
	rootPath := helpers.ReadEnv(kCertRootEnvKey)
	domain := helpers.ReadEnv(kCertDomainEnvKey)

	certDir := path.Join(rootPath, "config", "live", domain)
	tlsCert := path.Join(certDir, "fullchain.pem")
	tlsKey := path.Join(certDir, "privkey.pem")

	return tlsCert, tlsKey
}
