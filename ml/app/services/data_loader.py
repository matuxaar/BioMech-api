import psycopg2
import numpy as np
from app.config import settings


def fetch_samples(session_ids: list[str]) -> np.ndarray:
    conn = psycopg2.connect(settings.database_url)
    cur = conn.cursor()

    placeholders = ",".join(["%s"] * len(session_ids))
    query = f"""
        SELECT channel_1, channel_2, channel_3, channel_4,
               channel_5, channel_6, channel_7, channel_8
        FROM emg_samples
        WHERE session_id IN ({placeholders})
        ORDER BY timestamp ASC
    """
    cur.execute(query, session_ids)
    rows = cur.fetchall()
    cur.close()
    conn.close()

    if not rows:
        return np.empty((0, settings.n_channels))

    return np.array(rows, dtype=np.float32)
