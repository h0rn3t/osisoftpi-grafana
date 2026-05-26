# Datasource PI Web API для Grafana

Плагін дає доступ до даних **OSIsoft PI** та **PI AF** через **PI Web API**.

Поточна версія у цьому репозиторії: **5.3.1** (форк [h0rn3t/osisoftpi-grafana](https://github.com/h0rn3t/osisoftpi-grafana)).

![Огляд системи](https://github.com/GridProtectionAlliance/osisoftpi-grafana/raw/master/docs/img/system_overview.png)

## Вимоги

- Grafana **≥ 10.1.0** (Grafana 8.x / 9.x не підтримуються з версії 5.0.0)
- Доступний **PI Web API** (рекомендовано окремий екземпляр, якщо використовуєте PI Vision / PI Coresight)
- Для самопідписаних сертифікатів PI Web API — увімкніть відповідну опцію TLS у налаштуваннях datasource

---

## Встановлення

### Варіант 1 — з каталогу Grafana (офіційна збірка)

Підходить, якщо потрібна стабільна версія з [Grafana Marketplace](https://grafana.com/grafana/plugins/gridprotectionalliance-osisoftpi-datasource/) (зазвичай **5.2.x**, без WebSocket streaming з 5.3.0):

```bash
grafana-cli plugins install gridprotectionalliance-osisoftpi-datasource
sudo systemctl restart grafana-server
```

Шлях після встановлення (Linux):

```text
/var/lib/grafana/plugins/gridprotectionalliance-osisoftpi-datasource/
```

### Варіант 2 — з цього форку (версія 5.3.x, streaming)

Офіційний каталог ще може не містити 5.3.x. Для збірки з форку:

#### 1. Клонування та збірка

```bash
git clone https://github.com/h0rn3t/osisoftpi-grafana.git
cd osisoftpi-grafana
git checkout v5.3.1   # або гілка master

# Frontend
npm install
npm run build

# Backend (потрібні Go та mage)
go install github.com/magefile/mage@latest
export PATH="$PATH:$(go env GOPATH)/bin"
mage -v build:linux      # для Linux / Docker
# mage -v build:darwin   # для macOS
# mage -v build:windows  # для Windows
```

Або одразу всі платформи:

```bash
npm run build:backend
```

Готова папка плагіна: **`dist/`** (містить `plugin.json`, `module.js`, `gpx_osipiwebapi_*`).

#### 2. Копіювання в Grafana

```bash
sudo mkdir -p /var/lib/grafana/plugins
sudo rm -rf /var/lib/grafana/plugins/gridprotectionalliance-osisoftpi-datasource
sudo cp -r dist /var/lib/grafana/plugins/gridprotectionalliance-osisoftpi-datasource
sudo chown -R grafana:grafana /var/lib/grafana/plugins/gridprotectionalliance-osisoftpi-datasource
```

> **macOS (Homebrew):**  
> `cp -r dist "$(brew --prefix)/var/lib/grafana/plugins/gridprotectionalliance-osisoftpi-datasource"`

#### 3. Дозвіл на непідписаний плагін

Збірка з вихідного коду не підписана Grafana. Додайте в `grafana.ini`:

```ini
[plugins]
allow_loading_unsigned_plugins = gridprotectionalliance-osisoftpi-datasource
```

Або змінну середовища:

```bash
export GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=gridprotectionalliance-osisoftpi-datasource
```

#### 4. Перезапуск Grafana

```bash
sudo systemctl restart grafana-server
# або: docker restart <ім'я-контейнера-grafana>
```

У **Administration → Plugins** має з’явитися **OSIsoft-PI** (версія 5.3.1).

### Варіант 3 — Grafana у Docker

Змонтуйте локальну папку `dist` як volume плагіна:

```yaml
volumes:
  - ./dist:/var/lib/grafana/plugins/gridprotectionalliance-osisoftpi-datasource:ro
environment:
  GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS: gridprotectionalliance-osisoftpi-datasource
```

Для локальної розробки в репозиторії:

```bash
npm run build
mage -v build:linux
docker compose up --build
```

Grafana буде на http://localhost:3000.

### Варіант 4 — оновлення вже встановленого плагіна

1. Зупиніть Grafana (або контейнер).
2. Збережіть резервну копію старої папки плагіна.
3. Зберіть нову версію (`git pull`, `npm run build`, `mage build:…`).
4. Замініть вміст  
   `/var/lib/grafana/plugins/gridprotectionalliance-osisoftpi-datasource/`  
   вмістом з `dist/`.
5. Перезапустіть Grafana.
6. Перевірте версію в **Connections → Data sources → OSIsoft-PI** або в списку плагінів.

Datasource-и в дашбордах зазвичай зберігаються; змінюється лише код плагіна.

---

## Налаштування datasource

1. **Connections → Data sources → Add data source** → **OSIsoft-PI**.
2. URL PI Web API, наприклад `https://pivision.company.com/piwebapi`.
3. Режим доступу: рекомендовано **Server (proxy)**.
4. За потреби — **Basic** auth у PI Web API та облікові дані в datasource.
5. Для **5.3.x** — увімкніть **Enable Streaming Support**, якщо потрібні оновлення в реальному часі через WebSocket.

Документація PI Web API: [docs.osisoft.com — PI Web API](https://docs.osisoft.com/bundle/pi-web-api).

---

## Запити до PI Asset Framework (AF)

![Елементи та атрибути](https://github.com/GridProtectionAlliance/osisoftpi-grafana/raw/master/docs/img/elements_and_attributes.png)

1. Перемикач **PI Point Search** має бути вимкнений.
2. У **Element** натисніть **Select AF Database** і оберіть базу.
3. **Select AF Element** — оберіть елемент AF (ієрархію за потреби).
   - Якщо сегмент UI не з’явився — натисніть **+** у **Attributes**, щоб оновити інтерфейс.
4. У **Attributes** натисніть **+** і оберіть атрибут (можна фільтрувати введенням імені).

---

## Запити до PI Data Server (PI Points)

![Запит PI Point](https://github.com/GridProtectionAlliance/osisoftpi-grafana/raw/master/docs/img/pi_point_query.png)

1. Увімкніть **PI Point Search**.
2. **Select Dataserver** — оберіть сервер даних.
3. У **PI Points** натисніть **+** і введіть ім’я тега (регістр не важливий: `sinusoid` = `SINUSOID`).

---

## Streaming (версія 5.3.x)

Плагін може оновлювати графіки **без повного перезавантаження панелі**: спочатку завантажується історія (HTTP), далі — live-точки через streaming.

### Транспорти в плагіні

Налаштування **Streaming transport** у datasource:

| Режим | Ендпоінт PI Web API | Опис |
|-------|---------------------|------|
| **WebSocket (channel)** — за замовчуванням | `streamsets/channel` | Push-оновлення, одне WebSocket-з’єднання на batch тегів |
| **HTTP plot** | `streamsets/plot` | Періодичний GET (коли WebSocket недоступний або блокується проксі) |

### Де в PI Web API є WebSocket

У **PI Web API WebSocket використовується лише в Channels** (протокол `ws://` / `wss://`, не `https://`). Інші методи (`plot`, `value`, `recorded`) — звичайний HTTP.

| # | Шаблон URL (WebSocket) | Призначення |
|---|------------------------|-------------|
| 1 | `wss://{host}/piwebapi/streams/{webId}/channel` | Один stream (PI point або AF attribute з Data Reference) |
| 2 | `wss://{host}/piwebapi/streamsets/{webId}/channel` | Усі stream-атрибути Element / Event Frame / Attribute (батьківський webId) |
| 3 | `wss://{host}/piwebapi/streamsets/channel?webId={id1}&webId={id2}&...` | Довільний набір тегів (використовує цей плагін для кількох PI points) |

**Корисні параметри** (для всіх варіантів channel):

| Параметр | Опис |
|----------|------|
| `includeInitialValues=true` | Одразу після підключення надіслати поточне значення |
| `heartbeatRate={n}` | Порожні повідомлення для keep-alive (кратне інтервалу опиту PI) |

**Приклад** (один PI point, як у плагіні за замовчуванням):

```text
wss://mypiportal.example.com/piwebapi/streamsets/channel?webId=F1DPdfhknGLNe0exLjI6f9NlpAtVEDAAU...&includeInitialValues=true
```

**Діагностика активних підключень** (HTTP, не відкриває WebSocket):

```text
GET https://{host}/piwebapi/channels/instances
```

Документація AVEVA: [Channels](https://docs.aveva.com/bundle/pi-web-api-reference/page/help/topics/channels.html).

### Що не є WebSocket

| Ендпоінт | Протокол | Примітка |
|----------|----------|----------|
| `streamsets/plot`, `value`, `recorded`, `interpolated`, `summary` | HTTP GET | Історія та зрізи; плагін використовує `plot` для першого завантаження |
| **Stream Updates** (`/streams/{webId}/updates`, `/streamsets/updates`) | HTTP POST + GET | Інкрементальні зміни без WebSocket (окремий механізм PI Web API) |

Якщо потрібен саме URL на кшталт `streamsets/plot?startTime=...&intervals=3000` — це **не** live WebSocket; для streaming у плагіні оберіть **WebSocket (channel)** або режим **HTTP plot** (опит plot у фоні).

### Швидкий старт

1. **Connections → Data sources → OSIsoft-PI** → увімкніть **Enable Streaming Support**.
2. **Streaming transport** → `WebSocket (streamsets/channel)` (або `HTTP plot`, якщо WS недоступний).
3. У запиті панелі:
   - увімкніть **Is Pi Point?**;
   - увімкніть **Enable Streaming**;
   - вимкніть Summary, Interpolated, Recorded values, Expression.
4. У **Query options** панелі:
   - **Max data points** — для історії та (у режимі plot) кількість інтервалів;
   - **Min interval** — очікувана частота оновлень (наприклад `4s`, `10s`).

Streaming працює лише для **PI points** у режимі звичайного часового ряду (plot). Не застосовується до summary, interpolated, recorded values, expression та AF attributes.

### Переваги streaming та приклади

| Перевага | Без streaming | З streaming (channel) |
|----------|---------------|------------------------|
| Затримка | Залежить від **Refresh** дашборду (наприклад 30 с) | Оновлення при зміні значення в PI (секунди та менше) |
| Навантаження на PI | Повний HTTP-запит `plot` кожного refresh | Один WebSocket + лише змінені значення |
| Кілька тегів на панелі | Окремі HTTP batch-запити | Одне з’єднання `streamsets/channel?webId=...&webId=...` |

**Приклад 1 — оперативний моніторинг одного тега**

- Панель: Time series.
- PI Point: `RDC CENTRAL REGION.4.33651.VALUE` (або ваш шлях).
- Datasource: **Enable Streaming Support** + **Enable Streaming** у запиті.
- **Min interval**: `2s`, refresh дашборду можна поставити рідше (наприклад 1 хв) — графік оновлюватиметься між refresh завдяки WebSocket.

**Приклад 2 — кілька PI points на одній панелі**

- Додайте кілька запитів (Query A, B, C) з увімкненим streaming.
- Плагін згрупує webId в одне WebSocket-з’єднання (менше навантаження, ніж окремий channel на тег).

**Приклад 3 — увімкнення streaming змінною дашборду**

У datasource залиште **Enable Streaming** вимкненим у запиті, але вкажіть **Streaming variable**: `$streaming`

Створіть змінну дашборду:

- Ім’я: `streaming`
- Тип: Custom
- Значення: `true` / `false`

Коли `$streaming=true`, плагін увімкне streaming для запиту без редагування кожної панелі.

**Приклад 4 — сервер без WebSocket (проксі, стара ОС)**

- **Streaming transport** → `HTTP plot (streamsets/plot)`.
- Плагін періодично викликає той самий `plot`, що й для історії, і додає нові точки на графік.
- **Min interval** задає період опиту (1–30 с). Навантаження на PI вище, ніж у channel — використовуйте лише за потреби.

**Приклад 5 — історія + live на одному графіку**

1. Перше відкриття панелі: HTTP `streamsets/plot` за обраний часовий діапазон (повна історія).
2. Далі: channel надсилає лише нові значення (`IncludeDataOnly` у Grafana).
3. При зміні часового діапазону Grafana робить новий QueryData — реєструється новий streaming-канал.

### Обмеження PI Web API Channels

- Потрібен **PI Web API ≥ 2015 R3** і ОС із підтримкою WebSocket.
- Підключення через **`wss://`**, не `https://.../channel` (помилка «The protocol of the request is not supported»).
- Дуже велика кількість `webId` в одному URL може обмежуватися; альтернатива — кілька WebSocket або AF Data Pipe.
- Для проксі/firewall переконайтеся, що дозволено upgrade до WebSocket на шляху до PI Web API.

### Поведінка плагіна при перевантаженні

Якщо підписник не встигає обробляти оновлення, повідомлення можуть відкидатися після короткої спроби доставки (до 50 ms); у логах Grafana з’явиться попередження з `dropped=true`. Збільште **Min interval** або зменште кількість streamable-запитів на панелі.

---

## Змінні шаблону (template variables)

Підтримуються дочірні елементи AF. Потрібен JSON у запиті змінної, наприклад:

```json
{"path": "PISERVER\\DatabaseName\\ElementNameWithChildren"}
```

![Налаштування змінної](https://github.com/GridProtectionAlliance/osisoftpi-grafana/raw/master/docs/img/template_setup_1.png)

---

## Event Frames та анотації

Datasource підтримує **AF Event Frames** як анотації на графіках.

![Event Frame](https://github.com/GridProtectionAlliance/osisoftpi-grafana/raw/master/docs/img/event_frame.png)

У запиті анотації вкажіть категорію Event Frame; підтримуються колір і regex для імені.

![Анотації](https://github.com/GridProtectionAlliance/osisoftpi-grafana/raw/master/docs/img/annotations.png)

---

## Усунення несправностей

| Проблема | Що перевірити |
|----------|----------------|
| Плагін не видно в Grafana | `allow_loading_unsigned_plugins`, права на папку, перезапуск |
| `plugin unavailable` / помилка backend | Наявність `gpx_osipiwebapi_linux_amd64` (або для вашої ОС) у папці плагіна |
| 401 / 403 до PI | Basic auth, права користувача в PI Web API |
| Streaming не працює | **Enable Streaming Support**, Grafana ≥ 10.1, PI point, **Enable Streaming**, не summary/interpolated/recorded |
| WebSocket «protocol not supported» | URL має бути `wss://`, не `https://`; перевірте версію PI Web API та ОС |
| Streaming не працює за проксі | Спробуйте **Streaming transport → HTTP plot** або налаштуйте WebSocket upgrade на проксі |
| «Обрізана» історія на графіку | Збільшити **Max data points** і **Min interval** |
| Повідомлення `dropped=true` у логах | Зменшити частоту тегів / збільшити **Min interval** |

Логи Grafana (Linux): `journalctl -u grafana-server -f` або логи контейнера Docker.

---

## Розробка та тести

```bash
npm run dev          # frontend у режимі watch
go test ./pkg/plugin/...
```

---

## Торгові марки

Усі назви продуктів, логотипи та бренди належать відповідним власникам.  
**OSIsoft**, логотип OSIsoft та **PI Web API** — торгові марки [AVEVA Group plc](https://www.aveva.com/en/legal/osisoft-terms-and-conditions/).

Оригінальний проєкт: [GridProtectionAlliance/osisoftpi-grafana](https://github.com/GridProtectionAlliance/osisoftpi-grafana).
