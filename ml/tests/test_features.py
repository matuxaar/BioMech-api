"""Tests for feature extraction."""

import numpy as np
import pytest
from app.services.features import extract_features, sliding_window


def test_extract_features_basic():
    window = np.zeros((100, 8), dtype=np.float32)
    window[50:60] = 1.0
    features = extract_features(window)
    assert features.shape == (40,)
    assert features.dtype == np.float32
    assert features[0] > 0


def test_extract_features_nan_handling():
    window = np.zeros((100, 8), dtype=np.float32)
    window[0, 0] = np.nan
    features = extract_features(window)
    assert not np.any(np.isnan(features))


def test_sliding_window():
    data = np.random.randn(200, 8).astype(np.float32)
    features = sliding_window(data, 50, 25)
    assert len(features) == 7
    assert features.shape[1] == 40


def test_sliding_window_small_data():
    data = np.random.randn(30, 8).astype(np.float32)
    features = sliding_window(data, 50, 25)
    assert len(features) == 0


def test_sliding_window_nan_handling():
    data = np.random.randn(200, 8).astype(np.float32)
    data[0, 0] = np.nan
    features = sliding_window(data, 50, 25)
    assert not np.any(np.isnan(features))


@pytest.mark.parametrize("window_size,stride,expected", [
    (100, 50, 3),
    (100, 100, 1),
])
def test_sliding_window_params(window_size, stride, expected):
    data = np.random.randn(300, 8).astype(np.float32)
    features = sliding_window(data, window_size, stride)
    assert len(features) == expected
