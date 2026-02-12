output "container_name" {
  description = "Nginx container name (hostname for other containers)"
  value       = docker_container.nginx.name
}