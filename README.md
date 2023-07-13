# Работа компьютерного клуба

## Входные данные

необходимо заполнить input.txt необходимой информацией для работы приложения:
- **<количество столов в компьютерном клубе>** количество столов в виде целого положительного числа.

- **<время начала работы> <время окончания работы>** задается время начала и окончания работы компьютерного клуба, разделенные пробелом.

- **<стоимость часа в компьютерном клубе>** стоимость часа в компьютерном клубе в виде целого положительного числа.

- **<время события 1> <идентификатор события 1> <тело события 1>** далее перечисляются события в соотвествии с форматом.

## Запуск приложения

Необходимо собрать образ Docker, выполнив команду

```sh
docker build -t myapp .
```
После завершения сборки образа Docker возможен запуск программы, при передачи пути к входному файлу в качестве аргумента:
```sh
docker run -v input.txt -e INPUT_FILE=/app/input.txt myapp
```

## Результат работы

На первой строке выводится время начала работы.

Далее перечислены все события, произошедшие за рабочий день (входящие и исходящие), каждое на отдельной строке.

После списка событий на отдельной строке выводится время окончания работы.

Для каждого стола на отдельной строке выведены через пробел следующие параметры: Номер стола, Выручка за день и Время, которое он был занят в течение рабочего дня.

Последней строкой выводится общая выручка клуба за день.
