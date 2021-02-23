package duplosdk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DuploEcacheInstance is a Duplo SDK object that represents an RDS instance
type DuploEcacheInstance struct {
	// NOTE: The TenantID field does not come from the backend - we synthesize it
	TenantID string `json:"-,omitempty"`

	// NOTE: The Name field does not come from the backend - we synthesize it
	Name string `json:"Name"`

	Identifier             string `json:"Identifier"`
	Arn                    string `json:"Arn"`
	CacheType              int    `json:"CacheType,omitempty"`
	Size                   string `json:"Size,omitempty"`
	Replicas               int    `json:"Replicas,omitempty"`
	EnableEncryptionAtRest bool   `json:"EnableEncryptionAtRest,omitempty"`
	InstanceStatus         string `json:"InstanceStatus,omitempty"`
}

// DuploEcacheInstanceSchema returns a Terraform resource schema for an ECS Service
func DuploEcacheInstanceSchema() *map[string]*schema.Schema {
	return &map[string]*schema.Schema{
		"tenant_id": {
			Type:     schema.TypeString,
			Optional: false,
			Required: true,
			ForceNew: true, //switch tenant
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"identifier": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"arn": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cache_type": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  0,
		},
		"size": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"replicas": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"enable_encryption_at_rest": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"instance_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

/*************************************************
 * API CALLS to duplo
 */

// EcacheInstanceCreate creates an ECS service via the Duplo API.
func (c *Client) EcacheInstanceCreate(tenantID string, duploObject *DuploEcacheInstance) (*DuploEcacheInstance, error) {
	return c.EcacheInstanceCreateOrUpdate(tenantID, duploObject, false)
}

// EcacheInstanceUpdate updates an ECS service via the Duplo API.
func (c *Client) EcacheInstanceUpdate(tenantID string, duploObject *DuploEcacheInstance) (*DuploEcacheInstance, error) {
	return c.EcacheInstanceCreateOrUpdate(tenantID, duploObject, true)
}

// EcacheInstanceCreateOrUpdate creates or updates an ECS service via the Duplo API.
func (c *Client) EcacheInstanceCreateOrUpdate(tenantID string, duploObject *DuploEcacheInstance, updating bool) (*DuploEcacheInstance, error) {

	// Build the request
	verb := "POST"
	if updating {
		verb = "PUT"
	}
	rqBody, err := json.Marshal(&duploObject)
	if err != nil {
		log.Printf("[TRACE] EcacheInstanceCreateOrUpdate 1 JSON gen : %s", err.Error())
		return nil, err
	}
	url := fmt.Sprintf("%s/v2/subscriptions/%s/ECacheDBInstance", c.HostURL, tenantID)
	log.Printf("[TRACE] EcacheInstanceCreate 2 : %s <= %s", url, rqBody)
	req, err := http.NewRequest(verb, url, strings.NewReader(string(rqBody)))
	if err != nil {
		log.Printf("[TRACE] EcacheInstanceCreateOrUpdate 3 HTTP builder : %s", err.Error())
		return nil, err
	}

	// Call the API and get the response
	body, err := c.doRequest(req)
	if err != nil {
		log.Printf("[TRACE] EcacheInstanceCreateOrUpdate 4 HTTP %s : %s", verb, err.Error())
		return nil, err
	}
	bodyString := string(body)
	log.Printf("[TRACE] EcacheInstanceCreateOrUpdate 4 HTTP RESPONSE : %s", bodyString)

	// Handle the response
	rpObject := DuploEcacheInstance{}
	if bodyString == "" {
		log.Printf("[TRACE] EcacheInstanceCreateOrUpdate 5 NO RESULT : %s", bodyString)
		return nil, err
	}
	err = json.Unmarshal(body, &rpObject)
	if err != nil {
		log.Printf("[TRACE] EcacheInstanceCreateOrUpdate 6 JSON parse : %s", err.Error())
		return nil, err
	}
	return &rpObject, nil
}

// EcacheInstanceDelete deletes an ECS service via the Duplo API.
func (c *Client) EcacheInstanceDelete(id string) (*DuploEcacheInstance, error) {
	idParts := strings.SplitN(id, "/", 5)
	tenantID := idParts[2]
	name := idParts[4]

	// Build the request
	url := fmt.Sprintf("%s/v2/subscriptions/%s/ECacheDBInstance/duplo-%s", c.HostURL, tenantID, name)
	log.Printf("[TRACE] EcacheInstanceGet 1 : %s", url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("[TRACE] EcacheInstanceGet 2 HTTP builder : %s", err.Error())
		return nil, err
	}

	// Call the API and get the response
	body, err := c.doRequest(req)
	bodyString := string(body)
	if err != nil {
		log.Printf("[TRACE] EcacheInstanceGet 3 HTTP DELETE : %s", err.Error())
		return nil, err
	}
	log.Printf("[TRACE] EcacheInstanceGet 4 HTTP RESPONSE : %s", bodyString)

	// Parse the response into a duplo object
	duploObject := DuploEcacheInstance{}
	if bodyString == "" {
		// tolerate an empty response from DELETE
		duploObject.Name = name
	} else {
		err = json.Unmarshal(body, &duploObject)
		if err != nil {
			log.Printf("[TRACE] EcacheInstanceGet 5 JSON PARSE : %s", bodyString)
			return nil, err
		}
	}

	// Fill in the tenant ID and return the object
	duploObject.TenantID = tenantID
	return &duploObject, nil
}

// EcacheInstanceGet retrieves an RDS instance via the Duplo API.
func (c *Client) EcacheInstanceGet(id string) (*DuploEcacheInstance, error) {
	idParts := strings.SplitN(id, "/", 5)
	tenantID := idParts[2]
	name := idParts[4]

	// Build the request
	url := fmt.Sprintf("%s/v2/subscriptions/%s/ECacheDBInstance/duplo-%s", c.HostURL, tenantID, name)
	log.Printf("[TRACE] EcacheInstanceGet 1 : %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[TRACE] EcacheInstanceGet 2 HTTP builder : %s", err.Error())
		return nil, err
	}

	// Call the API and get the response
	body, err := c.doRequest(req)
	if err != nil {
		log.Printf("[TRACE] EcacheInstanceGet 3 HTTP GET : %s", err.Error())
		return nil, err
	}
	bodyString := string(body)
	log.Printf("[TRACE] EcacheInstanceGet 4 HTTP RESPONSE : %s", bodyString)

	// Parse the response into a duplo object, detecting a missing object
	if bodyString == "null" {
		return nil, nil
	}
	duploObject := DuploEcacheInstance{}
	err = json.Unmarshal(body, &duploObject)
	if err != nil {
		log.Printf("[TRACE] EcacheInstanceGet 5 JSON PARSE : %s", bodyString)
		return nil, err
	}

	// Fill in the tenant ID and the name and return the object
	duploObject.TenantID = tenantID
	duploObject.Name = name
	return &duploObject, nil
}

/*************************************************
 * DATA CONVERSIONS to/from duplo/terraform
 */

// EcacheInstanceFromState converts resource data respresenting an RDS instance to a Duplo SDK object.
func EcacheInstanceFromState(d *schema.ResourceData) (*DuploEcacheInstance, error) {
	duploObject := new(DuploEcacheInstance)

	// First, convert things into simple scalars
	duploObject.Name = d.Get("name").(string)
	duploObject.Identifier = d.Get("identifier").(string)
	duploObject.Arn = d.Get("arn").(string)
	duploObject.CacheType = d.Get("cache_type").(int)
	duploObject.Size = d.Get("size").(string)
	duploObject.Replicas = d.Get("replicas").(int)
	duploObject.EnableEncryptionAtRest = d.Get("enable_encryption_at_rest").(bool)
	duploObject.InstanceStatus = d.Get("instance_status").(string)

	return duploObject, nil
}

// EcacheInstanceToState converts a Duplo SDK object respresenting an RDS instance to terraform resource data.
func EcacheInstanceToState(duploObject *DuploEcacheInstance, d *schema.ResourceData) map[string]interface{} {
	if duploObject == nil {
		return nil
	}
	jsonData, _ := json.Marshal(duploObject)
	log.Printf("[TRACE] duplo-EcacheInstanceToState ******** 1: INPUT <= %s ", jsonData)

	jo := make(map[string]interface{})

	// First, convert things into simple scalars
	jo["tenant_id"] = duploObject.TenantID
	jo["name"] = duploObject.Name
	jo["identifier"] = duploObject.Identifier
	jo["arn"] = duploObject.Arn
	jo["cache_type"] = duploObject.CacheType
	jo["size"] = duploObject.Size
	jo["replicas"] = duploObject.Replicas
	jo["enable_encryption_at_rest"] = duploObject.EnableEncryptionAtRest
	jo["instance_status"] = duploObject.InstanceStatus

	jsonData2, _ := json.Marshal(jo)
	log.Printf("[TRACE] duplo-EcacheInstanceToState ******** 2: OUTPUT => %s ", jsonData2)

	return jo
}
