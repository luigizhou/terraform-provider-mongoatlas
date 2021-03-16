package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

const testVpcPeering string = "test-vpcpeering"

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"mongoatlas": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("MONGOATLAS_USERNAME"); v == "" {
		t.Fatal("MONGOATLAS_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("MONGOATLAS_APIKEY"); v == "" {
		t.Fatal("MONGOATLAS_APIKEY must be set for acceptance tests")
	}
	if v := os.Getenv("MONGOATLAS_GROUPID"); v == "" {
		t.Fatal("MONGOATLAS_GROUPID must be set for acceptance tests")
	}
	if v := os.Getenv("MONGOATLAS_VPCID"); v == "" {
		t.Fatal("MONGOATLAS_VPCID must be set for acceptance tests")
	}
	if v := os.Getenv("MONGOATLAS_AWSACCOUNTID"); v == "" {
		t.Fatal("MONGOATLAS_AWSACCOUNTID must be set for acceptance tests")
	}
}
