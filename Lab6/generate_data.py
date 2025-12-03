"""
Генератор тестовых данных для Лабораторной работы №6
Автор: Игнатьев Дмитрий Сергеевич, ИСТмд-11

Генерирует два файла с участниками олимпиад для проверки всех сценариев:
- Студент только в одной олимпиаде (получает "д") 
- Студент в двух олимпиадах без призов (получает "3")
- Студент в двух олимпиадах с одним призом (получает "4")
- Студент в двух олимпиадах с двумя призами (получает "5")
"""

import random

# Списки фамилий для генерации
surnames = [
    "Иванов", "Петров", "Сидоров", "Кузнецов", "Смирнов",
    "Попов", "Васильев", "Новиков", "Федоров", "Волков",
    "Лебедев", "Козлов", "Соколов", "Егоров", "Орлов",
    "Андреев", "Назаров", "Зайцев", "Павлов", "Макаров",
    "Игнатьев", "Степанов", "Романов", "Виноградов", "Богданов",
    "Филиппов", "Комаров", "Дмитриев", "Алексеев", "Антонов",
    "Тимофеев", "Николаев", "Максимов", "Артемьев", "Григорьев",
    "Яковлев", "Михайлов", "Давыдов", "Герасимов", "Тарасов",
    "Белов", "Поляков", "Медведев", "Борисов", "Абрамов",
    "Власов", "Голубев", "Беляев", "Тихонов", "Фомин",
]

def generate_test_data():
    """Генерирует тестовые данные, покрывающие все сценарии"""
    
    file1_lines = []
    file2_lines = []
    
    # Сценарий 1: Только в первой олимпиаде (должно быть "д")
    for i in range(20):
        surname = surnames[i]
        prize = random.choice(['', '1', '2', '3'])  # Может быть с призом или без
        line = f"{surname} {prize}".strip() if prize else surname
        file1_lines.append(line)
    
    # Сценарий 2: Только во второй олимпиаде (должно быть "д")
    for i in range(20, 30):
        surname = surnames[i]
        prize = random.choice(['', '1', '2', '3'])
        line = f"{surname} {prize}".strip() if prize else surname
        file2_lines.append(line)
    
    # Сценарий 3: В обеих олимпиадах без призов (должно быть "3")
    for i in range(30, 40):
        surname = surnames[i]
        file1_lines.append(surname)
        file2_lines.append(surname)
    
    # Сценарий 4: В обеих олимпиадах с одним призом (должно быть "4")
    for i in range(40, 45):
        surname = surnames[i]
        prize = random.choice(['1', '2', '3'])
        file1_lines.append(f"{surname} {prize}")
        file2_lines.append(surname)
    
    for i in range(45, 50):
        surname = surnames[i]
        prize = random.choice(['1', '2', '3'])
        file1_lines.append(surname)
        file2_lines.append(f"{surname} {prize}")
    
    # Сценарий 5: В обеих олимпиадах с двумя призами (должно быть "5")
    remaining = surnames[50:]
    for surname in remaining[:min(30, len(remaining))]:
        prize1 = random.choice(['1', '2', '3'])
        prize2 = random.choice(['1', '2', '3'])
        file1_lines.append(f"{surname} {prize1}")
        file2_lines.append(f"{surname} {prize2}")
    
    # Перемешиваем для реалистичности
    random.shuffle(file1_lines)
    random.shuffle(file2_lines)
    
    # Записываем в файлы
    with open('olymp1.txt', 'w', encoding='utf-8') as f:
        f.write('\n'.join(file1_lines))
    
    with open('olymp2.txt', 'w', encoding='utf-8') as f:
        f.write('\n'.join(file2_lines))
    
    print("✅ Тестовые файлы сгенерированы:")
    print(f"   olymp1.txt - {len(file1_lines)} записей")
    print(f"   olymp2.txt - {len(file2_lines)} записей")
    print(f"\nОжидаемые сценарии:")
    print(f"   'd' (только одна олимпиада): ~50 студентов")
    print(f"   '3' (две без призов): ~10 студентов")
    print(f"   '4' (две с одним призом): ~10 студентов")
    print(f"   '5' (две с двумя призами): ~30 студентов")

if __name__ == "__main__":
    generate_test_data()
