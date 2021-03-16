package main

import (
	"fmt"
	"strings"
)

func validateDiskSizeGB(v interface{}, k string) (ws []string, errors []error) {
	value := v.(float64)
	if value < 10 || value > 16384 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid value between 10 and 16386",
			k))
		return
	}
	return
}

func validateProviderName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "AWS" {
		errors = append(errors, fmt.Errorf(
			"%q must be AWS as it is currently the only supported cloud provider",
			k))
		return
	}
	return
}

func validateNumShards(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 12 {
		errors = append(errors, fmt.Errorf(
			"%q must be a value between 1 and 12",
			k))
		return
	}
	return
}

func validateInstanceSizeName(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToUpper(v.(string))
	if !stringInSlice(value, []string{"M10", "M20", "M30", "M40", "M50", "M60", "M100"}) {
		errors = append(errors, fmt.Errorf(
			"%q is invalid. Valid values are M10,M20,M30,M40,M50,M60,M100",
			k))
		return
	}

	return

}

func validateRegionName(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToUpper(v.(string))

	if !stringInSlice(value, []string{"AP_SOUTHEAST_2", "EU_WEST_1", "US_EAST_1", "US_WEST_2"}) {
		errors = append(errors, fmt.Errorf(
			"%q is invalid. Valid values are AP_SOUTHEAST_2, EU_WEST1, US_EAST_1, US_WEST_2",
			k))
		return
	}
	return
}

func validateReplicationFactor(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	if value != 3 && value != 5 && value != 7 {
		errors = append(errors, fmt.Errorf(
			"%q must be 3, 5 or 7. Other values are invalid",
			k))
		return
	}
	return
}

func validateRoleName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if !stringInSlice(value, []string{"atlasAdmin", "backup", "clusterMonitor", "dbAdmin", "dbAdminAnyDatabase", "enableSharding", "read", "readAnyDatabase", "readWrite", "readWriteAnyDatabase"}) {
		errors = append(errors, fmt.Errorf(
			"%q is invalid. Valid values are atlasAdmin, backup, clusterMonitor, dbAdmin, dbAdminAnydatabase, enableSharding, read, readAnyDatabase, readWrite, readWriteAnyDatabase",
			k))
		return
	}
	return

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
