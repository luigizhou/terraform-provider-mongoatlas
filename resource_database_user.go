package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"io/ioutil"
)

type DatabaseUser struct {
	DatabaseName string  `json:"databaseName,omitempty"`
	Username     string  `json:"username,omitempty"`
	Roles        *[]Role `json:"roles,omitempty"` // check if this is right or not
	Password     string  `json:"password,omitempty"`
}

type Role struct {
	DatabaseName string `json:"databaseName,omitempty"`
	RoleName     string `json:"roleName,omitempty"`
}

//TODO: need to write validation for rolename
func resourceDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseUserCreate,
		Update: resourceDatabaseUserUpdate,
		Read:   resourceDatabaseUserRead,
		Delete: resourceDatabaseUserDelete,
		Schema: map[string]*schema.Schema{
			"groupId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"databaseName": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"databaseName": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"roleName": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateRoleName,
						},
					},
				},
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Sensitive: true,
			},
		},
	}
}

func newDatabaseUser(d *schema.ResourceData) *DatabaseUser {

	var roles []Role

	rolesInterface := d.Get("roles").([]interface{})
	for _, roleInterface := range rolesInterface {
		rolesMap := roleInterface.(map[string]interface{})
		role := Role{
			DatabaseName: rolesMap["databaseName"].(string),
			RoleName:     rolesMap["roleName"].(string),
		}
		roles = append(roles, role)
	}


	databaseuser := &DatabaseUser{
		DatabaseName: d.Get("databaseName").(string),
		Username:     d.Get("username").(string),
		Roles:        &roles,
		Password:     d.Get("password").(string),
	}

	return databaseuser

}

func resourceDatabaseUserCreate(d *schema.ResourceData, m interface{}) error {
	// Instantiate the MongoatlasClient. Source in client.go, it handles all the HTTP REST API communication
	client := m.(*MongoatlasClient)

	databaseuser := newDatabaseUser(d)

	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(databaseuser)

	log.Printf("Sending %s \n", jsonpayload)

	// communication with API commence here
	databaseuser_req, err := client.Post(fmt.Sprintf("groups/%s/databaseUsers",
		d.Get("groupId").(string),
	), jsonpayload)

	// if statuscode is different than 201, nothing is persisted in the tfstate file, but what happened on mongo atlas is not checked.
	// TODO: handle other status code
	if databaseuser_req.StatusCode != 201 {
		body, err := ioutil.ReadAll(databaseuser_req.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Failed to create database user. Got the following response body %s", string(body))
	}

	decoder := json.NewDecoder(databaseuser_req.Body)
	err = decoder.Decode(&databaseuser)
	if err != nil {
		return err
	}

	log.Printf("Received %s \n", databaseuser_req.Body)

	

	// The following statement saves set the data that will be saved in the .tfstate file
	d.SetId(databaseuser.Username)
	d.Set("username", databaseuser.Username)
	d.Set("databaseName", databaseuser.DatabaseName)

	d.Set("password", databaseuser.Password)

	var s []map[string]interface{}
	for _, t := range *databaseuser.Roles {
		mapping := map[string]interface{}{
			"databaseName": t.DatabaseName,
			"roleName":     t.RoleName,
		}

		s = append(s, mapping)
	}

	if err := d.Set("roles", s); err != nil {
		return err
	}

	return resourceDatabaseUserRead(d, m)
}

func resourceDatabaseUserRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*MongoatlasClient)

	log.Printf("%s", d.Get("username").(string))
	log.Printf("%s", d.Get("groupId").(string))

	databaseuser_req, err := client.Get(fmt.Sprintf("groups/%s/databaseUsers/admin/%s",
		d.Get("groupId").(string),
		d.Get("username").(string),
	))

	if err != nil {
		return err
	}

	var databaseuser DatabaseUser

	decoder := json.NewDecoder(databaseuser_req.Body)
	err = decoder.Decode(&databaseuser)
	if err != nil {
		return err
	}
	log.Printf("Received %s \n", databaseuser_req.Body)

	d.Set("username", databaseuser.Username)
	d.Set("databaseName", databaseuser.DatabaseName)

	var s []map[string]interface{}
	for _, t := range *databaseuser.Roles {
		mapping := map[string]interface{}{
			"databaseName": t.DatabaseName,
			"roleName":     t.RoleName,
		}

		log.Printf("[DEBUG] mongoatlas roles - adding role mapping: %v", mapping)
		s = append(s, mapping)
	}
	if err := d.Set("roles", s); err != nil {
		return err
	}

	return nil
}

func resourceDatabaseUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*MongoatlasClient)

	databaseuser := DatabaseUser{}

	if d.HasChange("roles") || d.HasChange("password") {
		var roles []Role
		rolesInterface := d.Get("roles").([]interface{})
		for i, roleInterface := range rolesInterface {
			log.Printf("%+v", i)
			rolesMap := roleInterface.(map[string]interface{})
			role := Role{
				DatabaseName: rolesMap["databaseName"].(string),
				RoleName:     rolesMap["roleName"].(string),
			}
			roles = append(roles, role)
		}
		databaseuser.Roles = &roles
		if d.HasChange("password") {
			databaseuser.Password = d.Get("password").(string)
		}
	}

	
	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(databaseuser)

	log.Printf("Sending %s \n", jsonpayload)

	databaseuser_req, err := client.Patch(fmt.Sprintf("groups/%s/databaseUsers/admin/%s",
		d.Get("groupId").(string),
		d.Get("username").(string),
	), jsonpayload)

	if err != nil {
		return err
	}

	if databaseuser_req.StatusCode == 200 {
		decoder := json.NewDecoder(databaseuser_req.Body)
		err = decoder.Decode(&databaseuser)
		if err != nil {
			return err
		}
	} else {
		body, err := ioutil.ReadAll(databaseuser_req.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Failed to patch database user: %s", string(body))
	}

	return resourceDatabaseUserRead(d, m)
}

func resourceDatabaseUserDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*MongoatlasClient)

	delete_response, err := client.Delete(fmt.Sprintf("groups/%s/databaseUsers/admin/%s",
		d.Get("groupId"),
		d.Get("username"),
	))

	if err != nil {
		return err
	}

	// should handle 404 when resource is not found. In that case, resource is already deleted
	if delete_response.StatusCode != 200 {
		return fmt.Errorf("Failed to delete the database user. Got status code %d", delete_response.StatusCode)
	}

	return nil
}
