// da rivedere, non c'e' la delete sulle API

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"io/ioutil"
)

type Container struct {
	Id             string `json:"id,omitempty"`
	GroupId        string `json:"groupId,omitempty"`
	AtlasCidrBlock string `json:"atlasCidrBlock,omitempty"`
	ProviderName   string `json:"providerName,omitempty"`
	RegionName     string `json:"regionName,omitempty"`
	VpcId          string `json:"vpcId,omitempty"`
	IsProvisioned    string `json:"isProvisioned,omitempty"`
}

func resourceContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerCreate,
		Update: resourceContainerUpdate,
		Read:   resourceContainerRead,
		Delete: resourceContainerDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"groupId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"atlasCidrBlock": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"providerName": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"regionName": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpcId": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"isProvisioned": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func newContainer(d *schema.ResourceData) *Container {
	container := &Container{
		AtlasCidrBlock: d.Get("atlasCidrBlock").(string),
		ProviderName:   d.Get("providerName").(string),
		RegionName:     d.Get("regionName").(string),
	}

	return container

}

func resourceContainerCreate(d *schema.ResourceData, m interface{}) error {
	// Instantiate the MongoatlasClient. Source in client.go, it handles all the HTTP REST API communication
	client := m.(*MongoatlasClient)

	container := newContainer(d)
	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(container)

	log.Printf("Sending %s \n", jsonpayload)

	// communication with API commence here
	container_req, err := client.Post(fmt.Sprintf("groups/%s/containers",
		d.Get("groupId").(string),
	), jsonpayload)

	if container_req.StatusCode != 201 {
		body, err := ioutil.ReadAll(container_req.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Failed to create container - got following response body: %s", string(body))
	}

	decoder := json.NewDecoder(container_req.Body)
	err = decoder.Decode(&container)
	if err != nil {
		return err
	}

	log.Printf("Received %s \n", container_req.Body)
	

	// The following statement saves set the data that will be saved in the .tfstate file
	d.SetId(container.Id)
	d.Set("id", container.Id)
	d.Set("atlasCidrBlock", container.AtlasCidrBlock)
	d.Set("providerName", container.ProviderName)
	d.Set("regionName", container.RegionName)
	d.Set("vpcId", container.VpcId)
	d.Set("isProvisioned", container.IsProvisioned)

	return resourceContainerRead(d, m)

}

func resourceContainerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*MongoatlasClient)

	log.Printf("%s", d.Get("id").(string))
	log.Printf("%s", d.Get("groupId").(string))

	container_req, err := client.Get(fmt.Sprintf("groups/%s/containers/%s",
		d.Get("groupId").(string),
		d.Get("id").(string),
	))

	if err != nil {
		return err
	}

	var container Container

	decoder := json.NewDecoder(container_req.Body)
	err = decoder.Decode(&container)
	if err != nil {
		return err
	}

	d.Set("id", container.Id)
	d.Set("atlasCidrBlock", container.AtlasCidrBlock)
	d.Set("providerName", container.ProviderName)
	d.Set("regionName", container.RegionName)
	d.Set("vpcId", container.VpcId)
	d.Set("isProvisioned", container.IsProvisioned)

	return nil
}

func resourceContainerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceContainerCreate(d, m)
}

func resourceContainerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
