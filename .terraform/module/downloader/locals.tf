locals {
  container_name = var.name
  name           = "${var.name}-${var.slug}"
  labels = {
    app  = var.name
    slug = var.slug
  }
}
