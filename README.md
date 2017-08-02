# vRealize
Terraform vRealize provider

## Introduction
This provider uses the https://github.com/stefan-caraiman/govrealize client library to implement creation, reading and deletion vRealize machines, based on the https://github.com/sky-uk/govrealize one

## Usage
Build the provider
```bash
go build -o terraform-provider-vrealize github.com/stefan-caraiman/terraform-provider-vrealize && cp terraform-provider-vrealize /usr/local/terraform/
```
Configure vRealize provider in a .tf file
```golang
provider "vrealize" {
	username = "xxxxxxxx"
	password = "xxxxxxxx"
	tenant =   "vsphere"
	server = "server.com"
}
```
Define resource
```golang
resource "vrealize_machine" "test" {
    catalogItemRefId = "xxxxxxx-xxxx-xxxxx-xxxx-xxxxxxxxx"
    tenantRef = "vsphere.local"
    subTenantRef = "xxxxxxx-xxxx-xxxxx-xxxx-xxxxxxxxx"
		requestData = {
			key = "provider-provisioningGroupId"
			value = "xxxxxxx-xxxx-xxxxx-xxxx-xxxxxxxxx"
		}
		requestData = {
			key = "provider-VirtualMachine.CPU.Count"
			value = 1
		}
		requestData = {
			key = "provider-VirtualMachine.Memory.Size"
			value = 1024
		}
		requestData = {
			key = "provider-Cafe.Shim.VirtualMachine.Description"
			value = "Test API request"
		}
		requestData = {
			key = "reasons"
			value = "Test reason"
		}
}
```
Execute provisioners
```golang
connection {
  user     = "user"
  password = "password"
}

provisioner "file" {
  source      = "test.sh"
  destination = "/tmp/terraform.sh"
}

provisioner "remote-exec" {
  inline = [
    "chmod +x /tmp/terraform.sh",
    "/tmp/terraform.sh",
  ]
}
```
