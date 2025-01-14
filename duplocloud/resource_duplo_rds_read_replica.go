package duplocloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"terraform-provider-duplocloud/duplosdk"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func rdsReadReplicaSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tenant_id": {
			Description:  "The GUID of the tenant that the RDS read replica will be created in.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true, //switch tenant
			ValidateFunc: validation.IsUUID,
		},
		"name": {
			Description: "The short name of the RDS read replica.  Duplo will add a prefix to the name.  You can retrieve the full name from the `identifier` attribute.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 63-MAX_DUPLO_NO_HYPHEN_LENGTH),
				validation.StringMatch(regexp.MustCompile(`^[a-z0-9-]*$`), "Invalid RDS read replica name"),
				validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "RDS read replica name cannot end with a hyphen"),
				validation.StringDoesNotMatch(regexp.MustCompile(`--`), "RDS read replica name cannot contain two hyphens"),

				// NOTE: some validations are moot, because Duplo provides a prefix and suffix for the name:
				//
				// - First character must be a letter
				//
				// Because Duplo automatically prefixes names, it is impossible to break any of the rules in the above bulleted list.
			),
		},
		"cluster_identifier": {
			Description: "The full name of the RDS Cluster.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"size": {
			Description: "The type of the RDS read replica.\n" +
				"See AWS documentation for the [available instance types](https://aws.amazon.com/rds/instance-types/).",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^db\.`), "RDS read replica types must start with 'db.'"),
		},
		"identifier": {
			Description: "The full name of the RDS read replica.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"arn": {
			Description: "The ARN of the RDS read replica.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"endpoint": {
			Description: "The endpoint of the RDS read replica.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"host": {
			Description: "The DNS hostname of the RDS read replica.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"port": {
			Description: "The listening port of the RDS read replica.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"engine": {
			Description: "The numerical index of database engine to be used the for the RDS read replica.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"engine_version": {
			Description: "The database engine version to be used the for the RDS read replica.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"encrypt_storage": {
			Description: "Whether or not to encrypt the RDS instance storage.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"enable_logging": {
			Description: "Whether or not to enable the RDS instance logging. This setting is not applicable for document db cluster instance.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"multi_az": {
			Description: "Specifies if the RDS instance is multi-AZ.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"kms_key_id": {
			Description: "The globally unique identifier for the key.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"replica_status": {
			Description: "The current status of the RDS read replica.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

// SCHEMA for resource crud
func resourceDuploRdsReadReplica() *schema.Resource {
	return &schema.Resource{
		Description: "`duplocloud_rds_read_replica` manages an AWS RDS read replica in Duplo.",

		ReadContext:   resourceDuploRdsReadReplicaRead,
		CreateContext: resourceDuploRdsReadReplicaCreate,
		UpdateContext: resourceDuploRdsReadReplicaUpdate,
		DeleteContext: resourceDuploRdsReadReplicaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: rdsReadReplicaSchema(),
	}
}

/// READ resource
func resourceDuploRdsReadReplicaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[TRACE] resourceDuploRdsReadReplicaRead ******** start")

	// Get the object from Duplo, detecting a missing object
	c := m.(*duplosdk.Client)
	duplo, err := c.RdsInstanceGet(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if duplo == nil {
		d.SetId("")
		return nil
	}

	// Convert the object into Terraform resource data
	jo := rdsReadReplicaToState(duplo, d)
	for key := range jo {
		d.Set(key, jo[key])
	}
	d.SetId(fmt.Sprintf("v2/subscriptions/%s/RDSDBInstance/%s", duplo.TenantID, duplo.Name))

	log.Printf("[TRACE] resourceDuploRdsReadReplicaRead ******** end")
	return nil
}

/// CREATE resource
func resourceDuploRdsReadReplicaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[TRACE] resourceDuploRdsReadReplicaCreate ******** start")
	tenantID := d.Get("tenant_id").(string)
	// Convert the Terraform resource data into a Duplo object
	duplo, err := rdsReadReplicaFromState(d)
	if err != nil {
		return diag.Errorf("Internal error: %s", err)
	}

	// Post the object to Duplo
	c := m.(*duplosdk.Client)

	// Get RDS writer instance
	idParts := strings.SplitN(duplo.ClusterIdentifier, "-cluster", 2)
	name := strings.TrimPrefix(idParts[0], "duplo")
	duploWriterInstance, err := c.RdsInstanceGetByName(tenantID, name)
	if err != nil {
		return diag.FromErr(err)
	}
	duplo.Identifier = duplo.Name
	duplo.Engine = duploWriterInstance.Engine
	duplo.Cloud = duploWriterInstance.Cloud
	id := fmt.Sprintf("v2/subscriptions/%s/RDSDBInstance/%s", tenantID, duplo.Name)

	// Validate the RDS instance.
	errors := validateRdsInstance(duplo)
	if len(errors) > 0 {
		return errorsToDiagnostics(fmt.Sprintf("Cannot create RDS DB read replica: %s: ", id), errors)
	}

	_, err = c.RdsInstanceCreate(tenantID, duplo)
	if err != nil {
		return diag.Errorf("Error creating RDS DB read replica '%s': %s", id, err)
	}
	d.SetId(id)

	// Wait up to 60 seconds for Duplo to be able to return the instance details.
	diags := waitForResourceToBePresentAfterCreate(ctx, d, "RDS DB Read Replica", id, func() (interface{}, duplosdk.ClientError) {
		return c.RdsInstanceGet(id)
	})
	if diags != nil {
		return diags
	}

	// Wait for the instance to become available.
	err = rdsReadReplicaWaitUntilAvailable(ctx, c, id, d.Timeout("create"))
	if err != nil {
		return diag.Errorf("Error waiting for RDS DB read replica '%s' to be available: %s", id, err)
	}

	diags = resourceDuploRdsReadReplicaRead(ctx, d, m)

	log.Printf("[TRACE] resourceDuploRdsReadReplicaCreate ******** end")
	return diags
}

/// UPDATE resource
func resourceDuploRdsReadReplicaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

/// DELETE resource
func resourceDuploRdsReadReplicaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[TRACE] resourceDuploRdsReadReplicaDelete ******** start")

	// Delete the object from Duplo
	c := m.(*duplosdk.Client)
	id := d.Id()
	_, err := c.RdsInstanceDelete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	diags := waitForResourceToBeMissingAfterDelete(ctx, d, "RDS DB Read Replica", id, func() (interface{}, duplosdk.ClientError) {
		return c.RdsInstanceGet(id)
	})

	// Wait 1 more minute to deal with consistency issues.
	if diags == nil {
		time.Sleep(time.Minute)
	}

	log.Printf("[TRACE] resourceDuploRdsReadReplicaDelete ******** end")
	return diags
}

// RdsInstanceWaitUntilAvailable waits until an RDS instance is available.
//
// It should be usable both post-creation and post-modification.
func rdsReadReplicaWaitUntilAvailable(ctx context.Context, c *duplosdk.Client, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			"processing", "backing-up", "backtracking", "configuring-enhanced-monitoring", "configuring-iam-database-auth", "configuring-log-exports", "creating",
			"maintenance", "modifying", "moving-to-vpc", "rebooting", "renaming",
			"resetting-master-credentials", "starting", "stopping", "storage-optimization", "upgrading",
		},
		Target:       []string{"available"},
		MinTimeout:   10 * time.Second,
		PollInterval: 30 * time.Second,
		Timeout:      timeout,
		Refresh: func() (interface{}, string, error) {
			resp, err := c.RdsInstanceGet(id)
			if err != nil {
				return 0, "", err
			}
			if resp.InstanceStatus == "" {
				resp.InstanceStatus = "processing"
			}
			return resp, resp.InstanceStatus, nil
		},
	}
	log.Printf("[DEBUG] RdsInstanceWaitUntilAvailable (%s)", id)
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// RdsInstanceFromState converts resource data respresenting an RDS read replica to a Duplo SDK object.
func rdsReadReplicaFromState(d *schema.ResourceData) (*duplosdk.DuploRdsInstance, error) {
	duploObject := new(duplosdk.DuploRdsInstance)
	duploObject.Name = d.Get("name").(string)
	duploObject.Identifier = d.Get("name").(string)
	duploObject.SizeEx = d.Get("size").(string)
	duploObject.ClusterIdentifier = d.Get("cluster_identifier").(string)
	return duploObject, nil
}

// RdsInstanceToState converts a Duplo SDK object respresenting an RDS instance to terraform resource data.
func rdsReadReplicaToState(duploObject *duplosdk.DuploRdsInstance, d *schema.ResourceData) map[string]interface{} {
	if duploObject == nil {
		return nil
	}
	jsonData, _ := json.Marshal(duploObject)
	log.Printf("[TRACE] duplo-RdsInstanceToState ******** 1: INPUT <= %s ", jsonData)

	jo := make(map[string]interface{})

	// First, convert things into simple scalars
	jo["tenant_id"] = duploObject.TenantID
	jo["name"] = duploObject.Name
	jo["identifier"] = duploObject.Identifier
	jo["arn"] = duploObject.Arn
	jo["endpoint"] = duploObject.Endpoint
	if duploObject.Endpoint != "" {
		uriParts := strings.SplitN(duploObject.Endpoint, ":", 2)
		jo["host"] = uriParts[0]
		if len(uriParts) == 2 {
			jo["port"], _ = strconv.Atoi(uriParts[1])
		}
	}
	jo["engine"] = duploObject.Engine
	jo["engine_version"] = duploObject.EngineVersion
	jo["size"] = duploObject.SizeEx
	jo["encrypt_storage"] = duploObject.EncryptStorage
	jo["kms_key_id"] = duploObject.EncryptionKmsKeyId
	jo["enable_logging"] = duploObject.EnableLogging
	jo["multi_az"] = duploObject.MultiAZ
	jo["replica_status"] = duploObject.InstanceStatus

	jsonData2, _ := json.Marshal(jo)
	log.Printf("[TRACE] duplo-RdsInstanceToState ******** 2: OUTPUT => %s ", jsonData2)

	return jo
}
