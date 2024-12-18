package main

import (
	"fmt"
	"net/http"

	auth "github.com/mowazzem/tote/auth/internal"
	"github.com/mowazzem/tote/pkg/router"
)

func main() {
	r := router.New()
	gac := auth.NewGithubAuthClient()

	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Path))
	})

	r.HandleFunc("GET /auth/github", gac.GithubLogin)
	r.HandleFunc("GET /auth/github/callback", gac.GithubCallback)
	fmt.Println("auth server running")
	http.ListenAndServe(":8765", r)
}
