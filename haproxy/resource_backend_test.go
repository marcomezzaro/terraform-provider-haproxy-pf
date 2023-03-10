package haproxy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBackendResource(t *testing.T) {
	backendName1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	backendName2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
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
				`, backendName1, backendName1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_backend.%s", backendName1), "name", backendName1),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_backend.%s", backendName1), "id"),
				),
			},

			// ImportState testing
			{
				ResourceName:      fmt.Sprintf("haproxy-pf_backend.%s", backendName1),
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
				`, backendName1, backendName1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_backend.%s", backendName1), "name", backendName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_backend.%s", backendName1), "mode", "tcp"),
				),
			},
			// add two resources at the same time
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_backend" "%s" {
					name = "%s"
					balance = "roundrobin"
					mode = "tcp"
				}
				resource "haproxy-pf_backend" "%s" {
					name = "%s"
					balance = "roundrobin"
					mode = "tcp"
				}
				`, backendName1, backendName1, backendName2, backendName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_backend.%s", backendName1), "name", backendName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_backend.%s", backendName1), "mode", "tcp"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_backend.%s", backendName2), "name", backendName2),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_backend.%s", backendName2), "mode", "tcp"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
