package haproxy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServerTemplateResource(t *testing.T) {
	backendName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	serverTemplateName1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	serverTemplateName2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			// TODO: remove providerConfig and use global env
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_backend" "%s" {
					name = "%s"
					balance = "leastconn"
					mode    = "tcp"
				}
				resource "haproxy-pf_server_template" "%s" {
					prefix = "%s"
					fqdn = "www.google.com"
					num_or_range = "1-3"
					port = 9999
					parent_name = "%s"
					check = "enabled"
					resolvers = "myresolver"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				`, backendName, backendName, serverTemplateName1, serverTemplateName1, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "prefix", serverTemplateName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "port", "9999"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "fqdn", "www.google.com"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_backend" "%s" {
					name = "%s"
					balance = "leastconn"
					mode    = "tcp"
				}
				resource "haproxy-pf_server_template" "%s" {
					prefix = "%s"
					fqdn = "www.google.com"
					num_or_range = "1-3"
					port = 8888
					parent_name = "%s"
					check = "disabled"
					resolvers = "myresolver"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				`, backendName, backendName, serverTemplateName1, serverTemplateName1, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "prefix", serverTemplateName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "port", "8888"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "check", "disabled"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_backend" "%s" {
					name = "%s"
					balance = "leastconn"
					mode    = "tcp"
				}
				resource "haproxy-pf_server_template" "%s" {
					prefix = "%s"
					fqdn = "www.google.com"
					num_or_range = "1-3"
					port = 8888
					parent_name = "%s"
					check = "disabled"
					resolvers = "myresolver"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				resource "haproxy-pf_server_template" "%s" {
					prefix = "%s"
					fqdn = "www.google.it"
					num_or_range = "1-5"
					port = 7777
					parent_name = "%s"
					check = "enabled"
					resolvers = "myresolver"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				`, backendName, backendName, serverTemplateName1, serverTemplateName1, backendName, backendName, serverTemplateName2, serverTemplateName2, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "prefix", serverTemplateName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "port", "8888"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "fqdn", "www.google.com"),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName1), "id"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName2), "name", serverTemplateName2),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName2), "port", "7777"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName2), "fqdn", "www.google.it"),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_server_template.%s", serverTemplateName2), "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
