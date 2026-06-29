terraform {
  required_version = ">= 1.6"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
  backend "gcs" {
    bucket = "biomech-tf-state"
    prefix = "desertacia"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# --- Service enablement ---
resource "google_project_service" "services" {
  for_each = toset([
    "cloudrun.googleapis.com",
    "sqladmin.googleapis.com",
    "storage.googleapis.com",
    "artifactregistry.googleapis.com",
    "iam.googleapis.com",
    "cloudbuild.googleapis.com",
    "secretmanager.googleapis.com",
  ])
  service            = each.value
  disable_on_destroy = false
}

# --- Cloud SQL (PostgreSQL) ---
resource "google_sql_database_instance" "main" {
  name             = "biomech-db"
  database_version = "POSTGRES_17"
  region           = var.region
  depends_on       = [google_project_service.services]

  settings {
    tier              = "db-custom-2-7680"
    disk_size         = 20
    disk_type         = "SSD"
    disk_autoresize   = true
    disk_autoresize_limit = 200
    availability_type = "ZONAL"

    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_id
      require_ssl     = true
    }

    backup_configuration {
      enabled                        = true
      start_time                     = "03:00"
      point_in_time_recovery_enabled = true
      transaction_log_retention_days = 7
      backup_retention_settings {
        retained_backups = 30
        retention_unit   = "COUNT"
      }
    }
  }

  deletion_protection = false
}

resource "google_sql_database" "db" {
  name     = "desertacia"
  instance = google_sql_database_instance.main.name
}

resource "google_sql_user" "app" {
  name     = "app"
  instance = google_sql_database_instance.main.name
  password = random_password.db_password.result
}

resource "random_password" "db_password" {
  length  = 24
  special = false
}

# --- Artifact Registry ---
resource "google_artifact_registry_repository" "main" {
  location      = var.region
  repository_id = "cloud-run-source-deploy"
  format        = "DOCKER"
  depends_on    = [google_project_service.services]
}

# --- Cloud Run: API ---
resource "google_cloud_run_v2_service" "api" {
  name     = "biomech-api"
  location = var.region
  depends_on = [google_project_service.services]

  template {
    scaling {
      min_instance_count = 1
      max_instance_count = 10
    }
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/cloud-run-source-deploy/api:latest"
      env {
        name  = "SERVER_PORT"
        value = "8080"
      }
      env {
        name  = "DATABASE_URL"
        value = "postgres://${google_sql_user.app.name}:${random_password.db_password.result}@localhost:5432/${google_sql_database.db.name}?sslmode=disable"
      }
      env {
        name  = "ML_SERVICE_URL"
        value = "https://${google_cloud_run_v2_service.ml.uri}"
      }
      env {
        name  = "INTERNAL_API_KEY"
        value = random_password.internal_api_key.result
      }
      env {
        name  = "FIREBASE_CREDENTIALS"
        value = "/secrets/firebase-service-account.json"
      }
      env {
        name  = "UPLOADS_DIR"
        value = "/data/uploads"
      }
      env {
        name  = "SERVER_READ_TIMEOUT"
        value = "30s"
      }
      env {
        name  = "SERVER_WRITE_TIMEOUT"
        value = "60s"
      }
      ports {
        container_port = 8080
      }
      resources {
        limits = {
          cpu    = "2"
          memory = "1Gi"
        }
      }
      startup_probe {
        http_get {
          path = "/health"
        }
        initial_delay_seconds = 10
        period_seconds        = 10
        failure_threshold     = 6
      }
      liveness_probe {
        http_get {
          path = "/health"
        }
        period_seconds    = 30
        failure_threshold = 3
      }
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].containers[0].image,
    ]
  }
}

# --- Cloud Run: ML Service ---
resource "google_cloud_run_v2_service" "ml" {
  name     = "biomech-ml"
  location = var.region
  depends_on = [google_project_service.services]

  template {
    scaling {
      min_instance_count = 0
      max_instance_count = 3
    }
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/cloud-run-source-deploy/ml:latest"
      env {
        name  = "DATABASE_URL"
        value = "postgresql://${google_sql_user.app.name}:${random_password.db_password.result}@localhost:5432/${google_sql_database.db.name}"
      }
      env {
        name  = "GCS_BUCKET"
        value = var.project_id
      }
      env {
        name  = "MODELS_DIR"
        value = "/tmp/models"
      }
      env {
        name  = "BACKEND_CALLBACK_URL"
        value = "https://${google_cloud_run_v2_service.api.uri}/api/v1/training/jobs"
      }
      env {
        name  = "BACKEND_API_KEY"
        value = random_password.internal_api_key.result
      }
      env {
        name  = "LOG_LEVEL"
        value = "INFO"
      }
      ports {
        container_port = 8080
      }
      resources {
        limits = {
          cpu    = "2"
          memory = "4Gi"
        }
      }
      startup_probe {
        http_get {
          path = "/health"
        }
        initial_delay_seconds = 30
        period_seconds        = 10
        failure_threshold     = 6
      }
      liveness_probe {
        http_get {
          path = "/health"
        }
        period_seconds    = 30
        failure_threshold = 3
      }
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].containers[0].image,
    ]
  }
}

# --- Cloud SQL Proxy sidecar for Cloud Run ---
# Cloud Run connects via the VPC connector to Cloud SQL's private IP
# No sidecar needed — Cloud Run + VPC + private IP Cloud SQL is the recommended pattern

# --- VPC Access Connector (for Cloud Run -> Cloud SQL) ---
resource "google_vpc_access_connector" "main" {
  name          = "biomech-vpc-conn"
  region        = var.region
  network       = var.vpc_name
  ip_cidr_range = "10.8.0.0/28"
  machine_type  = "e2-micro"
  min_instances = 2
  max_instances = 10
}

# --- IAM: Allow unauthenticated invocations ---
resource "google_cloud_run_v2_service_iam_member" "api_public" {
  name     = google_cloud_run_v2_service.api.name
  location = google_cloud_run_v2_service.api.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloud_run_v2_service_iam_member" "ml_public" {
  name     = google_cloud_run_v2_service.ml.name
  location = google_cloud_run_v2_service.ml.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# --- Internal API Key secret ---
resource "random_password" "internal_api_key" {
  length  = 32
  special = false
}

resource "google_secret_manager_secret" "internal_api_key" {
  secret_id = "internal-api-key"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "internal_api_key" {
  secret      = google_secret_manager_secret.internal_api_key.id
  secret_data = random_password.internal_api_key.result
}

# --- Firebase service account key placeholder ---
# Place the actual Firebase service account JSON at:
#   gcloud secrets create firebase-service-account --data-file=./firebase-key.json
# resource "google_secret_manager_secret" "firebase_sa" {
#   secret_id = "firebase-service-account"
#   replication { auto {} }
# }

# --- Outputs ---
output "api_url" {
  value = google_cloud_run_v2_service.api.uri
}

output "ml_url" {
  value = google_cloud_run_v2_service.ml.uri
}

output "db_instance_name" {
  value = google_sql_database_instance.main.name
}

output "db_password" {
  value     = random_password.db_password.result
  sensitive = true
}

output "internal_api_key" {
  value     = random_password.internal_api_key.result
  sensitive = true
}
