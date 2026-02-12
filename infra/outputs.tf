output "nginx_url" {
  description = "Nginx URL (monitor target)"
  value       = "http://localhost:8080"
}

output "monitor_container" {
  description = "cloudops-monitor container name"
  value       = module.app.container_name
}