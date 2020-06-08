package dbinterface

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
)

// A PostsModifier is an object that handles database operations relating to
// journal posts.
type PostsModifier interface {
	CreatePost(post JPost) (string, error)
	GetLastPosts(userID string, n int64, startFrom *JPost) (posts []JPost, lastEvaluated *JPost, err error)
	UpdatePost(post JPost) error
	DeletePost(id, userID string) error
}

// JPost represents a journal post, storing its ID along with its content.
type JPost struct {
	ID     string `dynamo:"id" json:"id"`
	UserID string `dynamo:"userId" json:"userId"`
	// Content is the text content of the post, or, if the post was encrypted, a
	// string of unicode code points representing the binary encrypted data of the post.
	Content   string `dynamo:"content,omitempty" json:"content,omitempty"`
	CreatedAt int64  `dynamo:"createdAt" json:"createdAt"`
	UpdatedAt int64  `dynamo:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	Encrypted bool   `dynamo:"encrypted,omitempty" json:"encrypted,omitempty"`
	// IV is the initialization vector used when encrypting this post client-side.
	// If encryption was not enabled for this post, it will be empty.
	IV string `dynamo:"iv,omitempty" json:"iv,omitempty"`
}

// HashKey fulfills the dynamo.Keyed interface.
func (p *JPost) HashKey() interface{} {
	return p.ID
}

// RangeKey fulfills the dynamo.Keyed interface.
func (p *JPost) RangeKey() interface{} {
	return p.UserID
}

// ToDynamoPagingKey returns a dynamo.PagingKey, for use in query.StartFrom(),
// from the key information in a JPost object.
func (p JPost) ToDynamoPagingKey() dynamo.PagingKey {
	return dynamo.PagingKey{
		"id": &dynamodb.AttributeValue{
			S: &p.ID,
		},
		"userId": &dynamodb.AttributeValue{
			S: &p.UserID,
		},
		"createdAt": &dynamodb.AttributeValue{
			N: aws.String(fmt.Sprint(p.CreatedAt)),
		},
	}
}

// Posts is a PostModifier for use with a Dynamo backend.
type Posts struct {
	Dynamo                 DynamoHandler
	UsercreatedAtIndexName string
}

// CreatePost adds a new journal post. This should always be used to create Posts,
// even though it is just a wrapper around SetPost, as it is where IDs are
// created and set.
func (p *Posts) CreatePost(post JPost) (string, error) {
	post.ID = p.Dynamo.newID()
	return post.ID, p.UpdatePost(post)
}

// GetLastPosts lists the most recent n posts, with the specified starting key.
func (p *Posts) GetLastPosts(userID string, n int64, startFrom *JPost) ([]JPost, *JPost, error) {
	db, err := p.Dynamo.getClient()
	if err != nil {
		return nil, nil, err
	}

	index, ok := p.Dynamo.Indexes[p.UsercreatedAtIndexName]
	if !ok {
		return nil, nil, fmt.Errorf("Index %s not found in index configuration", p.UsercreatedAtIndexName)
	}

	result := []JPost{}

	table := db.Table(p.Dynamo.TableName)
	query := table.Get(index.PrimaryKey, userID).Index(p.UsercreatedAtIndexName)
	query = query.Order(dynamo.Descending).Limit(n)

	if startFrom != nil {
		query = query.StartFrom(startFrom.ToDynamoPagingKey())
	}

	lastEvaluated, err := query.AllWithLastEvaluatedKey(&result)
	if err != nil {
		return nil, nil, err
	}

	if lastEvaluated == nil {
		return result, nil, nil
	}

	// If we can get the full and correct key information from lastEvaluated,
	// return it as a post with fields set correctly. If there are any errors,
	// skip returning the last evaluated key.
	lastEvaluatedID, idOK := lastEvaluated["id"]
	lastEvaluatedUserID, uidOK := lastEvaluated["userId"]
	lastEvaluatedTimestamp, timeOK := lastEvaluated["createdAt"]

	if idOK && uidOK && timeOK {
		lastPost := JPost{}
		lastPost.ID = *lastEvaluatedID.S
		lastPost.UserID = *lastEvaluatedUserID.S
		lastPost.CreatedAt, err = strconv.ParseInt(*lastEvaluatedTimestamp.N, 10, 64)
		if err != nil {
			fmt.Printf("Error during conversion of lastEvaluatedTimestamp: %v\n", err)
		} else {
			return result, &lastPost, nil
		}
	}

	return result, nil, nil
}

// UpdatePost writes the given post to the database. Unlike CreatePost,
// it requires the passed argument to already have an ID set.
func (p *Posts) UpdatePost(post JPost) error {
	db, err := p.Dynamo.getClient()
	if err != nil {
		return err
	}

	// Confirm that the post object has an ID and UserID set
	if post.ID == "" {
		return errors.New("ID not set in SetPost request")
	}
	if post.UserID == "" {
		return errors.New("UserID not set in SetPost request")
	}

	table := db.Table(p.Dynamo.TableName)
	err = table.Put(post).Run()
	return err
}

// DeletePost removes a post specified by ID.
func (p *Posts) DeletePost(id, userID string) error {
	db, err := p.Dynamo.getClient()
	if err != nil {
		return err
	}

	table := db.Table(p.Dynamo.TableName)
	err = table.Delete(p.Dynamo.PrimaryKey, id).Range(p.Dynamo.SortKey, userID).Run()
	return err
}
