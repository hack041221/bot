terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.7.0"
    }
  }
  required_version = ">= 1.0"
}
