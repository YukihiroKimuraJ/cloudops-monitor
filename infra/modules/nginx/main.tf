# ネットワーク（コンテナ同士が通信するため）
resource "docker_network" "monitor_net" {
  name = "cloudops-monitor-net"
}

# nginxモジュール（監視対象）
module "nginx" {
  source = "./modules/nginx"

  network_name   = docker_network.monitor_net.name
  container_name = "monitor-nginx"
}

# cloudops-monitorモジュール
module "app" {
  source = "./modules/app"

  network_name  = docker_network.monitor_net.name
  nginx_host    = module.nginx.container_name
  monitor_image = var.monitor_image
}