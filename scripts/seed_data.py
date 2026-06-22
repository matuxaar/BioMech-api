"""
Генератор тестовых EMG-данных.

Создаёт пользователя, устройства, EMG-сессии и синтетические сэмплы
для всех 5 классов жестов: rest, fist, open, pinch, point.

Использование:
    pip install requests numpy
    python scripts/seed_data.py [--api-url http://localhost:8080]
"""

import argparse
import math
import random
from datetime import datetime, timezone

import requests

API_URL = "http://localhost:8080"
USER_EMAIL = "test@desertacia.dev"
USER_PASSWORD = "testpassword123"

GESTURES = ["rest", "fist", "open", "pinch", "point"]
N_SAMPLES_PER_SESSION = 500
N_CHANNELS = 8
SAMPLING_RATE = 100  # Hz


def generate_emg_samples(gesture: str, n: int) -> list[dict]:
    samples = []
    t0 = datetime.now(timezone.utc).timestamp()

    # Параметры сигнала для каждого жеста
    params = {
        "rest":  {"amp": 0.01, "freq": 5,   "noise": 0.005},
        "fist":  {"amp": 0.8,  "freq": 50,  "noise": 0.05},
        "open":  {"amp": 0.5,  "freq": 30,  "noise": 0.04},
        "pinch": {"amp": 0.6,  "freq": 60,  "noise": 0.03},
        "point": {"amp": 0.4,  "freq": 40,  "noise": 0.04},
    }
    p = params[gesture]

    for i in range(n):
        t = i / SAMPLING_RATE
        channels = []
        for ch in range(N_CHANNELS):
            phase = random.uniform(0, 2 * math.pi)
            value = (
                p["amp"] * math.sin(2 * math.pi * p["freq"] * t + phase)
                + p["amp"] * 0.3 * math.sin(2 * math.pi * p["freq"] * 2 * t + phase)
                + random.gauss(0, p["noise"])
            )
            channels.append(round(value, 6))

        ts = datetime.fromtimestamp(t0 + i / SAMPLING_RATE, tz=timezone.utc).isoformat()
        samples.append({
            "timestamp": ts,
            "channel_1": channels[0],
            "channel_2": channels[1],
            "channel_3": channels[2],
            "channel_4": channels[3],
            "channel_5": channels[4],
            "channel_6": channels[5],
            "channel_7": channels[6],
            "channel_8": channels[7],
        })

    return samples


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--api-url", default=API_URL)
    args = parser.parse_args()
    base = args.api_url

    session = requests.Session()

    # --- 1. Регистрация / логин ---
    print("[1/6] Регистрация пользователя...")
    r = session.post(f"{base}/api/v1/auth/register", json={
        "email": USER_EMAIL,
        "password": USER_PASSWORD,
    })
    if r.status_code == 201:
        token = r.json()["access_token"]
        print("      Пользователь создан")
    elif r.status_code == 409:
        print("      Пользователь уже существует, логинимся...")
        r = session.post(f"{base}/api/v1/auth/login", json={
            "email": USER_EMAIL,
            "password": USER_PASSWORD,
        })
        token = r.json()["access_token"]
    else:
        print(f"      Ошибка: {r.status_code} {r.text}")
        return

    session.headers.update({"Authorization": f"Bearer {token}"})

    # --- 2. Создать устройство (сенсор) ---
    print("[2/6] Создание устройства-сенсора...")
    r = session.post(f"{base}/api/v1/devices", json={
        "type": "sensor",
        "name": "EMG Armband v1",
        "hw_version": "1.0.0",
    })
    if r.status_code == 201:
        sensor_id = r.json()["id"]
        print(f"      Сенсор: {sensor_id}")
    else:
        print(f"      Ошибка: {r.status_code} {r.text}")
        return

    # --- 3. Создать устройство (протез) ---
    print("[3/6] Создание устройства-протеза...")
    r = session.post(f"{base}/api/v1/devices", json={
        "type": "prosthetic",
        "name": "Bionic Hand Pro",
        "hw_version": "2.0.0",
    })
    if r.status_code == 201:
        prosthetic_id = r.json()["id"]
        print(f"      Протез: {prosthetic_id}")
    else:
        print(f"      Ошибка: {r.status_code} {r.text}")
        return

    # --- 4-5. Для каждого жеста: создать сессию + загрузить сэмплы ---
    print("[4/6] Создание EMG-сессий...")
    session_ids = {}
    for gesture in GESTURES:
        r = session.post(f"{base}/api/v1/emg/sessions", json={
            "device_id": sensor_id,
            "label": gesture,
        })
        if r.status_code == 201:
            sid = r.json()["id"]
            session_ids[gesture] = sid
            print(f"      {gesture}: {sid}")
        else:
            print(f"      Ошибка: {r.status_code} {r.text}")
            return

    print("[5/6] Генерация и загрузка сэмплов...")
    for gesture in GESTURES:
        sid = session_ids[gesture]
        samples = generate_emg_samples(gesture, N_SAMPLES_PER_SESSION)

        r = session.post(f"{base}/api/v1/emg/sessions/{sid}/samples/batch", json={
            "samples": samples,
        })
        if r.status_code == 201:
            print(f"      {gesture}: {len(samples)} сэмплов загружено")
        else:
            print(f"      Ошибка: {r.status_code} {r.text}")
            return

    # --- 6. Закончить сессии ---
    print("[6/6] Завершение сессий...")
    for gesture in GESTURES:
        sid = session_ids[gesture]
        r = session.post(f"{base}/api/v1/emg/sessions/{sid}/end")
        if r.status_code == 200:
            print(f"      {gesture}: завершена")
        else:
            print(f"      Ошибка: {r.status_code} {r.text}")

    print("\nГотово! Данные загружены:")
    print(f"  Пользователь: {USER_EMAIL} / {USER_PASSWORD}")
    print(f"  Сессий: {len(session_ids)}")
    print(f"  Сэмплов: {len(GESTURES) * N_SAMPLES_PER_SESSION}")
    print(f"\nТеперь можно запустить обучение:")
    print(f"  curl -X POST {base}/api/v1/training/jobs \\")
    print(f"    -H 'Authorization: Bearer {token}' \\")
    print(f"    -H 'Content-Type: application/json' \\")
    print(f"    -d '{{\"session_ids\": [\"{session_ids['fist']}\", \"{session_ids['open']}\"]}}'")


if __name__ == "__main__":
    main()
