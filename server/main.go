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
	authHook := auth.AuthorizerHook(signer)
	api := &dispatch.API{}
	api.AddEndpoint("POST/api/register", loginManager.SignupUser)
	api.AddEndpoint("POST/api/login", loginManager.AuthenticateUser)
	api.AddEndpoint("GET/api/logout", loginManager.LogoutUser)
	api.AddEndpoint("POST/api/posts/new", createPost, authHook)
	api.AddEndpoint("GET/api/post/{id}", getPost, authHook)
	api.AddEndpoint("POST/api/post/{id}", updatePost, authHook)
	api.AddEndpoint("DELETE/api/post/{id}", deletePost, authHook)
	api.AddEndpoint("GET/api/posts/all", listPosts, authHook)
	http.HandleFunc("/api/", api.GetHandler())
	http.Handle("/", http.FileServer(http.Dir("../web/dist/")))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
