package haproxy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBackendResource(t *testing.T) {
	backendName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_backend" "%s" {
					name = "%s"
					balance = "roundrobin"
					mode = "http"
				}
				`, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_backend.%s", backendName), "name", backendName),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_backend.%s", backendName), "id"),
				),
			},

			// ImportState testing
			{
				ResourceName:      fmt.Sprintf("haproxy-pf_backend.%s", backendName),
				ImportState:       true,
				ImportStateVerify: true,
			},

			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_backend" "%s" {
					name = "%s"
					balance = "roundrobin"
					mode = "tcp"
				}
				`, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_backend.%s", backendName), "name", backendName),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_backend.%s", backendName), "mode", "tcp"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
