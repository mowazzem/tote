package internal

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const gh_client_secret = "f7cb7623a37a9d0cb8343f18e8c166c3cb6eb446"
const gh_client_id = "Ov23liDuUep8YZXEPFQ4"

type githubAuth struct {
	conf *oauth2.Config
}

func NewGithubAuthClient() *githubAuth {
	return &githubAuth{
		conf: &oauth2.Config{
			ClientID:     gh_client_id,
			ClientSecret: gh_client_secret,
			Endpoint:     github.Endpoint,
		},
	}
}

func generateRandomString() (string, error) {
	const ln = 32
	rb := make([]byte, ln)
	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(rb), nil
}

type stateObj struct {
	verifier    string
	redirectUrl string
}

func (ga *githubAuth) GithubLogin(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	redirectUrl := r.URL.Query().Get("redirectUrl")
	if redirectUrl == "" {
		http.Error(w, "invalid redirectUrl", http.StatusUnauthorized)
		return
	}
	redirectUrl = redirectUrl[1 : len(redirectUrl)-1] // removing the double quotes
	verifier := oauth2.GenerateVerifier()
	state, err := generateRandomString()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	stateObj := stateObj{
		verifier:    verifier,
		redirectUrl: redirectUrl,
	}
	cacheClient.Set(state, stateObj)
	expireAfter := time.After(time.Minute * 3) // delete the generated code after 3 minutes
	go func(after <-chan time.Time, k string) {
		select {
		case <-after:
			cacheClient.Delete(k)
		}
	}(expireAfter, state)
	url := ga.conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce, oauth2.S256ChallengeOption(verifier))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (ga *githubAuth) GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	cacheState := cacheClient.Get(state)
	stateObj, ok := cacheState.(stateObj)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "invalid auth state")
		return
	}
	token, err := ga.conf.Exchange(r.Context(), code, oauth2.VerifierOption(stateObj.verifier))
	if err != nil {
		fmt.Fprintln(w, "invalid auth code")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, stateObj.redirectUrl+"?access_token="+token.AccessToken, http.StatusPermanentRedirect)
}
