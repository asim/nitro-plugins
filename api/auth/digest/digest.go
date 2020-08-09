// Package digest provides digest auth using htdigest
package digest

import (
	"net/http"

	"github.com/abbot/go-http-auth"
)

type Digest struct {
	File    string
	Realm   string
	Secret  auth.SecretProvider
	Checker *auth.DigestAuth
}

func (d *Digest) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, info := d.Checker.CheckAuth(r)
		if len(u) == 0 {
			d.Checker.RequireAuth(w, r)
			return
		}

		if info != nil {
			w.Header().Set("Authentication-Info", *info)
		}
		h.ServeHTTP(w, r)
	})
}

func New(file, realm string) *Digest {
	secret := auth.HtdigestFileProvider(file)

	return &Digest{
		File:    file,
		Realm:   realm,
		Secret:  secret,
		Checker: auth.NewDigestAuthenticator(realm, secret),
	}
}
