package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/olafal0/dispatch"
	"github.com/sethvargo/go-diceware/diceware"
)

// JPost represents a journal post, storing its ID along with its content.
type JPost struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	// Content is the text content of the post, or, if the post was encrypted, a
	// string of unicode code points representing the binary encrypted data of the post.
	Content   string `json:"content,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`
	Encrypted bool   `json:"encrypted"`
	// IV is the initialization vector used when encrypting this post client-side.
	// If encryption was not enabled for this post, it will be empty.
	IV string `json:"iv,omitempty"`
}

// NewPost represents input from the client requesting a new post. This can also
// be used for updating existing posts.
type NewPost struct {
	Content   string `json:"content,omitempty"`
	Encrypted bool   `json:"encrypted"`
	IV        string `json:"iv,omitempty"`
}

func getPostIDs(username string) ([]string, error) {
	postIDs := make([]string, 0, 16)
	err := db.GetObject(fmt.Sprintf("posts_%s", username), "id_list", &postIDs)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return postIDs, nil
}

func createPost(newPost NewPost, ctx *dispatch.Context) (id string, err error) {
	// Generate random title-cased list of words to use as the ID
	// TODO: remove assumption that generated ID is unique
	username := ctx.Claims.Subject
	postIDs, err := getPostIDs(username)
	if err != nil {
		return "", err
	}

	words, _ := diceware.Generate(3)
	id = strings.ReplaceAll(strings.Title(strings.Join(words, " ")), " ", "")
	post := JPost{
		ID:        id,
		Content:   newPost.Content,
		CreatedAt: time.Now().UnixNano() / 1e6,
		Username:  ctx.Claims.Subject,
		Encrypted: newPost.Encrypted,
		IV:        newPost.IV,
	}
	err = db.SetObject(fmt.Sprintf("posts_%s", username), id, post)
	if err != nil {
		return "", err
	}
	postIDs = append(postIDs, id)
	err = db.SetObject(fmt.Sprintf("posts_%s", username), "id_list", postIDs)
	if err != nil {
		return "", err
	}
	return id, nil
}

func getPost(ctx *dispatch.Context) (*JPost, error) {
	username := ctx.Claims.Subject
	result := JPost{}
	err := db.GetObject(fmt.Sprintf("posts_%s", username), ctx.PathVars["id"], &result)
	return &result, err
}

func updatePost(newPost NewPost, ctx *dispatch.Context) (err error) {
	username := ctx.Claims.Subject
	postID := ctx.PathVars["id"]

	post := JPost{}
	err = db.GetObject(fmt.Sprintf("posts_%s", username), postID, &post)
	if err != nil {
		return err
	}

	post.Content = newPost.Content
	post.UpdatedAt = time.Now().UnixNano() / 1e6
	post.Encrypted = newPost.Encrypted
	post.IV = newPost.IV

	err = db.SetObject(fmt.Sprintf("posts_%s", username), postID, post)
	if err != nil {
		return err
	}
	return nil
}

func listPosts(ctx *dispatch.Context) ([]JPost, error) {
	username := ctx.Claims.Subject
	postIDs, err := getPostIDs(username)
	if err != nil {
		return nil, err
	}

	postsPerPage := 25
	posts := make([]JPost, 0, postsPerPage)

	// postIDs will be most recent first
	var searchSet []string
	// If we have less than postsPerPage, return all
	if len(postIDs) <= postsPerPage {
		searchSet = postIDs
	} else {
		// otherwise, return the last postsPerPage
		searchSet = postIDs[len(postIDs)-postsPerPage:]
	}
	for _, id := range searchSet {
		retrievedPost := JPost{}
		err := db.GetObject(fmt.Sprintf("posts_%s", username), id, &retrievedPost)
		if err != nil {
			return nil, err
		}
		posts = append(posts, retrievedPost)
	}
	return posts, nil
}

func deletePost(ctx *dispatch.Context) error {
	username := ctx.Claims.Subject
	postIDs, err := getPostIDs(username)
	if err != nil {
		return err
	}
	postID := ctx.PathVars["id"]
	err = db.DeleteObject(fmt.Sprintf("posts_%s", username), postID)
	if err != nil {
		return err
	}

	postIndex, err := func() (int, error) {
		for i, v := range postIDs {
			if v == postID {
				return i, nil
			}
		}
		return 0, errors.New("Post not found for ID")
	}()
	if err != nil {
		return err
	}

	postIDs = append(postIDs[:postIndex], postIDs[postIndex+1:]...)
	err = db.SetObject(fmt.Sprintf("posts_%s", username), "id_list", postIDs)
	return err
}
