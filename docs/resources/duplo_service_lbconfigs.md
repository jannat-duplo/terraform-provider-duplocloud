---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_duplo_service_lbconfigs Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_duplo_service_lbconfigs manages load balancer configuration(s) for a container-based service in Duplo.
  NOTE: For Amazon ECS services, see the duplocloud_ecs_service resource.
---

# duplocloud_duplo_service_lbconfigs (Resource)

`duplocloud_duplo_service_lbconfigs` manages load balancer configuration(s) for a container-based service in Duplo.

NOTE: For Amazon ECS services, see the `duplocloud_ecs_service` resource.

## Example Usage

```terraform
resource "duplocloud_tenant" "myapp" {
  account_name = "myapp"
  plan_id      = "default"
}

# Deploy NGINX using Duplo's native container agent, and configure a load balancer.
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **lbconfigs** (Block List, Min: 1) (see [below for nested schema](#nestedblock--lbconfigs))
- **replication_controller_name** (String) The name of the duplo service.
- **tenant_id** (String) The GUID of the tenant that hosts the duplo service.

### Optional

- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **wait_until_ready** (Boolean) Whether or not to wait until Duplo considers all of the load balancers ready Defaults to `true`.

### Read-Only

- **arn** (String) The load balancer ARN.
- **id** (String) The ID of this resource.
- **status** (String) The load balancer's current status.

<a id="nestedblock--lbconfigs"></a>
### Nested Schema for `lbconfigs`

Required:

- **external_port** (Number) The frontend port associated with this load balancer configuration.
- **lb_type** (Number) The numerical index of the type of load balancer configuration to create.
Should be one of:

   - `0` : ELB (Classic Load Balancer)
   - `1` : ALB (Application Load Balancer)
   - `2` : Health-check Only (No Load Balancer)
   - `3` : K8S Service w/ Cluster IP (No Load Balancer)
   - `4` : K8S Service w/ Node Port (No Load Balancer)
   - `5` : Azure Shared Application Gateway
   - `6` : NLB (Network Load Balancer)
   - `7` : Target Group Only
- **port** (String) The backend port associated with this load balancer configuration.
- **protocol** (String) The backend protocol associated with this load balancer configuration.

Optional:

- **certificate_arn** (String) The ARN of an ACM certificate to associate with this load balancer.  Only applicable for HTTPS.
- **external_traffic_policy** (String) Only for K8S Node Port (`lb_type = 4`) or load balancers in Kubernetes.  Set the kubernetes service `externalTrafficPolicy` attribute.
- **extra_selector_label** (Block List) Only for K8S services or load balancers in Kubernetes.  Sets an additional selector label to narrow which pods can receive traffic. (see [below for nested schema](#nestedblock--lbconfigs--extra_selector_label))
- **health_check_url** (String) The health check URL to associate with this load balancer configuration.
- **host_name** (String) (Azure Only) Set only if Azure Shared Application Gateway is used (`lb_type = 5`).
- **is_internal** (Boolean) Whether or not to create an internal load balancer.
- **is_native** (Boolean)

Read-Only:

- **backend_protocol_version** (String)
- **cloud_name** (String) The name of the cloud load balancer (if applicable).
- **dns_name** (String) The DNS name of the cloud load balancer (if applicable).
- **frontend_ip** (String)
- **host_port** (Number) The automatically assigned host port.
- **index** (Number) The load balancer Index.
- **is_infra_deployment** (Boolean)
- **name** (String) The name of the duplo service.
- **replication_controller_name** (String) The name of the duplo service.
- **target_group_arn** (String) The ARN of the Target Group to which to route traffic.

<a id="nestedblock--lbconfigs--extra_selector_label"></a>
### Nested Schema for `lbconfigs.extra_selector_label`

Required:

- **key** (String)
- **value** (String)



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing service's load balancer configurations
#  - *TENANT_ID* is the tenant GUID
#  - *NAME* is the name of the service
#
terraform import duplocloud_duplo_service_lbconfigs.myservice v2/subscriptions/*TENANT_ID*/ServiceLBConfigsV2/*NAME*
```
