import os
import json
import tempfile

import numpy as np
import tensorflow as tf
from sklearn.preprocessing import StandardScaler
from sklearn.model_selection import train_test_split

from app.config import settings
from app.services.features import sliding_window
from app.services.storage import save_model


def build_model(input_dim: int, n_classes: int) -> tf.keras.Model:
    model = tf.keras.Sequential(
        [
            tf.keras.layers.Input(shape=(input_dim,)),
            tf.keras.layers.Dense(128, activation="relu"),
            tf.keras.layers.Dropout(0.3),
            tf.keras.layers.Dense(64, activation="relu"),
            tf.keras.layers.Dropout(0.2),
            tf.keras.layers.Dense(n_classes, activation="softmax"),
        ]
    )
    model.compile(
        optimizer="adam",
        loss="sparse_categorical_crossentropy",
        metrics=["accuracy"],
    )
    return model


def train(
    data: np.ndarray,
    labels: np.ndarray,
    job_id: str,
    n_classes: int = 5,
) -> tuple[str, float]:
    features = sliding_window(data, settings.window_size, settings.window_size // 2)

    if len(features) == 0:
        raise ValueError(f"not enough samples (need at least {settings.window_size})")

    # Align labels with sliding windows: each window gets the label at its midpoint
    y = np.array([
        labels[min(start + settings.window_size // 2, len(labels) - 1)]
        for start in range(0, len(data) - settings.window_size + 1, settings.window_size // 2)
    ], dtype=np.int32)

    min_len = min(len(features), len(y))
    features, y = features[:min_len], y[:min_len]

    scaler = StandardScaler()
    features = scaler.fit_transform(features)

    X_train, X_test, y_train, y_test = train_test_split(
        features, y, test_size=0.2, random_state=42
    )

    model = build_model(features.shape[1], n_classes)

    model.fit(
        X_train,
        y_train,
        epochs=50,
        batch_size=32,
        validation_data=(X_test, y_test),
        verbose=1,
    )

    loss, accuracy = model.evaluate(X_test, y_test, verbose=0)

    with tempfile.TemporaryDirectory() as tmpdir:
        model_path = os.path.join(tmpdir, f"{job_id}.h5")
        model.save(model_path)

        scaler_data = {
            "mean": scaler.mean_.tolist(),
            "scale": scaler.scale_.tolist(),
            "n_features": features.shape[1],
            "n_classes": n_classes,
        }
        scaler_path = os.path.join(tmpdir, f"{job_id}_scaler.json")
        with open(scaler_path, "w") as f:
            json.dump(scaler_data, f)

        remote_model_path = save_model(model_path, job_id)
        save_model(scaler_path, f"{job_id}_scaler.json")

    return remote_model_path, float(accuracy)
