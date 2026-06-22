"""
EMG feature extraction utilities.

Extracts common time-domain features used in EMG pattern recognition:
- RMS (Root Mean Square)
- MAV (Mean Absolute Value)
- WL (Waveform Length)
- ZC (Zero Crossing)
- SSC (Slope Sign Change)
"""

import numpy as np


def extract_features(window: np.ndarray) -> np.ndarray:
    n_channels = window.shape[1]
    features = []

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
    features = []
    for start in range(0, len(data) - window_size + 1, stride):
        window = data[start : start + window_size]
        features.append(extract_features(window))
    return np.array(features)
