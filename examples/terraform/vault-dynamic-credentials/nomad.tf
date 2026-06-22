resource "nomad_acl_policy" "allocation_observer" {
  name = "allocation-observer"

  rules_hcl = <<-HCL
    namespace "*" {
      policy       = "read"
      capabilities = ["read-job"]
    }
  HCL
}

