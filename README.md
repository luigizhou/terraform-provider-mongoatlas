# DEPRECATED

This is a deprecated repo.

The code is over 3 years old and an official terraform provider for mongo atlas is already available.


<details><summary>Deprecated README</summary>
<p>

# Terraform Mongo Atlas Provider

## Implemented Features:
- [x] Containers
- [x] VPC peering 
- [x] Clusters 
- [x] DB Users
- [x] Group IP Whitelist
- [ ] Alert Configurations 

## Building: 
```
$ make
```

### Testing:
```
$ export MONGOATLAS_USERNAME=mongoatlas_username
$ export MONGOATLAS_APIKEY=xxxxxxxxx
$ export MONGOATLAS_VPCID=xxxxxx
$ export MONGOATLAS_AWSACCOUNTID=xxxxxx
$ export MONGOATLAS_GROUPID=xxxxxx
$ make test
```

a terraform-provider-mongoatlas binary will be built

## Configuration:

main.tf
```
provider "mongoatlas" {
    username = "user@example.com"
    apiKey = "XXXXXXXXXXXXXXXX"
}

resource "mongoatlas_vpc_peering" "test" {
    "groupId"= "0000000000000000000000"
    "vpcId"= "vpc-123456789"
    "awsAccountId" = "01234567890"
    "routeTableCidrBlock" = "10.0.0.0/24"
} 

resource "mongoatlas_cluster" "terratest1" {
    groupId = "xxxxxxxxxxxxxxxxxxxxx"
    name = "terratest1"
    backupEnabled = false
    instanceSizeName = "M10"
    providerName = "AWS"
    regionName = "EU_WEST_1" 
    diskSizeGB = 10
    replicationFactor = 3 
    encryptEBSVolume = false
    diskIOPS = 120
}


resource "mongoatlas_groupip_whitelist" "test_ipwhitelist" {
    cidrBlock = "1.2.3.4/32"
    groupId = "xxxxxxx"
    comment = "terraform test"
}


resource "mongoatlas_database_user" "test_user" {
    databaseName = "admin"
    username = "terratest"
    password = "test"
    groupId = "xxxx"
    roles = [
        {
            databaseName = "admin"
            roleName = "readWriteAnyDatabase"
        },
        {
            databaseName = "admin"
            roleName = "backup"
        }
    ]   
}

```


## KNOWN ISSUE
- MongoDB Version is currently not settable. It defaults to 3.2.8 and it's a value returned from MongoDB Atlas
- Only one provider available: AWS

## TODO:

- [x] Write acceptance tests for VPC Peering
- [x] Acceptance Tests for Cluster
- [x] Acceptance Tests for Group Ip Whitelist
- [x] Acceptance Tests for Database User
- [x] Must check how Container works with mongodb atlas
- [ ] Complete TODO throughout the written code
- [ ] Check if schema are correctly defined


</p>
</details>
