import os
import sys
import threading

from flask import Flask, render_template, request, jsonify
import requests

import weather
import advisor

# при сборке в exe (pyinstaller) шаблоны лежат во временной папке
if getattr(sys, "frozen", False):
    app = Flask(__name__,
                template_folder=os.path.join(sys._MEIPASS, "templates"))
else:
    app = Flask(__name__)

# греем модель в фоне, чтобы первый запрос пользователя не ждал минуту
threading.Thread(target=advisor.warmup, daemon=True).start()


@app.route("/")
def index():
    return render_template("index.html")


@app.route("/api/advice")
def api_advice():
    city_name = (request.args.get("city") or "").strip()
    if not city_name:
        return jsonify({"error": "Укажите город"}), 400

    try:
        city = weather.find_city(city_name)
    except weather.CityNotFound:
        return jsonify({"error": f"Город «{city_name}» не найден"}), 404
    except requests.RequestException:
        return jsonify({"error": "Сервис геокодинга недоступен"}), 502

    try:
        wx = weather.get_forecast(city["lat"], city["lon"])
    except requests.RequestException:
        return jsonify({"error": "Сервис погоды недоступен"}), 502

    try:
        advice = advisor.get_advice(f"{city['name']} ({city['country']})", wx)
    except requests.RequestException:
        return jsonify({"error": "Локальная модель не отвечает. "
                                 "Проверьте, что Ollama запущена."}), 502

    return jsonify({
        "city": city,
        "weather": wx,
        "advice": advice,
    })


if __name__ == "__main__":
    app.run(host="127.0.0.1", port=5000, debug=False)
