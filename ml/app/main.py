from fastapi import FastAPI
from app.router import train, predict
from app.config import settings

app = FastAPI(
    title="Desertacia ML Service",
    description="EMG signal processing and muscle contraction classification",
    version="0.1.0",
)

app.include_router(train.router)
app.include_router(predict.router)


@app.get("/health")
def health():
    db_ok = False
    try:
        import psycopg2
        conn = psycopg2.connect(settings.database_url)
        conn.close()
        db_ok = True
    except Exception:
        pass
    return {"status": "ok" if db_ok else "degraded", "database": "connected" if db_ok else "unavailable"}
