# --------------------------------------
# IAM
# --------------------------------------
resource "yandex_iam_service_account" "hb_s3_sa" {
  name = "hb-s3-sa"
}

resource "yandex_resourcemanager_folder_iam_member" "hb_s3_sa_editor" {
  folder_id = var.yc.folder_id
  role      = "storage.editor"
  member    = "serviceAccount:${yandex_iam_service_account.hb_s3_sa.id}"
}

resource "yandex_iam_service_account_static_access_key" "hb_s3_sa_static_key" {
  service_account_id = yandex_iam_service_account.hb_s3_sa.id
}

# --------------------------------------
# STORAGE
# --------------------------------------
# resource "yandex_storage_bucket" "hb_sitemaps" {
#   access_key    = yandex_iam_service_account_static_access_key.hb_s3_sa_static_key.access_key
#   secret_key    = yandex_iam_service_account_static_access_key.hb_s3_sa_static_key.secret_key
#   bucket        = var.yc.sitemaps_storage.bucket
#   max_size      = 0
#   force_destroy = false
# }

# resource "yandex_storage_bucket" "hb_storage" {
#   access_key    = yandex_iam_service_account_static_access_key.hb_s3_sa_static_key.access_key
#   secret_key    = yandex_iam_service_account_static_access_key.hb_s3_sa_static_key.secret_key
#   bucket        = var.yc.storage.bucket
#   max_size      = 0
#   force_destroy = false
# }
