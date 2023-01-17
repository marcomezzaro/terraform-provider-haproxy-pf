resource "haproxy-pf_frontend" "fe" {
name = "f1"
maxconn = 200
}