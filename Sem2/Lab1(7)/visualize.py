import json
import matplotlib.pyplot as plt
import numpy as np

def visualize():
    try:
        with open('results.json', 'r') as f:
            data = json.load(f)
    except FileNotFoundError:
        print("results.json не найден")
        return

    algos = sorted(list(set(d['algo'] for d in data)))
    ns = sorted(list(set(d['n'] for d in data)))

    # График времени
    plt.figure(figsize=(12, 7))
    for algo in algos:
        # Берем только успешные попытки для времени
        subset = [d for d in data if d['algo'] == algo and d['success']]
        if not subset:
            continue
        xs = [d['n'] for d in subset]
        ys = [d['time_ms'] for d in subset]
        plt.plot(xs, ys, marker='o', label=algo)
    
    plt.title('Сравнение времени выполнения (мс)')
    plt.xlabel('Размер доски (N)')
    plt.ylabel('Время (мс)')
    plt.legend()
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.savefig('time_comparison.png')
    plt.close()

    # График итераций
    plt.figure(figsize=(12, 7))
    for algo in algos:
        subset = [d for d in data if d['algo'] == algo]
        xs = [d['n'] for d in subset]
        ys = [d['iterations'] for d in subset]
        # Если неуспешно - рисуем другим стилем или помечаем
        plt.plot(xs, ys, marker='s', label=algo)
    
    plt.title('Сравнение количества итераций')
    plt.xlabel('Размер доски (N)')
    plt.ylabel('Итерации (шкала логарифмическая)')
    plt.yscale('log') # Логарифмическая шкала для итераций
    plt.legend()
    plt.grid(True, which="both", linestyle='--', alpha=0.5)
    plt.savefig('iter_comparison.png')
    plt.close()
    
    print("Графики обновлены: time_comparison.png, iter_comparison.png")

if __name__ == "__main__":
    visualize()
