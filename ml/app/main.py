from fastapi import FastAPI
from app.router import train, predict

app = FastAPI(
    title="Desertacia ML Service",
    description="EMG signal processing and muscle contraction classification",
    version="0.1.0",
)

app.include_router(train.router)
app.include_router(predict.router)


@app.get("/health")
def health():
    return {"status": "ok"}
