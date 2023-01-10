package haproxy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFrontendResource(t *testing.T) {
	backendName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	frontendName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy_frontend" "%s" {
					name = "%s"
					maxconn = 2000
					mode = "http"
					http_connection_mode = ""
				}
				`, frontendName, frontendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy_frontend.%s", frontendName), "name", frontendName),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy_frontend.%s", frontendName), "id"),
				),
			},

			// ImportState testing
			{
				ResourceName:      fmt.Sprintf("haproxy_frontend.%s", frontendName),
				ImportState:       true,
				ImportStateVerify: true,
			},

			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy_backend" "%s" {
					name = "%s"
					balance = "roundrobin"
					mode = "tcp"
				}
				resource "haproxy_frontend" "%s" {
					name = "%s"
					maxconn = 2000
					mode = "tcp"
					http_connection_mode = ""
					default_backend = "%s"
					depends_on = [
						haproxy_backend.%s
					]
				}
				`, backendName, backendName, frontendName, frontendName, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy_frontend.%s", frontendName), "name", frontendName),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy_frontend.%s", frontendName), "mode", "tcp"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy_frontend.%s", frontendName), "default_backend", backendName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
