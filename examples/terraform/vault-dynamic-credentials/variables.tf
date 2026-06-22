variable "consul_address" {
  description = "Consul HTTP API address."
  type        = string
  default     = "127.0.0.1:8500"
}

variable "consul_token" {
  description = "Consul token used by Terraform and by Vault's Consul secrets backend."
  type        = string
  sensitive   = true
  default     = null
}

variable "nomad_address" {
  description = "Nomad HTTP API address."
  type        = string
  default     = "http://127.0.0.1:4646"
}

variable "nomad_token" {
  description = "Nomad token used by Terraform and by Vault's Nomad secrets backend."
  type        = string
  sensitive   = true
  default     = null
}

variable "vault_address" {
  description = "Vault HTTP API address."
  type        = string
  default     = "http://127.0.0.1:8200"
}

variable "vault_token" {
  description = "Vault token used to configure policies, secrets engines, and KV data."
  type        = string
  sensitive   = true
  default     = null
}

variable "vault_kv_mount" {
  description = "KV v2 mount used for the cleanup HTTP access token."
  type        = string
  default     = "secret"
}

variable "cleanup_access_token_length" {
  description = "Length of the generated bearer token for the HTTP oneshot endpoint."
  type        = number
  default     = 32
}

