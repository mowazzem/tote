package internal

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func AuthProxy(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse("http://localhost:8765/")
	if err != nil {
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	fmt.Println("request found")
	fmt.Println("request found"
	proxy.ServeHTTP(w, r)
}
