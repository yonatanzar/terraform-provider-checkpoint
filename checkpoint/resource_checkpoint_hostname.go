package checkpoint

import (
	"fmt"
	checkpoint "github.com/Checkpoint/api_go_sdk/APIFiles"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)


func resourceHostname() *schema.Resource {
	return &schema.Resource{
		Create: createHostname,
		Read:   readHostname,
		Update: updateHostname,
		Delete: deleteHostname,
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Required: true,
				Description: "interface name",
			},
		},
	}
}

func hostnameParseSchemaToMap(d *schema.ResourceData) map[string]interface{} {
	hostnameMap := make(map[string]interface{})

	if val, ok := d.GetOk("name"); ok {
		hostnameMap["name"] = val.(string)
	}

	return hostnameMap
}

func createHostname(d *schema.ResourceData, m interface{}) error {
	log.Println("Enter createHostname...")
	client := m.(*checkpoint.ApiClient)
	payload := hostnameParseSchemaToMap(d)
	log.Println(payload)
	setPIRes, _ := client.ApiCall("set-hostname",payload,client.GetSessionID(),true,false)
	if !setPIRes.Success {
		return fmt.Errorf(setPIRes.ErrorMsg)
	}

	// Set Schema UID = Object key
	d.SetId(setPIRes.GetData()["name"].(string))

	log.Println("Exit createHostname...")
	return readHostname(d, m)
}

func readHostname(d *schema.ResourceData, m interface{}) error {
	log.Println("Enter readHostname...")
	client := m.(*checkpoint.ApiClient)
	payload := map[string]interface{}{}
	showHostnameRes, _ := client.ApiCall("show-hostname",payload,client.GetSessionID(),true,false)
	if !showHostnameRes.Success {
		// Handle deletion of an object from other clients - Object not found
		if objectNotFound(showHostnameRes.GetData()["code"].(string)) {
			d.SetId("") // Destroy resource
			return nil
		}
		return fmt.Errorf(showHostnameRes.ErrorMsg)
	}
	hostnameJson := showHostnameRes.GetData()
	log.Println(hostnameJson)

	_ = d.Set("name", hostnameJson["name"].(string))

	log.Println("Exit readHostname...")
	return nil
}

func updateHostname(d *schema.ResourceData, m interface{}) error {
	log.Println("Enter updateHostname...")
	client := m.(*checkpoint.ApiClient)
	payload := hostnameParseSchemaToMap(d)
	setNetworkRes, _ := client.ApiCall("set-hostname",payload,client.GetSessionID(),true,false)
	if !setNetworkRes.Success {
		return fmt.Errorf(setNetworkRes.ErrorMsg)
	}
	log.Println("Exit updateHostname...")
	return readHostname(d, m)
}

func deleteHostname(d *schema.ResourceData, m interface{}) error {
	log.Println("Enter deleteHostname...")
	d.SetId("") // Destroy resource
	log.Println("Exit deleteHostname...")
	return nil
}
