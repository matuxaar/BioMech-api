"""Tests for pydantic schemas."""

import pytest
from pydantic import ValidationError
from app.models.schemas import (
    EMGSample,
    StreamSample,
    PredictRequest,
    PredictResponse,
    TrainRequest,
    TrainResponse,
    StreamPredictRequest,
    StreamPredictResponse,
)
from datetime import datetime, timezone


class TestEMGSample:
    def test_valid(self):
        s = EMGSample(
            id="s1",
            session_id="ses1",
            timestamp=datetime.now(timezone.utc),
            channel_1=0.5,
            channel_2=-0.3,
            channel_3=0.0,
            channel_4=0.1,
            channel_5=-0.2,
            channel_6=0.4,
            channel_7=0.0,
            channel_8=-0.1,
        )
        assert s.channel_1 == 0.5

    def test_out_of_range(self):
        with pytest.raises(ValidationError):
            EMGSample(
                id="s1",
                session_id="ses1",
                timestamp=datetime.now(timezone.utc),
                channel_1=999,
                channel_2=0,
                channel_3=0,
                channel_4=0,
                channel_5=0,
                channel_6=0,
                channel_7=0,
                channel_8=0,
            )


class TestStreamSample:
    def test_alias(self):
        s = StreamSample(ch1=0.5, ch2=0.1, ch3=0, ch4=0, ch5=0, ch6=0, ch7=0, ch8=0)
        assert s.channel_1 == 0.5

    def test_full_names(self):
        s = StreamSample(channel_1=0.5, channel_2=0.1, channel_3=0, channel_4=0, channel_5=0, channel_6=0, channel_7=0, channel_8=0)
        assert s.channel_1 == 0.5


class TestPredictRequest:
    def test_min_samples(self):
        with pytest.raises(ValidationError):
            PredictRequest(samples=[])

    def test_valid(self):
        s = EMGSample(
            id="s1", session_id="ses1",
            timestamp=datetime.now(timezone.utc),
            channel_1=0, channel_2=0, channel_3=0, channel_4=0,
            channel_5=0, channel_6=0, channel_7=0, channel_8=0,
        )
        req = PredictRequest(samples=[s])
        assert len(req.samples) == 1


class TestTrainRequest:
    def test_valid(self):
        req = TrainRequest(job_id="job1", session_ids=["ses1", "ses2"])
        assert req.job_id == "job1"

    def test_empty_session_ids(self):
        with pytest.raises(ValidationError):
            TrainRequest(job_id="job1", session_ids=[])


class TestResponses:
    def test_predict_response(self):
        r = PredictResponse(predictions=["fist", "open"])
        assert len(r.predictions) == 2

    def test_stream_response(self):
        r = StreamPredictResponse(predictions=["rest"])
        assert r.predictions[0] == "rest"

    def test_train_response(self):
        r = TrainResponse(job_id="j1", status="completed", accuracy=0.95, model_path="/tmp/model.h5")
        assert r.accuracy == 0.95
