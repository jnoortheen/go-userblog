package actions_test

import (
	"net/http"
	"muserblog/models"
	"fmt"
	"strconv"
	"github.com/satori/go.uuid"
)

var (
	likePostUrl = "/like"
)

func (as *ActionSuite) countLikes() int {
	return as.Count(models.Like{})
}
func (as *ActionSuite) Test_LikePost() {
	//data for test
	_, post := createPost(as)
	prevCount := as.countLikes()

	// without login it has to return error
	res := as.HTML(likePostUrl).Post(map[string]string{"post_id": post.ID.String()})
	as.Equal(http.StatusUnauthorized, res.Code)

	// login a user
	user := userForTest()
	signinUser(as, user)

	// like post
	res = as.HTML(likePostUrl).Post(map[string]string{"post_id": post.ID.String()})
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "count")
	as.Contains(res.Body.String(), strconv.Itoa(as.countLikes()))
	// check that it create a like record
	as.Equal(1, as.countLikes()-prevCount)
}

func (as *ActionSuite) Test_UnLikePost() {
	//data for test
	user, post := createPost(as)
	// create  a 'like' record
	like := &models.Like{UserID: user.ID, PostID: post.ID}
	as.DB.Create(like)
	prevCount := as.countLikes()
	// without login it has to return error
	res := as.HTML(likePostUrl).Post(map[string]string{"post_id": post.ID.String()})
	as.Equal(http.StatusUnauthorized, res.Code)

	// login user
	signinUser(as, user)

	// then dislike the post
	res = as.HTML(likePostUrl).Post(map[string]string{"post_id": post.ID.String()})
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "count")
	as.Contains(res.Body.String(), strconv.Itoa(as.countLikes()))
	as.Equal(1, prevCount-as.countLikes())
}
