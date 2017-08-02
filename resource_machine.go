package vrealize

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stefan-caraiman/govrealize"
)

func resourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceMachineCreate,
		Read:   resourceMachineRead,
		Delete: resourceMachineDelete,

		Schema: map[string]*schema.Schema{
			"catalog_item_ref_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant_ref": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sub_tenant_ref": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_data": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*govrealize.Client)

	entries := []govrealize.MachineRequestDataEntry{{
		Key: "provider-Cafe_Shim_VirtualMachine_NumberOfInstances",
		Value: govrealize.MachineRequestDataEntryValue{
			Type:  "integer",
			Value: "1",
		},
	}}

	requestDataCount := d.Get("request_data.#").(int)

	for i := 0; i < requestDataCount; i++ {
		prefix := fmt.Sprintf("request_data.%d", i)
		key := d.Get(prefix + ".key").(string)
		value := d.Get(prefix + ".value").(string)
		if key != "provider-Cafe.Shim.VirtualMachine.NumberOfInstances" {
			t := fmt.Sprintf("%T", value)
			mrde := govrealize.MachineRequestDataEntry{
				Key: key,
				Value: govrealize.MachineRequestDataEntryValue{
					Type:  t,
					Value: value,
				},
			}
			entries = append(entries, mrde)
		}
	}

	opts := &govrealize.MachineCreateRequest{
		Type: "CatalogItemRequest",
		CatalogItemRef: govrealize.MachineCreateRequestCatalogItemRef{
			ID: d.Get("catalog_item_ref_id").(string),
		},
		Organization: govrealize.MachineCreateRequestOrganization{
			TenantRef:    d.Get("tenant_ref").(string),
			SubtenantRef: d.Get("sub_tenant_ref").(string),
		},
		State:         "SUBMITTED",
		RequestNumber: 0,
		RequestData: govrealize.MachineCreateRequestRequestData{
			Entries: entries,
		},
	}

	log.Printf("[DEBUG] Machine create configuration: %#v", opts)

	machine, _, err := client.Machine.CreateMachine(opts)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating Machine: %s", err)
	}

	d.SetId(machine.ID)
	log.Printf("[INFO] Machine name: %s", machine.ID)

	if machine.IPAddress() != "" {
		ip := machine.IPAddress()
		name := machine.MachineName()
		d.Set("ip", ip)
		d.Set("name", name)
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": ip,
		})
	}

	return resourceMachineRead(d, m)
}

func resourceMachineRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*govrealize.Client)

	id := d.Id()

	if id == "" {
		return fmt.Errorf("invalid machine id: %v", id)
	}

	machine, resp, err := client.Machine.GetMachine(id)

	if err != nil {
		// If the record is somehow already destroyed, mark as
		// successfully gone
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("catalog_item_ref_id", machine.CatalogItem.ID)
	d.Set("tenant_ref", machine.Organization.TenantRef)
	d.Set("sub_tenant_ref", machine.Organization.SubtenantRef)
	d.Set("request_data", machine.ResourceData)
	if machine.IPAddress() != "" {
		ip := machine.IPAddress()
		name := machine.MachineName()
		d.Set("ip", ip)
		d.Set("name", name)
	}

	return nil
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*govrealize.Client)

	id := d.Id()

	if id == "" {
		return fmt.Errorf("invalid machine id: %v", id)
	}

	log.Printf("[INFO] Deleting machine: %s", d.Id())

	// Destroy the machine
	_, _, err := client.Machine.DestroyMachine(id)

	if err != nil {
		return fmt.Errorf("Error deleting machine: %s", err)
	}

	return nil
}
