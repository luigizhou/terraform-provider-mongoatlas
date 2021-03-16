package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMongoatlasDatabaseUser_basic(t *testing.T) {
	var databaseuser DatabaseUser

	testGroupId := os.Getenv("MONGOATLAS_GROUPID")

	testAccMongoatlasDatabaseUserConfig_basic := fmt.Sprintf(
		`resource "mongoatlas_database_user" "acceptancetest_databaseuser" {
			databaseName = "admin"
		    username = "acctest"
		    password = "test"
		    roles = [
		        {
		            databaseName = "admin"
		            roleName = "backup"
		        }
		    ]   
	    	groupId = "%s"
		}
	`, testGroupId)

	testAccMongoatlasDatabaseUserConfig_update := fmt.Sprintf(
		`resource "mongoatlas_database_user" "acceptancetest_databaseuser" {
	    	databaseName = "admin"
		    username = "acctest"
		    password = "test"
		    roles = [
		        {
		        	databaseName = "admin"
		        	roleName = "readAnyDatabase"
		        }
		    ]   
	    	groupId = "%s"
		}
	`, testGroupId)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoatlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMongoatlasDatabaseUserConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoatlasDatabaseUserExists("mongoatlas_database_user.acceptancetest_databaseuser", &databaseuser),
				),
			},

			resource.TestStep{
				Config: testAccMongoatlasDatabaseUserConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoatlasDatabaseUserExists("mongoatlas_database_user.acceptancetest_databaseuser", &databaseuser),
					resource.TestCheckResourceAttr(
						"mongoatlas_database_user.acceptancetest_databaseuser", "roles.0.roleName", "readAnyDatabase"),
				),
			},
		},
	})

}

func testAccCheckMongoatlasDatabaseUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*MongoatlasClient)
	rs, ok := s.RootModule().Resources["mongoatlas_database_user.acceptancetest_databaseuser"]
	if !ok {
		return fmt.Errorf("Not found %s", "mongoatlas_database_user.acceptancetest_databaseuser")
	}

	response, err := client.Get(fmt.Sprintf("groups/%s/databaseUsers/admin/%s", rs.Primary.Attributes["groupId"], rs.Primary.Attributes["username"]))

	if err != nil {
		return err
	}

	if response.StatusCode != 404 {
		return fmt.Errorf("DatabaseUser still exists")
	}

	return nil
}

func testAccCheckMongoatlasDatabaseUserExists(n string, databaseuser *DatabaseUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No databaseuser ID/name is set")
		}
		return nil
	}
}
