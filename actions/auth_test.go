package actions_test

import (
	"muserblog/models"
	"net/http"
	"github.com/markbates/pop/nulls"
)

func (as *ActionSuite) countUsers() int {
	return as.Count(models.User{})
}

const (
	signupPath = "/auth/signup"
	signinPath = "/auth/signin"
)

var (
	usrName = "test_1"
	usrEmail = "test@mail.com"
	usrPwd = "plain_password"
)

func userForTest() *models.User {
	return &models.User{Name:usrName, Email:nulls.NewString(usrEmail), Pwd:usrPwd}
}

func (as *ActionSuite) Test_GetAuthFormHanlers() {
	// test sign in
	res := as.HTML("/auth/signin").Get()
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
	as.Equal(1, as.countUsers() - prevCnt)

	// fails to create user with the same name
	res = as.HTML(signupPath).Post(user)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Equal(1, as.countUsers() - prevCnt)
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
