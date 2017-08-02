# Configure the vRealize Provider
provider "vrealize" {
  username = "username"
  password = "password"
  tenant   = "vsphere.local"
  server   = "example.com"
}

resource "vrealize_machine" "machine" {
  count            = 1
  catalog_item_ref_id = "c94fa0c3-4aed-43ce-b7a6-4163a07e4cd6"
  tenant_ref        = "vsphere.local"
  sub_tenant_ref     = "f04f060d-73be-48a3-b82c-20cb98efd2d2"

  request_data = {
    key   = "provider-Cafe.Shim.VirtualMachine.Description"
    value = "Test API request"
  }

  request_data = {
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
