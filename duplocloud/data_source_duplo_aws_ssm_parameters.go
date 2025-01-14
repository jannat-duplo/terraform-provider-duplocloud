package duplocloud

import (
	"fmt"
	"log"
	"terraform-provider-duplocloud/duplosdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source listing SSM parameters
func dataSourceAwsSsmParameters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSsmParametersRead,

		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allowed_pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

/// READ resource
func dataSourceSsmParametersRead(d *schema.ResourceData, m interface{}) error {
	tenantID := d.Get("tenant_id").(string)

	log.Printf("[TRACE] dataSourceSsmParametersRead(%s): start", tenantID)

	// List the secrets from Duplo.
	c := m.(*duplosdk.Client)
	rp, err := c.SsmParameterList(tenantID)
	if err != nil {
		return fmt.Errorf("failed to list SSM parameters: %s", err)
	}
	d.SetId(tenantID)

	// Set the Terraform resource data
	list := make([]map[string]interface{}, 0, len(*rp))
	for _, ssmParam := range *rp {
		list = append(list, map[string]interface{}{
			"tenant_id":          tenantID,
			"name":               ssmParam.Name,
			"type":               ssmParam.Type,
			"key_id":             ssmParam.KeyId,
			"description":        ssmParam.Description,
			"allowed_pattern":    ssmParam.AllowedPattern,
			"last_modified_user": ssmParam.LastModifiedUser,
			"last_modified_date": ssmParam.LastModifiedDate,
		})
	}

	d.Set("parameters", list)

	log.Printf("[TRACE] dataSourceSsmParametersRead(%s): end", tenantID)
	return nil
}
