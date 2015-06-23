// Basic authentication modified from source:
// https://github.com/goji/httpauth/blob/master/basic_auth.go

package main

import (
  "bytes"
  "crypto/sha256"
  "crypto/subtle"
  "encoding/base64"
  "fmt"
  "net/http"
  "strings"
)

type basicAuth struct {
  h    http.Handler
  opts AuthOptions
}

// AuthOptions stores the configuration for HTTP Basic Authentication.
//
// A http.Handler may also be passed to UnauthorizedHandler to override the
// default error handler if you wish to serve a custom template/response.
type AuthOptions struct {
  Realm               string
  User                string
  Password            string
  UnauthorizedHandler http.Handler
}

// Satisfies the http.Handler interface for basicAuth.
func (b basicAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // Check if we have a user-provided error handler, else set a default
  if b.opts.UnauthorizedHandler == nil {
    b.opts.UnauthorizedHandler = http.HandlerFunc(defaultUnauthorizedHandler)
  }

  // Check that the provided details match
  if b.authenticate(r) == false {
    b.requestAuth(w, r)
    return
  }

  // Call the next handler on success.
  b.h.ServeHTTP(w, r)
}

// authenticate retrieves and then validates the user:password combination provided in
// the request header. Returns 'false' if the user has not successfully authenticated.
func (b *basicAuth) authenticate(r *http.Request) bool {
  const basicScheme string = "Basic "

  // Confirm the request is sending Basic Authentication credentials.
  auth := r.Header.Get("Authorization")
  if !strings.HasPrefix(auth, basicScheme) {
    return false
  }

  // Get the plain-text username and password from the request
  // The first six characters are skipped - e.g. "Basic ".
  str, err := base64.StdEncoding.DecodeString(auth[len(basicScheme):])
  if err != nil {
    return false
  }

  // Split on the first ":" character only, with any subsequent colons assumed to be part
  // of the password. Note that the RFC2617 standard does not place any limitations on
  // allowable characters in the password.
  creds := bytes.SplitN(str, []byte(":"), 2)

  if len(creds) != 2 {
    return false
  }

  // Equalize lengths of supplied and required credentials
  // by hashing them
  givenUser := sha256.Sum256(creds[0])
  givenPass := sha256.Sum256(creds[1])
  requiredUser := sha256.Sum256([]byte(b.opts.User))
  requiredPass := sha256.Sum256([]byte(b.opts.Password))

  // Compare the supplied credentials to those set in our options
  if subtle.ConstantTimeCompare(givenUser[:], requiredUser[:]) == 1 &&
    subtle.ConstantTimeCompare(givenPass[:], requiredPass[:]) == 1 {
    return true
  }

  return false
}

// Require authentication, and serve our error handler otherwise.
func (b *basicAuth) requestAuth(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm=%q`, b.opts.Realm))
  b.opts.UnauthorizedHandler.ServeHTTP(w, r)
}

// defaultUnauthorizedHandler provides a default HTTP 401 Unauthorized response.
func defaultUnauthorizedHandler(w http.ResponseWriter, r *http.Request) {
  http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func BasicAuth(o AuthOptions) func(http.Handler) http.Handler {
  fn := func(h http.Handler) http.Handler {
    return basicAuth{h, o}
  }
  return fn
}