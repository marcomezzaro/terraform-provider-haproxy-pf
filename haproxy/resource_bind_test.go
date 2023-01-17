package haproxy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBindResource(t *testing.T) {
	frontendName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	bindName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_frontend" "%s" {
					name = "%s"
					maxconn = 2000
					mode = "http"
					http_connection_mode = ""
				}
				resource "haproxy-pf_bind" "%s" {
					name = "%s"
					address = "127.0.0.1"
					port = 9999
					parent_name = "%s"
					depends_on = [
						haproxy-pf_frontend.%s
					]
				}
				`, frontendName, frontendName, bindName, bindName, frontendName, frontendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName), "name", bindName),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName), "port", "9999"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName), "address", "127.0.0.1"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_bind.%s", bindName), "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      fmt.Sprintf("haproxy-pf_bind.%s", bindName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_frontend" "%s" {
					name = "%s"
					maxconn = 2000
					mode = "http"
					http_connection_mode = ""
				}
				resource "haproxy-pf_bind" "%s" {
					name = "%s"
					address = "127.0.0.1"
					port = 8888
					parent_name = "%s"
					depends_on = [
						haproxy-pf_frontend.%s
					]
				}
				`, frontendName, frontendName, bindName, bindName, frontendName, frontendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName), "name", bindName),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName), "port", "8888"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName), "address", "127.0.0.1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
