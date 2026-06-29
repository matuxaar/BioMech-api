# Terraform — Desertacia GCP Infrastructure

## Prerequisites

```bash
gcloud auth application-default login
gcloud config set project biomech-217fe
terraform init
```

## Usage

```bash
# Review planned changes
terraform plan

# Apply
terraform apply

# Destroy
terraform destroy
```

## Manual steps after `terraform apply`

1. **Firebase service account** — place the JSON secret:
   ```bash
   gcloud secrets create firebase-service-account \
     --data-file=./path/to/firebase-service-account.json
   ```

2. **Enable Cloud SQL Auth Proxy** (if not using VPC):
   Not needed if VPC connector + private IP Cloud SQL is used.

3. **Push initial Docker images** via Cloud Build:
   ```bash
   gcloud builds submit --config cloudbuild.yaml \
     --substitutions=_REGION=us-central1
   ```
