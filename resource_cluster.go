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

type Cluster struct {
	Name              string            `json:"name,omitempty"`
	BackupEnabled     *bool             `json:"backupEnabled,omitempty"`
	MongoDBMajorVersion string 			`json:"mongoDBMajorVersion,omitempty"`
	MongoDBVersion    string            `json:"mongoDBVersion,omitempty"`
	MongoURI          string            `json:"mongoURI,omitempty"`
	MongoURIUpdated   string            `json:"mongoURIUpdated,omitempty"`
	NumShards         int               `json:"numShards,omitempty"`
	ReplicationFactor int               `json:"replicationFactor,omitempty"`
	ProviderSettings  *ProviderSettings `json:"providerSettings,omitempty"`
	DiskSizeGB        float64           `json:"diskSizeGB,omitempty"`
	StateName         string            `json:"stateName,omitempty"`
}

type ProviderSettings struct {
	ProviderName     string `json:"providerName,omitempty"`
	RegionName       string `json:"regionName,omitempty"`
	InstanceSizeName string `json:"instanceSizeName,omitempty"`
	DiskIOPS         int    `json:"diskIOPS,omitempty"`
	EncryptEBSVolume *bool  `json:"encryptEBSVolume,omitempty"`
}

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Update: resourceClusterUpdate,
		Read:   resourceClusterRead,
		Delete: resourceClusterDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"groupId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"backupEnabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"diskSizeGB": &schema.Schema{
				Type:         schema.TypeFloat,
				ValidateFunc: validateDiskSizeGB,
				Optional:     true,
				Computed:     true,
			},
			"mongoDBMajorVersion": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"mongoDBVersion": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"mongoURIUpdated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"numShards": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateNumShards,
				Computed:     true,
			},
			"providerName": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateProviderName,
				Required:     true,
			},
			"diskIOPS": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"encryptEBSVolume": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instanceSizeName": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateInstanceSizeName,
				Required:     true,
			},
			"regionName": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateRegionName,
				Required:     true,
			},
			"replicationFactor": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateReplicationFactor,
				Computed:     true,
			},
			"stateName": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func newCluster(d *schema.ResourceData) (*Cluster, error) {
	// <--- START PROVIDER SETTINGS
	providerSettings := &ProviderSettings{
		InstanceSizeName: strings.ToUpper(d.Get("instanceSizeName").(string)),
		ProviderName:     d.Get("providerName").(string),
		RegionName:       d.Get("regionName").(string),
	}

	if attr, ok := d.GetOk("diskIOPS"); ok {
		providerSettings.DiskIOPS = attr.(int)
	}

	if attr, ok := d.GetOk("encryptEBSVolume"); ok {
		encryptEBSVolume := new(bool)
		*encryptEBSVolume = attr.(bool)
		providerSettings.EncryptEBSVolume = encryptEBSVolume
	}

	// <--- END PROVIDER SETTINGS
	// <--- START CLUSTER SETTINGS
	backupEnabled := new(bool)
	*backupEnabled = d.Get("backupEnabled").(bool)

	cluster := &Cluster{
		Name:             d.Get("name").(string),
		BackupEnabled:    backupEnabled,
		ProviderSettings: providerSettings,
		MongoDBMajorVersion: d.Get("mongoDBMajorVersion").(string),
	}

	if attr, ok := d.GetOk("numShards"); ok {
		cluster.NumShards = attr.(int)
	}

	if attr, ok := d.GetOk("replicationFactor"); ok {
		cluster.ReplicationFactor = attr.(int)
	}
	// <--- END CLUSTER
	if attr, ok := d.GetOk("diskSizeGB"); ok {
		cluster.DiskSizeGB = attr.(float64)
	}

	return cluster, nil
}

func resourceClusterCreate(d *schema.ResourceData, m interface{}) error {
	// Instantiate the MongoatlasClient. Source in client.go, it handles all the HTTP REST API communication
	client := m.(*MongoatlasClient)

	cluster, err := newCluster(d)

	if err != nil {
		return err
	}

	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)

	enc.Encode(cluster)

	log.Printf("Sending %s \n", jsonpayload)

	// communication with API commence here
	cluster_req, err := client.Post(fmt.Sprintf("groups/%s/clusters",
		d.Get("groupId").(string),
	), jsonpayload)


	if cluster_req.StatusCode != 201 {
		body, err := ioutil.ReadAll(cluster_req.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Failed to create cluster. Got the following response body %s", string(body))
	}

	decoder := json.NewDecoder(cluster_req.Body)
	err = decoder.Decode(&cluster)
	if err != nil {
		return err
	}

	log.Printf("Received %s \n", cluster_req.Body)

	// if statuscode is different than 201, nothing is persisted in the tfstate file, but what happened on mongo atlas is not checked.
	// TODO: handle other status code
	
	// The following statement saves set the data that will be saved in the .tfstate file

	d.SetId(cluster.Name)
	d.Set("name", cluster.Name)
	d.Set("backupEnabled", cluster.BackupEnabled)
	d.Set("diskSizeGB", cluster.DiskSizeGB)
	d.Set("mongoDBMajorVersion", cluster.MongoDBMajorVersion)
	d.Set("mongoDBVersion", cluster.MongoDBVersion)
	d.Set("mongoURIUpdated", cluster.MongoURIUpdated)
	d.Set("numShards", cluster.NumShards)
	d.Set("providerName", cluster.ProviderSettings.ProviderName)
	d.Set("diskIOPS", cluster.ProviderSettings.DiskIOPS)
	d.Set("encryptEBSVolume", cluster.ProviderSettings.EncryptEBSVolume)
	d.Set("instanceSizeName", cluster.ProviderSettings.InstanceSizeName)
	d.Set("regionName", cluster.ProviderSettings.RegionName)
	d.Set("replicationFactor", cluster.ReplicationFactor)
	d.Set("stateName", cluster.StateName)

	return resourceClusterRead(d, m)

}

func resourceClusterRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*MongoatlasClient)

	log.Printf("%s", d.Get("name").(string))
	log.Printf("%s", d.Get("groupId").(string))

	cluster_req, err := client.Get(fmt.Sprintf("groups/%s/clusters/%s",
		d.Get("groupId").(string),
		d.Get("name").(string),
	))

	if cluster_req.StatusCode == 404 {
		log.Printf("[DEBUG] %s no longer exist, so we'll drop it from the state", d.Get("name").(string))
		d.SetId("")
		return nil
	}


	if err != nil {
		return err
	}

	var cluster Cluster

	decoder := json.NewDecoder(cluster_req.Body)
	err = decoder.Decode(&cluster)
	if err != nil {
		return err
	}
	log.Printf("Received %s \n", cluster_req.Body)

	d.Set("name", cluster.Name)
	d.Set("backupEnabled", *cluster.BackupEnabled)
	d.Set("diskSizeGB", cluster.DiskSizeGB)
	d.Set("mongoDBMajorVersion", cluster.MongoDBMajorVersion)
	d.Set("mongoDBVersion", cluster.MongoDBVersion)
	d.Set("mongoURIUpdated", cluster.MongoURIUpdated)
	d.Set("numShards", cluster.NumShards)
	d.Set("providerName", cluster.ProviderSettings.ProviderName)
	d.Set("diskIOPS", cluster.ProviderSettings.DiskIOPS)
	d.Set("encryptEBSVolume", *cluster.ProviderSettings.EncryptEBSVolume)
	d.Set("instanceSizeName", cluster.ProviderSettings.InstanceSizeName)
	d.Set("regionName", cluster.ProviderSettings.RegionName)
	d.Set("replicationFactor", cluster.ReplicationFactor)
	d.Set("stateName", cluster.StateName)
	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, m interface{}) error {
	setProvider := false
	client := m.(*MongoatlasClient)

	cluster := Cluster{}
	providerSettings := ProviderSettings{}

	if d.HasChange("diskSizeGB") {
		cluster.DiskSizeGB = d.Get("diskSizeGB").(float64)
	}

	if d.HasChange("backupEnabled") {
		backupEnabled := new(bool)
		*backupEnabled = d.Get("backupEnabled").(bool)
		cluster.BackupEnabled = backupEnabled
	}

	if d.HasChange("mongoDBMajorVersion") {
		cluster.MongoDBMajorVersion = d.Get("mongoDBMajorVersion").(string)
	}

	if d.HasChange("numShard") {
		cluster.NumShards = d.Get("numShard").(int)
	}

	if d.HasChange("replicationFactor") {
		cluster.ReplicationFactor = d.Get("replicationFactor").(int)
	}

	if d.HasChange("regionName") {
		if !setProvider {
			setProvider = true
		}
		providerSettings.RegionName = d.Get("regionName").(string)
	}

	if d.HasChange("diskIOPS") {
		providerSettings.DiskIOPS = d.Get("diskIOPS").(int)
	}

	if d.HasChange("encryptEBSVolume") {
		if !setProvider {
			setProvider = true
		}
		encryptEBSVolume := new(bool)
		*encryptEBSVolume = d.Get("encryptEBSVolume").(bool)
		providerSettings.EncryptEBSVolume = encryptEBSVolume
	}

	if d.HasChange("instanceSizeName") {
		if !setProvider {
			setProvider = true
		}
		providerSettings.InstanceSizeName = d.Get("instanceSizeName").(string)
	}

	if setProvider {
		providerSettings.ProviderName = "AWS"
		cluster.ProviderSettings = &providerSettings
	}

	var jsonbuffer []byte

	jsonpayload := bytes.NewBuffer(jsonbuffer)
	enc := json.NewEncoder(jsonpayload)
	enc.Encode(cluster)

	log.Printf("Sending %s \n", jsonpayload)

	cluster_req, err := client.Patch(fmt.Sprintf("groups/%s/clusters/%s",
		d.Get("groupId").(string),
		d.Get("name").(string),
	), jsonpayload)

	if err != nil {
		return err
	}

	if cluster_req.StatusCode == 200 {
		decoder := json.NewDecoder(cluster_req.Body)
		err = decoder.Decode(&cluster)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Failed to patch: %d", cluster_req.StatusCode)
	}

	return resourceClusterRead(d, m)

}

func resourceClusterDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*MongoatlasClient)

	delete_response, err := client.Delete(fmt.Sprintf("groups/%s/clusters/%s",
		d.Get("groupId").(string),
		d.Get("name").(string),
	))

	log.Printf("%s", d.Get("name").(string))

	if err != nil {
		return err
	}

	// should handle 404 when resource is not found. In that case, resource is already deleted
	if delete_response.StatusCode != 202 {
		return fmt.Errorf("Failed to delete the cluster. Got status code %d", delete_response.StatusCode)
	}
	return nil
}
