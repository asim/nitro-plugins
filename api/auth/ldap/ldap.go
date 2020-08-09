// Package ldap provides ldap authentication
package ldap

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/ldap.v3"
)

type LDAP struct {
	URL    string
	Realm  string
	BaseDN string
}

func (l *LDAP) requireAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("WWW-Authenticate", `Basic realm="`+l.Realm+`"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf("%d %s\n", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))))
}

func (l *LDAP) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no url specified
		if len(l.URL) == 0 {
			h.ServeHTTP(w, r)
			return
		}

		// connect to ldap server
		c, err := ldap.DialURL(l.URL)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer c.Close()

		// get basic auth
		u, p, ok := r.BasicAuth()
		if !ok {
			l.requireAuth(w, r)
			return
		}

		// cn=user,ou=bla,dn=foo
		user := strings.Join([]string{fmt.Sprintf("cn=%s", u), l.BaseDN}, ",")

		// try bind to ldap server
		if err := c.Bind(user, p); err != nil {
			l.requireAuth(w, r)
			return
		}

		// serve http
		h.ServeHTTP(w, r)
	})
}

func New(uri, realm string) *LDAP {
	var baseDN string
	u, _ := url.Parse(uri)
	if u != nil && len(u.Path) > 1 && u.Path[0] == '/' {
		baseDN = u.Path[1:]
	}

	return &LDAP{
		URL:    uri,
		Realm:  realm,
		BaseDN: baseDN,
	}
}
