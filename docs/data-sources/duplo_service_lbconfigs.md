---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_duplo_service_lbconfigs Data Source - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_duplo_service_lbconfigs retrieves load balancer configuration(s) for container-based service(s) in Duplo.
  NOTE: For Amazon ECS services, see the duplocloud_ecs_services data source.
---

# duplocloud_duplo_service_lbconfigs (Data Source)

`duplocloud_duplo_service_lbconfigs` retrieves load balancer configuration(s) for container-based service(s) in Duplo.

NOTE: For Amazon ECS services, see the `duplocloud_ecs_services` data source.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **tenant_id** (String) The GUID of the tenant that hosts the duplo service.

### Optional

- **name** (String) The name of the duplo service.

### Read-Only

- **id** (String) The ID of this resource.
- **services** (List of Object) (see [below for nested schema](#nestedatt--services))

<a id="nestedatt--services"></a>
### Nested Schema for `services`

Read-Only:

- **arn** (String)
- **lbconfigs** (List of Object) (see [below for nested schema](#nestedobjatt--services--lbconfigs))
- **name** (String)
- **replication_controller_name** (String)
- **status** (String)

<a id="nestedobjatt--services--lbconfigs"></a>
### Nested Schema for `services.lbconfigs`

Read-Only:

- **backend_protocol_version** (String)
- **certificate_arn** (String)
- **cloud_name** (String)
- **dns_name** (String)
- **external_port** (Number)
- **external_traffic_policy** (String)
- **extra_selector_label** (List of Object) (see [below for nested schema](#nestedobjatt--services--lbconfigs--extra_selector_label))
- **frontend_ip** (String)
- **health_check_url** (String)
- **host_name** (String)
- **host_port** (Number)
- **index** (Number)
- **is_infra_deployment** (Boolean)
- **is_internal** (Boolean)
- **is_native** (Boolean)
- **lb_type** (Number)
- **name** (String)
- **port** (String)
- **protocol** (String)
- **replication_controller_name** (String)
- **target_group_arn** (String)

<a id="nestedobjatt--services--lbconfigs--extra_selector_label"></a>
### Nested Schema for `services.lbconfigs.extra_selector_label`

Read-Only:

- **key** (String)
- **value** (String)


