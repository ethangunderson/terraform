package heroku

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"heroku-go/v3"
)

func resourceHerokuSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceHerokuSpaceCreate,
		Read:   resourceHerokuSpaceRead,
		Delete: resourceHerokuSpaceDelete,
		Update: resourceHerokuSpaceUpdate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// "compliance": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	ForceNew: true,
			// 	Elem:     &schema.Schema{Type: schema.TypeString},
			// },

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceHerokuSpaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	opts := heroku.SpaceCreateOpts{}

	if v, ok := d.GetOk("name"); ok {
		vs := v.(string)
		log.Printf("[DEBUG] Space name: %s", vs)
		opts.Name = &vs
	}
	if v, ok := d.GetOk("organization"); ok {
		vs := v.(string)
		log.Printf("[DEBUG] Space organization: %s", vs)
		opts.Organization = &vs
	}
	if v, ok := d.GetOk("region"); ok {
		vs := v.(string)
		log.Printf("[DEBUG] Space region: %s", vs)
		opts.Region = &vs
	}
	// if v, ok := d.GetOk("compliance"); ok {
	// 	vs := v.(map[string]string)
	// 	log.Printf("[DEBUG] Space compliance: %s", vs)
	// 	opts.Compliance = &vs
	// }

	var spc *heroku.Space

	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		s, err := client.SpaceCreate(opts)
		if err != nil {
			log.Printf("[DEBUG] Error creating space: %s", err)
			return resource.NonRetryableError(err)
		}
		spc = s
		return nil
	})

	if err != nil {
		log.Printf("[DEBUG] Error creating space: %s", err)
		return err
	}

	d.Set("name", spc.Name)
	// d.Set("region", spc.Region)
	// d.Set("compliance", spc.Compliance)

	return resourceHerokuSpaceRead(d, meta)
}

func resourceHerokuSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	opts := heroku.SpaceUpdateOpts{}

	if v, ok := d.GetOk("name"); ok {
		vs := v.(string)
		log.Printf("[DEBUG] Space name: %s", vs)
		opts.Name = &vs
	}

	spc, err := client.SpaceUpdate(d.Get("id").(string), opts)
	if err != nil {
		return fmt.Errorf("Error updating space: %s", err)
	}

	d.Set("name", spc.Name)
	d.Set("region", spc.Region)
	// d.Set("compliance", spc.Compliance)

	return nil
}

func resourceHerokuSpaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	spc, err := client.SpaceInfo(d.Get("id").(string))
	if err != nil {
		return fmt.Errorf("Error retrieving space: %s", err)
	}

	d.Set("name", spc.Name)
	// d.Set("region", spc.Region)
	// d.Set("compliance", spc.Compliance)

	return nil
}

func resourceHerokuSpaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	log.Printf("[INFO] Deleting space: %s", d.Id())

	err := client.SpaceDelete(d.Get("id").(string))
	if err != nil {
		return fmt.Errorf("Error deleting space: %s", err)
	}

	return nil
}
