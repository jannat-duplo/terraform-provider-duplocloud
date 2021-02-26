package duplosdk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// DuploAwsKmsKey represents an AWS KMS key for a Duplo tenant
type DuploAwsKmsKey struct {
	// NOTE: The TenantID field does not come from the backend - we synthesize it
	TenantID string `json:"-,omitempty"`

	Arn         string `json:"Arn,omitempty"`
	KeyName     string `json:"KeyName,omitempty"`
	KeyID       string `json:"KeyId,omitempty"`
	KeyArn      string `json:"KeyArn,omitempty"`
	Description string `json:"Description,omitempty"`
}

// TenantGetPlanKmsKeys retrieves a list of the AWS KMS keys for a tenant via the Duplo API.
func (c *Client) TenantGetPlanKmsKeys(tenantID string) (*[]DuploAwsKmsKey, error) {

	// Format the URL
	url := fmt.Sprintf("%s/subscriptions/%s/GetPlanKmsKeys", c.HostURL, tenantID)
	log.Printf("[TRACE] duplo-TenantGetPlanKmsKeys 1 ********: %s ", url)

	// Get the list from Duplo
	req2, _ := http.NewRequest("GET", url, nil)
	body, err := c.doRequest(req2)
	if err != nil {
		log.Printf("[TRACE] duplo-TenantGetPlanKmsKeys 2 ********: %s", err.Error())
		return nil, err
	}
	bodyString := string(body)
	log.Printf("[TRACE] duplo-TenantGetPlanKmsKeys 3 ********: %s", bodyString)

	// Return it as a list.
	list := []DuploAwsKmsKey{}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}
	log.Printf("[TRACE] duplo-TenantGetPlanKmsKeys 4 ********: %d items", len(list))
	for i := range list {
		list[i].TenantID = tenantID
		list[i].Arn = list[i].KeyArn
		list[i].Description = list[i].KeyName
	}
	return &list, nil
}

// TenantGetTenantKmsKey retrieves a tenant specific AWS KMS keys via the Duplo API.
func (c *Client) TenantGetTenantKmsKey(tenantID string) (*DuploAwsKmsKey, error) {

	// Format the URL
	url := fmt.Sprintf("%s/subscriptions/%s/GetTenantKmsKey", c.HostURL, tenantID)
	log.Printf("[TRACE] duplo-TenantGetTenantKmsKey 1 ********: %s ", url)

	// Get the key from Duplo
	req2, _ := http.NewRequest("GET", url, nil)
	body, err := c.doRequest(req2)
	if err != nil {
		log.Printf("[TRACE] duplo-TenantGetTenantKmsKey 2 ********: %s", err.Error())
		return nil, err
	}
	bodyString := string(body)
	log.Printf("[TRACE] duplo-TenantGetTenantKmsKey 3 ********: %s", bodyString)

	// Return it as an object.
	kms := DuploAwsKmsKey{}
	err = json.Unmarshal(body, &kms)
	if err != nil {
		return nil, err
	}
	log.Printf("[TRACE] duplo-TenantGetTenantKmsKey 4 ********")
	kms.TenantID = tenantID
	kms.KeyArn = kms.Arn
	kms.KeyName = kms.Description
	return &kms, nil
}

// TenantGetAllKmsKeys retrieves a list of all AWS KMS keys usable by a tenant via the Duplo API.
func (c *Client) TenantGetAllKmsKeys(tenantID string) ([]DuploAwsKmsKey, error) {

	// Tenant specific key
	tenantKey, err := c.TenantGetTenantKmsKey(tenantID)
	if err != nil {
		return nil, err
	}

	// Plan keys
	planKeys, err := c.TenantGetPlanKmsKeys(tenantID)
	if err != nil {
		return nil, err
	}

	// Make a list of all keys that have an ID
	allKeys := make([]DuploAwsKmsKey, 1, len(*planKeys)+1)
	allKeys[0] = *tenantKey
	for _, key := range *planKeys {
		if key.KeyID != "" {
			allKeys = append(allKeys, key)
		}
	}

	return allKeys, nil
}

// TenantGetKmsKeyByName retrieves a KMS key with a specific name, that is usable by a tenant via the Duplo API.
func (c *Client) TenantGetKmsKeyByName(tenantID string, keyName string) (*DuploAwsKmsKey, error) {

	// Get all keys.
	allKeys, err := c.TenantGetAllKmsKeys(tenantID)
	if err != nil {
		return nil, err
	}

	// Find and return the key with the specific name.
	for _, key := range allKeys {
		if key.KeyName == keyName {
			return &key, nil
		}
	}

	// No key was found.
	return nil, nil
}

// TenantGetKmsKeyByID retrieves a KMS key with a specific ID, that is usable by a tenant via the Duplo API.
func (c *Client) TenantGetKmsKeyByID(tenantID string, keyID string) (*DuploAwsKmsKey, error) {

	// Get all keys.
	allKeys, err := c.TenantGetAllKmsKeys(tenantID)
	if err != nil {
		return nil, err
	}

	// Find and return the key with the specific name.
	for _, key := range allKeys {
		if key.KeyID == keyID {
			return &key, nil
		}
	}

	// No key was found.
	return nil, nil
}