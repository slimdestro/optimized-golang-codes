/**
	@ failSafe high performance proxy server with oAuth in Golang
	@ dlimdestro
*/
package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"log"
	"context"
	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
)

// ProxyServer struct
type ProxyServer struct {
	Proxy *httputil.ReverseProxy
	Config *ProxyConfig
}

// ProxyConfig struct
type ProxyConfig struct {
	Target *url.URL
	Scheme string
	Timeout time.Duration
	Auth bool
	OAuth2 bool
	MaxConcurrentRequests int
}

// NewProxyServer creates a new proxy server
func NewProxyServer(config *ProxyConfig) *ProxyServer {
	proxy := httputil.NewSingleHostReverseProxy(config.Target)
	proxy.Transport = &http.Transport{
		MaxConnsPerHost: config.MaxConcurrentRequests,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &ProxyServer{
		Proxy: proxy,
		Config: config,
	}
}

// ServeHTTP handles the incoming requests
func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), p.Config.Timeout)
	defer cancel()
	r = r.WithContext(ctx)
	if p.Config.Auth {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	if p.Config.OAuth2 {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	p.Proxy.ServeHTTP(w, r)
}

func main() {
	target, _ := url.Parse("http://example.com")
	proxyConfig := &ProxyConfig{
		Target: target,
		Scheme: "http",
		Timeout: 10 * time.Second,
		Auth: true,
		OAuth2: true,
		MaxConcurrentRequests: 1000000,
	}
	proxy := NewProxyServer(proxyConfig)
	r := mux.NewRouter()
	r.HandleFunc("/", proxy.ServeHTTP)
	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r)))
}