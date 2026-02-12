variable "network_name" {
  description = "Docker network name to attach nginx to"
  type        = string
}

variable "container_name" {
  description = "Name for the nginx container"
  type        = string
  default     = "monitor-nginx"
}

variable "external_port" {
  description = "Port exposed on localhost"
  type        = number
  default     = 8080
}