terraform {
  required_version = ">= 1.6.0"

  required_providers {
    consul = {
      source  = "hashicorp/consul"
      version = "~> 2.21"
    }

    nomad = {
      source  = "hashicorp/nomad"
      version = "~> 2.3"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }

    vault = {
      source  = "hashicorp/vault"
      version = "~> 4.3"
    }
  }
}

provider "consul" {
  address = var.consul_address
  token   = var.consul_token
}

provider "nomad" {
  address   = var.nomad_address
  secret_id = var.nomad_token
}

provider "random" {}

provider "vault" {
  address = var.vault_address
  token   = var.vault_token
}

