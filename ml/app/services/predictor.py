import json
import os

import numpy as np
import tensorflow as tf

from app.models.schemas import EMGSample
from app.services.features import extract_features
from app.config import settings

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
        h5_files = [f for f in os.listdir(self.model_dir) if f.endswith(".h5")]
        if not h5_files:
            raise RuntimeError("no trained model found")
        latest = max(h5_files, key=lambda f: os.path.getmtime(os.path.join(self.model_dir, f)))
        model_path = os.path.join(self.model_dir, latest)

        scaler_path = model_path.replace(".h5", "_scaler.json")
        if not os.path.exists(scaler_path):
            raise RuntimeError(f"scaler not found for {latest}")

        self._model = tf.keras.models.load_model(model_path)

        with open(scaler_path) as f:
            scaler_data = json.load(f)
            self._n_classes = scaler_data["n_classes"]

        self._model_path = latest.replace(".h5", "")

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

        if len(data) < settings.window_size:
            predictions = self._model.predict(data, verbose=0)
        else:
            window = data[-settings.window_size:]
            features = extract_features(window).reshape(1, -1)
            predictions = self._model.predict(features, verbose=0)

        predicted_indices = np.argmax(predictions, axis=1)

        return [
            GESTURE_LABELS[idx] if idx < len(GESTURE_LABELS) else "unknown"
            for idx in predicted_indices
        ]
