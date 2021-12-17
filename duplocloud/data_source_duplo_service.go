package duplocloud

import (
	"context"
	"fmt"

	"log"
	"terraform-provider-duplocloud/duplosdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func duploServiceComputedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"other_docker_host_config": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"other_docker_config": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"extra_config": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"allocation_tags": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"volumes": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"commands": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"cloud": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"agent_platform": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"replicas": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"replicas_matching_asg_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"docker_image": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"lb_synced_deployment": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"any_host_allowed": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"cloud_creds_from_k8s_service_account": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"tags": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     KeyValueSchema(),
		},
	}
}

func dataSourceDuploService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDuploServiceRead,
		Schema:      duploServiceComputedSchema(),
	}
}

func dataSourceDuploServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[TRACE] dataSourceDuploServiceRead: start")

	// Parse the identifying attributes
	tenantID := d.Get("tenant_id").(string)
	name := d.Get("name").(string)

	// Get the object from Duplo, detecting a missing object
	c := m.(*duplosdk.Client)
	duplo, err := c.ReplicationControllerGet(tenantID, name)
	if err != nil {
		return diag.Errorf("Unable to read tenant %s service '%s': %s", tenantID, name, err)
	}
	if duplo == nil {
		return diag.Errorf("Unable to read tenant %s service '%s': not found", tenantID, name)
	}
	d.SetId(fmt.Sprintf("%s/%s", tenantID, name))

	// Apply TF state
	d.Set("volumes", duplo.Volumes)
	d.Set("lb_synced_deployment", duplo.IsLBSyncedDeployment)
	d.Set("any_host_allowed", duplo.IsAnyHostAllowed)
	d.Set("cloud_creds_from_k8s_service_account", duplo.IsCloudCredsFromK8sServiceAccount)
	d.Set("agent_platform", duplo.AgentPlatform)
	d.Set("replicas_matching_asg_name", duplo.ReplicasMatchingAsgName)
	d.Set("replicas", duplo.Replicas)
	d.Set("cloud", duplo.Cloud)
	d.Set("tags", keyValueToState("tags", duplo.Tags))

	// If we have a pod template, read data from it
	if duplo.Template != nil {
		d.Set("other_docker_host_config", duplo.Template.OtherDockerHostConfig)
		d.Set("other_docker_config", duplo.Template.OtherDockerConfig)
		d.Set("allocation_tags", duplo.Template.AllocationTags)
		d.Set("extra_config", duplo.Template.ExtraConfig)
		d.Set("commands", duplo.Template.Commands)

		// If there is at least one container, get the first docker image from it.
		if duplo.Template.Containers != nil && len(*duplo.Template.Containers) > 0 {
			d.Set("docker_image", (*duplo.Template.Containers)[0].Image)
		}
	}

	log.Printf("[TRACE] dataSourceDuploServiceRead: end")
	return nil
}
