"""Tests for training model."""

import numpy as np
import pytest
from app.models.trainer import build_model


class TestBuildModel:
    def test_build_model(self):
        model = build_model(input_dim=40, n_classes=5)
        assert model.input_shape == (None, 40)
        assert model.output_shape == (None, 5)

    def test_build_model_different_params(self):
        model = build_model(input_dim=80, n_classes=3)
        assert model.input_shape == (None, 80)
        assert model.output_shape == (None, 3)
