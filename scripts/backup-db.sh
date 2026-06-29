#!/usr/bin/env bash
set -euo pipefail

# PostgreSQL backup script for Desertacia
# Usage: ./backup-db.sh [bucket-name]
# Requires: gcloud, pg_dump, gsutil

BUCKET="${1:-biomech-db-backups}"
PROJECT="${PROJECT_ID:-biomech-217fe}"
INSTANCE="${CLOUD_SQL_INSTANCE:-biomech-db}"
REGION="${REGION:-us-central1}"
DB_NAME="${DB_NAME:-desertacia}"
TIMESTAMP="$(date -u +%Y%m%d-%H%M%S)"
DUMP_FILE="${DB_NAME}-${TIMESTAMP}.sql.gz"

echo "=== Starting backup: ${DUMP_FILE} ==="

echo "1. Connecting via Cloud SQL Auth Proxy..."
CLOUDSQL="${PROJECT}:${REGION}:${INSTANCE}"

echo "2. Dumping database..."
gcloud sql export sql "${INSTANCE}" "gs://${BUCKET}/${DUMP_FILE}" \
  --database="${DB_NAME}" \
  --project="${PROJECT}" \
  --quiet

echo "3. Verifying backup..."
gsutil stat "gs://${BUCKET}/${DUMP_FILE}"

echo "4. Cleaning up backups older than 30 days..."
gsutil ls "gs://${BUCKET}/${DB_NAME}-*.sql.gz" | while read -r OLD_FILE; do
  CREATE_DATE=$(gsutil stat "$OLD_FILE" | grep "Creation time:" | awk '{print $3}')
  if [[ -n "$CREATE_DATE" ]]; then
    AGE_DAYS=$(( ($(date +%s) - $(date -d "$CREATE_DATE" +%s)) / 86400 ))
    if [[ $AGE_DAYS -gt 30 ]]; then
      echo "  Removing old backup: $(basename "$OLD_FILE") (${AGE_DAYS} days)"
      gsutil rm "$OLD_FILE"
    fi
  fi
done

echo "=== Backup complete: gs://${BUCKET}/${DUMP_FILE} ==="
