package mock

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceExample() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceExampleRead,
		Schema: map[string]*schema.Schema{
			"things": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceExampleRead(d *schema.ResourceData, m any) error {
	fmt.Printf("ResourceData: %+v\n", d)

	// We'll create some stubbed data as if we had made an actual API call and
	// got back a bunch of things.
	const jsonStream = `
		[
			{"id": 1, "version": "a"},
			{"id": 2, "version": "b"},
			{"id": 3, "version": "c"},
		]`

	// In order for us to store the returned data into terraform we need to
	// marshal the data into a format that matches what the schema expects.
	things := make([]map[string]any, 0)
	err := json.NewDecoder(strings.NewReader(jsonStream)).Decode(&things)
	if err != nil {
		return err
	}

	// We store our the data into terraform.
	if err := d.Set("things", things); err != nil {
		return err
	}

	// We don't have a unique ID for this data resource so we create one using a
	// timestamp format. I've seen people use a hash of the returned API data as
	// a unique key.
	//
	// NOTE:
	// That hashcode helper is no longer available! It has been moved into an
	// internal directory meaning it's not supposed to be consumed.
	//
	// Reference:
	// https://github.com/hashicorp/terraform-plugin-sdk/blob/master/internal/helper/hashcode/hashcode.go
	//
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
