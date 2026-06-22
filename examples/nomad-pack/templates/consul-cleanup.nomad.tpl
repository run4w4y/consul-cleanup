job [[ .my.job_name | quote ]] {
  datacenters = [[ .my.datacenters | toStringList ]]
  type = "service"

  group "consul-cleanup" {
    count = [[ .my.count ]]

    network {
      mode = "host"

      port "http" {
        static = [[ .my.port ]]
      }
    }

    task "consul-cleanup" {
      driver = "docker"

      vault {
        policies = ["default", "consul-cleanup"]
      }

      env {
        PORT = "${NOMAD_PORT_http}"
      }

      config {
        image = "[[ .my.docker_image_registry_url | trimSuffix "/" ]]/[[ .my.docker_image_name ]]:[[ .my.docker_image_tag ]]"

        auth = {
          [[ if empty .my.docker_image_username | not -]]
          username = [[ .my.docker_image_username | quote ]]
          [[- end ]]
          [[ if empty .my.docker_image_password | not -]]
          password = [[ .my.docker_image_password | quote ]]
          [[- end ]]
        }

        command = "serve"
        ports = ["http"]
      }

      template {
        data = <<EOH
{{ with secret "secret/data/consul-cleanup/root-credentials" }}
CLEANUP_ACCESS_TOKEN={{ .Data.data.access_token }}
{{ end }}
        EOH
        env = true
        destination = "secrets/access_token.env"
      }

      template {
        data = <<EOH
CONSUL_HTTP_ADDR={{ env "NOMAD_IP_http" }}:8500
{{ with secret "consul/creds/consul-cleanup" }}
CONSUL_HTTP_TOKEN={{ .Data.token }}
{{ end }}
NOMAD_ADDR=http://{{ env "NOMAD_IP_http" }}:4646/
{{ with secret "nomad/creds/consul-cleanup" }}
NOMAD_TOKEN={{ .Data.secret_id }}
{{ end }}
        EOH
        env = true
        destination = "secrets/credentials.env"
      }

      resources {
        cpu = [[ .my.resources.cpu ]]
        memory = [[ .my.resources.memory ]]
      }

      service {
        name = [[ .my.job_name | quote ]]
        port = "http"

        check {
          type = "http"
          path = "/api/v1/health"
          interval = "2s"
          timeout = "2s"
        }
      }
    }

    service {
      name = [[ .my.job_name | quote ]]
      port = "http"
    }

    restart {
      attempts = 10
      interval = "5m"
      delay = "25s"
      mode = "delay"
    }
  }

  update {
    max_parallel = 1
    min_healthy_time = "10s"
    healthy_deadline = "3m"
    auto_revert = false
    canary = 0
  }
}

