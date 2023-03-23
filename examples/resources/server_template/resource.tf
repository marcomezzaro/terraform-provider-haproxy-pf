resource "haproxy-pf_server_template" "server_template" {
  prefix = "demo"
  fqdn = "my-server.domain.com"
  num_or_range = "1-5"
  port = 9955
  check = "enabled"
  resolvers = "resolver"
  parent_name = "parent-name"
}