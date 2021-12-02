variable "telegram_token" { type = string }
variable "aws_secret_access_key" { type = string }
variable "aws_access_key_id" { type = string }
variable "aws_region" { type = string }
variable "aws_bucket" { type = string }
variable "job_url" { type = string }
variable "state_url" { type = string }

variable "slug" { type = string }
variable "image" { type = string }
variable "image_version" { type = string }
variable "namespace" { type = string }

variable "docker_registry_password" { type = string }
variable "docker_registry_server" { type = string }
variable "docker_registry_username" { type = string }

variable "k8s_host" { type = string }
variable "k8s_ca_crt" { type = string }
variable "k8s_token" { type = string }
