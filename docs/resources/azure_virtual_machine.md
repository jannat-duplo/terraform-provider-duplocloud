---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_azure_virtual_machine Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_azure_virtual_machine manages an Azure virtual machine in Duplo.
---

# duplocloud_azure_virtual_machine (Resource)

`duplocloud_azure_virtual_machine` manages an Azure virtual machine in Duplo.

## Example Usage

```terraform
resource "duplocloud_tenant" "myapp" {
  account_name = "myapp"
  plan_id      = "default"
}

resource "duplocloud_azure_virtual_machine" "az_vm" {
  tentenant_id  = duplocloud_tenant.myapp.tenant_id
  friendly_name = "test-vm"

  image_id       = "16.04-LTS;Canonical;UbuntuServer"
  capacity       = "Standard_D2s_v3"
  agent_platform = 0 # Duplo native container agent

  admin_username = "azureuser"
  admin_password = "Root!12345"
  disk_size_gb   = 50
  subnet_id      = "duploinfra-default"

  minion_tags {
    key   = "AllocationTags"
    value = "test-host"
  }

  tags {
    key   = "CreatedBy"
    value = "duplo"
  }

  tags {
    key   = "Owner"
    value = "duplo"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **admin_username** (String) Specifies the name of the local administrator account.
- **capacity** (String) Specifies the [size of the Virtual Machine](https://docs.microsoft.com/azure/virtual-machines/sizes-general). See also [Azure VM Naming Conventions](https://docs.microsoft.com/azure/virtual-machines/vm-naming-conventions).
- **friendly_name** (String) The short name of the host.
- **image_id** (String) The Image ID to use to create virtual machine. Provide id as semicolon separated string with sequence of sku, publisher and offer. For example, 16.04-LTS;Canonical;UbuntuServe
- **subnet_id** (String) Subnet ID which should be associated with the Virtual Machine.
- **tenant_id** (String) The GUID of the tenant that the host will be created in.

### Optional

- **admin_password** (String, Sensitive) The password associated with the local administrator account.
- **agent_platform** (Number) The numeric ID of the container agent pool that this host is added to. Defaults to `0`.
- **allocated_public_ip** (Boolean) Whether or not to allocate a public IP. Defaults to `false`.
- **base64_user_data** (String) Base64 encoded user data to associated with the host.
- **disk_size_gb** (Number) Specifies the size of the OS Disk in gigabytes
- **enable_log_analytics** (Boolean) Enable log analytics on virtual machine. Defaults to `false`.
- **encrypt_disk** (Boolean) Defaults to `false`.
- **is_minion** (Boolean) Defaults to `true`.
- **join_domain** (Boolean) Join a Windows Server virtual machine to an Azure Active Directory Domain Services. Defaults to `false`.
- **minion_tags** (Block List) A map of tags to assign to the resource. Example - `AllocationTags` can be passed as tag key with any value. (see [below for nested schema](#nestedblock--minion_tags))
- **tags** (Block List) (see [below for nested schema](#nestedblock--tags))
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **volume** (Block List) (see [below for nested schema](#nestedblock--volume))
- **wait_until_ready** (Boolean) Whether or not to wait until azure virtual machine to be ready, after creation. Defaults to `true`.

### Read-Only

- **id** (String) The ID of this resource.
- **instance_id** (String) The Azure Virtual Machine ID of the host.
- **status** (String) The current status of the host.
- **user_account** (String) The name of the tenant that the host will be created in.

<a id="nestedblock--minion_tags"></a>
### Nested Schema for `minion_tags`

Required:

- **key** (String)
- **value** (String)


<a id="nestedblock--tags"></a>
### Nested Schema for `tags`

Required:

- **key** (String)
- **value** (String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)


<a id="nestedblock--volume"></a>
### Nested Schema for `volume`

Optional:

- **iops** (Number)
- **name** (String)
- **size** (Number)
- **volume_id** (String)
- **volume_type** (String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing Azure Virtual Machine
#  - *TENANT_ID* is the tenant GUID
#  - *SHORT_NAME* is the short name of the Azure Virtual Machine
#
terraform import duplocloud_azure_virtual_machine.myvm *TENANT_ID*/*SHORT_NAME*
```