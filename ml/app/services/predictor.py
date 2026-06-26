import json
import os

import numpy as np
import tensorflow as tf

from app.models.schemas import EMGSample
from app.services.features import extract_features
from app.config import settings
from app.services.storage import download_model

GESTURE_LABELS = [
    "rest",
    "fist",
    "open",
    "pinch",
    "point",
]


class Predictor:
    def __init__(self, model_dir: str = settings.models_dir):
        self.model_dir = model_dir
        self._model = None
        self._scaler = None
        self._n_classes = settings.n_classes

    def _load_latest(self):
        if settings.gcs_bucket:
            from google.cloud import storage
            client = storage.Client()
            bucket = client.bucket(settings.gcs_bucket)
            blobs = list(bucket.list_blobs(prefix=settings.gcs_models_prefix + "/"))
            h5_blobs = [b for b in blobs if b.name.endswith(".h5")]
            if not h5_blobs:
                raise RuntimeError("no trained model found in GCS")
            latest = max(h5_blobs, key=lambda b: b.updated)
            latest_name = latest.name.split("/")[-1]
            model_local = download_model(f"gs://{settings.gcs_bucket}/{latest.name}")
            scaler_name = latest.name.replace(".h5", "_scaler.json")
            scaler_local = download_model(f"gs://{settings.gcs_bucket}/{scaler_name}")
        else:
            h5_files = [f for f in os.listdir(self.model_dir) if f.endswith(".h5")]
            if not h5_files:
                raise RuntimeError("no trained model found")
            latest_name = max(h5_files, key=lambda f: os.path.getmtime(os.path.join(self.model_dir, f)))
            model_local = os.path.join(self.model_dir, latest_name)
            scaler_local = model_local.replace(".h5", "_scaler.json")
            if not os.path.exists(scaler_local):
                raise RuntimeError(f"scaler not found for {latest_name}")

        self._model = tf.keras.models.load_model(model_local)
        with open(scaler_local) as f:
            scaler_data = json.load(f)
            self._n_classes = scaler_data.get("n_classes", settings.n_classes)

    def get_model(self):
        if self._model is None:
            self._load_latest()
        return self._model

    def predict(self, samples: list[EMGSample]) -> list[str]:
        if self._model is None:
            self._load_latest()

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

        if len(data) >= settings.window_size:
            features = extract_features(data[-settings.window_size:]).reshape(1, -1)
        elif len(data) >= 1:
            # Pad with zeros to reach window_size, then extract features
            padded = np.zeros((settings.window_size, settings.n_channels), dtype=np.float32)
            padded[-len(data):] = data
            features = extract_features(padded).reshape(1, -1)
        else:
            raise ValueError("at least one sample required for prediction")

        predictions = self._model.predict(features, verbose=0)
        predicted_indices = np.argmax(predictions, axis=1)

        return [
            GESTURE_LABELS[idx] if idx < len(GESTURE_LABELS) else "unknown"
            for idx in predicted_indices
        ]
