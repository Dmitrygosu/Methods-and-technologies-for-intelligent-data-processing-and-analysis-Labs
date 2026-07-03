# Общение с локальной моделью через Ollama.
# Модель meteo-sovetnik собрана из qwen2.5:3b через Modelfile (см. папку models),
# поэтому системный промпт здесь не нужен — он уже зашит в модель.

import json
import os

import requests

OLLAMA_URL = os.environ.get("OLLAMA_URL", "http://localhost:11434")
MODEL = os.environ.get("OLLAMA_MODEL", "meteo-sovetnik")

# схема ответа: Ollama умеет принуждать модель к валидному JSON
# по заданной json-схеме (structured outputs), чем и пользуемся
ADVICE_SCHEMA = {
    "type": "object",
    "properties": {
        "summary": {"type": "string"},
        "clothing": {"type": "array", "items": {"type": "string"}},
        "umbrella": {"type": "boolean"},
        "umbrella_reason": {"type": "string"},
        "tips": {"type": "array", "items": {"type": "string"}},
        "comfort": {"type": "integer", "minimum": 1, "maximum": 10},
    },
    "required": ["summary", "clothing", "umbrella", "umbrella_reason",
                 "tips", "comfort"],
}


def build_prompt(city, wx):
    now, today = wx["now"], wx["today"]
    return (
        f"Город: {city}.\n"
        f"Сейчас: {now['description']}, {now['temp']}°C "
        f"(ощущается как {now['feels_like']}°C), "
        f"влажность {now['humidity']}%, ветер {now['wind']} км/ч "
        f"(порывы до {now['gusts']} км/ч), осадки {now['precipitation']} мм.\n"
        f"Прогноз на сегодня: {today['description']}, "
        f"от {today['t_min']}°C до {today['t_max']}°C, "
        f"вероятность осадков {today['rain_prob']}%, "
        f"сумма осадков {today['rain_sum']} мм, "
        f"ветер до {today['wind_max']} км/ч, УФ-индекс {today['uv_max']}.\n"
        f"Составь совет по одежде на сегодня."
    )


def umbrella_needed(wx):
    # то же правило, что зашито в системный промпт модели
    return wx["today"]["rain_prob"] > 40 or wx["today"]["rain_sum"] > 1.0


# маленькая модель изредка предлагает в списке одежды нижнее бельё,
# зонт или странные длинные фразы — такие пункты бракуем
BAD_CLOTHING = ("бель", "трус", "лиф", "пижам", "купальн", "зонт")


def clothing_item_ok(item):
    low = item.lower()
    return not any(b in low for b in BAD_CLOTHING) and len(item.split()) <= 8


def clothing_ok(items):
    return 2 <= len(items) <= 6 and all(clothing_item_ok(i) for i in items)


def get_advice(city, wx):
    # ответ модели проверяется: решение про зонт должно совпадать с правилом,
    # список одежды — быть осмысленным; иначе повторный запрос строже
    need = umbrella_needed(wx)
    advice = None
    for temp in (0.3, 0.1):
        r = requests.post(f"{OLLAMA_URL}/api/chat", json={
            "model": MODEL,
            "messages": [{"role": "user", "content": build_prompt(city, wx)}],
            "format": ADVICE_SCHEMA,
            "stream": False,
            "keep_alive": "2h",
            "options": {"temperature": temp},
        }, timeout=180)
        r.raise_for_status()
        advice = json.loads(r.json()["message"]["content"])
        if (advice["umbrella"] == need
                and len(advice["summary"].split()) >= 4
                and clothing_ok(advice["clothing"])):
            return advice

    # модель дважды ошиблась — чиним ответ по данным
    advice["umbrella"] = need
    advice["umbrella_reason"] = (
        f"вероятность осадков {wx['today']['rain_prob']}%, "
        f"сумма {wx['today']['rain_sum']} мм")
    advice["clothing"] = [it for it in advice["clothing"]
                          if clothing_item_ok(it)] or ["одежда по сезону"]
    return advice


def warmup():
    # первая генерация после старта Ollama занимает под минуту —
    # модель поднимается в видеопамять; греем её заранее
    try:
        requests.post(f"{OLLAMA_URL}/api/generate", json={
            "model": MODEL, "prompt": "ок", "stream": False,
            "keep_alive": "2h", "options": {"num_predict": 1},
        }, timeout=180)
    except requests.RequestException:
        pass
