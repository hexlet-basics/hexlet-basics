# samber/oops: что попадает в Sentry и slog

**Вопрос.** Доходят ли стек и атрибуты `samber/oops` (v1.23.2) до событий
`sentry-go` (v0.50.0) и до записей `slog`? Нужен ли самописный адаптер и какого
он размера? Где вызывается `sentry.Init` и с какими опциями? (Тикет
[#881](https://github.com/hexlet-basics/hexlet-basics/issues/881), карта #880.)

**Дата.** 2026-10-17. **Статус.** Только исследование, продуктовый код не менялся.

**Источники.** Исходники закреплённых версий в module cache
(`~/go/pkg/mod/github.com/samber/oops@v1.23.2` — далее `oops/`,
`~/go/pkg/mod/github.com/getsentry/sentry-go@v0.50.0` — далее `sentry/`), файлы
репозитория, запуск. Сторонние статьи не использовались. Раздел «Проверено
запуском» опирается на выполненный код (Go 1.27.2, приложение в конце); раздел
«Прочитано в исходниках» — только на чтение кода, запуском не подтверждён, и это
помечено.

---

## Короткий вывод

1. **Стек доходит до Sentry без адаптера** (проверено запуском): `sentry-go` сам
   находит у `OopsError` метод `StackFrames()`. Но кадры в событии **искажены**
   (см. п. 3).
2. **Атрибуты `oops` в Sentry не доходят вообще** (проверено): `domain`, `tags`,
   `code`, `With(...)`, `hint` в событии отсутствуют. Остаётся только текст
   цепочки сообщений. В `slog` доходят **все** (проверено), стек и
   атрибуты вложены в группу `err`.
3. **Дефект кадров в Sentry** (проверено запуском): на обычной оптимизированной
   сборке Sentry показывает строку на 1 меньше (`repo:20` вместо `repo:21`) и
   лишние кадры уже инлайнёных функций (`db`, `errors.New`), которых в стеке в
   момент `Wrapf` не было. С `-gcflags=-l` кадры верные. Механизм (повторное
   разрешение уже разрешённых PC через `runtime.CallersFrames`) — гипотеза по
   коду, подтверждена только поведением.
4. **Нужен небольшой адаптер**, ~40 строк в одном месте: `BeforeSend` в
   `telemetry.NewSentryClient` (покрывает и HTTP, и River — оба пути идут через
   этот клиент). Прототип запущен и работает (ниже). Кадры чинить можно отдельно
   (опционально, см. §4).
5. `sentry.Init` в проекте **не вызывается**: используется `sentry.NewClient`,
   `AttachStacktrace` не задан (по умолчанию `false`) — для исключений это не
   влияет, см. §1.

---

## 1. Где инициализируется Sentry (прочитано в исходниках)

- `internal/telemetry/telemetry.go:22-27`: `sentry.NewClient(sentry.ClientOptions{Dsn, Environment, Release})`.
  Других опций нет: ни `AttachStacktrace`, ни `BeforeSend`, ни `Debug`.
  `sentry.Init` в `cmd/` и `internal/` не вызывается (поиск по `*.go`; клиент
  строится через DI, `internal/di/common_package.go:38`).
- Хаб на запрос: `internal/telemetry/http.go:18-25` (`NewHub` + `SetHubOnContext`
  + `sentryhttp` с `Repanic`).
- Две точки отправки исключений: `internal/handlers/errors.go:172`
  (`hub.CaptureException(err)` внутри `WithScope`, теги `trace_id`, `span_id`,
  `http.*` задаются руками в `report()`, строки 143-173) и
  `internal/jobs/errors.go:65` (`sentryClient.CaptureException(err, &EventHint{...}, scope)`).
- `AttachStacktrace` по документации поля: «attach stacktraces to pure capture
  message calls» (`sentry/client.go:139-141`); использован только в
  `EventFromMessage` (`sentry/client.go:897`). Для `CaptureException` стек берётся
  из самой ошибки (§2), так что отсутствие флага проблемой не является.
- `BeforeSend` вызывается на каждое событие с `hint.OriginalException`
  (`sentry/client.go:187-189`, `sentry/interfaces.go:730`).
- Логирование: `internal/logging/logging.go:12-17` — `tint.NewTextHandler`.
  Handler без `ReplaceAttr`, поэтому `slog.LogValuer` ошибки раскрывается как есть.
  Моста slog → Sentry (`sentryslog`) нет.
- Использование `oops` в проекте: только `oops.Wrapf` (80 вызовов) и
  `oops.Errorf` (14), без `In`/`Tags`/`With`/`Code` (поиск по `*.go`). То есть
  сейчас в `OopsError` нет атрибутов вообще, кроме `trace` и стека; все выводы
  про атрибуты — для будущего использования builder-API.

## 2. Как библиотеки стыкуются (прочитано в исходниках)

- `sentry` извлекает стек рефлексией: ищет метод `StackFrames`, затем
  `StackTrace`, затем `GetStackTracer` (`sentry/stacktrace.go:76-102`), а PC
  берёт из элементов среза по полю `PC`/`ProgramCounter`
  (`sentry/stacktrace.go:105-130`). `OopsError` имеет
  `StackFrames() []runtime.Frame` (`oops/error.go:544-559`), у `runtime.Frame`
  есть поле `PC`, так что стек находится. (Вопрос в тикете называл `Stacktrace()
  string` — он для Sentry не используется, достаточно `StackFrames()`.)
- Цепочку `sentry` обходит по `Unwrap()` (`sentry/exception.go:40-138`),
  по одному `Exception` на звено; `OopsError.Unwrap` — `oops/error.go:144`.
- В событие уходят только `Value`, `Type`, `Stacktrace`, `Mechanism`
  (`sentry/exception.go:92-96`). Интерфейса, через который ошибка могла бы
  отдать теги/контекст, у `sentry-go` нет — читать `OopsError` нечему.
- `OopsError.LogValue()` (`oops/error.go:585-696`) строит группу атрибутов:
  `message`, `err`, `code`, `time`, `duration`, `domain`, `tags`, `trace`,
  `hint`, `public`, `owner`, группы `context`, `user`, `tenant`, дамп
  request/response, `stacktrace`, `sources` (если не выключен
  `oops.SourceFragmentsHidden`; по умолчанию скрыт, `oops/error.go:22-24, 686`).
  Приёмник значения, поэтому `slog` вызывает его для `"err", err`.

## 3. Проверено запуском

Скрэтч-тест (приложение A) в репозитории, не закоммичен: oops-ошибка из
`oops.Wrapf(errors.New("db down"), ...)`, а также с
`oops.In("billing").Tags("payment","retry").Code("E42").With("user_id", 7).Hint("retry later").Wrapf(...)`;
отправка через `sentry.NewHub(client).CaptureException` с
`sentry.MockTransport` (как в `internal/handlers/errors_test.go:52`); `slog` —
JSON-handler и `logging.New` проекта (tint). Это тот же вызов
`CaptureException`, что в `report()`; сам `report()` я не запускал.

**Sentry-событие (`AttachStacktrace` false и true — результат одинаков):**

| Звено | `Type` | `Value` | кадров |
|---|---|---|---|
| 0 | `*errors.errorString` | `db down` | 0 |
| 1 | `oops.OopsError` | `load course go: db down` | 4 (ожидалось 2) |
| 2 (rich) | `oops.OopsError` | `charge: load course go: db down` | 2 |

- `event.Tags` — пустой. `Contexts` содержит только `device`, `os`, `runtime`,
  `trace`. Ни `domain`, ни `tags`, ни `code`, ни `user_id`, ни `hint`.
- Кадры звена 1, обычная сборка (кадры в порядке Sentry, от старого к новому):
  `TestFrames r2_test.go:15`, `repo r_test.go:20`, `db r_test.go:17`,
  `errors.New errors.go:65`. С `-gcflags=all=-l`: `TestFrames:15`, `repo:21`.
  `slog` для той же ошибки (`oops`-овский собственный разбор) даёт верное
  `repo():21 → TestR()`. То есть oops хранит верный стек, а Sentry
  пересчитывает PC заново и получает сдвиг на инлайненных вызовах.

**slog (JSON и tint):** у простого `Wrapf` в группе `err`: `err`, `message`,
`stacktrace` (формат `Oops: ... --- at file:line fn()`), `time`, `trace` (ULID
самого oops, **не** OTel `trace_id` — они разные; `report()` логирует
OTel-идентификаторы отдельно). У ошибки с builder-API дополнительно `code`,
`domain`, `hint`, `tags`, `context.user_id`. В tint это выглядит как
`err.domain=billing err.tags="[payment retry]" err.context.user_id=7 err.stacktrace="..."`
в одной строке (путь в `stacktrace` абсолютный, с домашним каталогом сборки).

**Прототип адаптера** (приложение A, `enrich`): `BeforeSend` берёт
`errors.As(hint.OriginalException, &oops.OopsError{})` и копирует `Domain()`,
`Tags()`, `Code()`, `Context()`, `Hint()`. Результат запуска:
`tags=map[oops.code:E42 oops.domain:billing oops.tag.0:payment oops.tag.1:retry]`,
`contexts.oops=map[hint:retry later user_id:7]`.

## 4. Вывод: адаптер и его размер

- **Нужен ли.** Для Sentry — да, если решено использовать builder-API `oops`
  (`In/Tags/With/Code/Hint`): иначе эти данные пропадают в Sentry, оставаясь
  только в логах. Если проект остаётся на `Wrapf`/`Errorf` (как сейчас), то
  адаптер даёт только цепочку и стек; атрибутов нет, и адаптер не нужен.
- **Для slog — не нужен**: `LogValuer` уже работает со стандартным handler'ом.
- **Размер.** ~40 строк: функция `BeforeSend(*Event, *EventHint) *Event` в
  `internal/telemetry`, подключаемая одной строкой в `ClientOptions`
  (`telemetry.go:23`). Одна точка покрывает HTTP и River
  (`handlers/errors.go:172`, `jobs/errors.go:65`), хэндлеры не меняются.
  Решения за дизайном: плоские теги Sentry (`oops.domain`, `oops.code`) против
  контекста; атрибуты из `With(...)` лучше класть в контекст, не в теги
  (теги индексируются, у них лимит длины и кардинальность); `public`/`user`/
  `request` не копировать без решения о PII.
- **Кадры.** Дефект из §3 можно чинить в том же `BeforeSend`:
  заново собрать `Stacktrace` звена из
  `oe.StackFrames()` через `sentry.NewFrame(runtime.Frame)` (~15 строк) и
  подставить в `event.Exception`. Этот вариант я **не запускал**; перед
  принятием его надо проверить. Альтернатива — принять сдвиг (стек сопоставим
  по функциям, строки на 1 меньше) и зафиксировать это в ADR.
- **Не проверено.** Поведение на `oops.Join`/`Recover`, реальный `report()`
  end-to-end, отправка в настоящий Sentry (только `MockTransport`), другие версии
  Go/сборки. Гипотеза про инлайнинг подтверждена одним сравнением
  «с `-l` / без `-l`».

## Приложение A. Воспроизведение

Файлы лежали в `internal/zzresearch/` (удалены перед коммитом); команда:
`go test ./internal/zzresearch/ -run 'TestR|TestFrames|TestAdapter' -v -count=1`,
вариант без инлайнинга: добавить `-gcflags=all=-l`. Суть:

```go
func repo() error { return oops.Wrapf(db(), "load course %s", "go") }
func rich() error {
	return oops.In("billing").Tags("payment", "retry").Code("E42").
		With("user_id", 7).Hint("retry later").Wrapf(repo(), "charge")
}
// Sentry: client := sentry.NewClient({Dsn, Transport: &sentry.MockTransport{}, BeforeSend: enrich})
//         sentry.NewHub(client, sentry.NewScope()).CaptureException(rich())
// slog:   slog.New(slog.NewJSONHandler(&buf, nil)).Error("request failed", "err", rich())
func enrich(ev *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	var oe oops.OopsError
	err, _ := hint.OriginalException.(error)
	if !errors.As(err, &oe) { return ev }
	// ev.Tags["oops.domain"]=oe.Domain(); oe.Tags(); oe.Code();
	// ev.Contexts["oops"]=sentry.Context{oe.Context()..., "hint": oe.Hint()}
	return ev
}
```
