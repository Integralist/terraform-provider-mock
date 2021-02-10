package mock

import (
	// Documentation:
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema
	//
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource:
// A 'thing' you create, and then manage (update/delete) via terraform.
//
// Data Source:
// Data you can get access to and reference within your resources.

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"foo": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MOCK_FOO", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			// Naming format...
			//
			// Map key: <provider>_<thing>
			// File:    resource_<thing>.go
			//
			// NOTE:
			// The map key is what's documented as the 'thing' a consumer of this
			// provider would add to their terraform HCL file.
			// e.g. resource "mock_example" "my_own_name_for_this" {...}
			//
			"mock_example": resourceExample(),
		},
		// DataSource is a subset of Resource.
		DataSourcesMap: map[string]*schema.Resource{
			// Naming format...
			//
			// Map key: <provider>_<thing>
			// File:    data_source_<thing>.go
			//
			// NOTE:
			// The map key is what's documented as the 'thing' a consumer of this
			// provider would add to their terraform HCL file.
			// e.g. data_source "mock_example" "my_own_name_for_this" {...}
			//
			"mock_example": dataSourceExample(),
		},

		// To configure the provider (i.e. create an API client)
		// then pass ConfigureFunc. The interface{} value returned by this function
		// is stored and passed into the subsequent resources as the meta
		// parameter (this includes Data Sources as they are subsets of Resources).
		//
		// Documentation:
		// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#ConfigureFunc
		// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#ConfigureContextFunc
	}
}
