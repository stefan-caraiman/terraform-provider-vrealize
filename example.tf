# Configure the vRealize Provider
provider "vrealize" {
  username = "username"
  password = "password"
  tenant   = "vsphere.local"
  server   = "example.com"
}

resource "vrealize_machine" "machine" {
  count            = 1
  catalogItemRefId = "c94fa0c3-4aed-43ce-b7a6-4163a07e4cd6"
  tenantRef        = "vsphere.local"
  subTenantRef     = "f04f060d-73be-48a3-b82c-20cb98efd2d2"

  requestData = {
    key   = "provider-Cafe.Shim.VirtualMachine.Description"
    value = "Test API request"
  }

  requestData = {
    key   = "reasons"
    value = "Test reason"
  }

  connection {
    user     = "user"
    password = "password"
  }

  provisioner "file" {
    source      = "test.sh"
    destination = "/tmp/example.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/example.sh",
      "/tmp/example.sh",
    ]
  }
}
