"""
EMG feature extraction utilities.
"""

import logging

import numpy as np

logger = logging.getLogger(__name__)


def extract_features(window: np.ndarray) -> np.ndarray:
    n_channels = window.shape[1]
    features = []

    if np.any(np.isnan(window)):
        logger.warning("NaN values detected in feature extraction window")
        window = np.nan_to_num(window)

    for ch in range(n_channels):
        x = window[:, ch]

        rms = np.sqrt(np.mean(x**2))
        mav = np.mean(np.abs(x))
        wl = np.sum(np.abs(np.diff(x)))
        zc = np.sum(np.diff(np.signbit(x).astype(int)) != 0)
        ssc = np.sum(
            np.diff(np.signbit(np.diff(x)).astype(int)) != 0
        )

        features.extend([rms, mav, wl, zc, ssc])

    return np.array(features, dtype=np.float32)


def sliding_window(
    data: np.ndarray, window_size: int, stride: int
) -> np.ndarray:
    if np.any(np.isnan(data)):
        logger.warning("NaN values in sliding_window input, replacing with 0")
        data = np.nan_to_num(data)

    features = []
    for start in range(0, len(data) - window_size + 1, stride):
        window = data[start : start + window_size]
        features.append(extract_features(window))
    return np.array(features)
