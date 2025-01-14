---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_duplo_service_params Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_duplo_service_lbconfigs manages additional configuration for a container-based service in Duplo.
  NOTE: For Amazon ECS services, see the duplocloud_ecs_service resource.
---

# duplocloud_duplo_service_params (Resource)

`duplocloud_duplo_service_lbconfigs` manages additional configuration for a container-based service in Duplo.

NOTE: For Amazon ECS services, see the `duplocloud_ecs_service` resource.

## Example Usage

```terraform
resource "duplocloud_tenant" "myapp" {
  account_name = "myapp"
  plan_id      = "default"
}

# Deploy NGINX using Duplo's native container agent, and configure additional load balancer settings.
resource "duplocloud_duplo_service" "myservice" {
  tenant_id = duplocloud_tenant.myapp.tenant_id

  name           = "myservice"
  agent_platform = 0 # Duplo native container agent
  docker_image   = "nginx:latest"
  replicas       = 1
}
resource "duplocloud_duplo_service_lbconfigs" "myservice" {
  tenant_id                   = duplocloud_duplo_service.myservice.tenant_id
  replication_controller_name = duplocloud_duplo_service.myservice.name

  lbconfigs {
    external_port    = 80
    health_check_url = "/"
    is_native        = false
    lb_type          = 1 # Application load balancer
    port             = "80"
    protocol         = "http"
  }
}
resource "duplocloud_duplo_service_params" "myservice" {
  tenant_id = duplocloud_duplo_service_lbconfigs.myservice.tenant_id

  replication_controller_name = duplocloud_duplo_service_lbconfigs.myservice.replication_controller_name
  dns_prfx                    = "myservice"
  drop_invalid_headers        = true
  enable_access_logs          = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **replication_controller_name** (String) The name of the duplo service.
- **tenant_id** (String) The GUID of the tenant that hosts the duplo service.

### Optional

- **dns_prfx** (String) The DNS prefix to assign to this service's load balancer.
- **drop_invalid_headers** (Boolean) Whether or not to drop invalid HTTP headers received by the load balancer.
- **enable_access_logs** (Boolean) Whether or not to enable access logs.  When enabled, Duplo will send access logs to a centralized S3 bucket per plan
- **http_to_https_redirect** (Boolean) Whether or not to enable http to https redirection.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **webaclid** (String) The ARN of a web application firewall to associate this load balancer.

### Read-Only

- **id** (String) The ID of this resource.
- **load_balancer_arn** (String) The load balancer ARN.
- **load_balancer_name** (String) The load balancer name.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing service's parameters
#  - *TENANT_ID* is the tenant GUID
#  - *NAME* is the name of the service
#
terraform import duplocloud_duplo_service_params.myservice v2/subscriptions/*TENANT_ID*/ReplicationControllerParamsV2/*NAME*
```
