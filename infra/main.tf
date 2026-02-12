# ネットワーク（コンテナ同士が通信するため）
resource "docker_network" "monitor_net" {
  name = "cloudops-monitor-net"
}

# nginx（監視対象）
resource "docker_container" "nginx" {
  name  = "monitor-nginx"
  image = docker_image.nginx.image_id

  networks_advanced {
    name = docker_network.monitor_net.name
  }

  ports {
    internal = 80
    external = 8080
  }

  upload {
    content = file("${path.module}/nginx-health.conf")
    file    = "/etc/nginx/conf.d/default.conf"
  }
}

# nginxイメージ（公式イメージ）
resource "docker_image" "nginx" {
  name = "nginx:alpine"
}

# cloudops-monitorイメージ（ローカルビルド）
resource "docker_image" "monitor" {
  name = var.monitor_image

  build {
    context    = abspath("${path.module}/..")
    dockerfile = "Dockerfile"
  }
}

# cloudops-monitorコンテナ
resource "docker_container" "monitor" {
  name  = "cloudops-monitor"
  image = docker_image.monitor.image_id

  networks_advanced {
    name = docker_network.monitor_net.name
  }

  entrypoint = ["sh", "-c"]
  command = [
    "while true; do ./cloudops-monitor -f /app/urls.txt -t 5 -c 5; sleep 60; done"
  ]

  upload {
    content = <<-EOT
      http://monitor-nginx:80/
      http://monitor-nginx:80/health
    EOT
    file    = "/app/urls.txt"
  }
}