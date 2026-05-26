# Datasource PI Web API для Grafana

Плагін дає доступ до даних **OSIsoft PI** та **PI AF** через **PI Web API**.

Поточна версія у цьому репозиторії: **5.3.0** (форк [h0rn3t/osisoftpi-grafana](https://github.com/h0rn3t/osisoftpi-grafana)).

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

### Варіант 2 — з цього форку (версія 5.3.0, streaming)

Офіційний каталог ще може не містити 5.3.0. Для збірки з форку:

#### 1. Клонування та збірка

```bash
git clone https://github.com/h0rn3t/osisoftpi-grafana.git
cd osisoftpi-grafana
git checkout v5.3.0   # або гілка master / feature/v5.3.0

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

У **Administration → Plugins** має з’явитися **OSIsoft-PI** (версія 5.3.0).

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
5. Для **5.3.0** — увімкніть **Enable Streaming Support**, якщо потрібні оновлення в реальному часі через WebSocket.

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

## Streaming (версія 5.3.0)

Потокові дані йдуть через **PI Web API streamsets/channel** (одне WebSocket-з’єднання на batch тегів).

1. У конфігурації datasource — **Enable Streaming Support**.
2. У редакторі запиту — увімкніть streaming для запиту (можна через dashboard variable).
3. Для швидких потоків у **Query options** панелі:
   - збільште **Max data points**;
   - **Min interval** — очікуваний інтервал оновлення (наприклад `4s`, `10s`).

Streaming не застосовується до запитів summary та інших несумісних режимів — плагін автоматично використовує звичайний HTTP-запит.

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
| Streaming не працює | **Enable Streaming Support**, Grafana ≥ 10.1, не summary-запит |
| «Обрізана» історія на графіку | Збільшити **Max data points** і **Min interval** |

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
