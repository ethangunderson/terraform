package heroku

import (
	"fmt"
	"testing"

	"github.com/cyberdelia/heroku-go/v3"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHerokuSpace_Basic(t *testing.T) {
	var space heroku.Space
	appName := fmt.Sprintf("tftest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHerokuspaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckHerokuSpaceConfig_basic(appName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHerokuspaceExists("heroku_space.foobar", &space),
					testAccCheckHerokuSpaceAttributes(&space),
					resource.TestCheckResourceAttr(
						"heroku_space.foobar", "url", "syslog://terraform.example.com:1234"),
					resource.TestCheckResourceAttr(
						"heroku_space.foobar", "app", appName),
				),
			},
		},
	})
}

func testAccCheckHerokuspaceAttributes(Space *heroku.Space) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if Space.Name != "syslog://terraform.example.com:1234" {
			return fmt.Errorf("Bad Name: %s", Space.Name)
		}

		return nil
	}
}
