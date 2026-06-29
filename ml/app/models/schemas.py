from pydantic import BaseModel, Field
from datetime import datetime
from typing import Optional


class TrainRequest(BaseModel):
    job_id: str = Field(..., min_length=1)
    session_ids: list[str] = Field(..., min_length=1)


class TrainResponse(BaseModel):
    job_id: str
    status: str
    accuracy: Optional[float] = 0.0
    model_path: str = ""


class EMGSample(BaseModel):
    id: str
    session_id: str
    timestamp: datetime
    channel_1: float = Field(..., ge=-10, le=10)
    channel_2: float = Field(..., ge=-10, le=10)
    channel_3: float = Field(..., ge=-10, le=10)
    channel_4: float = Field(..., ge=-10, le=10)
    channel_5: float = Field(..., ge=-10, le=10)
    channel_6: float = Field(..., ge=-10, le=10)
    channel_7: float = Field(..., ge=-10, le=10)
    channel_8: float = Field(..., ge=-10, le=10)
    metadata: str = ""


class PredictRequest(BaseModel):
    samples: list[EMGSample] = Field(..., min_length=1)


class PredictResponse(BaseModel):
    predictions: list[str]


class StreamSample(BaseModel):
    ch1: float = Field(..., ge=-10, le=10)
    ch2: float = Field(..., ge=-10, le=10)
    ch3: float = Field(..., ge=-10, le=10)
    ch4: float = Field(..., ge=-10, le=10)
    ch5: float = Field(..., ge=-10, le=10)
    ch6: float = Field(..., ge=-10, le=10)
    ch7: float = Field(..., ge=-10, le=10)
    ch8: float = Field(..., ge=-10, le=10)


class StreamPredictRequest(BaseModel):
    samples: list[StreamSample] = Field(..., min_length=1)


class StreamPredictResponse(BaseModel):
    predictions: list[str]
