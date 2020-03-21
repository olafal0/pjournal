package main

import (
	"log"
	"net/http"

	"github.com/olafal0/dispatch"
	"github.com/olafal0/dispatch/auth"
	"github.com/olafal0/dispatch/kvstore"
)

var db *kvstore.KeyValueDB
var signer *auth.TokenSigner
var loginManager *auth.LoginManager

func init() {
	var err error
	db, err = kvstore.NewDB("pjournal.db")
	if err != nil {
		log.Fatal(err)
	}
	signer = auth.NewTokenSigner("pjournal", []byte("secretly hidden word of wisdom")) // no, this key isn't in use
	loginManager = &auth.LoginManager{
		DB:    db,
		Token: signer,
	}
}

func main() {
	authHook := dispatch.AuthorizerHook(signer)
	api := &dispatch.API{}
	api.AddEndpoint("POST/register", loginManager.SignupUser)
	api.AddEndpoint("POST/login", loginManager.AuthenticateUser)
	api.AddEndpoint("POST/posts/new", createPost, authHook)
	api.AddEndpoint("GET/post/{id}", getPost, authHook)
	api.AddEndpoint("DELETE/post/{id}", deletePost, authHook)
	api.AddEndpoint("GET/posts/all", listPosts, authHook)
	http.HandleFunc("/", api.GetHandler())
	log.Fatal(http.ListenAndServe(":8000", nil))
}
