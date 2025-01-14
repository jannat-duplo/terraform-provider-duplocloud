---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_oci_containerengine_node_pool Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_oci_containerengine_node_pool manages an OCI container node pool in Duplo.
---

# duplocloud_oci_containerengine_node_pool (Resource)

`duplocloud_oci_containerengine_node_pool` manages an OCI container node pool in Duplo.

## Example Usage

```terraform
resource "duplocloud_tenant" "myapp" {
  account_name = "myapp"
  plan_id      = "default"
}

resource "duplocloud_oci_containerengine_node_pool" "myOciNodePool" {
  tenant_id     = duplocloud_tenant.myapp.tenant_id
  name          = "tf-test"
  node_shape    = "VM.Standard2.1"
  node_image_id = "ocid1.image.oc1.ap-mumbai-1.aaaaaaaagosxifkwha6a6pi2fxx4idf3te3icdsf7z6jar2sxls6xycnehna"

  initial_node_labels {
    key   = "allocationtags"
    value = "test"
  }
  node_config_details {
    size = 1

    placement_configs {
      availability_domain = "uwFr:AP-MUMBAI-1-AD-1"
      subnet_id           = "ocid1.subnet.oc1.ap-mumbai-1.aaaaaaaasz36nwww2zygjn7arpuq4fbz3z22kn6adlalldvld3b5nu6afuxa"
    }

    freeform_tags = {
      CreatedBy = "duplo"
    }
  }

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The name of the node pool.
- **node_shape** (String) The name of the node shape of the nodes in the node pool.
- **node_shape_config** (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--node_shape_config))
- **tenant_id** (String) The GUID of the tenant that the Node Pool resource in Oracle Cloud Infrastructure Container Engine service be created in.

### Optional

- **defined_tags** (Map of String) Defined tags for this resource. Each key is predefined and scoped to a namespace.
- **freeform_tags** (Map of String) Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
- **initial_node_labels** (Block List) (see [below for nested schema](#nestedblock--initial_node_labels))
- **kubernetes_version** (String) The version of Kubernetes to install on the nodes in the node pool.
- **node_config_details** (Block List, Max: 1) (see [below for nested schema](#nestedblock--node_config_details))
- **node_image_id** (String)
- **node_image_name** (String)
- **node_metadata** (Map of String)
- **node_source_details** (Block List, Max: 1) (see [below for nested schema](#nestedblock--node_source_details))
- **quantity_per_subnet** (Number)
- **ssh_public_key** (String)
- **subnet_ids** (Set of String)
- **system_tags** (Map of String)
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **wait_until_ready** (Boolean) Whether or not to wait until oci node pool to be ready, after creation. Defaults to `true`.

### Read-Only

- **cluster_id** (String) The OCID of the cluster to which this node pool is attached.
- **compartment_id** (String) The OCID of the compartment in which the node pool exists.
- **id** (String) The ID of this resource.
- **node_pool_id** (String) The OCID of the node pool.
- **nodes** (Block List) (see [below for nested schema](#nestedblock--nodes))

<a id="nestedblock--node_shape_config"></a>
### Nested Schema for `node_shape_config`

Optional:

- **memory_in_gbs** (Number)
- **ocpus** (Number)


<a id="nestedblock--initial_node_labels"></a>
### Nested Schema for `initial_node_labels`

Optional:

- **key** (String)
- **value** (String)


<a id="nestedblock--node_config_details"></a>
### Nested Schema for `node_config_details`

Required:

- **placement_configs** (Block List, Min: 1) (see [below for nested schema](#nestedblock--node_config_details--placement_configs))
- **size** (Number)

Optional:

- **defined_tags** (Map of String)
- **freeform_tags** (Map of String)
- **is_pv_encryption_in_transit_enabled** (Boolean)
- **kms_key_id** (String)
- **nsg_ids** (Set of String)

<a id="nestedblock--node_config_details--placement_configs"></a>
### Nested Schema for `node_config_details.placement_configs`

Required:

- **availability_domain** (String)
- **subnet_id** (String)

Optional:

- **capacity_reservation_id** (String)



<a id="nestedblock--node_source_details"></a>
### Nested Schema for `node_source_details`

Required:

- **image_id** (String)
- **source_type** (String)

Optional:

- **boot_volume_size_in_gbs** (String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)


<a id="nestedblock--nodes"></a>
### Nested Schema for `nodes`

Read-Only:

- **availability_domain** (String)
- **defined_tags** (Map of String)
- **fault_domain** (String)
- **freeform_tags** (Map of String)
- **id** (String) The ID of this resource.
- **kubernetes_version** (String)
- **lifecycle_details** (String)
- **name** (String)
- **node_pool_id** (String)
- **private_ip** (String)
- **public_ip** (String)
- **state** (String)
- **subnet_id** (String)
- **system_tags** (Map of String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing OCI Container Engine Node Pool
#  - *TENANT_ID* is the tenant GUID
#  - *SHORTNAME* is the short name of the cluster
#
terraform import duplocloud_oci_containerengine_node_pool.myOciNodePool *TENANT_ID*/*SHORTNAME*
```
