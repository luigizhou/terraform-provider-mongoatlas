package main

import (
	"bytes"
	"net/http"
	// "log"
	httpdigest "github.com/ryanjdew/http-digest-auth-client"
)

type MongoatlasClient struct {
	Username string
	ApiKey   string
}

func (c *MongoatlasClient) Get(endpoint string) (*http.Response, error) {
	dh := &httpdigest.DigestHeaders{}
	dh, err := dh.Auth(c.Username, c.ApiKey, "https://cloud.mongodb.com/api/atlas/v1.0/")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://cloud.mongodb.com/api/atlas/v1.0/"+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	dh.ApplyAuth(req)
	return client.Do(req)
}

func (c *MongoatlasClient) Post(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	dh := &httpdigest.DigestHeaders{}
	dh, err := dh.Auth(c.Username, c.ApiKey, "https://cloud.mongodb.com/api/atlas/v1.0/")
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://cloud.mongodb.com/api/atlas/v1.0/"+endpoint, jsonpayload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	dh.ApplyAuth(req)
	return client.Do(req)
}

func (c *MongoatlasClient) Patch(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	dh := &httpdigest.DigestHeaders{}
	dh, err := dh.Auth(c.Username, c.ApiKey, "https://cloud.mongodb.com/api/atlas/v1.0/")
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", "https://cloud.mongodb.com/api/atlas/v1.0/"+endpoint, jsonpayload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	dh.ApplyAuth(req)
	return client.Do(req)
}

func (c *MongoatlasClient) Delete(endpoint string) (*http.Response, error) {
	dh := &httpdigest.DigestHeaders{}
	dh, err := dh.Auth(c.Username, c.ApiKey, "https://cloud.mongodb.com/api/atlas/v1.0/")
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "https://cloud.mongodb.com/api/atlas/v1.0/"+endpoint, nil)
	if err != nil {
		return nil, err
	}
	dh.ApplyAuth(req)
	return client.Do(req)
}

/*
func (c *MongoatlasClient) Put(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
    client := &http.Client{}
    req, err := http.NewRequest("PUT", "https://cloud.mongodb.com/api/atlas/v1.0/"+endpoint, jsonpayload)
    if err != nil {
        return nil, err
    }
    req.SetBasicAuth(c.Username, c.ApiKey)
    req.Header.Add("content-type", "application/json")
    return client.Do(req)
}

func (c *MongoatlasClient) PutOnly(endpoint string) (*http.Response, error) {
    client := &http.Client{}
    req, err := http.NewRequest("PUT", "https://cloud.mongodb.com/api/atlas/v1.0/"+endpoint, nil)
    if err != nil {
        return nil, err
    }
    req.SetBasicAuth(c.Username, c.ApiKey)
    return client.Do(req)
}

func (c *MongoatlasClient) Delete(endpoint string) (*http.Response, error) {
    client := &http.Client{}
    req, err := http.NewRequest("DELETE", "https://cloud.mongodb.com/api/atlas/v1.0/"+endpoint, nil)
    if err != nil {
        return nil, err
    }
    req.SetBasicAuth(c.Username, c.ApiKey)
    return client.Do(req)
}
*/
