package actions_test

import (
	"testing"

	"github.com/gobuffalo/suite"
	"muserblog/actions"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	as := &ActionSuite{suite.NewAction(actions.App())}
	suite.Run(t, as)
}

func (as *ActionSuite) Count(model interface{}) int {
	cnt, err := as.DB.Count(model)
	as.NoError(err)
	return cnt
}
