package actions_test

import "muserblog/models"

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

func countPostsTo(as *ActionSuite, expectedCount int) {
	as.Equal(expectedCount, as.countPosts())
}

// create a new post and return
func createPost(as *ActionSuite) *models.Post {
	prevCount := as.countPosts()
	post := postForTest()
	as.NoError(as.DB.Create(post))
	// check exactly one new record is created
	as.Equal(as.countPosts()-prevCount, 1)
	return post
}

// return the First record that matches the query
func getPostFromDB(as *ActionSuite, post *models.Post) *models.Post {
	err := as.DB.First(post)
	as.NoError(err)
	as.NotZero(post.ID)
	as.NotZero(post.CreatedAt)
	return post
}

func (as *ActionSuite) Test_PostsResource_List() {
	res := as.HTML(postsListUrl).Get()

	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Posts")
}

func (as *ActionSuite) Test_PostsResource_Show() {
	post := createPost(as)
	res := as.HTML(postUrl, post.ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), post.Content)
}

func (as *ActionSuite) Test_PostsResource_New() {
	res := as.HTML(postsNewUrl).Get()
	as.Contains(res.Body.String(), "<h1>New Post</h1>")
}

func (as *ActionSuite) Test_PostsResource_Create() {
	countPostsTo(as, 0)
	post := postForTest()

	res := as.HTML(postsListUrl).Post(post)
	as.Equal(302, res.Code)

	countPostsTo(as, 1)

	post = getPostFromDB(as, post)
	as.Equal("Post 1", post.Title)
}

func (as *ActionSuite) Test_PostsResource_Edit() {
	post := createPost(as)
	res := as.HTML(postsEditUrl, post.ID).Get()
	as.Contains(res.Body.String(), "<h1>Edit Post</h1>")
	as.Contains(res.Body.String(), post.Title)
	as.Contains(res.Body.String(), post.Content)
}

func (as *ActionSuite) Test_PostsResource_Update() {
	post := createPost(as)

	post.Title = updatedPostTitle

	res := as.HTML(postUrl, post.ID).Put(post)
	as.Equal(302, res.Code)

	countPostsTo(as, 1)

	post = getPostFromDB(as, post)
	as.Equal(updatedPostTitle, post.Title)
}

func (as *ActionSuite) Test_PostsResource_Destroy() {
	post := createPost(as)

	res := as.HTML(postUrl, post.ID).Delete()
	as.Equal(302, res.Code)

	countPostsTo(as, 0)
}
