package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMongoatlasCluster_basic(t *testing.T) {
	var cluster Cluster

	testGroupId := os.Getenv("MONGOATLAS_GROUPID")

	testAccMongoatlasClusterConfig := fmt.Sprintf(
		`resource "mongoatlas_cluster" "acceptancetest_cluster" {
	    	groupId = "%s"
	    	name = "terratest3"
		    backupEnabled = false
		    instanceSizeName = "M10"
		    diskSizeGB = "10"
		    providerName = "AWS"
		    regionName = "EU_WEST_1"
		}
	`, testGroupId)

	testAccMongoatlasClusterConfig_update := fmt.Sprintf(
		`resource "mongoatlas_cluster" "acceptancetest_cluster" {
	    	groupId = "%s"
	    	name = "terratest3"
		    backupEnabled = false
		    instanceSizeName = "M10"
		    diskSizeGB = "11"
		    providerName = "AWS"
		    regionName = "EU_WEST_1"
		}
	`, testGroupId)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoatlasClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMongoatlasClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoatlasClusterExists("mongoatlas_cluster.acceptancetest_cluster", &cluster),
				),
			},

			resource.TestStep{
				Config: testAccMongoatlasClusterConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoatlasClusterExists("mongoatlas_cluster.acceptancetest_cluster", &cluster),
					resource.TestCheckResourceAttr(
						"mongoatlas_cluster.acceptancetest_cluster", "diskSizeGB", "11"),
				),
			},
		},
	})

}

func testAccCheckMongoatlasClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*MongoatlasClient)
	rs, ok := s.RootModule().Resources["mongoatlas_cluster.acceptancetest_cluster"]
	if !ok {
		return fmt.Errorf("Not found %s", "mongoatlas_cluster.acceptancetest_cluster")
	}

	response, err := client.Get(fmt.Sprintf("groups/%s/clusters/%s", rs.Primary.Attributes["groupId"], rs.Primary.Attributes["name"]))

	if err != nil {
		return err
	}

	if response.StatusCode != 404 {
		var cluster Cluster

		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&cluster)
		if err != nil {
			return err
		}

		if cluster.StateName != "DELETING" {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func TestAccMongoAtlasClusterDiskSizeGB_validation(t *testing.T) {
	cases := []struct {
		Value    float64
		ErrCount int
	}{
		{
			Value:    9,
			ErrCount: 1,
		},
		{
			Value:    10,
			ErrCount: 0,
		},
		{
			Value:    16384,
			ErrCount: 0,
		},
		{
			Value:    16385,
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateDiskSizeGB(tc.Value, "mongoatlas_cluster_disksizegb")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %+v Validation Error, Got %+v Validation Error for %+v VALUE", tc.ErrCount, len(errors), tc.Value)
		}
	}
}

func TestAccMongoAtlasClusterProviderName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "AWS",
			ErrCount: 0,
		},
		{
			Value:    "Azure",
			ErrCount: 1,
		},
		{
			Value:    "Anything Else",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateProviderName(tc.Value, "mongoatlas_cluster_providername")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %+v Validation Error, Got %+v Validation Error for %+v VALUE", tc.ErrCount, len(errors), tc.Value)
		}
	}
}

func TestAccMongoAtlasClusterNumShards_validation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{
			Value:    0,
			ErrCount: 1,
		},
		{
			Value:    6,
			ErrCount: 0,
		},
		{
			Value:    13,
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateNumShards(tc.Value, "mongoatlas_cluster_numshards")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %+v Validation Error, Got %+v Validation Error for %+v VALUE", tc.ErrCount, len(errors), tc.Value)
		}
	}
}

func TestAccMongoAtlasClusterInstanceSizeName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "M10",
			ErrCount: 0,
		},
		{
			Value:    "M20",
			ErrCount: 0,
		},
		{
			Value:    "M30",
			ErrCount: 0,
		},
		{
			Value:    "M40",
			ErrCount: 0,
		},
		{
			Value:    "M50",
			ErrCount: 0,
		},
		{
			Value:    "M60",
			ErrCount: 0,
		},
		{
			Value:    "M100",
			ErrCount: 0,
		},
		{
			Value:    "M0",
			ErrCount: 1,
		},
		{
			Value:    "m10",
			ErrCount: 0,
		},
		{
			Value:    "M15",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateInstanceSizeName(tc.Value, "mongoatlas_cluster_instancesizename")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %+v Validation Error, Got %+v Validation Error for %+v VALUE", tc.ErrCount, len(errors), tc.Value)
		}
	}
}

func TestAccMongoAtlasClusterRegionName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "AP_SOUTHEAST_2",
			ErrCount: 0,
		},
		{
			Value:    "EU_WEST_1",
			ErrCount: 0,
		},
		{
			Value:    "US_EAST_1",
			ErrCount: 0,
		},
		{
			Value:    "US_WEST_2",
			ErrCount: 0,
		},
		{
			Value:    "us-east-1",
			ErrCount: 1,
		},
		{
			Value:    "us-east-1",
			ErrCount: 1,
		},
		{
			Value:    "some non existent value",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateRegionName(tc.Value, "mongoatlas_cluster_regionname")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %+v Validation Error, Got %+v Validation Error for %+v VALUE", tc.ErrCount, len(errors), tc.Value)
		}
	}
}

func TestAccMongoAtlasClusterReplicationFactor_validation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{
			Value:    3,
			ErrCount: 0,
		},
		{
			Value:    5,
			ErrCount: 0,
		},
		{
			Value:    7,
			ErrCount: 0,
		},
		{
			Value:    11,
			ErrCount: 1,
		},
		{
			Value:    0,
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateReplicationFactor(tc.Value, "mongoatlas_cluster_replicationfactor")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %+v Validation Error, Got %+v Validation Error for %+v VALUE", tc.ErrCount, len(errors), tc.Value)
		}
	}
}

func testAccCheckMongoatlasClusterExists(n string, cluster *Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No cluster ID/name is set")
		}
		return nil
	}
}
