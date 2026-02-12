variable "network_name" {
  description = "Docker network name to attach monitor to"
  type        = string
}

variable "nginx_host" {
  description = "Nginx container hostname (for URL construction)"
  type        = string
}

variable "monitor_image" {
  description = "cloudops-monitor Docker image name"
  type        = string
}