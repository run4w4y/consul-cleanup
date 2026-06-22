variable "job_name" {
  description = "The name to use as the job name which overrides using the pack name"
  type        = string
  default     = "consul-cleanup"
}

variable "datacenters" {
  description = "A list of datacenters in the region which are eligible for task placement."
  type        = list(string)
  default     = ["dc1"]
}

variable "count" {
  description = "The number of replicas to be deployed"
  type        = number
  default     = 1
}

variable "docker_image_registry_url" {
  description = "URL to the Docker registry"
  type        = string
  default     = "ghcr.io/run4w4y"
}

variable "docker_image_name" {
  description = "Docker image name"
  type        = string
  default     = "consul-cleanup"
}

variable "docker_image_tag" {
  description = "Docker image tag"
  type        = string
  default     = "latest"
}

variable "docker_image_username" {
  description = "Username to log in to the Docker registry with"
  type        = string
  default     = null
}

variable "docker_image_password" {
  description = "Password to log in to the Docker registry with"
  type        = string
  default     = null
}

variable "resources" {
  description = "Resources assigned to each replica of the job"
  type = object({
    cpu    = number
    memory = number
  })
  default = {
    cpu    = 200,
    memory = 256
  }
}

variable "port" {
  description = "Port to bind the HTTP server to"
  type        = number
  default     = 8989
}

