package actions_test

import (
	"github.com/markbates/pop/nulls"
	"github.com/markbates/willie"
	"muserblog/models"
	"net/http"
)

func (as *ActionSuite) countUsers() int {
	return as.Count(models.User{})
}

const (
	signupPath  = "/auth/signup"
	signinPath  = "/auth/signin"
	signoutPath = "/auth/signout"
)

var (
	usrName  = "test_1"
	usrEmail = "test@mail.com"
	usrPwd   = "plain_password"
)

func userForTest() *models.User {
	return &models.User{Name: usrName, Email: nulls.NewString(usrEmail), Pwd: usrPwd}
}

func signinUser(as *ActionSuite, user *models.User) {
	// creates a new user record
	tUser := *user
	tUser.SaltPassword()
	as.NoError(as.DB.Create(&tUser))

	// ensure that session cookie is present
	res := as.HTML("/").Get()

	// logins that user
	res = as.HTML(signinPath).Post(user)
	as.Equal(http.StatusFound, res.Code)
}

func (as *ActionSuite) Test_GetAuthFormHanlers() {
	// test sign in
	res := as.HTML(signinPath).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Sign in")

	// test signup
	res = as.HTML(signupPath).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Sign Up")
}

func (as *ActionSuite) Test_UserSignUp() {
	prevCnt := as.countUsers()
	user := userForTest()

	// creates a new user successfully
	res := as.HTML(signupPath).Post(user)
	as.Equal(http.StatusFound, res.Code)
	as.Equal(1, as.countUsers()-prevCnt)

	// fails to create user with the same name
	res = as.HTML(signupPath).Post(user)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Equal(1, as.countUsers()-prevCnt)
}

func (as *ActionSuite) Test_UserSignIn() {
	user := userForTest()
	user.SaltPassword()
	as.NoError(as.DB.Create(user))

	res := as.HTML(signinPath).Post(userForTest())
	as.Equal(http.StatusFound, res.Code)

	// wrong credentials
	user = userForTest()
	user.Name = "test_uk"
	res = as.HTML(signinPath).Post(user)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
}

func (as *ActionSuite) Test_UserSigninProcess() {
	var res *willie.Response

	// create session cookie
	res = as.HTML("/").Get()

	// create new user by creating an account
	res = as.HTML(signupPath).Post(userForTest())
	as.Equal(http.StatusFound, res.Code)

	// signin user : this will create auth_token and add to the cookie
	res = as.HTML(signinPath).Post(userForTest())

	// go to the posts page and check that the user is signed in
	res = as.HTML(postsListUrl).Get()
	as.Contains(res.Body.String(), "Signout")

	// user clicks the signout link
	res = as.HTML(signoutPath).Get()
	// it redirects to signin page
	as.Equal(res.Location(), "/auth/signin")

	res = as.HTML(res.Location()).Get()
	as.Contains(res.Body.String(), "Signin")
}
