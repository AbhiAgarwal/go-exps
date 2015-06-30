package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"
	"fmt"
	"io/ioutil"

  "golang.org/x/oauth2"
  "golang.org/x/oauth2/facebook"
	"github.com/gorilla/context"
)

var facebookconf = &oauth2.Config{
  ClientID:     "509424122546702",
  ClientSecret: "e69efc763ed78440c9566bd81568077e",
  RedirectURL:  "http://lvh.me:8080/facebook",
  Scopes: []string{},
  Endpoint: facebook.Endpoint,
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				WriteError(w, ErrInternalServer)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func acceptHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Header.Get("Accept") != "application/vnd.api+json" {
			WriteError(w, ErrNotAcceptable)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func contentTypeHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Header.Get("Content-Type") != "application/vnd.api+json" {
			WriteError(w, ErrUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func bodyHandler(v interface{}) func(http.Handler) http.Handler {
	t := reflect.TypeOf(v)

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			val := reflect.New(t).Interface()
			err := json.NewDecoder(r.Body).Decode(val)

			if err != nil {
				WriteError(w, ErrBadRequest)
				return
			}

			if next != nil {
				context.Set(r, "body", val)
				next.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(fn)
	}

	return m
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
  url := facebookconf.AuthCodeURL("state")
  http.Redirect(w, r, url, 301)
}

func facebookHandler(w http.ResponseWriter, r *http.Request) {
	if val, ok := r.URL.Query()["code"]; ok && len(val) > 0 {
		tok, err := facebookconf.Exchange(oauth2.NoContext, val[0])
		if err != nil {
        fmt.Println("err is", err)
    }
    fmt.Println("https://graph.facebook.com/v2.3/me?access_token=" + tok.AccessToken)
    response, err := http.Get("https://graph.facebook.com/v2.3/me?access_token=" + tok.AccessToken)
    defer response.Body.Close()
    contents, err := ioutil.ReadAll(response.Body)
    fmt.Print(string(contents))
	}
}
