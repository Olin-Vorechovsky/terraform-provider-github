package github

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGithubTeamMembership() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGithubTeamMembershipRead,

		Schema: map[string]*schema.Schema{
			"teamslug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGithubTeamMembershipRead(d *schema.ResourceData, meta interface{}) error {
	teamslug := d.Get("teamslug").(string)

	username := d.Get("username").(string)
	log.Printf("[INFO] Refreshing GitHub team membership: %s:%s", teamslug, username)

	client := meta.(*Owner).v3client

	orgName := meta.(*Owner).name
	if configuredOrg := d.Get("organization").(string); configuredOrg != "" {
		orgName = configuredOrg
	}

	ctx := context.Background()

	membership, resp, err := client.Teams.GetTeamMembershipBySlug(ctx, orgName, teamslug, username)

	if err != nil {
		return err
	}

	d.SetId(buildTwoPartID(teamslug, username))

	d.Set("username", username)
	d.Set("role", membership.GetRole())
	d.Set("etag", resp.Header.Get("ETag"))
	return nil
}
