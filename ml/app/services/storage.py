import os
import tempfile

from google.cloud import storage

from app.config import settings

_client = None


def get_client() -> storage.Client:
    global _client
    if _client is None:
        _client = storage.Client()
    return _client


def save_model(local_path: str, job_id: str) -> str:
    if not settings.gcs_bucket:
        os.makedirs(settings.models_dir, exist_ok=True)
        dest = os.path.join(settings.models_dir, os.path.basename(local_path))
        os.replace(local_path, dest)
        return dest

    bucket = get_client().bucket(settings.gcs_bucket)
    blob_name = f"{settings.gcs_models_prefix}/{os.path.basename(local_path)}"
    blob = bucket.blob(blob_name)
    blob.upload_from_filename(local_path)
    return f"gs://{settings.gcs_bucket}/{blob_name}"


def download_model(remote_path: str) -> str:
    if remote_path.startswith("gs://"):
        parts = remote_path.replace("gs://", "").split("/", 1)
        bucket_name, blob_name = parts[0], parts[1]
        bucket = get_client().bucket(bucket_name)
        blob = bucket.blob(blob_name)

        os.makedirs(settings.models_dir, exist_ok=True)
        local_path = os.path.join(settings.models_dir, os.path.basename(blob_name))
        blob.download_to_filename(local_path)
        return local_path

    if os.path.exists(remote_path):
        return remote_path

    local_path = os.path.join(settings.models_dir, os.path.basename(remote_path))
    if os.path.exists(local_path):
        return local_path

    raise FileNotFoundError(f"model not found: {remote_path}")
