module "bot" {
  source          = "./module/bot"
  namespace       = var.namespace
  slug            = var.slug
  docker-registry = local.docker-registry
  image           = {
    name = "${var.image}/app"
    tag  = var.image_version
  }
  name            = "bot"
  app             = {
    telegram_token = var.telegram_token
    aws_secret_key = var.aws_secret_key
    aws_access_key = var.aws_access_key
    aws_region     = var.aws_region
    aws_bucket     = var.aws_bucket
    job_url        = var.job_url
    state_url      = var.state_url
  }
}

module "downloader" {
  source          = "./module/downloader"
  namespace       = var.namespace
  slug            = var.slug
  docker-registry = local.docker-registry
  image           = {
    name = "${var.image}/downloader"
    tag  = var.image_version
  }
  name            = "downloader"
  app             = {
    telegram_token = var.telegram_token
    aws_secret_key = var.aws_secret_key
    aws_access_key = var.aws_access_key
    aws_region     = var.aws_region
    aws_bucket     = var.aws_bucket
    job_url        = var.job_url
    state_url      = var.state_url
    frame_url      = var.frame_url
    audio_url      = var.audio_url
  }
}
