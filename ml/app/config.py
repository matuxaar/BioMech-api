import os

GESTURE_LABELS = (
    "rest",
    "fist",
    "open",
    "pinch",
    "point",
)


class Settings:
    database_url: str = os.getenv("DATABASE_URL", "")
    models_dir: str = os.getenv("MODELS_DIR", "/tmp/models")
    gcs_bucket: str = os.getenv("GCS_BUCKET", "")
    gcs_models_prefix: str = os.getenv("GCS_MODELS_PREFIX", "models")
    n_channels: int = int(os.getenv("N_CHANNELS", "8"))
    window_size: int = int(os.getenv("WINDOW_SIZE", "100"))
    n_classes: int = int(os.getenv("N_CLASSES", "5"))
    max_predict_batch: int = int(os.getenv("MAX_PREDICT_BATCH", "10000"))
    backend_callback_url: str = os.getenv("BACKEND_CALLBACK_URL", "")
    backend_api_key: str = os.getenv("BACKEND_API_KEY", "")


settings = Settings()
