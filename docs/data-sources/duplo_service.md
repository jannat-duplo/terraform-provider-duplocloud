---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_duplo_service Data Source - terraform-provider-duplocloud"
subcategory: ""
description: |-
  
---

# duplocloud_duplo_service (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String)
- **tenant_id** (String)

### Read-Only

- **agent_platform** (Number)
- **allocation_tags** (String)
- **any_host_allowed** (Boolean)
- **cloud** (Number)
- **cloud_creds_from_k8s_service_account** (Boolean)
- **commands** (List of String)
- **docker_image** (String)
- **extra_config** (String)
- **force_stateful_set** (Boolean)
- **hpa_specs** (String)
- **id** (String) The ID of this resource.
- **is_daemonset** (Boolean)
- **lb_synced_deployment** (Boolean)
- **other_docker_config** (String)
- **other_docker_host_config** (String)
- **replica_collocation_allowed** (Boolean)
- **replicas** (Number)
- **replicas_matching_asg_name** (String)
- **tags** (List of Object) (see [below for nested schema](#nestedatt--tags))
- **volumes** (String)

<a id="nestedatt--tags"></a>
### Nested Schema for `tags`

Read-Only:

- **key** (String)
- **value** (String)


