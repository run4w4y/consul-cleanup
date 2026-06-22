job_name    = "consul-cleanup"
datacenters = ["dc1"]

port  = 8989
count = 1

resources = {
  cpu    = 200
  memory = 256
}

