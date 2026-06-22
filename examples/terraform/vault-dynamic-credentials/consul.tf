resource "consul_acl_policy" "consul_cleanup" {
  name = "consul-cleanup"

  rules = <<-HCL
    service_prefix "" {
      policy = "write"
    }

    node_prefix "" {
      policy = "write"
    }
  HCL
}

