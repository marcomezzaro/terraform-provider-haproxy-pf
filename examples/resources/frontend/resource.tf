resource "haproxy_frontend" "fe" {
name = "f1"
maxconn = 200
}