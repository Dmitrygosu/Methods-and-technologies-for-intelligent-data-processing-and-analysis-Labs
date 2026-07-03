# Клиент Open-Meteo: геокодинг города + текущая погода и прогноз на день.
# API публичный, ключ не нужен, лимитов для учебных целей хватает с запасом.

import time

import requests

GEO_URL = "https://geocoding-api.open-meteo.com/v1/search"

# основной и резервный серверы прогноза: у них одинаковый API, но живут они
# на разных адресах — в некоторых сетях (например, в универе) основной
# домен не резолвится, а резервный работает
FORECAST_URLS = [
    "https://api.open-meteo.com/v1/forecast",
    "https://historical-forecast-api.open-meteo.com/v1/forecast",
]

# расшифровка кодов погоды WMO, которые возвращает Open-Meteo
WMO_CODES = {
    0: "ясно", 1: "преимущественно ясно", 2: "переменная облачность",
    3: "пасмурно", 45: "туман", 48: "изморозь",
    51: "лёгкая морось", 53: "морось", 55: "сильная морось",
    56: "ледяная морось", 57: "сильная ледяная морось",
    61: "небольшой дождь", 63: "дождь", 65: "ливень",
    66: "ледяной дождь", 67: "сильный ледяной дождь",
    71: "небольшой снег", 73: "снег", 75: "сильный снегопад",
    77: "снежная крупа", 80: "кратковременный дождь",
    81: "ливневый дождь", 82: "сильный ливень",
    85: "снежный заряд", 86: "сильный снежный заряд",
    95: "гроза", 96: "гроза с градом", 99: "сильная гроза с градом",
}


class CityNotFound(Exception):
    pass


# координаты города не меняются, поэтому геокодинг кешируется без TTL
# на весь срок жизни процесса — экономит запрос при повторном вводе города
_geo_cache = {}

# прогноз меняется со временем, поэтому короткий TTL: спасает от повторного
# похода в API при случайных повторных кликах или демонстрации подряд
_FORECAST_TTL = 300  # 5 минут
_forecast_cache = {}


def find_city(name):
    key = name.strip().lower()
    if key in _geo_cache:
        return _geo_cache[key]

    r = requests.get(GEO_URL, params={
        "name": name, "count": 1, "language": "ru", "format": "json",
    }, timeout=10)
    r.raise_for_status()
    results = r.json().get("results")
    if not results:
        raise CityNotFound(name)
    c = results[0]
    city = {
        "name": c["name"],
        "country": c.get("country", ""),
        "lat": c["latitude"],
        "lon": c["longitude"],
    }
    _geo_cache[key] = city
    return city


def get_forecast(lat, lon):
    key = (round(lat, 2), round(lon, 2))
    cached = _forecast_cache.get(key)
    if cached and time.time() - cached[0] < _FORECAST_TTL:
        return cached[1]
    params = {
        "latitude": lat,
        "longitude": lon,
        "current": "temperature_2m,apparent_temperature,relative_humidity_2m,"
                   "precipitation,weather_code,wind_speed_10m,wind_gusts_10m",
        "daily": "temperature_2m_max,temperature_2m_min,"
                 "precipitation_probability_max,precipitation_sum,weather_code,"
                 "wind_speed_10m_max,uv_index_max",
        "forecast_days": 1,
        "timezone": "auto",
    }
    last_err = None
    # два прохода по списку серверов на случай разовых сетевых сбоев
    for url in list(FORECAST_URLS) * 2:
        try:
            r = requests.get(url, params=params, timeout=6)
            r.raise_for_status()
            # запоминаем рабочий сервер, чтобы не ждать таймаут каждый раз
            if FORECAST_URLS[0] != url:
                FORECAST_URLS.remove(url)
                FORECAST_URLS.insert(0, url)
            break
        except requests.RequestException as e:
            last_err = e
    else:
        raise last_err
    data = r.json()

    cur = data["current"]
    day = data["daily"]
    forecast = {
        "now": {
            "temp": cur["temperature_2m"],
            "feels_like": cur["apparent_temperature"],
            "humidity": cur["relative_humidity_2m"],
            "precipitation": cur["precipitation"],
            "wind": cur["wind_speed_10m"],
            "gusts": cur["wind_gusts_10m"],
            "description": WMO_CODES.get(cur["weather_code"], "н/д"),
        },
        "today": {
            "t_min": day["temperature_2m_min"][0],
            "t_max": day["temperature_2m_max"][0],
            "rain_prob": day["precipitation_probability_max"][0],
            "rain_sum": day["precipitation_sum"][0],
            "wind_max": day["wind_speed_10m_max"][0],
            "uv_max": day["uv_index_max"][0],
            "description": WMO_CODES.get(day["weather_code"][0], "н/д"),
        },
    }
    _forecast_cache[key] = (time.time(), forecast)
    return forecast
