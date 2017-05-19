package actions_test

import (
	"fmt"
	"muserblog/models"
	"net/http"
)

var (
	postTitle        = "Post 1"
	postContent      = "Post 1 content"
	updatedPostTitle = "updated title for post"

	postsListUrl = "/posts"
	postsNewUrl  = postsListUrl + "/new"
	postUrl      = postsListUrl + "/%s"
	postsEditUrl = postUrl + "/edit"
)

func postForTest() *models.Post {
	return &models.Post{Title: postTitle, Content: postContent}
}

func (as *ActionSuite) countPosts() int {
	return as.Count(models.Post{})
}

// create a new post and return
func createPost(as *ActionSuite) (*models.User, *models.Post) {
	user := userForTest()
	userCopy := *user
	user.SaltPassword()
	as.NoError(as.DB.Create(user))

	prevCount := as.countPosts()
	post := postForTest()
	post.UserID = user.ID
	as.NoError(as.DB.Create(post))

	// check exactly one new record is created
	as.Equal(as.countPosts()-prevCount, 1)
	userCopy.ID = user.ID
	return &userCopy, post
}

// return the First record that matches the query
func getPostFromDB(as *ActionSuite, post *models.Post) *models.Post {
	err := as.DB.First(post)
	as.NoError(err)
	as.NotZero(post.ID)
	as.NotZero(post.CreatedAt)
	return post
}

func user2ForTest(as *ActionSuite) *models.User {
	user := userForTest()
	user.Name = "test_2"

	copyUser := *user

	user.SaltPassword()
	as.NoError(as.DB.Create(user))

	return &copyUser
}

func (as *ActionSuite) Test_PostsResource_List() {
	res := as.HTML(postsListUrl).Get()

	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Posts")
}

func (as *ActionSuite) Test_PostsResource_Show() {
	// initially create a post
	user, post := createPost(as)
	url := fmt.Sprintf(postsEditUrl, post.ID)
	editLink := fmt.Sprintf(`href="%s"`, url)

	// case1: before logging in
	res := as.HTML(postUrl, post.ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), post.Content)
	as.NotContains(res.Body.String(), editLink)

	// signin with author
	signinUser(as, user)

	// case2: after signin post with edit button
	res = as.HTML(postUrl, post.ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), editLink)

	// case3: no edit option for no author
	as.HTML(signoutPath).Get()
	user2 := user2ForTest(as)
	signinUser(as, user2)
	res = as.HTML(postUrl, post.ID).Get()
	as.Equal(200, res.Code)
	as.NotContains(res.Body.String(), editLink)
}

func (as *ActionSuite) Test_PostsResource_New() {
	// without login it has to redirect to signin page
	res := as.HTML(postsNewUrl).Get()
	as.Equal(http.StatusFound, res.Code)

	// login a user
	user := userForTest()
	signinUser(as, user)

	// now it should return post new page
	res = as.HTML(postsNewUrl).Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "<h1>New Post</h1>")
}

func (as *ActionSuite) Test_PostsResource_Create() {
	// without login it has to redirect to signin page
	res := as.HTML(postsNewUrl).Get()
	as.Equal(http.StatusFound, res.Code)

	// initial
	prevCount := as.countPosts()
	post := postForTest()

	// login a user
	user := userForTest()
	signinUser(as, user)

	// create post
	res = as.HTML(postsListUrl).Post(post)
	as.Equal(302, res.Code)
	as.Equal(1, as.countPosts()-prevCount)

	// check post exists in db
	post = getPostFromDB(as, post)
	user.WithName(as.DB)
	as.Equal("Post 1", post.Title)
	as.Equal(post.UserID, user.ID)
}

func (as *ActionSuite) Test_PostsResource_Edit() {
	// initial
	user, post := createPost(as)

	// without login it has to redirect to signin page
	res := as.HTML(postsEditUrl, post.ID).Get()
	as.Equal(http.StatusFound, res.Code)

	// login as author
	signinUser(as, user)

	// get to edit post page
	res = as.HTML(postsEditUrl, post.ID).Get()
	as.Contains(res.Body.String(), "<h1>Edit Post</h1>")
	as.Contains(res.Body.String(), post.Title)
	as.Contains(res.Body.String(), post.Content)
}

func (as *ActionSuite) Test_PostsResource_Update() {
	//initial
	user, post := createPost(as)
	prevCount := as.countPosts()

	// without login it has to redirect to signin page
	res := as.HTML(postUrl, post.ID).Put(post)
	as.Equal(http.StatusFound, res.Code)
	as.Equal(res.Location(), signinPath)

	// login a user
	signinUser(as, user)

	//modified post
	post.Title = updatedPostTitle
	// edit post
	res = as.HTML(postUrl, post.ID).Put(post)
	as.Equal(302, res.Code)
	as.Equal(res.Location(), fmt.Sprintf(postUrl, post.ID))
	as.Equal(0, as.countPosts()-prevCount)

	// check post has been updated
	post = getPostFromDB(as, post)
	user.WithName(as.DB)
	as.Equal(updatedPostTitle, post.Title)
	as.Equal(user.ID, post.UserID)

	//	don't allow editing from other users
	user2 := user2ForTest(as)
	as.HTML(signoutPath).Get()
	signinUser(as, user2)
	res = as.HTML(postUrl, post.ID).Put(post)
	as.Equal(http.StatusUnauthorized, res.Code)
}

func (as *ActionSuite) Test_PostsResource_Destroy() {
	user, post := createPost(as)

	//	login as another user
	user2 := user2ForTest(as)
	signinUser(as, user2)
	res := as.HTML(postUrl, post.ID).Delete()
	as.Equal(http.StatusUnauthorized, res.Code)
	as.HTML(signoutPath).Get()

	// login as author
	signinUser(as, user)

	res = as.HTML(postUrl, post.ID).Delete()
	as.Equal(302, res.Code)
}
