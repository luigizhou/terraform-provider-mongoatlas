package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"io/ioutil"
)

type VpcPeering struct {
	VpcId               string `json:"vpcId,omitempty"`
	AwsAccountId        string `json:"awsAccountId,omitempty"`
	RouteTableCidrBlock string `json:"routeTableCidrBlock,omitempty"`
	Id                  string `json:"id,omitempty"`
	ConnectionId        string `json:"connectionId,omitempty"`
	StatusName          string `json:"statusName,omitempty"`
	ErrorStateName      string `json:"errorStateName,omitempty"` // TODO: Must handle ErrorStateName as it is empty if everything goes smoothly
}

func resourceVpcPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcPeeringCreate,
		Update: resourceVpcPeeringUpdate,
		Read:   resourceVpcPeeringRead,
		Delete: resourceVpcPeeringDelete,
		Schema: map[string]*schema.Schema{
			"groupId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpcId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"awsAccountId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"routeTableCidrBlock": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"connectionId": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"statusName": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"errorStateName": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func newVpcPeering(d *schema.ResourceData) *VpcPeering {
	vpcpeering := &VpcPeering{
		VpcId:               d.Get("vpcId").(string),
		AwsAccountId:        d.Get("awsAccountId").(string),
		RouteTableCidrBlock: d.Get("routeTableCidrBlock").(string),
	}

	return vpcpeering

}

func resourceVpcPeeringCreate(d *schema.ResourceData, m interface{}) error {
	// Instantiate the MongoatlasClient. Source in client.go, it handles all the HTTP REST API communication
	client := m.(*MongoatlasClient)

	vpcpeering := newVpcPeering(d)
	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(vpcpeering)

	log.Printf("Sending %s \n", jsonpayload)

	// communication with API commence here
	vpcpeering_req, err := client.Post(fmt.Sprintf("groups/%s/peers",
		d.Get("groupId").(string),
	), jsonpayload)

	decoder := json.NewDecoder(vpcpeering_req.Body)
	err = decoder.Decode(&vpcpeering)
	if err != nil {
		return err
	}

	log.Printf("Received %s \n", vpcpeering_req.Body)

	// if statuscode is different than 201, nothing is persisted in the tfstate file, but what happened on mongo atlas is not checked.
	// TODO: handle other status code
	if vpcpeering_req.StatusCode != 201 {
		body, err := ioutil.ReadAll(vpcpeering_req.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Failed to create vpc peering got status code %d", string(body))
	}

	// The following statement saves set the data that will be saved in the .tfstate file
	d.SetId(vpcpeering.Id)
	d.Set("id", vpcpeering.Id)
	d.Set("awsAccountId", vpcpeering.AwsAccountId)
	d.Set("connectionId", vpcpeering.ConnectionId) //TODO: understand why connectionId returns empty | It get updates on read
	d.Set("routeTableCidrBlock", vpcpeering.RouteTableCidrBlock)
	d.Set("statusName", vpcpeering.StatusName)
	d.Set("errorStateName", vpcpeering.ErrorStateName)

	return resourceVpcPeeringRead(d, m)

}

func resourceVpcPeeringRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*MongoatlasClient)

	log.Printf("%s", d.Get("id").(string))
	log.Printf("%s", d.Get("groupId").(string))

	vpcpeering_req, err := client.Get(fmt.Sprintf("groups/%s/peers/%s",
		d.Get("groupId").(string),
		d.Get("id").(string),
	))

	if err != nil {
		return err
	}

	var vpcpeering VpcPeering

	decoder := json.NewDecoder(vpcpeering_req.Body)
	err = decoder.Decode(&vpcpeering)
	if err != nil {
		return err
	}
	log.Printf("Received %s \n", vpcpeering_req.Body)

	d.Set("awsAccountId", vpcpeering.AwsAccountId)
	d.Set("routeTableCidrBlock", vpcpeering.RouteTableCidrBlock)
	d.Set("id", vpcpeering.Id)
	d.Set("connectionId", vpcpeering.ConnectionId)
	d.Set("statusName", vpcpeering.StatusName)
	d.Set("errorStateName", vpcpeering.ErrorStateName)
	return nil
}

func resourceVpcPeeringUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*MongoatlasClient)

	vpcpeering := VpcPeering{}

	// TODO: VPCID and AWSACCOUNTID are handled together. On their own the changes are not applied
	if d.HasChange("vpcId") {
		vpcpeering.VpcId = d.Get("vpcId").(string)
	}
	if d.HasChange("awsAccountId") {
		vpcpeering.AwsAccountId = d.Get("awsAccountId").(string)
	}
	if d.HasChange("routeTableCidrBlock") {
		vpcpeering.RouteTableCidrBlock = d.Get("routeTableCidrBlock").(string)
	}

	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(vpcpeering)

	log.Printf("Sending %s \n", jsonpayload)

	vpcpeering_req, err := client.Patch(fmt.Sprintf("groups/%s/peers/%s",
		d.Get("groupId").(string),
		d.Get("id").(string),
	), jsonpayload)

	if err != nil {
		return err
	}

	if vpcpeering_req.StatusCode == 200 {
		decoder := json.NewDecoder(vpcpeering_req.Body)
		err = decoder.Decode(&vpcpeering)
		if err != nil {
			return err
		}
	} else {
		body, err := ioutil.ReadAll(vpcpeering_req.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Failed to patch: %s", string(body))
	}

	return resourceVpcPeeringRead(d, m)

	return nil
}

func resourceVpcPeeringDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*MongoatlasClient)

	delete_response, err := client.Delete(fmt.Sprintf("groups/%s/peers/%s",
		d.Get("groupId").(string),
		d.Get("id").(string),
	))

	if err != nil {
		return err
	}

	// should handle 404 when resource is not found. In that case, resource is already deleted
	if delete_response.StatusCode != 202 {
		return fmt.Errorf("Failed to delete the vpc peering. Got status code %d", delete_response.StatusCode)
	}
	return nil
}
