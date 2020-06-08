package main

import (
	"server/config"
	"server/dbinterface"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/flick-web/dispatch"
)

var db dbinterface.DatabaseModifier
var conf *config.PjournalConfig

func main() {
	conf = config.NewFromEnv()
	db = dbinterface.NewModifier(conf)

	api := &dispatch.API{}
	api.AddEndpoint("POST/posts/new", createPost)
	api.AddEndpoint("GET/posts", getPosts)
	api.AddEndpoint("POST/post/{id}", updatePost)
	api.AddEndpoint("DELETE/post/{id}", deletePost)

	lambda.Start(api.LambdaProxy)
}
