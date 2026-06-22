from pydantic import BaseModel
from datetime import datetime
from typing import Optional


class TrainRequest(BaseModel):
    job_id: str
    session_ids: list[str]


class TrainResponse(BaseModel):
    job_id: str
    status: str
    accuracy: float
    model_path: str


class EMGSample(BaseModel):
    id: str
    session_id: str
    timestamp: datetime
    channel_1: float
    channel_2: float
    channel_3: float
    channel_4: float
    channel_5: float
    channel_6: float
    channel_7: float
    channel_8: float
    metadata: str = ""


class PredictRequest(BaseModel):
    samples: list[EMGSample]


class PredictResponse(BaseModel):
    predictions: list[str]
