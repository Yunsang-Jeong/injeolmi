variable "context" {
  description = "The name of service."
  type        = map(string)
  default = {
    delimiter = "-"
  }
}

variable "service_name" {
  description = "The name of service."
  type        = string
  default     = "ingeolmi"
}

variable "gitlab_token" {
  description = "The token in gitlab."
  type        = string
  sensitive   = true
}

variable "gitlab_webhook_secret" {
  description = "The webhook secret in gitlab."
  type        = string
  sensitive   = true
}
