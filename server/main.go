package main

import (
	"log"
	"net/http"
	"os"

	"github.com/flick-web/auth"
	"github.com/flick-web/dispatch"
	"github.com/flick-web/kvstore"
)

var db kvstore.KeyValueStore
var signer *auth.TokenSigner
var loginManager *auth.LoginManager

var contentDir = os.Getenv("CONTENT_DIR")
var dbFile = os.Getenv("DB_FILE")

func init() {
	if dbFile == "" {
		dbFile = "pjournal.db"
	}

	if contentDir == "" {
		contentDir = "../web/dist"
	}

	log.Printf("Using CONTENT_DIR=%s DB_FILE=%s\n", contentDir, dbFile)

	var err error
	db, err = kvstore.NewSqliteDB(dbFile)
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
	api := &dispatch.API{}
	api.AddEndpoint("POST/api/register", loginManager.SignupUser)
	api.AddEndpoint("POST/api/login", loginManager.AuthenticateUser)
	api.AddEndpoint("GET/api/logout", loginManager.LogoutUser)
	api.AddEndpoint("POST/api/posts/new", createPost, loginManager.AuthorizerHook)
	api.AddEndpoint("GET/api/post/{id}", getPost, loginManager.AuthorizerHook)
	api.AddEndpoint("POST/api/post/{id}", updatePost, loginManager.AuthorizerHook)
	api.AddEndpoint("DELETE/api/post/{id}", deletePost, loginManager.AuthorizerHook)
	api.AddEndpoint("GET/api/posts/all", listPosts, loginManager.AuthorizerHook)
	http.HandleFunc("/api/", api.GetHandler())
	http.Handle("/", http.FileServer(http.Dir(contentDir)))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
