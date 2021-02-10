package mock

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExample() *schema.Resource {
	return &schema.Resource{
		// CRUD operations
		//
		// When terraform reads the state file, if a resource doesn't exist, then a
		// CREATE operation will be started (otherwise an UPDATE operation).
		//
		Create: resourceCreate,
		Read:   resourceRead,
		Update: resourceUpdate,
		Delete: resourceDelete,

		// Resource Schema
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// This attribute is 'required' meaning the consumer of this provider
			// will need to define the values expected when writing their terraform
			// HCL code.
			"foo": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bar": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceCreate(d *schema.ResourceData, m interface{}) error {
	// If this were a real provider, then we'd have an API client that would be
	// creating, reading, updating, deleting (i.e. CRUD) data.
	//
	// c := m.(*some.APIClient)

	foo := d.Get("foo").([]interface{})
	fmt.Printf("foo: %+v\n", foo)

	// We could do things with foo, e.g. loop over its elements and build up a
	// data structure to be used as input to an API client to create.
	//
	// Remember that "foo" was defined in the schema as "required" meaning the
	// consumer of this provider would need to provide the values associated with
	// the foo schema.

	// Let's pretend we made an API call and in the response we got back there
	// was an ID we could use as a unique key in our terraform state file so it
	// was able to track the resource.
	//
	// NOTE:
	// The mere existence of the ID and lack of error means terraform will
	// presume the CREATE operation was successful and store the provided "foo"
	// data the user provided into the local state file.
	d.SetId("123")

	// We do a READ operation to be sure we get the latest state stored locally.
	//
	// NOTE:
	// In this example, instead of doing a READ and then returning nil at the end
	// of the create function, we instead return the result of the READ. Meaning
	// if the READ fails, then we'll fail the CREATE. It's up to you whether
	// that's something you want to do as a READ could fail due to a network
	// issue and not necessarily mean there was an error with the CREATE
	// operation (which itself would have caused an error earlier and failed the
	// CREATE any way). This way we're ensuring the local state is up-to-date and
	// doesn't need a refresh.
	return resourceRead(d, m)
}

func resourceRead(d *schema.ResourceData, m interface{}) error {
	// We get the ID we set into terraform state after we had initially created
	// the resource.
	resourceID := d.Id()
	fmt.Println("resourceID:", resourceID)

	// Imagine we made an API call to get the latest version of the resource.
	//
	// In the returned data structure we might have a list of "bar" and we
	// could iterate over its elements and do something like:
	//
	// if err := d.Set("foo", <...>); err != nil {
	// 	return err
	// }
	//
	// Where <...> would be a flattened version of the list of bar data where
	// the returned type is `[]interface{}`. Flattened isn't a great description
	// for why we do this. It's more just a 'marshal' of the data into a format
	// that terraform can work with.

	// TODO: have an in-memory map that I can pull fake data from.

	return nil
}

func resourceUpdate(d *schema.ResourceData, m interface{}) error {
	// We get the ID we set into terraform state after we had initially created
	// the resource.
	resourceID := d.Id()
	fmt.Println("resourceID:", resourceID)

	if d.HasChange("foo") {
		foo := d.Get("foo").([]interface{})
		fmt.Printf("foo: %+v\n", foo)

		// Imagine we made an API call to update the given resource.
		//
		// We'd do this by iterating over the foo we pulled out of our terraform
		// state and coercing them into a type of map[string]interface{}
		//
		// e.g.
		//
		// for _, f := range foo {
		// 	i := f.(map[string]interface{})
		//
		// 	t := i["bar"].([]interface{})[0]
		// 	bar := t.(map[string]interface{})
		//
		//  ...constructing data structure to pass to API...
		//
		//  We might assign values to the data structure like:
		//
		//  bar["id"].(int)).
		// }

		// TODO: update "version" to be 2

		// Again, we do a READ operation to be sure we get the latest state stored locally.
		//
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceRead(d, m)
}

func resourceDelete(d *schema.ResourceData, m interface{}) error {
	// We get the ID we set into terraform state after we had initially created
	// the resource.
	resourceID := d.Id()
	fmt.Println("resourceID:", resourceID)

	// Imagine we use resourceID to issue a DELETE API call.

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return nil
}
