from pydantic import BaseModel, Field
from datetime import datetime
from typing import Optional

from app.config import settings


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
    samples: list[EMGSample] = Field(..., min_length=1, max_length=settings.max_predict_batch)


class PredictResponse(BaseModel):
    predictions: list[str]


class StreamSample(BaseModel):
    channel_1: float = Field(..., ge=-10, le=10, alias="ch1")
    channel_2: float = Field(..., ge=-10, le=10, alias="ch2")
    channel_3: float = Field(..., ge=-10, le=10, alias="ch3")
    channel_4: float = Field(..., ge=-10, le=10, alias="ch4")
    channel_5: float = Field(..., ge=-10, le=10, alias="ch5")
    channel_6: float = Field(..., ge=-10, le=10, alias="ch6")
    channel_7: float = Field(..., ge=-10, le=10, alias="ch7")
    channel_8: float = Field(..., ge=-10, le=10, alias="ch8")

    class Config:
        populate_by_name = True


class StreamPredictRequest(BaseModel):
    samples: list[StreamSample] = Field(..., min_length=1, max_length=settings.max_predict_batch)


class StreamPredictResponse(BaseModel):
    predictions: list[str]
