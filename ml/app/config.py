import os
from dotenv import load_dotenv

load_dotenv()


class Settings:
    database_url: str = os.getenv(
        "DATABASE_URL",
        "postgresql://postgres:postgres@localhost:5432/desertacia",
    )
    models_dir: str = os.getenv("MODELS_DIR", "/app/models")
    n_channels: int = 8
    window_size: int = 100
    n_classes: int = 5


settings = Settings()
