package actions_test

import (
	"muserblog/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (as *ActionSuite) countComments() int {
	return as.Count(models.Comment{})
}

func (as *ActionSuite) commentPath(post *models.Post) string {
	return fmt.Sprintf("/%s/comments", post.ID.String())
}
func (as *ActionSuite) commentUpdatePath(post *models.Post, comment *models.Comment) string {
	return fmt.Sprintf("/%s/comments/%s", post.ID.String(), comment.ID.String())
}

func (as *ActionSuite) Test_CommentsResource_List() {
	user, post := createPost(as)
	contents := []string{"content 1", "content 2", "content 3"}
	for _, content := range contents {
		as.DB.Create(&models.Comment{Content: content, UserID: user.ID, PostID: post.ID})
	}
	resp := as.JSON(as.commentPath(post)).Get()
	comments := &models.Comments{}

	if err := json.NewDecoder(resp.Body).Decode(comments); err != nil {
		as.Fail("failed to decode")
	}
	as.Equal(len(*comments), len(contents))
}

func (as *ActionSuite) Test_CommentsResource_Create() {
	// initial variables
	prevCount := as.countComments()
	user, post := createPost(as)
	content := "content"
	comment := &models.Comment{Content: content}

	// without signin
	resp := as.JSON(as.commentPath(post)).Post(comment)
	as.Equal(http.StatusUnauthorized, resp.Code)

	signinUser(as, user)

	// create post
	resp = as.JSON(as.commentPath(post)).Post(comment)
	as.Equal(201, resp.Code)
	as.Equal(1, as.countComments()-prevCount)

	// check post exists in db
	as.DB.First(comment)
	as.Equal(content, comment.Content)
	as.Equal(comment.UserID, user.ID)
	as.Equal(comment.PostID, post.ID)
}

func (as *ActionSuite) Test_CommentsResource_Update() {
	// initial variables
	user, post := createPost(as)
	content := "content"
	comment := &models.Comment{Content: content, UserID: user.ID, PostID: post.ID}
	as.DB.Create(comment)
	prevCount := as.countComments()

	// without signin
	resp := as.JSON(as.commentUpdatePath(post, comment)).Put(comment)
	as.Equal(http.StatusUnauthorized, resp.Code)

	//after signin
	signinUser(as, user)

	// update post
	comment.Content = "updated content"
	resp = as.JSON(as.commentUpdatePath(post, comment)).Put(comment)
	as.Equal(200, resp.Code)
	as.Equal(0, as.countComments()-prevCount)

	// get comment from db
	as.DB.Find(comment, comment.ID)
	as.Equal("updated content", comment.Content)
	as.Equal(comment.UserID, user.ID)
	as.Equal(comment.PostID, post.ID)
}

func (as *ActionSuite) Test_CommentsResource_Destroy() {
	as.Fail("Not Implemented!")
}
