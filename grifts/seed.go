package grifts

import (
	"muserblog/models"

	"github.com/markbates/grift/grift"
)

var _ = grift.Add("db:seed", func(c *grift.Context) error {
	// Add DB seeding stuff here
	user := models.NewUser("user_1", "password", "user1@mail.com")
	models.DB.Create(user)

	user = models.NewUser("user_1", "password", "user1@mail.com")
	models.DB.Create(user)

	return nil
})
