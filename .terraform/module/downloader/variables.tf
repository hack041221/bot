variable "name" { type = string }
variable "namespace" { type = string }
variable "slug" { type = string }

variable "image" {
  type = object({
    name = string
    tag  = string
  })
}

variable "deployment" {
  type    = object({
    replicas = number,
    requests = object({
      cpu    = string
      memory = string
    })
    limits   = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    replicas = 1,
    requests = {
      cpu    = "2"
      memory = "512Mi"
    }
    limits   = {
      cpu    = "2"
      memory = "512Mi"
    }
  }
}

variable "app" {
  type = object({
    telegram_token = string
    aws_secret_key = string
    aws_access_key = string
    aws_region     = string
    aws_bucket     = string
    job_url        = string
    state_url      = string
    frame_url      = string
    audio_url      = string
  })
}
