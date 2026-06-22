terraform {
  required_version = ">= 1.6.0"

  required_providers {
    consul = {
      source  = "hashicorp/consul"
      version = "~> 2.21"
    }
  }
}

provider "consul" {
  address = var.consul_address
  token   = var.consul_token
}

