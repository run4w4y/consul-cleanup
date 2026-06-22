variable "consul_address" {
  description = "Consul HTTP API address."
  type        = string
  default     = "127.0.0.1:8500"
}

variable "consul_token" {
  description = "Consul token used by Terraform to create the cleanup policy and token."
  type        = string
  sensitive   = true
  default     = null
}

variable "token_description" {
  description = "Description for the generated Consul ACL token."
  type        = string
  default     = "consul-cleanup"
}

variable "token_local" {
  description = "Whether the generated Consul token is local to the datacenter."
  type        = bool
  default     = true
}

