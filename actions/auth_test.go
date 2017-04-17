package actions_test

func (as *ActionSuite) Test_GetAuthFormHanlers() {
	// test sign in
	res := as.HTML("/auth/signin").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Sign in")

	//	test signup
	res = as.HTML("/auth/signup").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Sign Up")
}
