output "nginx_url" {
  description = "Nginx URL (監視対象)"
  value       = "http://localhost:8080"
}

output "monitor_container" {
  description = "cloudops-monitor コンテナ名"
  value       = docker_container.monitor.name
}