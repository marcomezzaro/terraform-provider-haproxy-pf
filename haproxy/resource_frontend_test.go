package haproxy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFrontendResource(t *testing.T) {
	backendName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	frontendName1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	frontendName2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
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
				`, frontendName1, frontendName1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName1), "name", frontendName1),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName1), "id"),
				),
			},

			// ImportState testing
			{
				ResourceName:      fmt.Sprintf("haproxy-pf_frontend.%s", frontendName1),
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
				resource "haproxy-pf_frontend" "%s" {
					name = "%s"
					maxconn = 2000
					mode = "tcp"
					http_connection_mode = ""
					default_backend = "%s"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				`, backendName, backendName, frontendName1, frontendName1, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName1), "name", frontendName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName1), "mode", "tcp"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName1), "default_backend", backendName),
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
				resource "haproxy-pf_frontend" "%s" {
					name = "%s"
					maxconn = 2000
					mode = "tcp"
					http_connection_mode = ""
					default_backend = "%s"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				resource "haproxy-pf_frontend" "%s" {
					name = "%s"
					maxconn = 2000
					mode = "tcp"
					http_connection_mode = ""
					default_backend = "%s"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				`, backendName, backendName, frontendName1, frontendName1, backendName, backendName, frontendName2, frontendName2, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName1), "name", frontendName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName1), "mode", "tcp"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName1), "default_backend", backendName),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName2), "name", frontendName2),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName2), "mode", "tcp"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_frontend.%s", frontendName2), "default_backend", backendName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
