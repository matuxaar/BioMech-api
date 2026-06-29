"""Tests for config."""

from app.config import Settings


class TestSettings:
    def test_defaults(self):
        s = Settings()
        assert s.n_channels == 8
        assert s.window_size == 100
        assert s.n_classes == 5
        assert s.max_predict_batch == 10000

    def test_custom_values(self, monkeypatch):
        monkeypatch.setenv("N_CHANNELS", "16")
        monkeypatch.setenv("WINDOW_SIZE", "200")
        s = Settings()
        assert s.n_channels == 16
        assert s.window_size == 200
