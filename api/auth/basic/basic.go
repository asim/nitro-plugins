// Package basic provides basic auth using htpasswd
package basic

import (
	"fmt"
	"net/http"

	"github.com/abbot/go-http-auth"
)

type Basic struct {
	File   string
	Realm  string
	Secret auth.SecretProvider
}

func (b *Basic) requireAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("WWW-Authenticate", `Basic realm="`+b.Realm+`"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf("%d %s\n", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))))
}

func (b *Basic) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			b.requireAuth(w, r)
			return
		}
		secret := b.Secret(u, b.Realm)
		// no secret
		if len(secret) == 0 {
			b.requireAuth(w, r)
			return
		}

		if !auth.CheckSecret(p, secret) {
			b.requireAuth(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func New(file, realm string) *Basic {
	return &Basic{
		File:   file,
		Realm:  realm,
		Secret: auth.HtpasswdFileProvider(file),
	}
}
