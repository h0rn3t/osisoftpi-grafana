# Журнал змін (українською)

## 5.3.1 — 2026-05-26

Патч-реліз: посилення надійності WebSocket streaming з 5.3.0. Схема JSON datasource не змінювалась.

### Додано

- Обмежене перепідключення WebSocket з exponential backoff (до 5 спроб, 1–30 с)
- Очищення застарілих записів реєстру streaming-каналів
- Тести eligibility, відключення, `Dispose()` та race для subscribe/register

### Змінено

- Streaming увімкнюється лише перемикачем **Enable Streaming Support** (experimental не потрібен)
- Streaming лише для **PI points**; виключено summary, interpolated, recorded values, expression та AF attributes
- Захист `channelConstruct` mutex; підтримка шляху `ds/<uid>/<uuid>`
- При переповненні буфера — коротка блокуюча відправка (50 ms), потім drop з `dropped=true` у логах
- Файл `steam.go` перейменовано на `stream.go`

### Виправлено

- Підписники більше не зависають після падіння WebSocket read loop
- `Dispose()` коректно закриває WebSocket і очищає maps
- Вимкнення experimental більше не вимикає streaming у UI

### Міграція

- Увімкніть **Enable Streaming Support** у налаштуваннях datasource
- Панелі з streaming на summary/interpolated тощо автоматично повертаються до HTTP

## 5.3.0

- Підтримка PI Web API `streamsets/channel` (WebSocket streaming)
- Одне WebSocket-з’єднання на batch тегів
- Перемикачі streaming у datasource та Query Editor

Повний журнал англійською: [CHANGELOG.md](./CHANGELOG.md).
