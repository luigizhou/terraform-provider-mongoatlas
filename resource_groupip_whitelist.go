package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
	"io/ioutil"
)

type GroupipWhitelist struct {
	CidrBlock string `json:"cidrBlock,omitempty"`
	IpAddress string `json:"ipAddress,omitempty"`
	GroupId   string `json:"groupId,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

type RequestObject struct {
	Requestarray []GroupipWhitelist
}

func resourceGroupipWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupipWhitelistCreate,
		Update: resourceGroupipWhitelistUpdate,
		Read:   resourceGroupipWhitelistRead,
		Delete: resourceGroupipWhitelistDelete,
		Schema: map[string]*schema.Schema{
			"cidrBlock": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipAddress": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"groupId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGroupipWhitelistCreate(d *schema.ResourceData, m interface{}) error {
	// Instantiate the MongoatlasClient. Source in client.go, it handles all the HTTP REST API communication
	var isip bool
	client := m.(*MongoatlasClient)

	groupipwhitelist := &GroupipWhitelist{}

	if attr, ok := d.GetOk("ipAddress"); ok {
		groupipwhitelist.IpAddress = attr.(string)
		isip = true
	} else if attr, ok := d.GetOk("cidrBlock"); ok {
		groupipwhitelist.CidrBlock = attr.(string)
		isip = false
	}

	if attr, ok := d.GetOk("comment"); ok {
		groupipwhitelist.Comment = attr.(string)
	}

	requestobject := &RequestObject{
		Requestarray: []GroupipWhitelist{*groupipwhitelist},
	}

	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(requestobject.Requestarray)

	log.Printf("Sending %s \n", jsonpayload)

	// communication with API commence here
	groupipwhitelist_req, err := client.Post(fmt.Sprintf("groups/%s/whitelist",
		d.Get("groupId").(string),
	), jsonpayload)
	if err != nil {
		return err
	}

	// if statuscode is different than 201, nothing is persisted in the tfstate file, but what happened on mongo atlas is not checked.
	// TODO: handle other status code
	if groupipwhitelist_req.StatusCode != 201 {
		body, err := ioutil.ReadAll(groupipwhitelist_req.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Failed to create vpc peering got status code %d", string(body))
	}

	if isip {
		d.SetId(d.Get("ipAddress").(string) + "/32")
		d.Set("cidrBlock", d.Get("ipAddress").(string)+"/32")
		d.Set("ipAddress", d.Get("ipAddress").(string))
	} else {
		d.SetId(d.Get("cidrBlock").(string))
		d.Set("cidrBlock", d.Get("cidrBlock").(string))
	}

	return resourceGroupipWhitelistRead(d, m)
}

func resourceGroupipWhitelistRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*MongoatlasClient)

	log.Printf("%s", d.Get("groupId").(string))
	address := strings.Replace(d.Get("cidrBlock").(string), "/", "%2F", -1)

	groupipwhitelist_req, err := client.Get(fmt.Sprintf("groups/%s/whitelist/%s",
		d.Get("groupId").(string),
		address,
	))

	if err != nil {
		return err
	}

	var groupipwhitelist GroupipWhitelist

	decoder := json.NewDecoder(groupipwhitelist_req.Body)
	err = decoder.Decode(&groupipwhitelist)
	if err != nil {
		return err
	}

	log.Printf("Received %s \n", groupipwhitelist_req.Body)

	d.Set("cidrBlock", groupipwhitelist.CidrBlock)
	d.Set("ipAddress", groupipwhitelist.IpAddress)
	d.Set("Comment", groupipwhitelist.Comment)

	return nil
}

func resourceGroupipWhitelistUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceGroupipWhitelistDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*MongoatlasClient)

	address := strings.Replace(d.Get("cidrBlock").(string), "/", "%2F", -1)

	delete_response, err := client.Delete(fmt.Sprintf("groups/%s/whitelist/%s",
		d.Get("groupId").(string),
		address,
	))

	if err != nil {
		return err
	}

	// should handle 404 when resource is not found. In that case, resource is already deleted
	if delete_response.StatusCode != 200 {
		return fmt.Errorf("Failed to delete the vpc peering. Got status code %d", delete_response.StatusCode)
	}
	return nil
}
