## Приложение Этикеток
### git https://github.com/boszhansale/scale_printer
### написан go 1.20 и с использованием fyne
обновить версию го нельзя (дальше windows 7 не поддерживается)
### Приложение состоит из двух частей
- retail штучный
- weight весовой

### для компляции нужен докер и taskfile

## команды
- task r  компилирует  retail app
- task w  компилирует  weight app

## Структура проекта
-  cmd главные запускаемые файлы и все верстки здесь
-  internal основной код
-  config  модуль конфига
-  logger модуль лога который пишут в лог файл app.log
-  models модели для связки с сервером
-  repository для связи с сервером и сохранение данных в файл data.json
-  services/label/ генератор этикет
-  services/printer/  управление принтером
-  services/scale/  управление весов
-  services/zpl/ генератор zpl кода
-  vendor библиотеки go
-  .env конф файл 

