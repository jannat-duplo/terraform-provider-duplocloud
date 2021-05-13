package duplocloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"terraform-provider-duplocloud/duplosdk"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ecacheInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tenant_id": {
			Description: "The GUID of the tenant that the elasticache instance will be created in.",
			Type:        schema.TypeString,
			Optional:    false,
			Required:    true,
			ForceNew:    true, //switch tenant
		},
		"name": {
			Description: "The short name of the elasticache instance.  Duplo will add a prefix to the name.  You can retrieve the full name from the `identifier` attribute.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"identifier": {
			Description: "The full name of the elasticache instance.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"arn": {
			Description: "The ARN of the elasticache instance.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"endpoint": {
			Description: "The endpoint of the elasticache instance.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"host": {
			Description: "The DNS hostname of the elasticache instance.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"port": {
			Description: "The listening port of the elasticache instance.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"cache_type": {
			Description: "The numerical index of elasticache instance type.\n" +
				"Should be one of:\n\n" +
				"   - `0` : Redis\n" +
				"   - `1` : Memcache\n",
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  0,
		},
		"size": {
			Description: "The instance type of the elasticache instance.\n" +
				"See AWS documentation for the [available instance types](https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/CacheNodes.SupportedTypes.html).",
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"replicas": {
			Description: "The number of replicas to create.",
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     1,
		},
		"encryption_at_rest": {
			Description: "Enables encryption-at-rest.",
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},
		"encryption_in_transit": {
			Description: "Enables encryption-in-transit.",
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},
		"instance_status": {
			Description: "The status of the elasticache instance.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

// SCHEMA for resource crud
func resourceDuploEcacheInstance() *schema.Resource {
	return &schema.Resource{
		Description: "`duplocloud_ecache_instance` manages an ElastiCache instance in Duplo.",

		ReadContext:   resourceDuploEcacheInstanceRead,
		CreateContext: resourceDuploEcacheInstanceCreate,
		DeleteContext: resourceDuploEcacheInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: ecacheInstanceSchema(),
	}
}

/// READ resource
func resourceDuploEcacheInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Parse the identifying attributes
	id := d.Id()
	tenantID, name, err := parseECacheInstanceIdParts(id)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[TRACE] resourceDuploEcacheInstanceRead(%s, %s): start", tenantID, name)

	// Get the object from Duplo, detecting a missing object
	c := m.(*duplosdk.Client)
	duplo, err := c.EcacheInstanceGet(tenantID, name)
	if duplo == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	// Convert the object into Terraform resource data
	flattenEcacheInstance(duplo, d)

	log.Printf("[TRACE] resourceDuploEcacheInstanceRead(%s, %s): end", tenantID, name)
	return nil
}

/// CREATE resource
func resourceDuploEcacheInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error

	tenantID := d.Get("tenant_id").(string)

	log.Printf("[TRACE] resourceDuploEcacheInstanceCreate(%s): start", tenantID)

	duplo := expandEcacheInstance(d)
	duplo.Identifier = duplo.Name
	id := fmt.Sprintf("v2/subscriptions/%s/ECacheDBInstance/%s", tenantID, duplo.Name)

	// Post the object to Duplo
	c := m.(*duplosdk.Client)
	_, err = c.EcacheInstanceCreate(tenantID, duplo)
	if err != nil {
		return diag.Errorf("Error updating ECache instance '%s': %s", id, err)
	}
	d.SetId(id)

	// Wait up to 60 seconds for Duplo to be able to return the instance details.
	diags := waitForResourceToBePresentAfterCreate(ctx, d, "ECache instance", id, func() (interface{}, duplosdk.ClientError) {
		return c.EcacheInstanceGet(tenantID, duplo.Name)
	})
	if diags != nil {
		return diags
	}

	// Wait for the instance to become available.
	err = ecacheInstanceWaitUntilAvailable(ctx, c, tenantID, duplo.Name)
	if err != nil {
		return diag.Errorf("Error waiting for ECache instance '%s' to be available: %s", id, err)
	}

	diags = resourceDuploEcacheInstanceRead(ctx, d, m)
	log.Printf("[TRACE] resourceDuploEcacheInstanceCreate(%s, %s): end", tenantID, duplo.Name)
	return diags
}

/// DELETE resource
func resourceDuploEcacheInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Parse the identifying attributes
	id := d.Id()
	tenantID, name, err := parseECacheInstanceIdParts(id)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[TRACE] resourceDuploEcacheInstanceDelete(%s, %s): start", tenantID, name)

	// Delete the object from Duplo
	c := m.(*duplosdk.Client)
	err = c.EcacheInstanceDelete(tenantID, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait up to 60 seconds for Duplo to show the object as deleted.
	diag := waitForResourceToBeMissingAfterDelete(ctx, d, "ECache instance", id, func() (interface{}, duplosdk.ClientError) {
		return c.EcacheInstanceGet(tenantID, name)
	})
	log.Printf("[TRACE] resourceDuploEcacheInstanceDelete(%s, %s): end", tenantID, name)
	return diag
}

/*************************************************
 * DATA CONVERSIONS to/from duplo/terraform
 */

// expandEcacheInstance converts resource data respresenting an ECache instance to a Duplo SDK object.
func expandEcacheInstance(d *schema.ResourceData) *duplosdk.DuploEcacheInstance {
	return &duplosdk.DuploEcacheInstance{
		Name:                d.Get("name").(string),
		Identifier:          d.Get("identifier").(string),
		Arn:                 d.Get("arn").(string),
		Endpoint:            d.Get("endpoint").(string),
		CacheType:           d.Get("cache_type").(int),
		Size:                d.Get("size").(string),
		Replicas:            d.Get("replicas").(int),
		EncryptionAtRest:    d.Get("encryption_at_rest").(bool),
		EncryptionInTransit: d.Get("encryption_in_transit").(bool),
		InstanceStatus:      d.Get("instance_status").(string),
	}
}

// flattenEcacheInstance converts a Duplo SDK object respresenting an ECache instance to terraform resource data.
func flattenEcacheInstance(duplo *duplosdk.DuploEcacheInstance, d *schema.ResourceData) {
	d.Set("tenant_id", duplo.TenantID)
	d.Set("name", duplo.Name)
	d.Set("identifier", duplo.Identifier)
	d.Set("arn", duplo.Arn)
	d.Set("endpoint", duplo.Endpoint)
	if duplo.Endpoint != "" {
		uriParts := strings.SplitN(duplo.Endpoint, ":", 2)
		d.Set("host", uriParts[0])
		if len(uriParts) == 2 {
			port, _ := strconv.Atoi(uriParts[1])
			d.Set("port", port)
		}
	}
	d.Set("cache_type", duplo.CacheType)
	d.Set("size", duplo.Size)
	d.Set("replicas", duplo.Replicas)
	d.Set("encryption_at_rest", duplo.EncryptionAtRest)
	d.Set("encryption_in_transit", duplo.EncryptionInTransit)
	d.Set("instance_status", duplo.InstanceStatus)
}

// ecacheInstanceWaitUntilAvailable waits until an ECache instance is available.
//
// It should be usable both post-creation and post-modification.
func ecacheInstanceWaitUntilAvailable(ctx context.Context, c *duplosdk.Client, tenantID, name string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"processing", "creating", "modifying", "rebooting cluster nodes", "snapshotting"},
		Target:       []string{"available"},
		MinTimeout:   10 * time.Second,
		PollInterval: 30 * time.Second,
		Timeout:      20 * time.Minute,
		Refresh: func() (interface{}, string, error) {
			resp, err := c.EcacheInstanceGet(tenantID, name)
			if err != nil {
				return 0, "", err
			}
			if resp.InstanceStatus == "" {
				resp.InstanceStatus = "processing"
			}
			return resp, resp.InstanceStatus, nil
		},
	}
	log.Printf("[DEBUG] EcacheInstanceWaitUntilAvailable (%s, %s)", tenantID, name)
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func parseECacheInstanceIdParts(id string) (tenantID, name string, err error) {
	idParts := strings.SplitN(id, "/", 5)
	if len(idParts) == 5 {
		tenantID, name = idParts[2], idParts[4]
	} else {
		err = fmt.Errorf("invalid resource ID: %s", id)
	}
	return
}
