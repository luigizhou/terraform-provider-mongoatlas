package main

import (
	"fmt"
	"os"
	"testing"

	"encoding/json"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMongoatlasVpcpeering_basic(t *testing.T) {
	var vpcpeering VpcPeering

	testGroupId := os.Getenv("MONGOATLAS_GROUPID")
	testVpcId := os.Getenv("MONGOATLAS_VPCID")
	testAwsAccountId := os.Getenv("MONGOATLAS_AWSACCOUNTID")

	testAccMongoatlasVpcpeeringConfig := fmt.Sprintf(
		`resource "mongoatlas_vpc_peering" "acceptancetest_vpcpeering" {
	    	groupId= "%s"
	    	vpcId= "%s"
	    	awsAccountId = "%s"
	    	routeTableCidrBlock = "10.230.8.0/24"
		}
	`, testGroupId, testVpcId, testAwsAccountId)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoatlasVpcpeeringDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMongoatlasVpcpeeringConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoatlasVpcpeeringExists("mongoatlas_vpc_peering.acceptancetest_vpcpeering", &vpcpeering),
				),
			},
		},
	})

}

func testAccCheckMongoatlasVpcpeeringDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*MongoatlasClient)
	rs, ok := s.RootModule().Resources["mongoatlas_vpc_peering.acceptancetest_vpcpeering"]
	if !ok {
		return fmt.Errorf("Not found %s", "mongoatlas_vpc_peering.acceptancetest_vpcpeering")
	}

	response, err := client.Get(fmt.Sprintf("groups/%s/peers/%s", rs.Primary.Attributes["groupId"], rs.Primary.Attributes["id"]))

	if err != nil {
		return err
	}

	if response.StatusCode != 404 {
		var vpcpeering VpcPeering

		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&vpcpeering)
		if err != nil {
			return err
		}

		if vpcpeering.StatusName != "TERMINATING" {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccCheckMongoatlasVpcpeeringExists(n string, vpcpeering *VpcPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No vpc peering ID is set")
		}
		return nil
	}
}
