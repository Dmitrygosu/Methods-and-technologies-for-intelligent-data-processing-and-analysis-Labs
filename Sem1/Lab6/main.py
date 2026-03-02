"""
Лабораторная работа №6
Обработка списков студентов-участников олимпиад

Автор: Игнатьев Дмитрий Сергеевич
Группа: ИСТмд-11

Задача: Объединить данные из двух файлов с участниками олимпиад и вывести
результат по следующей логике:
- 'д' - участвовал только в одной олимпиаде
- '3' - участвовал в двух олимпиадах и нигде не занял призовое место
- '4' - участвовал в двух олимпиадах и занял призовое место только в одной
- '5' - участвовал в двух олимпиадах и занял призовые места в обеих
"""

def read_olympiad_file(filename):
    """
    Читает файл с участниками олимпиады.
    
    Возвращает словарь: {фамилия: bool (есть ли приз)}
    Приз определяется наличием цифры 1, 2 или 3 после фамилии.
    """
    students = {}
    
    with open(filename, 'r', encoding='utf-8') as f:
        for line in f:
            line = line.strip()
            if not line:
                continue
            
            parts = line.split()
            surname = parts[0]
            
            # Проверяем наличие призового места (1, 2 или 3)
            has_prize = len(parts) > 1 and parts[1] in ['1', '2', '3']
            
            # Если студент уже есть, обновляем статус приза (OR логика)
            if surname in students:
                students[surname] = students[surname] or has_prize
            else:
                students[surname] = has_prize
    
    return students


def evaluate_students(file1, file2):
    """
    Обрабатывает два файла с участниками олимпиад и выводит результат.
    
    Логика оценивания:
    - Только в одной олимпиаде -> 'д'
    - В обеих, оба без призов -> '3'
    - В обеих, один с призом -> '4'
    - В обеих, оба с призами -> '5'
    """
    # Читаем оба файла
    olymp1 = read_olympiad_file(file1)
    olymp2 = read_olympiad_file(file2)
    
    # Объединяем все фамилии
    all_students = set(olymp1.keys()) | set(olymp2.keys())
    
    # Обрабатываем каждого студента
    results = []
    for surname in sorted(all_students):
        in_olymp1 = surname in olymp1
        in_olymp2 = surname in olymp2
        
        # Если только в одной олимпиаде
        if in_olymp1 and not in_olymp2:
            result = 'д'
        elif in_olymp2 and not in_olymp1:
            result = 'д'
        else:
            # В обеих олимпиадах - считаем призы
            prize1 = olymp1.get(surname, False)
            prize2 = olymp2.get(surname, False)
            prize_count = int(prize1) + int(prize2)
            
            if prize_count == 0:
                result = '3'
            elif prize_count == 1:
                result = '4'
            else:  # prize_count == 2
                result = '5'
        
        results.append((surname, result))
    
    return results


def main():
    """Главная функция - запускает обработку и выводит результат"""
    
    print("=" * 70)
    print("Лабораторная работа №6: Обработка списков участников олимпиад")
    print("Автор: Игнатьев Д.С., группа ИСТмд-11")
    print("=" * 70)
    print()
    
    # Обрабатываем данные
    results = evaluate_students('olymp1.txt', 'olymp2.txt')
    
    # Выводим результат
    print(f"{'Фамилия':<20} | Результат")
    print("-" * 35)
    
    for surname, grade in results:
        print(f"{surname:<20} | {grade}")
    
    # Статистика
    print()
    print("=" * 70)
    print("Статистика:")
    print("-" * 35)
    
    stats = {}
    for _, grade in results:
        stats[grade] = stats.get(grade, 0) + 1
    
    for grade in ['д', '3', '4', '5']:
        count = stats.get(grade, 0)
        print(f"  {grade} - {count} студентов")
    
    print(f"\nВсего обработано: {len(results)} студентов")
    print("=" * 70)


if __name__ == "__main__":
    main()
