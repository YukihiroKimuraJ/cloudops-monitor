resource "docker_image" "monitor" {
  name = var.monitor_image

  build {
    context    = abspath("${path.module}/../../..")
    dockerfile  = "Dockerfile"
  }
}

resource "docker_container" "monitor" {
  name  = "cloudops-monitor"
  image = docker_image.monitor.image_id

  networks_advanced {
    name = var.network_name
  }

  entrypoint = ["sh", "-c"]
  command = [
    "while true; do ./cloudops-monitor -f /app/urls.txt -t 5 -c 5; sleep 60; done"
  ]

  upload {
    content = <<-EOT
      http://${var.nginx_host}:80/
      http://${var.nginx_host}:80/health
    EOT
    file    = "/app/urls.txt"
  }
}