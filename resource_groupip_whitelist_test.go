package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMongoatlasGroupipWhitelist_basic(t *testing.T) {
	var groupipwhitelist GroupipWhitelist

	testGroupId := os.Getenv("MONGOATLAS_GROUPID")

	testAccMongoatlasGroupipWhitelistConfig := fmt.Sprintf(
		`resource "mongoatlas_groupip_whitelist" "acceptancetest_groupipwhitelist" {
			cidrBlock = "1.2.3.4/32"
		    groupId = "%s"
		    comment = "terraform test"
		}
	`, testGroupId)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoatlasGroupipWhitelistDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMongoatlasGroupipWhitelistConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoatlasGroupipWhitelistExists("mongoatlas_groupip_whitelist.acceptancetest_groupipwhitelist", &groupipwhitelist),
				),
			},
		},
	})

}

func testAccCheckMongoatlasGroupipWhitelistDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*MongoatlasClient)
	rs, ok := s.RootModule().Resources["mongoatlas_groupip_whitelist.acceptancetest_groupipwhitelist"]
	if !ok {
		return fmt.Errorf("Not found %s", "mongoatlas_groupip_whitelist.acceptancetest_groupipwhitelist")
	}

	address := strings.Replace(rs.Primary.Attributes["cidrBlock"], "/", "%2F", -1)

	response, err := client.Get(fmt.Sprintf("groups/%s/whitelist/%s", rs.Primary.Attributes["groupId"], address))

	if err != nil {
		return err
	}

	if response.StatusCode != 404 {
		return fmt.Errorf("GroupipWhitelist still exists")
	}

	return nil
}

func testAccCheckMongoatlasGroupipWhitelistExists(n string, groupipwhitelist *GroupipWhitelist) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Group IP Whitelist ID is set")
		}
		return nil
	}
}
