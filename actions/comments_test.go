package actions_test

import (
	"muserblog/models"
	"encoding/json"
	"net/http"
)

func commentForTest(user *models.User, post *models.Post) *models.Comment {
	return &models.Comment{Content: postContent, UserID: user.ID, PostID: post.ID}
}

func (as *ActionSuite) countComments() int {
	return as.Count(models.Comment{})
}

func (as *ActionSuite) Test_CommentsResource_List() {
	user, post := createPost(as)
	contents := []string{"content 1", "content 2", "content 3"}
	for _, content := range contents {
		as.DB.Create(&models.Comment{Content: content, UserID: user.ID, PostID: post.ID})
	}
	resp := as.JSON("/comments").Get()
	comments := &models.Comments{}

	if err := json.NewDecoder(resp.Body).Decode(comments); err != nil {
		as.Fail("failed to decode")
	}
	as.Equal(len(*comments), len(contents))
}

func (as *ActionSuite) Test_CommentsResource_Create() {
	// initial variables
	user, post := createPost(as)
	comment := commentForTest(user, post)
	prevCount := as.countComments()

	// without signin
	resp := as.JSON("/comments").Post(comment)
	as.Equal(http.StatusUnauthorized, resp.Code)

	signinUser(as, user)

	// create post
	resp = as.JSON("/comments").Post(comment)
	as.Equal(201, resp.Code)
	as.Equal(1, as.countPosts()-prevCount)

	// check post exists in db
	as.DB.First(comment)
	user.WithName(as.DB)
	as.Equal(postContent, comment.Content)
	as.Equal(comment.UserID, user.ID)
	as.Equal(comment.PostID, post.ID)
}

func (as *ActionSuite) Test_CommentsResource_Update() {
	as.Fail("Not Implemented!")
}

//func (as *ActionSuite) Test_CommentsResource_Destroy() {
//	as.Fail("Not Implemented!")
//}
