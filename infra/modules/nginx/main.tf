resource "docker_image" "nginx" {
  name         = "nginx:alpine"
  keep_locally = true
}

resource "docker_container" "nginx" {
  name  = var.container_name
  image = docker_image.nginx.image_id

  networks_advanced {
    name = var.network_name
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