package main

import (
	"fmt"
	"net/http"
	"strings"

	gw "github.com/mowazzem/tote/agw/internal"
	"github.com/mowazzem/tote/pkg/router"
)

func main() {
	r := router.New()
	r.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/auth/") {
			gw.AuthProxy(w, r)
		}
	})
	fmt.Println("Gateway running")
	http.ListenAndServe(":8760", r)

}
