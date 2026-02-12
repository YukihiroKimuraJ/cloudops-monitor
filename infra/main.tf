resource "docker_network" "monitor_net" {
  name = "cloudops-monitor-net"
}

module "nginx" {
  source = "./modules/nginx"

  network_name   = docker_network.monitor_net.name
  container_name = "monitor-nginx"
}

module "app" {
  source = "./modules/app"

  network_name  = docker_network.monitor_net.name
  nginx_host    = module.nginx.container_name
  monitor_image = var.monitor_image
}
