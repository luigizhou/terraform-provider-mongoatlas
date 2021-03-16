package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Required:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("MONGOATLAS_USERNAME", nil),
			},
			"apiKey": {
				Required:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("MONGOATLAS_APIKEY", nil),
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"mongoatlas_vpc_peering":       resourceVpcPeering(),
			"mongoatlas_cluster":           resourceCluster(),
			"mongoatlas_database_user":     resourceDatabaseUser(),
			"mongoatlas_groupip_whitelist": resourceGroupipWhitelist(),
			"mongoatlas_container":         resourceContainer(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := &MongoatlasClient{
		Username: d.Get("username").(string),
		ApiKey:   d.Get("apiKey").(string),
	}

	return client, nil
}
