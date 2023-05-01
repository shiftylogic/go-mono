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
	"crypto/subtle"
	"errors"
	"log"
	"net/url"
)

const (
	kClientID    = "24C4853F-9398-41BD-B155-3333F181066B"
	kRedirectURI = "https://local.vroov.com:9443/auth-cb"

	kUser = "dude@intfoo.com"
	kPwd  = "1234test"
)

var (
	kBadUserPasswordError = errors.New("invalid user or password")
)

type fixedAuthorizer struct{}

func (v *fixedAuthorizer) ValidateClient(cid, redir string) bool {
	redir, err := url.QueryUnescape(redir)
	if err != nil {
		log.Printf("Failed to unescape redirect URI - %v", err)
		return false
	}

	res := subtle.ConstantTimeCompare([]byte(kClientID), []byte(cid))
	res += subtle.ConstantTimeCompare([]byte(kRedirectURI), []byte(redir))

	return res == 2
}

func (v *fixedAuthorizer) ValidateUser(user, pwd, scope string) (string, error) {
	res := subtle.ConstantTimeCompare([]byte(kUser), []byte(user))
	res += subtle.ConstantTimeCompare([]byte(kPwd), []byte(pwd))

	if res != 2 {
		return "", kBadUserPasswordError
	}

	return "1", nil
}
