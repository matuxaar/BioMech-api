import json
import logging
import os
import threading

import numpy as np
import tensorflow as tf
from sklearn.preprocessing import StandardScaler

from app.models.schemas import EMGSample
from app.services.features import extract_features
from app.config import settings, GESTURE_LABELS
from app.services.storage import download_model

logger = logging.getLogger(__name__)


class Predictor:
    def __init__(self, model_dir: str = settings.models_dir):
        self.model_dir = model_dir
        self._model = None
        self._scaler = None
        self._n_classes = settings.n_classes
        self._fallback = True
        self._lock = threading.Lock()
        self._gcs_lock = threading.Lock()
        self._gcs_client = None

    def _get_gcs_client(self):
        with self._gcs_lock:
            if self._gcs_client is None:
                from google.cloud import storage
                self._gcs_client = storage.Client()
            return self._gcs_client

    def _load_latest(self):
        if settings.gcs_bucket:
            client = self._get_gcs_client()
            bucket = client.bucket(settings.gcs_bucket)
            blobs = list(bucket.list_blobs(prefix=settings.gcs_models_prefix + "/"))
            h5_blobs = [b for b in blobs if b.name.endswith(".h5")]
            if not h5_blobs:
                logger.warning("no model found in GCS bucket %s", settings.gcs_bucket)
                self._fallback = True
                return
            latest = max(h5_blobs, key=lambda b: b.updated)
            model_local = download_model(f"gs://{settings.gcs_bucket}/{latest.name}")
            root, _ = os.path.splitext(model_local)
            scaler_local = root + "_scaler.json"
            if not os.path.exists(scaler_local):
                scaler_name = latest.name.replace(".h5", "_scaler.json", 1)
                scaler_local = download_model(f"gs://{settings.gcs_bucket}/{scaler_name}")
        else:
            h5_files = [f for f in os.listdir(self.model_dir) if f.endswith(".h5")]
            if not h5_files:
                logger.warning("no model found in local dir %s", self.model_dir)
                self._fallback = True
                return
            latest_name = max(h5_files, key=lambda f: os.path.getmtime(os.path.join(self.model_dir, f)))
            model_local = os.path.join(self.model_dir, latest_name)
            root, _ = os.path.splitext(model_local)
            scaler_local = root + "_scaler.json"

        self._fallback = False
        logger.info("loading model from %s", model_local)
        self._model = tf.keras.models.load_model(model_local)
        if os.path.exists(scaler_local):
            with open(scaler_local) as f:
                scaler_data = json.load(f)
                self._scaler = StandardScaler()
                self._scaler.mean_ = np.array(scaler_data["mean"])
                self._scaler.scale_ = np.array(scaler_data["scale"])
                self._scaler.n_features_in_ = scaler_data["n_features"]
                self._n_classes = scaler_data.get("n_classes", settings.n_classes)
            logger.info("scaler loaded from %s", scaler_local)
        else:
            logger.warning("scaler not found at %s, predictions will skip scaling", scaler_local)

    def predict(self, samples: list[EMGSample]) -> list[str]:
        with self._lock:
            if self._model is None:
                self._load_latest()

        if self._fallback:
            logger.error("no model available, returning fallback predictions")
            raise RuntimeError("model not available - no trained model found")

        data = np.array(
            [
                [
                    s.channel_1,
                    s.channel_2,
                    s.channel_3,
                    s.channel_4,
                    s.channel_5,
                    s.channel_6,
                    s.channel_7,
                    s.channel_8,
                ]
                for s in samples
            ],
            dtype=np.float32,
        )

        if np.any(np.isnan(data)):
            raise ValueError("samples contain NaN values")

        if len(data) >= settings.window_size:
            features = extract_features(data[-settings.window_size:]).reshape(1, -1)
        elif len(data) >= 1:
            padded = np.zeros((settings.window_size, settings.n_channels), dtype=np.float32)
            padded[-len(data):] = data
            features = extract_features(padded).reshape(1, -1)
        else:
            raise ValueError("at least one sample required for prediction")

        if self._scaler is not None:
            features = self._scaler.transform(features)

        predictions = self._model.predict(features, verbose=0)
        predicted_indices = np.argmax(predictions, axis=1)

        return [
            GESTURE_LABELS[idx] if idx < len(GESTURE_LABELS) else "unknown"
            for idx in predicted_indices
        ]
