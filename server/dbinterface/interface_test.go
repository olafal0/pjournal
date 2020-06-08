package dbinterface_test

import (
	"os"
	"server/config"
	database "server/dbinterface"
	"testing"
	"time"
)

// Since these tests depend on the state of an AWS stack, they should not be run
// unless specified. Set TEST_DYNAMO=true to run these tests.
func TestMain(m *testing.M) {
	testEnv := os.Getenv("TEST_DYNAMO")
	if testEnv == "1" || testEnv == "true" {
		os.Exit(m.Run())
	} else {
		os.Exit(0)
	}
}

func TestPostsInterface(t *testing.T) {
	conf := &config.PjournalConfig{}
	conf.PostsTableName = "pjournal-development-PostsTable-198CAU4BK07B4"

	now := time.Now()

	post1 := database.JPost{
		UserID:    "userId",
		CreatedAt: now.UnixNano() / 1e6,
	}

	post2 := database.JPost{
		UserID:    "userId",
		CreatedAt: now.Add(time.Hour).UnixNano() / 1e6,
	}

	db := database.NewModifier(conf)
	postID1, err := db.CreatePost(post1)
	if err != nil {
		t.Error(err)
	}
	postID2, err := db.CreatePost(post2)
	if err != nil {
		t.Error(err)
	}

	post1.ID = postID1
	post2.ID = postID2

	resultPosts, lastEval, err := db.GetLastPosts("userId", 10, nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if lastEval != nil {
		t.Error("Unexpected lastEval value")
	}
	if resultPosts[0].ID != postID2 || resultPosts[1].ID != postID1 {
		t.Errorf(
			"Post IDS did not match: %v and %v (expected %v and %v)\n",
			resultPosts[0].ID, resultPosts[1].ID,
			postID2, postID1,
		)
	}

	// Use most recent post as start key, which should mean we only receive the
	// other, older post
	resultPosts, lastEval, err = db.GetLastPosts("userId", 10, &post2)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(resultPosts) != 1 {
		t.Errorf("Wrong number of result posts with start key: %v", len(resultPosts))
	}
	if resultPosts[0].ID != postID1 {
		t.Errorf("Post ID did not match: %v (expected %v)\n", resultPosts[0].ID, postID1)
	}

	err = db.DeletePost(postID1, "userId")
	if err != nil {
		t.Error(err)
	}

	err = db.DeletePost(postID2, "userId")
	if err != nil {
		t.Error(err)
	}
}
