resource "kubernetes_secret" "app" {
  metadata {
    name      = local.name
    namespace = var.namespace
  }

  data = {
    TELEGRAM_TOKEN        = var.app.telegram_token
    AWS_SECRET_ACCESS_KEY = var.app.aws_secret_access_key
    AWS_ACCESS_KEY_ID     = var.app.aws_access_key_id
    AWS_REGION            = var.app.aws_region
    AWS_BUCKET            = var.app.aws_bucket
    JOB_URL               = var.app.job_url
    STATE_URL             = var.app.state_url
  }
}
