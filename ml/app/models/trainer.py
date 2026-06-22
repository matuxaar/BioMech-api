import os
import json

import numpy as np
import tensorflow as tf
from sklearn.preprocessing import LabelEncoder, StandardScaler
from sklearn.model_selection import train_test_split

from app.config import settings
from app.services.features import sliding_window


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
    labels: np.ndarray | None,
    job_id: str,
    n_classes: int = 5,
) -> tuple[str, float]:
    features = sliding_window(data, settings.window_size, settings.window_size // 2)

    if labels is None:
        labels = np.zeros(len(features), dtype=np.int32)

    min_len = min(len(features), len(labels))
    features, labels = features[:min_len], labels[:min_len]

    scaler = StandardScaler()
    features = scaler.fit_transform(features)

    X_train, X_test, y_train, y_test = train_test_split(
        features, labels, test_size=0.2, random_state=42
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

    os.makedirs(settings.models_dir, exist_ok=True)
    model_path = os.path.join(settings.models_dir, f"{job_id}.h5")
    model.save(model_path)

    scaler_path = os.path.join(settings.models_dir, f"{job_id}_scaler.json")
    with open(scaler_path, "w") as f:
        json.dump(
            {
                "mean": scaler.mean_.tolist(),
                "scale": scaler.scale_.tolist(),
                "n_features": features.shape[1],
                "n_classes": n_classes,
            },
            f,
        )

    return model_path, float(accuracy)
