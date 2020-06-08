package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"server/dbinterface"

	"github.com/aws/aws-lambda-go/events"
	"github.com/flick-web/dispatch"
)

var errorMissingID = errors.New("no id specified")
var errorInternal = errors.New(http.StatusText(http.StatusInternalServerError))

func ctxToUserID(ctx events.APIGatewayProxyRequestContext) (id string, err error) {
	// AWS wraps claims.sub inside of an interface{}, inside of a map[string]interface{},
	// inside of another map[string]interface{}. This means there are four different
	// points at which this function might panic if we don't explicitly error check.
	// This is ridiculous, so just assume everything will go fine and recover in
	// case of panic.
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ctxToUserID panic: %v\n", r)
			err = errorInternal
			return
		}
	}()
	claims := ctx.Authorizer["claims"]
	claimsMap := claims.(map[string]interface{})
	userID := claimsMap["sub"]
	id = userID.(string)
	return id, nil
}

func createPost(post *dbinterface.JPost, ctx *dispatch.Context) (id string, err error) {
	userID, err := ctxToUserID(ctx.LambdaRequest.RequestContext)
	if err != nil {
		return "", err
	}
	post.UserID = userID
	post.CreatedAt = time.Now().UnixNano() / 1e6
	return db.CreatePost(*post)
}

type postsList struct {
	Posts         []dbinterface.JPost `json:"posts"`
	LastEvaluated *dbinterface.JPost  `json:"lastEvaluated,omitempty"`
}

func getPosts(ctx *dispatch.Context) (*postsList, error) {
	// Get number of posts to limit query to
	var nPosts int64
	nPostsStr, ok := ctx.LambdaRequest.QueryStringParameters["n"]
	if ok {
		var err error
		nPosts, err = strconv.ParseInt(nPostsStr, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	// Get user ID from authorization context
	userID, err := ctxToUserID(ctx.LambdaRequest.RequestContext)
	if err != nil {
		return nil, err
	}

	// Get fromId and fromTime, parameters that are used to determine the last
	// evaluated item for pagination purposes
	var startFromPost *dbinterface.JPost
	startFromID, idOK := ctx.LambdaRequest.QueryStringParameters["fromId"]
	startFromTimestampStr, timeOK := ctx.LambdaRequest.QueryStringParameters["fromTime"]
	if idOK && timeOK {
		startFromTimestamp, err := strconv.ParseInt(startFromTimestampStr, 10, 64)
		if err != nil {
			fmt.Printf("Error parsing fromTime query parameter: %v\n", err)
		} else {
			startFromPost = &dbinterface.JPost{
				ID:        startFromID,
				UserID:    userID,
				CreatedAt: startFromTimestamp,
			}
		}
	}

	posts, lastEval, err := db.GetLastPosts(userID, nPosts, startFromPost)
	if err != nil {
		return nil, err
	}
	postsResult := postsList{
		Posts:         posts,
		LastEvaluated: lastEval,
	}
	return &postsResult, nil
}

func updatePost(post dbinterface.JPost, ctx *dispatch.Context) error {
	id, ok := ctx.PathVars["id"]
	if !ok || id == "undefined" || id == "null" {
		return errorMissingID
	}
	post.ID = id
	userID, err := ctxToUserID(ctx.LambdaRequest.RequestContext)
	if err != nil {
		return err
	}
	post.UserID = userID
	post.UpdatedAt = time.Now().UnixNano() / 1e6
	return db.UpdatePost(post)
}

func deletePost(ctx *dispatch.Context) error {
	id, ok := ctx.PathVars["id"]
	if !ok {
		return errorMissingID
	}
	userID, err := ctxToUserID(ctx.LambdaRequest.RequestContext)
	if err != nil {
		return err
	}
	return db.DeletePost(id, userID)
}
